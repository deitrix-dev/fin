package fin

import (
	"context"
	"iter"
	"time"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Store interface {
	Account(ctx context.Context, id string) (Account, error)
	Accounts(ctx context.Context, query AccountsQuery) (Page[Account], error)
	CreateAccount(ctx context.Context, a Account) error
	UpdateAccount(ctx context.Context, a Account) error
	DeleteAccount(ctx context.Context, id string) error

	Payment(ctx context.Context, id string) (Payment, error)
	Payments(ctx context.Context, q PaymentsQuery) (Page[Payment], error)
	CreatePayment(ctx context.Context, p Payment) error
	UpdatePayment(ctx context.Context, p Payment) error
	DeletePayment(ctx context.Context, id string) error

	RecurringPayment(ctx context.Context, id string) (RecurringPayment, error)
	RecurringPayments(ctx context.Context, filter RecurringPaymentFilter) ([]RecurringPayment, error)
	CreateRecurringPayment(ctx context.Context, rp RecurringPayment) error
	UpdateRecurringPayment(ctx context.Context, rp RecurringPayment) error
	DeleteRecurringPayment(ctx context.Context, id string) error
}

type AccountsQuery struct {
	Filter AccountFilter `json:"filter"`
	Offset uint          `json:"offset"`
	Limit  uint          `json:"limit"`
}

func (q AccountsQuery) Validate() error {
	return v.ValidateStruct(&q,
		v.Field(&q.Filter),
		v.Field(&q.Offset, v.Min(0)),
		v.Field(&q.Limit, v.Required, v.Min(1), v.Max(100)),
	)
}

func (q AccountsQuery) WithPage(offset, limit uint) AccountsQuery {
	q.Offset = offset
	q.Limit = limit
	return q
}

type AccountFilter struct {
	Search string `json:"search"`
}

func (f AccountFilter) Validate() error {
	return v.ValidateStruct(&f,
		v.Field(&f.Search, v.Length(0, 100)),
	)
}

type RecurringPaymentFilter struct {
	Search string `json:"search"`
}

func (f RecurringPaymentFilter) Validate() error {
	return v.ValidateStruct(&f,
		v.Field(&f.Search, v.Length(0, 100)),
	)
}

// PaymentsQuery holds options and filters for querying payments.
type PaymentsQuery struct {
	Filter                PaymentFilter `json:"filter"`
	RecurringPaymentsOnly bool          `json:"recurringPaymentsOnly"`
	PaymentsOnly          bool          `json:"paymentsOnly"`
	Offset                uint          `json:"offset"`
	Limit                 uint          `json:"limit"`
}

func (q PaymentsQuery) WithPage(offset, limit uint) PaymentsQuery {
	q.Offset = offset
	q.Limit = limit
	return q
}

func (q PaymentsQuery) Validate() error {
	return v.ValidateStruct(&q,
		v.Field(&q.Filter),
		v.Field(&q.Offset, v.Min(0)),
		v.Field(&q.Limit, v.Required, v.Min(1), v.Max(100)),
	)
}

type PaymentFilter struct {
	After      time.Time `json:"after"`
	Before     time.Time `json:"before"`
	Search     string    `json:"search"`
	AccountIDs []string  `json:"accountIds"`
}

func (f PaymentFilter) Validate() error {
	return v.ValidateStruct(&f,
		v.Field(&f.After,
			v.When(!f.After.IsZero() && !f.Before.IsZero(), v.By(func(min any) error {
				return v.Validate(min, v.Min(f.Before).Error("must be no later than before"))
			})),
		),
		v.Field(&f.AccountIDs, v.Each(v.Required, is.UUIDv4)),
	)
}

type Page[T any] struct {
	Total   uint `json:"total"`
	Results []T  `json:"results"`
}

// PageIter can be used to iterate over all pages of a paginated query.
func PageIter[T any, Q interface{ WithPage(offset, limit uint) Q }](
	ctx context.Context,
	query Q,
	pageSize uint,
	queryFn func(context.Context, Q) (Page[T], error),
) iter.Seq2[T, error] {
	var total uint
	return func(yield func(T, error) bool) {
		for offset := uint(0); ; offset += pageSize {
			if total > 0 && offset >= total {
				return
			}
			page, err := queryFn(ctx, query.WithPage(offset, pageSize))
			if err != nil {
				if !yield(*new(T), err) {
					return
				}
			}
			if len(page.Results) == 0 {
				return
			}
			total = page.Total
			for _, result := range page.Results {
				if !yield(result, nil) {
					return
				}
			}
		}
	}
}
