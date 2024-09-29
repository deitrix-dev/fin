package page

import (
	. "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

func RecurringPayments() Node {
	return Layout("Recurring Payments",
		Div(Class("flex overflow-auto"), hx.Get("/api/recurring-payments?referer=true"), hx.Trigger("load")),
	)
}
