package api

import (
	"iter"
	"net/http"
	"time"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/date"
	"github.com/deitrix/fin/pkg/form"
	"github.com/deitrix/fin/pkg/iterx"
	"github.com/deitrix/fin/pkg/pointer"
	"github.com/deitrix/fin/ui"
	"github.com/deitrix/fin/ui/components"
)

type PaymentsInputs struct {
	Filter             string
	Search             string
	Offset             uint
	RecurringPaymentID string
	Source             string
}

func PaymentsFields(in *PaymentsInputs) form.Fields {
	return form.Fields{
		"paymentFilter":    form.String(&in.Filter),
		"paymentSearch":    form.String(&in.Search),
		"offset":           form.Uint(&in.Offset),
		"recurringPayment": form.String(&in.RecurringPaymentID),
		"source":           form.String(&in.Source),
	}
}

func Payments(svc *fin.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in PaymentsInputs
		fields := PaymentsFields(&in)
		if err := form.Decode(refererQuery(r), fields); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var paymentIter iter.Seq2[fin.Payment, error]
		if in.RecurringPaymentID == "" {
			var query fin.PaymentsQuery
			query.Filter.Search = in.Search
			switch in.Filter {
			case "paymentsOnly":
				query.PaymentsOnly = true
			case "recurringPaymentsOnly":
				query.RecurringPaymentsOnly = true
			}
			paymentIter = svc.Payments(r.Context(), query)
		} else {
			rp, err := svc.RecurringPayment(r.Context(), in.RecurringPaymentID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			paymentIter = iterx.WithNilErr(rp.PaymentsSince(date.Midnight(time.Now())))
		}

		payments, err := iterx.CollectNErr(iterx.SkipErr(paymentIter, int(in.Offset)), 26)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var nextPage *uint
		if len(payments) == 26 {
			nextPage = pointer.To(in.Offset + 25)
			payments = payments[:25]
		}

		w.Header().Set("HX-Replace-URL", hxReplaceURL(r, fields, "paymentFilter", "paymentSearch"))

		ui.Render(w, r, components.Payments(components.PaymentsInputs{
			Header:      "Upcoming Payments",
			Payments:    payments,
			FetchURL:    "/api/payments?" + form.Encode(fields).Encode(),
			NextPage:    nextPage,
			Search:      in.RecurringPaymentID == "",
			Description: in.RecurringPaymentID == "",
			OOB:         in.Source != "",
			OOBSearch:   in.Source != "" && in.Source != "paymentSearch",
			OOBFilter:   in.Source != "" && in.Source != "paymentFilter",
			Exclude:     nil,
			Filter:      in.Filter,
			Query:       in.Search,
		}))
	}
}
