package components

import (
	"github.com/deitrix/fin"
	. "github.com/deitrix/fin/pkg/gomponents/ext"
	"github.com/deitrix/fin/pkg/murl"
	s "github.com/deitrix/fin/ui/components/styled"
	. "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

type PaymentsInputs struct {
	Header      string
	Payments    []fin.Payment
	FetchURL    string
	NextPage    *uint
	Search      bool
	Description bool
	OOB         bool
	OOBSearch   bool
	OOBFilter   bool
	Exclude     []string
	Filter      string
	Query       string
}

func Payments(in PaymentsInputs) Node {
	return Article(Class("flex flex-col flex-1 m-0 bg-white p-4 border-2 border-solid border-gray-300"),
		Div(Class("flex items-center mb-4 gap-4"),
			H2(Class("m-0"), Text(in.Header)),
			If(in.Search,
				Div(Class("flex gap-2 flex-grow justify-end"),
					Select(ID("paymentFilter"), Class("px-2"), Name("paymentFilter"),
						Option(Value(""), Text("All"), If(in.Filter == "", Selected())),
						Option(Value("paymentsOnly"), Text("Payments only"), If(in.Filter == "paymentsOnly", Selected())),
						Option(Value("recurringPaymentsOnly"), Text("Recurring payments only"), If(in.Filter == "recurringPaymentsOnly", Selected())),
						hx.Get(murl.Mutate(in.FetchURL,
							murl.RemoveQuery("paymentFilter"),
							murl.AddQuery("source", "paymentFilter"),
						)), hx.Trigger("change"),
						hx.Target("#paymentsContainer"), hx.Select("#paymentsContainer"),
						hx.Swap("outerHTML scroll:top"),
						If(in.OOBFilter, hx.SwapOOB("outerHTML")),
					),
					Input(ID("paymentSearch"), Class("px-2 flex-grow max-w-[400px]"), Type("search"), AutoComplete("off"), Name("paymentSearch"), Placeholder("Search"),
						hx.Get(murl.Mutate(in.FetchURL,
							murl.RemoveQuery("paymentSearch"),
							murl.AddQuery("source", "paymentSearch"),
						)), hx.Trigger("input changed, search"),
						hx.Target("#paymentsContainer"), hx.Select("#paymentsContainer"),
						hx.Swap("outerHTML scroll:top"),
						If(in.OOBSearch, hx.SwapOOB("outerHTML")),
						Value(in.Query),
					),
					s.Link(s.Primary.Sm(), Href("/payments/new"), Text("Create")),
				),
			),
		),
		Div(Class("overflow-y-auto flex-shrink text-center"),
			s.Table(
				s.Tr(
					s.Th(Text("Date")),
					If(in.Description, s.Th(Text("Description"))),
					s.Th(Text("Account")),
					s.Th(Text("Amount")),
					s.Th(Text("Actions")),
				),
				TBody(
					ID("paymentsContainer"),
					Map(in.Payments, func(payment fin.Payment) Node {
						return s.Tr(
							s.Td(Textf("%s", payment.Date.Format("Mon 2 Jan 2006"))),
							If(in.Description, s.Td(Text(payment.Description))),
							s.Td(Text(payment.AccountID)),
							s.Td(Textf("£%.2f", float64(payment.Amount)/100)),
							s.Td(Div(Class("flex gap-3 justify-center"),
								Iff(payment.ID != nil, func() Node {
									return Group{
										s.Link(s.Primary.Text(), Href("/payments/"+*payment.ID), Text("edit")),
										s.Link(s.Danger.Text(), Href("/payments/"+*payment.ID+"/delete"), Text("delete"),
											Confirm("Are you sure you want to delete this payment")),
									}
								}),
							)),
						)
					}),
				),
			),
		),
		Div(Class("flex flex-col mt-2"),
			ID("loadMore"),
			If(in.OOB, hx.SwapOOB("outerHTML")),
			Iff(in.NextPage != nil, func() Node {
				return s.Link(s.Primary.Bordered(),
					Href("#"),
					hx.Swap("beforeend"),
					hx.Get(murl.Mutate(in.FetchURL, murl.AddQuery(
						"offset", *in.NextPage,
						"source", "loadMore",
					))),
					hx.Select("#paymentsContainer>tr"),
					hx.Target("#paymentsContainer"),
					Text("Load more"),
				)
			}),
		),
	)
}
