package handlers

import (
	"net/http"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/ui"
	"github.com/deitrix/fin/ui/page"
	"github.com/go-chi/chi/v5"
)

func RecurringPayment(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		rp, err := store.RecurringPayment(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ui.Render(w, r, page.RecurringPayment(rp))
	}
}

func RecurringPaymentUpdateForm(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rp, err := store.RecurringPayment(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ui.Render(w, r, page.RecurringPaymentForm(rp))
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
		w.Header().Set("HX-Trigger", "reload")
		ui.Render(w, r, page.RecurringPaymentForm(rp))
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
