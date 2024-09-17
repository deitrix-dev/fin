package fin

import (
	"cmp"
	"fmt"

	"github.com/deitrix/fin/pkg/pointer"
	"github.com/rickb777/date"
)

type Payment struct {
	Date               date.Date         `json:"date"`
	Amount             int               `json:"amount"`
	AccountID          string            `json:"accountId"`
	Account            *Account          `json:"account,omitempty"`
	RecurringPaymentID string            `json:"recurringPaymentId,omitempty"`
	RecurringPayment   *RecurringPayment `json:"recurringPayment,omitempty"`
}

func (a Payment) Compare(b Payment) int {
	return cmp.Or(
		cmp.Compare(a.Date.DaysSinceEpoch(), b.Date.DaysSinceEpoch()),
		cmp.Compare(pointer.Zero(a.RecurringPayment).Name, pointer.Zero(b.RecurringPayment).Name),
	)
}

func (a Payment) AmountGBP() string {
	return fmt.Sprintf("£%.2f", float64(a.Amount)/100)
}
