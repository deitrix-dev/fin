package api

import (
	"net/http"

	"github.com/deitrix/fin/auth"
	"github.com/deitrix/fin/web/components"
)

func HeaderUser(simulateUser string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile, ok := auth.ProfileFromContext(r.Context())
		if !ok {
			render(w, r, components.HeaderUser(simulateUser))
			return
		}
		render(w, r, components.HeaderUser(profile["email"].(string)))
	}
}
