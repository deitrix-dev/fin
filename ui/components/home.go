package components

import (
	. "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

func Home(recurringPaymentSearch string) Node {
	return Layout("Fin",
		Div(Class("grid grid-cols-3 grid-rows flex-1 overflow-auto gap-2"),
			Div(Class("flex overflow-auto"), hx.Get("/api/recurring-payments?recurringPaymentSearch="+recurringPaymentSearch), hx.Trigger("load")),
			Div(Class("col-span-2 flex overflow-auto"), hx.Get("/api/payments"), hx.Trigger("load")),
		),
	)
}
