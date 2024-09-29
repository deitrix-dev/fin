package handlers

import (
	"net/http"

	"github.com/deitrix/fin/ui"
	"github.com/deitrix/fin/ui/page"
)

func RecurringPayments(w http.ResponseWriter, r *http.Request) {
	ui.Render(w, r, page.RecurringPayments())
}
