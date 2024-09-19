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

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

var mySQL = goqu.Dialect("mysql")

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) RecurringPayment(ctx context.Context, id string) (fin.RecurringPayment, error) {
	return sqlg.Select(ctx, s.db, scanRecurringPayment, selectRecurringPayments.
		Where(goqu.C("id").Eq(id)).
		Limit(1))
}

func (s *Store) RecurringPayments(ctx context.Context) ([]fin.RecurringPayment, error) {
	return sqlg.SelectAll(ctx, s.db, scanRecurringPayment, selectRecurringPayments)
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

func jsonColumn[T any](dst *T) jsonCol[T] {
	return jsonCol[T]{dst}
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

func recurringPaymentRow(rp fin.RecurringPayment) goqu.Record {
	return goqu.Record{
		"id":        rp.ID,
		"name":      rp.Name,
		"enabled":   rp.Enabled,
		"debt":      rp.Debt,
		"schedules": jsonColumn(&rp.Schedules),
	}
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
