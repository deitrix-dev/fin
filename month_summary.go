package fin

import (
	"time"
)

type MonthSummary struct {
	Month      time.Time `json:"month"`
	Payments   []Payment `json:"payments"`
	Income     int       `json:"income"`
	Bills      int       `json:"bills"`
	Debts      int       `json:"debts"`
	Spending   int       `json:"spending"`
	Disposable int       `json:"disposable"`
}
