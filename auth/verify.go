package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
)

type Profile map[string]any

type profileKey struct{}

func ProfileFromContext(ctx context.Context) (Profile, bool) {
	p, ok := ctx.Value(profileKey{}).(Profile)
	return p, ok
}

type Config struct {
	Enabled  bool   `json:"enabled"`
	ClientID string `json:"clientId"`
	Domain   string `json:"domain"`
}

// Verify is a middleware to verify a CF Access token
func Verify(conf Config) func(http.Handler) http.Handler {
	if !conf.Enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	certsURL := fmt.Sprintf("%s/cdn-cgi/access/certs", conf.Domain)
	config := &oidc.Config{ClientID: conf.ClientID}
	keySet := oidc.NewRemoteKeySet(context.Background(), certsURL)
	verifier := oidc.NewVerifier(conf.Domain, keySet, config)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Make sure that the incoming request has our token header
			//  Could also look in the cookies for CF_AUTHORIZATION
			accessJWT := r.Header.Get("Cf-Access-Jwt-Assertion")
			if accessJWT == "" {
				slog.ErrorContext(r.Context(), "no Cf-Access-Jwt-Assertion header found")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Verify the access token
			idToken, err := verifier.Verify(r.Context(), accessJWT)
			if err != nil {
				slog.ErrorContext(r.Context(), "invalid token", "err", err)
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			var profile Profile
			if err := idToken.Claims(&profile); err != nil {
				slog.ErrorContext(r.Context(), "error unmarshalling id token claims", "err", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), profileKey{}, profile)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
