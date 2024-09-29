package date

import "time"

func Midnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func Month(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func MonthDays(m time.Month) int {
	switch m {
	case time.February:
		return 28
	case time.April, time.June, time.September, time.November:
		return 30
	default:
		return 31
	}
}
