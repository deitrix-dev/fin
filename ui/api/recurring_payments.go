package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/iterx"
	"github.com/deitrix/fin/ui"
	"github.com/deitrix/fin/ui/components"
)

func RecurringPayments(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("recurringPaymentSearch")

		rps, err := iterx.CollectErr(fin.PageIter(fin.RecurringPaymentsQuery{}, 100, store.RecurringPayments))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if query != "" {
			var filtered []fin.RecurringPayment
			for _, rp := range rps {
				if strings.Contains(strings.ToLower(rp.Name), strings.ToLower(query)) {
					filtered = append(filtered, rp)
				}
			}
			rps = filtered
		}

		u, _ := url.Parse(r.Referer())
		q := u.Query()
		if query != "" {
			q.Set("recurringPaymentSearch", query)
		} else {
			q.Del("recurringPaymentSearch")
		}
		u.RawQuery = q.Encode()
		w.Header().Set("HX-Replace-URL", u.String())

		ui.Render(w, r, components.RecurringPayments(rps, query))
	}
}
