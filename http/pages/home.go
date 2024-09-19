package pages

import (
	"net/http"
	"strings"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/pointer"
	"github.com/deitrix/fin/web/page"
	"github.com/rickb777/date"
)

func Home(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rps, err := store.RecurringPayments(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		payments := fin.PaymentsSinceN(rps, date.Today(), 11)
		var nextPayments *int
		if len(payments) == 11 {
			nextPayments = pointer.To(20)
			payments = payments[:10]
		}
		data := page.HomeData{
			RecurringPayments: rps,
			Payments:          payments,
			PaymentsState: page.PaymentsState{
				CurrentPage: 10,
			},
			NextPayments: nextPayments,
		}
		if q := r.URL.Query().Get("q"); q != "" {
			var filtered []fin.RecurringPayment
			for _, rp := range rps {
				if strings.Contains(strings.ToLower(rp.Name), strings.ToLower(q)) {
					filtered = append(filtered, rp)
				}
			}
			data.RecurringPayments = filtered
		}
		render(w, r, page.Home(data))
	}
}
