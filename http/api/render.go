package api

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

func render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	if err := component.Render(r.Context(), w); err != nil {
		slog.ErrorContext(r.Context(), "error rendering page", err)
	}
}
