package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/form"
	"github.com/deitrix/fin/ui"
	"github.com/deitrix/fin/ui/api"
	"github.com/deitrix/fin/ui/components"
	"github.com/go-chi/chi/v5"
)

func ScheduleForm(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rps := fin.PageIter(fin.RecurringPaymentsQuery{}, 100, store.RecurringPayments)
		accountsSet := make(map[string]struct{})
		for rp, err := range rps {
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			for _, sch := range rp.Schedules {
				accountsSet[sch.AccountID] = struct{}{}
			}
		}
		accounts := make([]string, 0, len(accountsSet))
		for acc := range accountsSet {
			accounts = append(accounts, acc)
		}
		sort.Strings(accounts)
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
		ui.Render(w, r, components.Layout("New schedule",
			components.ScheduleForm(accounts, rp, schedule, index),
		))
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
