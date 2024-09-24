package api

import (
	"fmt"
	"iter"
	"net/http"
	"strconv"
	"time"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/iterx"
	"github.com/deitrix/fin/pkg/pointer"
	"github.com/deitrix/fin/ui"
	"github.com/deitrix/fin/ui/components"
)

func Payments(store fin.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var offset int
		if s := r.URL.Query().Get("offset"); s != "" {
			var err error
			offset, err = strconv.Atoi(s)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		q := r.URL.Query().Get("q")
		recurringPaymentID := r.URL.Query().Get("recurringPayment")

		var paymentIter iter.Seq[fin.Payment]
		if recurringPaymentID == "" {
			rps, err := iterx.CollectErr(fin.PageIter(fin.RecurringPaymentsQuery{
				Filter: fin.RecurringPaymentFilter{
					Search: q,
				},
			}, 100, store.RecurringPayments))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			paymentIter = fin.PaymentsSince(rps, time.Now())
		} else {
			rp, err := store.RecurringPayment(r.Context(), recurringPaymentID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			paymentIter = rp.PaymentsSince(time.Now())
		}

		payments := iterx.CollectN(iterx.Skip(paymentIter, offset), 26)
		var nextPage *int
		if len(payments) == 26 {
			nextPage = pointer.To(offset + 25)
			payments = payments[:25]
		}

		fetchURL := fmt.Sprintf("/api/payments?oob=true&q=%s&offset=%d&recurringPayment=%s", q, offset+25, recurringPaymentID)
		ui.Render(w, r, components.Payments(components.PaymentsInputs{
			Header:      "Upcoming Payments",
			Payments:    payments,
			FetchURL:    fetchURL,
			NextPage:    nextPage,
			Search:      recurringPaymentID == "",
			Description: recurringPaymentID == "",
			OOB:         r.URL.Query().Get("oob") == "true",
		}))
	}
}
