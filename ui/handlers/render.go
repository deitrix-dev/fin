package handlers

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	g "github.com/maragudk/gomponents"
)

func render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	if err := component.Render(r.Context(), w); err != nil {
		slog.ErrorContext(r.Context(), "error rendering page", err)
	}
}

func renderGomp(w http.ResponseWriter, r *http.Request, node g.Node) {
	if err := node.Render(w); err != nil {
		slog.ErrorContext(r.Context(), "error rendering page", err)
	}
}
