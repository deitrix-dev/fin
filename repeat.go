package fin

import (
	"cmp"
	"fmt"
	"iter"
	"slices"
	"strings"
	"time"

	"github.com/rickb777/date"
)

type Step string

const (
	Month Step = "month"
	Week  Step = "week"
	Day   Step = "day"
)

type Weekday string

const (
	Monday    Weekday = "monday"
	Tuesday   Weekday = "tuesday"
	Wednesday Weekday = "wednesday"
	Thursday  Weekday = "thursday"
	Friday    Weekday = "friday"
	Saturday  Weekday = "saturday"
	Sunday    Weekday = "sunday"
)

var mapWeekday = map[Weekday]time.Weekday{
	Monday:    time.Monday,
	Tuesday:   time.Tuesday,
	Wednesday: time.Wednesday,
	Thursday:  time.Thursday,
	Friday:    time.Friday,
	Saturday:  time.Saturday,
	Sunday:    time.Sunday,
}

func (w Weekday) String() string {
	return title(string(w))
}

func title(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

// Repeat examples:
//
//	Repeat{Every: Month, Day: 12}                          // 12th of every month
//	Repeat{Every: Week, Weekday: Monday}                   // Every Monday
//	Repeat{Every: Month, Day: 1, Multiplier: 12}           // New Year's Day
//	Repeat{Every: Month, Day: 1, Multiplier: 2}            // Every other month on the 1st
//	Repeat{Every: Month, Day: 1, Multiplier: 2, Offset: 1} // Every other month on the 1st, starting in February
//
// When using a multiplier > 1, dates are anchored to the first occurrence equal to or after the
// epoch date (01/01/1970). Offset can be used to shift the anchor date by the Every unit.
type Repeat struct {
	Every      Step    `json:"every,omitempty"`
	Day        int     `json:"day,omitempty"`
	Weekday    Weekday `json:"weekday,omitempty"`
	Multiplier int     `json:"multiplier,omitempty"`
	Offset     int     `json:"offset,omitempty"`
}

func (r Repeat) String() string {
	switch r.Every {
	case Month:
		if r.Multiplier > 1 {
			return fmt.Sprintf("every %d months on the %d%s", r.Multiplier, r.Day, dateOrdinal(r.Day))
		}
		return fmt.Sprintf("monthly on %s", dateOrdinal(r.Day))
	case Week:
		if r.Multiplier > 1 {
			return fmt.Sprintf("every %d weeks on %s", r.Multiplier, r.Weekday)
		}
		return fmt.Sprintf("weekly on %s", r.Weekday)
	}
	return ""
}

// add adds n steps to the given date and returns the new date. It assumes that the given date
// occurs on the repeat pattern.
func (r Repeat) add(d time.Time, step Step, n int) time.Time {
	switch step {
	case Month:
		months := dateMonths(d)
		months += n
		return monthsDate(months, cmp.Or(r.Day, 1))
	case Week:
		return d.AddDate(0, 0, 7*n)
	case Day:
		return d.AddDate(0, 0, n)
	}
	panic(fmt.Sprintf("invalid repeat step: %s", step))
}

// First returns the first occurrence of the repeat pattern after or equal to the given date.
func (r Repeat) First(since time.Time) time.Time {
	mul := max(r.Multiplier, 1)
	switch r.Every {
	case Day:
		days := int(since.Unix() / 86400)
		days = days - (days+mul-r.Offset)%mul
		d := time.Unix(int64(days*86400), 0).In(since.Location())
		if d.Before(since) {
			d = r.add(d, Day, mul)
		}
		return d
	case Week:
		weeks := int(since.Unix() / 86400 / 7)
		weeks = weeks - (weeks+mul-r.Offset)%mul
		d := time.Unix(int64(weeks*7*86400), 0)
		for d.Weekday() != mapWeekday[cmp.Or(r.Weekday, Monday)] {
			d = d.Add(1 * 24 * time.Hour)
		}
		if d.Before(since) {
			d = r.add(d, Week, mul)
		}
		return d
	case Month:
		months := dateMonths(since)
		months = months - (months+mul-r.Offset)%mul
		d := monthsDate(months, cmp.Or(r.Day, 1))
		if d.Before(since) {
			d = r.add(d, Month, mul)
		}
		return d
	}
	panic(fmt.Sprintf("invalid repeat step: %s", r.Every))
}

// DatesSince returns an iterator that yields dates that occur on the repeat pattern after or equal
// to the given date. Because there are infinitely many dates that can be yielded, the iterator
// should be used with a limit or a break condition.
func (r Repeat) DatesSince(since time.Time) iter.Seq[time.Time] {
	mul := max(r.Multiplier, 1)
	first := r.First(since)
	return func(yield func(time.Time) bool) {
		for d := first; ; d = r.add(d, r.Every, mul) {
			if !yield(d) {
				return
			}
		}
	}
}

// DatesUntil returns an iterator that yields dates that occur on the repeat pattern before or equal
// to the given date. Because there are infinitely many dates that can be yielded, the iterator
// should be used with a limit or a break condition.
//
// Dates are yielded in reverse order, starting from the last date that occurs on the repeat pattern
// before or equal to the given date.
func (r Repeat) DatesUntil(until time.Time) iter.Seq[time.Time] {
	mul := max(r.Multiplier, 1)
	first := r.First(until)
	return func(yield func(time.Time) bool) {
		for d := first; ; d = r.add(d, r.Every, -mul) {
			if d.After(until) {
				continue
			}
			if !yield(d) {
				return
			}
		}
	}
}

// DatesUntilN returns the first n dates that occur on the repeat pattern before or equal to the
// given date.
func (r Repeat) DatesUntilN(until time.Time, n int) []time.Time {
	var dates []time.Time
	for d := range r.DatesUntil(until) {
		dates = append(dates, d)
		if len(dates) >= n {
			break
		}
	}
	slices.Reverse(dates)
	return dates
}

// DatesSinceN returns the first n dates that occur on the repeat pattern after or equal to the
// given date.
func (r Repeat) DatesSinceN(since time.Time, n int) []time.Time {
	var dates []time.Time
	for d := range r.DatesSince(since) {
		dates = append(dates, d)
		if len(dates) >= n {
			break
		}
	}
	return dates
}

// DatesBetween returns all dates that occur on the repeat pattern between the given dates,
// inclusive.
func (r Repeat) DatesBetween(since, until time.Time) []time.Time {
	var dates []time.Time
	for d := range r.DatesSince(since) {
		if d.After(until) {
			break
		}
		dates = append(dates, d)
	}
	return dates
}

// normalizeMonthDay returns the day of the month normalized to the number of days in the month. For
// example, if the month is February and the day is 31, the day will be normalized to either 28 or
// 29 (depending on whether the year is a leap year).
func normalizeMonthDay(year int, month time.Month, day int) int {
	daysIn := date.DaysIn(year, month)
	if day > daysIn {
		return daysIn
	}
	return day
}

func dateMonths(d time.Time) int {
	return (d.Year()-1970)*12 + int(d.Month()) - 1
}

func monthsDate(m int, day int) time.Time {
	if m < 0 {
		year := 1970 + m/12 - 1
		month := time.Month(m%12 + 13)
		return time.Date(year, month, normalizeMonthDay(year, month, day), 0, 0, 0, 0, time.UTC)
	}
	year := 1970 + m/12
	month := time.Month(m%12 + 1)
	return time.Date(year, month, normalizeMonthDay(year, month, day), 0, 0, 0, 0, time.UTC)
}

func dateOrdinal(n int) string {
	if n >= 11 && n <= 13 {
		return fmt.Sprintf("%dth", n)
	}

	switch n % 10 {
	case 1:
		return fmt.Sprintf("%dst", n)
	case 2:
		return fmt.Sprintf("%dnd", n)
	case 3:
		return fmt.Sprintf("%drd", n)
	default:
		return fmt.Sprintf("%dth", n)
	}
}
