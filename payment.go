package fin

import (
	"cmp"
	"fmt"
	"time"

	"github.com/deitrix/fin/pkg/pointer"
)

type Payment struct {
	ID                 *string           `json:"id,omitempty"`
	Date               time.Time         `json:"date"`
	Amount             int               `json:"amount"`
	AccountID          string            `json:"accountId"`
	Account            *Account          `json:"account,omitempty"`
	RecurringPaymentID *string           `json:"recurringPaymentId,omitempty"`
	RecurringPayment   *RecurringPayment `json:"recurringPayment,omitempty"`
}

func (a Payment) Compare(b Payment) int {
	return cmp.Or(
		a.Date.Compare(b.Date),
		cmp.Compare(pointer.Zero(a.RecurringPayment).Name, pointer.Zero(b.RecurringPayment).Name),
	)
}

func (a Payment) AmountGBP() string {
	return fmt.Sprintf("£%.2f", float64(a.Amount)/100)
}
