package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/deitrix/fin"
	"github.com/deitrix/sqlg"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

var mySQL = goqu.Dialect("mysql")

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

var accountsTable = goqu.T("accounts")

var selectAccounts = mySQL.
	Select("id", "name").
	From(accountsTable)

func scanAccount(row sqlg.Row) (a fin.Account, err error) {
	return a, row.Scan(
		&a.ID,
		&a.Name,
	)
}

func (s *Store) Account(ctx context.Context, id string) (fin.Account, error) {
	return sqlg.Select(ctx, s.db, scanAccount, selectAccounts.
		Where(goqu.C("id").Eq(id)).
		Limit(1))
}

func accountFilter(filter fin.AccountFilter) []exp.Expression {
	var exprs []exp.Expression
	if filter.Search != "" {
		exprs = append(exprs, goqu.C("name").ILike("%"+filter.Search+"%"))
	}
	return exprs
}

func (s *Store) Accounts(ctx context.Context, query fin.AccountsQuery) (fin.Page[fin.Account], error) {
	return pageOf(ctx, s.db, scanAccount, selectAccounts.
		Where(accountFilter(query.Filter)...).
		Offset(query.Offset).
		Limit(query.Limit))
}

func accountRow(a fin.Account) goqu.Record {
	return goqu.Record{
		"id":   a.ID,
		"name": a.Name,
	}
}

func (s *Store) CreateAccount(ctx context.Context, a fin.Account) error {
	return sqlg.Exec(ctx, s.db, mySQL.
		Insert(accountsTable).
		Rows(accountRow(a)))
}

func (s *Store) UpdateAccount(ctx context.Context, a fin.Account) error {
	return sqlg.Exec(ctx, s.db, mySQL.
		Update(accountsTable).
		Set(accountRow(a)).
		Where(goqu.C("id").Eq(a.ID)))
}

func (s *Store) DeleteAccount(ctx context.Context, id string) error {
	return sqlg.Exec(ctx, s.db, mySQL.
		Delete(accountsTable).
		Where(goqu.C("id").Eq(id)))
}

var recurringPaymentsTable = goqu.T("recurring_payments")

var selectRecurringPayments = mySQL.
	Select("id", "name", "enabled", "debt", "schedules").
	From(recurringPaymentsTable)

func scanRecurringPayment(row sqlg.Row) (rp fin.RecurringPayment, err error) {
	return rp, row.Scan(
		&rp.ID,
		&rp.Name,
		&rp.Enabled,
		&rp.Debt,
		jsonColumn(&rp.Schedules),
	)
}

func (s *Store) RecurringPayment(ctx context.Context, id string) (fin.RecurringPayment, error) {
	return sqlg.Select(ctx, s.db, scanRecurringPayment, selectRecurringPayments.
		Where(goqu.C("id").Eq(id)).
		Limit(1))
}

func recurringPaymentFilter(filter fin.RecurringPaymentFilter) []exp.Expression {
	var exprs []exp.Expression
	if filter.Search != "" {
		exprs = append(exprs, goqu.C("name").ILike("%"+filter.Search+"%"))
	}
	return exprs
}

func (s *Store) RecurringPayments(ctx context.Context, filter fin.RecurringPaymentFilter) ([]fin.RecurringPayment, error) {
	return sqlg.SelectAll(ctx, s.db, scanRecurringPayment, selectRecurringPayments.
		Where(recurringPaymentFilter(filter)...).
		Order(goqu.C("name").Asc()))
}

func recurringPaymentRow(rp fin.RecurringPayment) goqu.Record {
	return goqu.Record{
		"id":        rp.ID,
		"name":      rp.Name,
		"enabled":   rp.Enabled,
		"debt":      rp.Debt,
		"schedules": jsonColumn(&rp.Schedules),
	}
}

func (s *Store) CreateRecurringPayment(ctx context.Context, rp fin.RecurringPayment) error {
	return sqlg.Exec(ctx, s.db, mySQL.
		Insert(recurringPaymentsTable).
		Rows(recurringPaymentRow(rp)))
}

func (s *Store) UpdateRecurringPayment(ctx context.Context, rp fin.RecurringPayment) error {
	return sqlg.Exec(ctx, s.db, mySQL.
		Update(recurringPaymentsTable).
		Set(recurringPaymentRow(rp)).
		Where(goqu.C("id").Eq(rp.ID)))
}

func (s *Store) DeleteRecurringPayment(ctx context.Context, id string) error {
	return sqlg.Exec(ctx, s.db, mySQL.
		Delete(recurringPaymentsTable).
		Where(goqu.C("id").Eq(id)))
}

var paymentsTable = goqu.T("payments")

var selectPayments = mySQL.
	Select("id", "description", "date", "amount", "debt", "account_id", "recurring_payment_id").
	From(paymentsTable)

