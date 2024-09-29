package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/form"
	"github.com/deitrix/fin/pkg/pointer"
	"github.com/deitrix/fin/ui"
	"github.com/deitrix/fin/ui/page"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PaymentFormInput struct {
	fin.Payment
	ID           string
	FormPriority bool
	Amount       float64
	NewAccount   string
}

func (in *PaymentFormInput) Process() {
	in.Payment.Amount = int(in.Amount * 100)
	if in.ID != "" {
		in.Payment.ID = &in.ID
	}
	if in.NewAccount != "" {
		in.Payment.AccountID = in.NewAccount
	}
}

func PaymentFormFields(in *PaymentFormInput) form.Fields {
	return form.Fields{
		"id":           form.String(&in.ID),
		"description":  form.String(&in.Description),
		"date":         form.Time(&in.Date, "2006-01-02"),
		"amount":       form.Float(&in.Amount),
		"debt":         form.Bool(&in.Debt),
		"account":      form.String(&in.AccountID),
		"newAccount":   form.String(&in.NewAccount),
		"formPriority": form.Bool(&in.FormPriority),
	}
}

func PaymentForm(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := accountList(r.Context(), store)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var in PaymentFormInput
		if err := form.Decode(r.URL.Query(), PaymentFormFields(&in)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		in.Process()

		if in.FormPriority {
			// Pass the form data back to the

		} else if id := chi.URLParam(r, "id"); id != "" {
			var err error
			in.Payment, err = store.Payment(r.Context(), id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			in.Payment = fin.Payment{
				Date:   time.Now(),
				Amount: 100,
			}
			if len(accounts) != 0 {
				in.Payment.AccountID = accounts[0]
			}
		}

		ui.Render(w, r, page.PaymentForm(accounts, in.Payment))
	}
}

func PaymentHandleForm(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var in PaymentFormInput
		if err := form.Decode(r.Form, PaymentFormFields(&in)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		in.Process()

		var persist func(context.Context, fin.Payment) error
		if id := chi.URLParam(r, "id"); id != "" {
			in.Payment.ID = &id
			persist = store.UpdatePayment
		} else {
			in.Payment.ID = pointer.To(uuid.NewString())
			persist = store.CreatePayment
		}
		if err := persist(r.Context(), in.Payment); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func PaymentHandleDelete(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}
		if err := store.DeletePayment(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
