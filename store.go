package fin

import "context"

type Store interface {
	RecurringPayment(ctx context.Context, id string) (RecurringPayment, error)
	RecurringPayments(ctx context.Context) ([]RecurringPayment, error)
	CreateRecurringPayment(ctx context.Context, rp RecurringPayment) error
	UpdateRecurringPayment(ctx context.Context, rp RecurringPayment) error
	DeleteRecurringPayment(ctx context.Context, id string) error
}