func scanPayment(row sqlg.Row) (p fin.Payment, err error) {
	return p, row.Scan(
		&p.ID,
		&p.Description,
		&p.Date,
		&p.Amount,
		&p.Debt,
		&p.AccountID,
		&p.RecurringPaymentID,
	)
}

func (s *Store) Payment(ctx context.Context, id string) (fin.Payment, error) {
	return sqlg.Select(ctx, s.db, scanPayment, selectPayments.
		Where(goqu.C("id").Eq(id)).
		Limit(1))
}

func paymentFilter(filter fin.PaymentFilter) []exp.Expression {
	var exprs []exp.Expression
	if !filter.After.IsZero() {
		exprs = append(exprs, goqu.C("date").Gt(filter.After))
	}
	if !filter.Before.IsZero() {
		exprs = append(exprs, goqu.C("date").Lt(filter.Before))
	}
	if filter.Search != "" {
		exprs = append(exprs, goqu.C("description").ILike("%"+filter.Search+"%"))
	}
	if len(filter.AccountIDs) > 0 {
		exprs = append(exprs, goqu.C("account_id").In(filter.AccountIDs))
	}
	return exprs
}

func (s *Store) Payments(ctx context.Context, q fin.PaymentsQuery) (fin.Page[fin.Payment], error) {
	return pageOf(ctx, s.db, scanPayment, selectPayments.
		Where(paymentFilter(q.Filter)...).
		Order(
			goqu.C("date").Asc(),
			goqu.C("description").Asc(),
		).
		Offset(q.Offset).
		Limit(q.Limit))
}

func paymentRow(p fin.Payment) goqu.Record {
	return goqu.Record{
		"id":                   p.ID,
		"description":          p.Description,
		"date":                 p.Date,
		"amount":               p.Amount,
		"debt":                 p.Debt,
		"account_id":           p.AccountID,
		"recurring_payment_id": p.RecurringPaymentID,
	}
}

func (s *Store) CreatePayment(ctx context.Context, p fin.Payment) error {
	return sqlg.Exec(ctx, s.db, mySQL.
		Insert(paymentsTable).
		Rows(paymentRow(p)))
}

func (s *Store) UpdatePayment(ctx context.Context, p fin.Payment) error {
	return sqlg.Exec(ctx, s.db, mySQL.
		Update(paymentsTable).
		Set(paymentRow(p)).
		Where(goqu.C("id").Eq(p.ID)))
}

func (s *Store) DeletePayment(ctx context.Context, id string) error {
	return sqlg.Exec(ctx, s.db, mySQL.
		Delete(paymentsTable).
		Where(goqu.C("id").Eq(id)))
}

func jsonColumn[T any](dst *T) jsonCol[T] {
	return jsonCol[T]{dst}
}

func pageOf[T any](ctx context.Context, q sqlg.Queryable, scan sqlg.ScanFunc[T], ds *goqu.SelectDataset) (fin.Page[T], error) {
	result, err := sqlg.SelectAll(ctx, q, scan, ds)
	if err != nil {
		return fin.Page[T]{}, err
	}

	total, err := total(ctx, q, ds)
	if err != nil {
		return fin.Page[T]{}, fmt.Errorf("failed to get total: %w", err)
	}

	return fin.Page[T]{
		Total:   total,
		Results: result,
	}, nil
}

func total(ctx context.Context, q sqlg.Queryable, ds *goqu.SelectDataset) (uint, error) {
	scan := func(row sqlg.Row) (total uint, _ error) {
		return total, row.Scan(&total)
	}
	total, err := sqlg.Select(ctx, q, scan, mySQL.
		Select(goqu.COUNT(goqu.Star())).
		From(ds.GroupBy().ClearOffset().ClearLimit()))
	if err != nil {
		return 0, err
	}
	return total, nil
}

type jsonCol[T any] struct {
	dst *T
}

func (j jsonCol[T]) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, j.dst)
	case string:
		return json.Unmarshal([]byte(src), j.dst)
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
}

func (j jsonCol[T]) Value() (driver.Value, error) {
	bs, err := json.Marshal(j.dst)
	return bs, err
}

func records[T any](s []T, fn func(T) goqu.Record) []goqu.Record {
	records := make([]goqu.Record, len(s))
	for i, x := range s {
		records[i] = fn(x)
	}
	return records
}

func nullIfZero[T comparable](v *T) *nullZero[T] {
	return &nullZero[T]{
		v: v,
	}
}

type nullZero[T comparable] struct {
	v *T
}

func (v *nullZero[T]) Value() (driver.Value, error) {
	if v.v == nil {
		return nil, nil
	}
	var zero T
	if *v.v == zero {
		return nil, nil
	}
	return *v.v, nil
}

func (v *nullZero[T]) Scan(src any) error {
	var null sql.Null[T]
	if err := null.Scan(src); err != nil {
		return err
	}
	if null.Valid {
		*v.v = null.V
	}
	return nil
}
