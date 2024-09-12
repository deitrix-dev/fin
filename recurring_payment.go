package fin

import (
	"github.com/rickb777/date"
)

type PaymentSchedule struct {
	StartDate    *date.Date `json:"startDate,omitempty"`
	EndDate      *date.Date `json:"endDate,omitempty"`
	PaymentCount *int       `json:"paymentCount,omitempty"`
	Repeat       Repeat     `json:"repeat"`
	Amount       int        `json:"amount"`
	AccountID    string     `json:"accountId"`
	Account      *Account   `json:"account,omitempty"`
}

// RecurringPayment represents a payment that recurs at a regular interval. For example, a monthly
// subscription to Spotify, or a weekly payment to a babysitter.
//
// Recurring payments are made up of one or more payment schedules. This allows for a recurring
// payment to have differing amounts, accounts, or payment dates throughout the lifetime of the
// recurring payment. An example of this could be a personal loan, where the first payment is
type RecurringPayment struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Enabled   bool              `json:"enabled"`
	Debt      bool              `json:"debt"`
	Schedules []PaymentSchedule `json:"schedules"`
}
