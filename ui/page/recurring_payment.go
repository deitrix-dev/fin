package page

import (
	"fmt"
	"strconv"

	"github.com/deitrix/fin"
	. "github.com/deitrix/fin/pkg/gomponents/ext"
	"github.com/deitrix/fin/pkg/stringsx"
	s "github.com/deitrix/fin/ui/component/styled"
	lucide "github.com/eduardolat/gomponents-lucide"
	. "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

func RecurringPayment(rp fin.RecurringPayment) Node {
	return Layout(rp.Name,
		Div(Class("grid flex-1 overflow-auto gap-2"),
			Style("grid-template-columns: repeat(2, minmax(min-content, 1fr))"),
			Article(Class("flex flex-col flex-1 m-0 bg-white p-4 border-2 border-solid border-gray-300"),
				Div(Class("flex justify-between items-center"),
					H2(Class("m-0"), Text(rp.Name)),
					Div(Class("flex gap-4 items-center"),
						s.Link(s.Primary.Text(), Href("/recurring-payments/"+rp.ID+"/schedules/new"), Text("New schedule")),
						s.Link(s.Danger.Bordered(), Href("/recurring-payments/"+rp.ID+"/delete"), Text("Delete"),
							Confirm("Are you sure you want to delete this recurring payment")),
					),
				),
				Sub(Class("text-gray-500 mt-4"), Text("Recurring Payment")),
				RecurringPaymentForm(rp),
				Sub(Class("text-gray-500 mt-4 mb-2"), Text("Schedules")),
				s.Table(
					s.Tr(
						s.Th(Text("When")),
						s.Th(Text("Recurrence")),
						s.Th(Text("Amount")),
						s.Th(Text("Account")),
						s.Th(Text("Actions")),
					),
					Group(MapIndex(rp.Schedules, func(i int, ps fin.PaymentSchedule) Node {
						var when string
						repeat := stringsx.UpperFirst(ps.Repeat.String())
						switch {
						case ps.StartDate == nil && ps.EndDate == nil:
							when = "Forever"
						case ps.StartDate != nil && ps.EndDate != nil:
							if ps.StartDate.Equal(*ps.EndDate) {
								when = ps.StartDate.Format("2 Jan 2006")
								repeat = "Once"
							} else {
								when = fmt.Sprintf("%s - %s", ps.StartDate.Format("2 Jan 2006"), ps.EndDate.Format("2 Jan 2006"))
							}
						case ps.StartDate != nil:
							when = fmt.Sprintf("From %s", ps.StartDate.Format("2 Jan 2006"))
						case ps.EndDate != nil:
							when = fmt.Sprintf("Until %s", ps.EndDate.Format("2 Jan 2006"))
						}
						return s.Tr(
							s.Td(Text(when)),
							s.Td(Text(repeat)),
							s.Td(Text(fin.FormatCurrencyGBP(ps.Amount))),
							s.Td(Text(ps.AccountID)),
							s.Td(
								Div(Class("flex justify-center items-center gap-3"),
									s.Link(s.Primary.Text(), Href("/recurring-payments/"+rp.ID+"/schedules/"+strconv.Itoa(i)),
										lucide.PenLine(Class("size-6 cursor-pointer text-blue-500")),
										Title("Edit")),
									A(Href("/recurring-payments/"+rp.ID+"/schedules/"+strconv.Itoa(i)+"/delete"),
										lucide.Trash2(Class("size-6 cursor-pointer text-red-500")),
										Confirm("Are you sure you want to delete this payment schedule"),
										Title("Delete")),
								),
							),
						)
					})),
				),
			),
			Div(Class("flex overflow-auto"),
				hx.Get("/api/payments?recurringPayment="+rp.ID),
				hx.Trigger("load, reload from:body"),
			),
		),
	)
}

func RecurringPaymentForm(rp fin.RecurringPayment) Node {
	return Form(Class("flex"),
		hx.Post("/recurring-payments/"+rp.ID+"/form"),
		hx.Trigger("change"),
		hx.Swap("outerHTML"),
		Div(Class("mt-2 grid grid-cols-2 font-normal"),
			Label(For("enabled"), Strong(Text("Enabled"))),
			Input(Type("checkbox"), ID("enabled"), Name("enabled"), If(rp.Enabled, Checked())),
			Label(For("debt"), Strong(Text("Debt"))),
			Input(Type("checkbox"), ID("debt"), Name("debt"), If(rp.Debt, Checked())),
		),
	)
}
