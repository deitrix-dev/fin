package api

import (
	"net/http"
	"time"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/date"
	"github.com/deitrix/fin/ui"
	"github.com/deitrix/fin/ui/component"
)

func MonthSummaries(svc *fin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		summaries, err := svc.MonthSummaries(r.Context(), fin.MonthSummariesQuery{
			After: date.Midnight(time.Now()),
			Limit: 12,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ui.Render(w, r, component.MonthSummaries(summaries))
	}
}
