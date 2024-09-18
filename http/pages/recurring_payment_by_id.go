package pages

import (
	"net/http"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/auth"
	"github.com/deitrix/fin/web/page"
	"github.com/go-chi/chi/v5"
	"github.com/rickb777/date"
)

func RecurringPaymentByID(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var email string
		profile, ok := auth.ProfileFromContext(r.Context())
		if ok {
			email = profile["email"].(string)
		}
		id := chi.URLParam(r, "id")
		rp, err := store.RecurringPayment(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		since := date.Today()
		if q := r.URL.Query().Get("since"); q != "" {
			var err error
			since, err = date.ParseISO(q)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		payments := rp.PaymentsSinceN(since, 6)
		var loadMoreSince *date.Date
		if len(payments) == 6 {
			loadMoreSince = &payments[5].Date
			payments = payments[:5]
		}
		render(w, r, page.RecurringPayment(email, rp, payments, loadMoreSince))
	}
}
