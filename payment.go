package fin

import (
	"cmp"
	"fmt"
	"time"
)

type Payment struct {
	ID                 *string           `json:"id,omitempty"`
	Description        string            `json:"description"`
	Date               time.Time         `json:"date"`
	Amount             int               `json:"amount"`
	Debt               bool              `json:"debt"`
	AccountID          string            `json:"accountId"`
	Account            *Account          `json:"account,omitempty"`
	RecurringPaymentID *string           `json:"recurringPaymentId,omitempty"`
	RecurringPayment   *RecurringPayment `json:"recurringPayment,omitempty"`
}

func (a Payment) Compare(b Payment) int {
	return cmp.Or(
		a.Date.Compare(b.Date),
		cmp.Compare(a.Description, b.Description),
	)
}

func (a Payment) AmountGBP() string {
	return fmt.Sprintf("£%.2f", float64(a.Amount)/100)
}
