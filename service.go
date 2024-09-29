package fin

import (
	"context"
	"iter"
	"slices"
	"time"

	"github.com/deitrix/fin/pkg/date"
	"github.com/deitrix/fin/pkg/iterx"
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
		if query.Filter.After.IsZero() {
			query.Filter.After = date.Midnight(time.Now())
		}

		var iters []iter.Seq2[Payment, error]

		if !query.PaymentsOnly {
			rps, err := s.Store.RecurringPayments(ctx, RecurringPaymentFilter{Search: query.Filter.Search})
			if err != nil {
				yield(Payment{}, err)
				return
			}
			iters = append(iters, iterx.WithNilErr(PaymentsSince(rps, query.Filter.After)))
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

type MonthSummariesQuery struct {
	After  time.Time `json:"after"`
	Offset uint      `json:"offset"`
	Limit  uint      `json:"limit"`
}

func (s *Service) MonthSummaries(ctx context.Context, query MonthSummariesQuery) ([]MonthSummary, error) {
	if query.Limit == 0 {
		query.Limit = 12
	}

	paymentIter := s.Payments(ctx, PaymentsQuery{Filter: PaymentFilter{After: query.After}})
	monthPayments := make(map[time.Time][]Payment)
	for p, err := range paymentIter {
		if err != nil {
			return nil, err
		}
		month := date.Month(p.Date)
		if _, ok := monthPayments[month]; !ok && len(monthPayments) >= int(query.Limit) {
			break
		}
		monthPayments[month] = append(monthPayments[month], p)
	}

	var summaries []MonthSummary
	for month, payments := range monthPayments {
		summary := MonthSummary{
			Month:    month,
			Income:   500000,                               // TODO hardcoded for now
			Spending: date.MonthDays(month.Month()) * 4000, // TODO hardcoded for now
		}
		for _, p := range payments {
			if p.Debt {
				summary.Debts += p.Amount
			} else {
				summary.Bills += p.Amount
			}
		}
		summary.Disposable = summary.Income - summary.Bills - summary.Debts - summary.Spending
		summaries = append(summaries, summary)
	}

	slices.SortFunc(summaries, func(a, b MonthSummary) int {
		return a.Month.Compare(b.Month)
	})

	return summaries, nil
}
