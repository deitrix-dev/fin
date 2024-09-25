package components

import (
	"github.com/deitrix/fin"
	. "github.com/deitrix/fin/pkg/gomponents/ext"
	s "github.com/deitrix/fin/ui/components/styled"
	. "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

func RecurringPayments(recurringPayments []fin.RecurringPayment, search string) Node {
	return Article(Class("flex flex-col flex-1 m-0 bg-white p-4 border-2 border-solid border-gray-300"),
		Div(Class("flex justify-between items-center mb-4 gap-4"),
			H2(Class("m-0"), Text("Recurring Payments")),
			Div(Class("flex gap-2 flex-grow justify-end"),
				Input(
					Class("px-2 flex-grow max-w-[400px]"),
					Type("search"),
					ID("search"),
					Name("recurringPaymentSearch"),
					Placeholder("Search"),
					Value(search),
					hx.Get("/api/recurring-payments"),
					hx.Trigger("input changed"),
					hx.Target("#recurringPaymentsTable"),
					hx.Select("#recurringPaymentsTable"),
					hx.Swap("outerHTML scroll:top"),
				),
				s.Link(s.Primary.Sm(), Href("/create"), Text("Create")),
			),
		),
		Div(Class("overflow-y-auto flex-shrink"),
			s.Table(ID("recurringPaymentsTable"),
				s.Tr(
					s.Th(Text("Name")),
					s.Th(Text("Next Payment")),
					s.Th(Text("Actions")),
				),
				Map(recurringPayments, func(rp fin.RecurringPayment) Node {
					return s.Tr(
						s.Td(s.Link(s.Primary.Text(), Href("/recurring-payments/"+rp.ID), Text(rp.Name))),
						s.Td(Iff(rp.NextPayment() != nil, func() Node {
							np := rp.NextPayment()
							return Textf("%s on %s", np.AmountGBP(), np.Date.Format("2 Jan"))
						})),
						s.Td(s.Link(s.Danger.Text(), Href("/recurring-payments/"+rp.ID+"/delete"), Text("delete"),
							Confirm("Are you sure you want to delete this recurring payment?"))),
					)
				}),
			),
		),
	)
}
