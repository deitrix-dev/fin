package cache

import (
	"context"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/cache"
	"github.com/mitchellh/hashstructure/v2"
)

type Store struct {
	fin.Store
	accountCache           cache.Cache[string, fin.Account]
	accountsCache          cache.Cache[uint64, fin.Page[fin.Account]]
	paymentCache           cache.Cache[string, fin.Payment]
	paymentsCache          cache.Cache[uint64, fin.Page[fin.Payment]]
	recurringPaymentCache  cache.Cache[string, fin.RecurringPayment]
	recurringPaymentsCache cache.Cache[uint64, []fin.RecurringPayment]
}

func NewStore(s fin.Store) *Store {
	return &Store{
		Store: s,
	}
}

func (s *Store) Account(ctx context.Context, id string) (fin.Account, error) {
	return s.accountCache.GetFunc(id, func() (fin.Account, error) {
		return s.Store.Account(ctx, id)
	})
}

func (s *Store) Accounts(ctx context.Context, query fin.AccountsQuery) (fin.Page[fin.Account], error) {
	return s.accountsCache.GetFunc(hash(query), func() (fin.Page[fin.Account], error) {
		return s.Store.Accounts(ctx, query)
	})
}

func (s *Store) CreateAccount(ctx context.Context, a fin.Account) error {
	defer s.accountCache.Clear()
	defer s.accountsCache.Clear()
	return s.Store.CreateAccount(ctx, a)
}

func (s *Store) UpdateAccount(ctx context.Context, a fin.Account) error {
	defer s.accountCache.Clear()
	defer s.accountsCache.Clear()
	return s.Store.UpdateAccount(ctx, a)
}

func (s *Store) DeleteAccount(ctx context.Context, id string) error {
	defer s.accountCache.Clear()
	defer s.accountsCache.Clear()
	return s.Store.DeleteAccount(ctx, id)
}

func (s *Store) Payment(ctx context.Context, id string) (fin.Payment, error) {
	return s.paymentCache.GetFunc(id, func() (fin.Payment, error) {
		return s.Store.Payment(ctx, id)
	})
}

func (s *Store) Payments(ctx context.Context, q fin.PaymentsQuery) (fin.Page[fin.Payment], error) {
	return s.paymentsCache.GetFunc(hash(q), func() (fin.Page[fin.Payment], error) {
		return s.Store.Payments(ctx, q)
	})
}

func (s *Store) CreatePayment(ctx context.Context, p fin.Payment) error {
	defer s.paymentCache.Clear()
	defer s.paymentsCache.Clear()
	return s.Store.CreatePayment(ctx, p)
}

func (s *Store) UpdatePayment(ctx context.Context, p fin.Payment) error {
	defer s.paymentCache.Clear()
	defer s.paymentsCache.Clear()
	return s.Store.UpdatePayment(ctx, p)
}

func (s *Store) DeletePayment(ctx context.Context, id string) error {
	defer s.paymentCache.Clear()
	defer s.paymentsCache.Clear()
	return s.Store.DeletePayment(ctx, id)
}

func (s *Store) RecurringPayment(ctx context.Context, id string) (fin.RecurringPayment, error) {
	return s.recurringPaymentCache.GetFunc(id, func() (fin.RecurringPayment, error) {
		return s.Store.RecurringPayment(ctx, id)
	})
}

func (s *Store) RecurringPayments(ctx context.Context, filter fin.RecurringPaymentFilter) ([]fin.RecurringPayment, error) {
	return s.recurringPaymentsCache.GetFunc(hash(filter), func() ([]fin.RecurringPayment, error) {
		return s.Store.RecurringPayments(ctx, filter)
	})
}

func (s *Store) CreateRecurringPayment(ctx context.Context, rp fin.RecurringPayment) error {
	defer s.recurringPaymentCache.Clear()
	defer s.recurringPaymentsCache.Clear()
	return s.Store.CreateRecurringPayment(ctx, rp)
}

func (s *Store) UpdateRecurringPayment(ctx context.Context, rp fin.RecurringPayment) error {
	defer s.recurringPaymentCache.Clear()
	defer s.recurringPaymentsCache.Clear()
	return s.Store.UpdateRecurringPayment(ctx, rp)
}

func (s *Store) DeleteRecurringPayment(ctx context.Context, id string) error {
	defer s.recurringPaymentCache.Clear()
	defer s.recurringPaymentsCache.Clear()
	return s.Store.DeleteRecurringPayment(ctx, id)
}

func hash(v any) uint64 {
	h, err := hashstructure.Hash(v, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}
	return h
}
