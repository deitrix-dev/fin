package handlers

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/deitrix/fin"
	"github.com/go-chi/chi/v5"
)

func ScheduleDelete(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		index, err := strconv.Atoi(chi.URLParam(r, "index"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rp, err := store.RecurringPayment(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if index < 0 || index >= len(rp.Schedules) {
			http.Error(w, "index out of bounds", http.StatusBadRequest)
			return
		}
		rp.Schedules = slices.Delete(rp.Schedules, index, index+1)
		if err := store.UpdateRecurringPayment(r.Context(), rp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/recurring-payments/"+rp.ID, http.StatusSeeOther)
	}
}
