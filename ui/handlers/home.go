package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/iterx"
	"github.com/deitrix/fin/pkg/pointer"
	"github.com/deitrix/fin/web/page"
)

func Home(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rps, err := iterx.CollectErr(fin.PageIter(fin.RecurringPaymentsQuery{}, 100, store.RecurringPayments))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		payments := fin.PaymentsSinceN(rps, time.Now(), 11)
		var nextPayments *int
		if len(payments) == 11 {
			nextPayments = pointer.To(20)
			payments = payments[:10]
		}
		data := page.HomeData{
			RecurringPayments: rps,
			Payments:          payments,
			NextPayments:      nextPayments,
			Query:             r.URL.Query().Get("q"),
		}
		pushURL := r.URL.Query().Get("pushUrl") == "true"
		if data.Query != "" {
			var filtered []fin.RecurringPayment
			for _, rp := range rps {
				if strings.Contains(strings.ToLower(rp.Name), strings.ToLower(data.Query)) {
					filtered = append(filtered, rp)
				}
			}
			data.RecurringPayments = filtered
		}
		if pushURL {
			u := r.URL
			q := u.Query()
			q.Del("pushUrl")
			if q.Get("q") == "" {
				q.Del("q")
			}
			u.RawQuery = q.Encode()

			w.Header().Set("HX-Replace-URL", u.String())
		}
		render(w, r, page.Home(data))
	}
}
