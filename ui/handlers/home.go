package handlers

import (
	"net/http"

	"github.com/deitrix/fin/ui"
	"github.com/deitrix/fin/ui/components"
)

func Home(w http.ResponseWriter, r *http.Request) {
	ui.Render(w, r, components.Home(r.URL.Query().Get("recurringPaymentSearch")))
}
