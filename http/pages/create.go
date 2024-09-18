package pages

import (
	"net/http"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/auth"
	"github.com/deitrix/fin/web/page"
	"github.com/google/uuid"
)

func Create(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var email string
		profile, ok := auth.ProfileFromContext(r.Context())
		if ok {
			email = profile["email"].(string)
		}
		render(w, r, page.RecurringPaymentCreate(email))
	}
}

func CreatePOST(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rp := fin.RecurringPayment{
			ID:      uuid.NewString(),
			Name:    r.Form.Get("name"),
			Enabled: true,
			Debt:    r.Form.Get("debt") == "on",
		}
		if err := store.CreateRecurringPayment(r.Context(), rp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
