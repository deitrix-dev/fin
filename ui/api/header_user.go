package api

import (
	"net/http"

	"github.com/deitrix/fin/auth"
	"github.com/deitrix/fin/ui"
	"github.com/deitrix/fin/ui/components"
)

func HeaderUser(simulateUser string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile, ok := auth.ProfileFromContext(r.Context())
		if !ok {
			if simulateUser != "" {
				ui.Render(w, r, components.HeaderUser(simulateUser))
			}
			return
		}
		ui.Render(w, r, components.HeaderUser(profile["email"].(string)))
	}
}
