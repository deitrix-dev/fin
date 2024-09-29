package handlers

import (
	"context"
	"net/http"
	"sort"
	"strconv"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/form"
	"github.com/deitrix/fin/ui"
	"github.com/deitrix/fin/ui/api"
	"github.com/deitrix/fin/ui/page"
	"github.com/go-chi/chi/v5"
)

func accountList(ctx context.Context, store fin.Store) ([]string, error) {
	rps, err := store.RecurringPayments(ctx, fin.RecurringPaymentFilter{})
	if err != nil {
		return nil, err
	}
	payments := fin.PageIter(ctx, fin.PaymentsQuery{}, 100, store.Payments)
	accountsSet := make(map[string]struct{})
	for _, rp := range rps {
		for _, sch := range rp.Schedules {
			accountsSet[sch.AccountID] = struct{}{}
		}
	}
	for p, err := range payments {
		if err != nil {
			return nil, err
		}
		accountsSet[p.AccountID] = struct{}{}
	}
	accounts := make([]string, 0, len(accountsSet))
	for acc := range accountsSet {
		accounts = append(accounts, acc)
	}
	sort.Strings(accounts)
	return accounts, nil
}

func ScheduleForm(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := accountList(r.Context(), store)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rp, err := store.RecurringPayment(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		index := -1
		if s := chi.URLParam(r, "index"); s != "" {
			var err error
			index, err = strconv.Atoi(s)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		var in api.ScheduleFormInput
		if err := form.Decode(r.URL.Query(), api.ScheduleForm(&in)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var schedule fin.PaymentSchedule
		if in.FormPriority {
			schedule = in.Schedule()
		} else if index > -1 {
			schedule = rp.Schedules[index]
			schedule.Repeat.Day = min(max(schedule.Repeat.Day, 1), 31)
			schedule.Repeat.Multiplier = max(schedule.Repeat.Multiplier, 1)
			if schedule.Repeat.Weekday == "" {
				schedule.Repeat.Weekday = fin.Monday
			}
		} else {
			schedule = fin.PaymentSchedule{
				Repeat: fin.Repeat{
					Every:      fin.Month,
					Day:        1,
					Weekday:    fin.Monday,
					Multiplier: 1,
				},
				Amount: 100,
			}
			if len(accounts) > 0 {
				schedule.AccountID = accounts[0]
			}
		}
		ui.Render(w, r, page.ScheduleForm(accounts, rp, schedule, index))
	}
}

func ScheduleHandleForm(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var in api.ScheduleFormInput
		if err := form.Decode(r.Form, api.ScheduleForm(&in)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rp, err := store.RecurringPayment(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		index := -1
		if s := chi.URLParam(r, "index"); s != "" {
			var err error
			index, err = strconv.Atoi(s)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		if index > -1 {
			rp.Schedules[index] = in.Schedule()
		} else {
			rp.Schedules = append(rp.Schedules, in.Schedule())
		}
		if err := store.UpdateRecurringPayment(r.Context(), rp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/recurring-payments/"+rp.ID, http.StatusSeeOther)
	}
}
