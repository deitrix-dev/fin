package api

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/form"
	"github.com/deitrix/fin/ui"
	"github.com/deitrix/fin/ui/components"
)

type RecurringPaymentsInputs struct {
	Search string
}

func RecurringPaymentsFields(in *RecurringPaymentsInputs) form.Fields {
	return form.Fields{
		"recurringPaymentSearch": form.String(&in.Search),
	}
}

func RecurringPayments(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in RecurringPaymentsInputs
		if err := form.Decode(refererQuery(r), RecurringPaymentsFields(&in)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rps, err := store.RecurringPayments(r.Context(), fin.RecurringPaymentFilter{Search: in.Search})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Replace-URL", hxReplaceURL(r, RecurringPaymentsFields(&in), "recurringPaymentSearch"))

		ui.Render(w, r, components.RecurringPayments(rps, in.Search))
	}
}

func refererQuery(r *http.Request) url.Values {
	q := r.URL.Query()
	if r.URL.Query().Get("referer") == "true" {
		ref, _ := url.Parse(r.Referer())
		q = ref.Query()
		for k, v := range r.URL.Query() {
			q[k] = v
		}
	}
	return q
}

func hxReplaceURL(r *http.Request, fields form.Fields, keys ...string) string {
	u, _ := url.Parse(r.Referer())
	q := u.Query()
	for k, v := range form.Encode(fields) {
		if slices.Contains(keys, k) {
			q[k] = v
		}
	}
	for k, v := range fields {
		if slices.Contains(keys, k) {
			if len(v.EncodeForm()) == 0 {
				q.Del(k)
			}
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}
