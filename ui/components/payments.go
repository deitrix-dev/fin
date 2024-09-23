package components

import (
	"fmt"
	"net/url"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/pointer"
	s "github.com/deitrix/fin/ui/components/styled"
	. "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

type PaymentsInputs struct {
	Header      string
	Payments    []fin.Payment
	FetchURL    string
	NextPage    *int
	Search      bool
	Description bool
	OOB         bool
}

func Payments(in PaymentsInputs) Node {
	return Article(Class("flex flex-col flex-1 m-0 bg-white p-4 border-2 border-solid border-gray-300"),
		Div(Class("flex justify-between items-center mb-4 gap-4"),
			H2(Class("m-0"), Text(in.Header)),
			If(in.Search,
				Div(Class("flex gap-2"),
					Input(Type("search"), AutoComplete("off"), Name("q"), Placeholder("Search"),
						hx.Get(removeQuery(in.FetchURL, "q")), hx.Trigger("input changed, search"),
						hx.Target("#paymentsContainer"), hx.Select("#paymentsContainer"),
						hx.Swap("outerHTML scroll:top")),
					Form(Button(Action("/create"), Text("Create"))),
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
				),
				TBody(
					ID("paymentsContainer"),
					Group(Map(in.Payments, func(payment fin.Payment) Node {
						return s.Tr(
							s.Td(Textf("%s", payment.Date.Format("Mon 2 Jan 2006"))),
							If(in.Description, s.Td(Text(pointer.Zero(payment.RecurringPayment).Name))),
							s.Td(Text(payment.AccountID)),
							s.Td(Textf("£%.2f", float64(payment.Amount)/100)),
						)
					})),
				),
			),
		),
		Div(
			ID("loadMore"),
			If(in.OOB, hx.SwapOOB("outerHTML")),
			Iff(in.NextPage != nil, func() Node {
				return Button(Class("w-full p-2"),
					hx.Swap("beforeend"),
					hx.Get(addQuery(in.FetchURL, "offset", *in.NextPage)),
					hx.Select("#paymentsContainer>tr"),
					hx.Target("#paymentsContainer"),
					Text("Load more"),
				)
			}),
		),
	)
}

func addQuery(rawURL, key string, value any) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Set(key, fmt.Sprint(value))
	u.RawQuery = q.Encode()
	return u.String()
}

func removeQuery(rawURL, key string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Del(key)
	u.RawQuery = q.Encode()
	return u.String()
}
