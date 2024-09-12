package fin

import "github.com/rickb777/date"

type Payment struct {
	Date      date.Date `json:"date"`
	Amount    int       `json:"amount"`
	AccountID string    `json:"accountId"`
	Account   *Account  `json:"account,omitempty"`
}
