package fin

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
)

func FormatCurrencyGBP(amount int) string {
	if amount < 0 {
		return fmt.Sprintf("-£%s", FormatCurrency(-amount))
	}
	return fmt.Sprintf("£%s", FormatCurrency(amount))
}

func FormatCurrency(amount int) string {
	if amount%100 == 0 {
		return humanize.Comma(int64(amount / 100))
	}
	if amount%10 == 0 {
		return strings.TrimSuffix(humanize.Commaf(float64(amount+1)/100), "1") + "0" // awful hack
	}
	return humanize.Commaf(float64(amount) / 100)
}
