package page

import (
	. "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

func Home() Node {
	return Layout("Fin",
		Div(Class("grid grid-rows-2 flex-1 overflow-auto gap-2"),
			Style("grid-template-columns: repeat(2, minmax(min-content, 1fr))"),
			Div(Class("flex overflow-auto"), hx.Get("/api/recurring-payments?referer=true"), hx.Trigger("load")),
			Div(Class("flex overflow-auto row-span-2"), hx.Get("/api/payments?referer=true"), hx.Trigger("load")),
			Div(Class("flex overflow-auto"), hx.Get("/api/month-summaries?referer=true"), hx.Trigger("load")),
		),
	)
}
