package ui

import (
	"log/slog"
	"net/http"

	g "github.com/maragudk/gomponents"
)

func Render(w http.ResponseWriter, r *http.Request, node g.Node) {
	if err := node.Render(w); err != nil {
		slog.ErrorContext(r.Context(), "error rendering page", err)
	}
}
