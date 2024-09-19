package pages

import (
	"net/http"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/web/page"
	"github.com/go-chi/chi/v5"
	"github.com/rickb777/date"
)

func RecurringPayment(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		render(w, r, page.RecurringPayment(rp, payments, loadMoreSince))
	}
}

func RecurringPaymentUpdateForm(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rp, err := store.RecurringPayment(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render(w, r, page.RecurringPaymentForm(rp))
	}
}

func RecurringPaymentHandleUpdateForm(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rp, err := store.RecurringPayment(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rp.Enabled = r.Form.Get("enabled") == "on"
		rp.Debt = r.Form.Get("debt") == "on"
		if err := store.UpdateRecurringPayment(r.Context(), rp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		render(w, r, page.RecurringPaymentForm(rp))
	}
}

func RecurringPaymentDelete(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := store.DeleteRecurringPayment(r.Context(), chi.URLParam(r, "id")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
