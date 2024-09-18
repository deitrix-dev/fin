package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/pointer"
	"github.com/deitrix/fin/web/page"
	"github.com/rickb777/date"
)

func Payments(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		size, err := strconv.Atoi(r.URL.Query().Get("size"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		q := r.URL.Query().Get("q")
		rps, err := store.RecurringPayments(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		payments := fin.PaymentsSinceNFilter(rps, date.Today(), size+1, func(payment fin.Payment) bool {
			if q == "" {
				return true
			}
			return strings.Contains(strings.ToLower(payment.RecurringPayment.Name), strings.ToLower(q))
		})
		var nextPage *int
		if len(payments) == size+1 {
			nextPage = pointer.To(size + 10)
			payments = payments[:size]
		}
		render(w, r, page.Payments(payments, page.PaymentsState{
			CurrentPage: size,
			Query:       q,
		}, nextPage))
	}
}
