package component

import (
	"github.com/deitrix/fin"
	s "github.com/deitrix/fin/ui/component/styled"
	. "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func MonthSummaries(summaries []fin.MonthSummary) Node {
	return Article(Class("flex flex-col flex-1 m-0 bg-white p-4 border-2 border-solid border-gray-300"),
		Div(Class("flex items-center mb-4 gap-4"),
			H2(Class("m-0"), Text("Month Summaries")),
		),
		Div(Class("overflow-y-auto flex-shrink text-center"),
			s.Table(
				s.Tr(
					s.Th(Text("Month")),
					s.Th(Text("Income")),
					s.Th(Text("Expenses")),
					s.Th(Text("Spending")),
					s.Th(Text("Disposable")),
				),
				TBody(
					Map(summaries, func(summary fin.MonthSummary) Node {
						return s.Tr(
							s.Td(Text(summary.Month.Format("Jan 2006"))),
							s.Td(Span(Class("text-green-600"), Text(fin.FormatCurrencyGBP(summary.Income)))),
							s.Td(Text(fin.FormatCurrencyGBP(summary.Bills+summary.Debts))),
							s.Td(Text(fin.FormatCurrencyGBP(summary.Spending))),
							s.Td(Text(fin.FormatCurrencyGBP(summary.Disposable))),
						)
					}),
				),
			),
		),
	)
}
