package fin

import (
	"fmt"
	"iter"
	"time"

	"github.com/deitrix/fin/pkg/iterx"
)

type PaymentSchedule struct {
	StartDate *time.Time `json:"startDate,omitempty"`
	EndDate   *time.Time `json:"endDate,omitempty"`
	Repeat    Repeat     `json:"repeat"`
	Amount    int        `json:"amount"`
	AccountID string     `json:"accountId"`
	Account   *Account   `json:"account,omitempty"`
}

func (s PaymentSchedule) AmountGBP() string {
	return fmt.Sprintf("£%.2f", float64(s.Amount)/100)
}

func (s PaymentSchedule) PaymentsSince(since time.Time) iter.Seq[Payment] {
	if s.StartDate != nil && s.StartDate.After(since) {
		since = *s.StartDate
	}
	return func(yield func(Payment) bool) {
		for d := range s.Repeat.DatesSince(since) {
			if s.EndDate != nil && d.After(*s.EndDate) {
				return
			}
			if !yield(Payment{
				Date:      d,
				Amount:    s.Amount,
				AccountID: s.AccountID,
			}) {
				return
			}
		}
	}
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

func (rp *RecurringPayment) NextPayment() *Payment {
	for payment := range rp.PaymentsSince(time.Now()) {
		return &payment
	}
	return nil
}

func (rp RecurringPayment) PaymentsSince(since time.Time) iter.Seq[Payment] {
	if !rp.Enabled {
		return iterx.Empty[Payment]()
	}
	seqs := make([]iter.Seq[Payment], len(rp.Schedules))
	for i, s := range rp.Schedules {
		seqs[i] = withRecurringPaymentSeq(&rp, s.PaymentsSince(since))
	}
	return iterx.JoinFunc(seqs, Payment.Compare)
}

func (rp RecurringPayment) PaymentsSinceN(since time.Time, n int) []Payment {
	return iterx.CollectN(rp.PaymentsSince(since), n)
}

func withRecurringPaymentSeq(rp *RecurringPayment, seq iter.Seq[Payment]) iter.Seq[Payment] {
	return func(yield func(Payment) bool) {
		for p := range seq {
			p.RecurringPaymentID = &rp.ID
			p.RecurringPayment = rp
			if !yield(p) {
				return
			}
		}
	}
}

func PaymentsSince(rps []RecurringPayment, since time.Time) iter.Seq[Payment] {
	seqs := make([]iter.Seq[Payment], len(rps))
	for i, rp := range rps {
		seqs[i] = rp.PaymentsSince(since)
	}
	return iterx.JoinFunc(seqs, Payment.Compare)
}

func PaymentsSinceN(rps []RecurringPayment, since time.Time, n int) []Payment {
	return iterx.CollectN(PaymentsSince(rps, since), n)
}

func PaymentsSinceNFilter(rps []RecurringPayment, since time.Time, n int, filter func(Payment) bool) []Payment {
	return iterx.CollectNFilter(PaymentsSince(rps, since), n, filter)
}
