package fin

import (
	"context"
	"iter"
	"time"

	"github.com/deitrix/fin/pkg/date"
	"github.com/deitrix/fin/pkg/iterx"
	"github.com/deitrix/fin/pkg/pointer"
)

type ServiceClient interface {
	Store
}

type Service struct {
	Store
}

func NewService(store Store) *Service {
	return &Service{Store: store}
}

func (s *Service) Payments(ctx context.Context, query PaymentsQuery) iter.Seq2[Payment, error] {
	return func(yield func(Payment, error) bool) {
		if query.Filter.After == nil {
			query.Filter.After = pointer.To(date.Midnight(time.Now()))
		}

		var iters []iter.Seq2[Payment, error]

		if !query.PaymentsOnly {
			rps, err := s.Store.RecurringPayments(ctx, RecurringPaymentFilter{Search: query.Filter.Search})
			if err != nil {
				yield(Payment{}, err)
				return
			}
			iters = append(iters, iterx.WithNilErr(PaymentsSince(rps, *query.Filter.After)))
		}

		if !query.RecurringPaymentsOnly {
			iters = append(iters, PageIter(ctx, query, 100, s.Store.Payments))
		}

		seq := iterx.JoinErr(Payment.Compare, iters...)

		for p, err := range seq {
			if !yield(p, err) {
				break
			}
		}
	}
}
