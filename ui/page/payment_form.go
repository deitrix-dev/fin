package page

import (
	"fmt"
	"slices"

	"github.com/deitrix/fin"
	. "github.com/deitrix/fin/pkg/gomponents/ext"
	s "github.com/deitrix/fin/ui/component/styled"
	. "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

func PaymentForm(accounts []string, payment fin.Payment) Node {
	createUpdate := "Create"
	if payment.ID != nil {
		createUpdate = "Update"
	}
	isNewAccount := !slices.Contains(accounts, payment.AccountID)
	var postURL string
	if payment.ID != nil {
		postURL = fmt.Sprintf("/payments/%s", *payment.ID)
	} else {
		postURL = "/payments/new"
	}
	refreshFormURL := postURL + "?formPriority=true"
	return Layout("Payment",
		Div(Class("flex flex-col justify-center items-center"),
			Article(Class("flex flex-col flex-1 m-0 bg-white p-8 border-2 border-solid border-gray-300 w-[500px]"),
				Div(Class("mb-8"),
					H1(Class("text-center m-0"), Text(createUpdate+" Payment")),
				),
				Form(ID("paymentForm"), Action(postURL), Method("post"),
					Div(Class("flex flex-col gap-4"),
						Input(Type("hidden"), ID("id"), Name("id"), Iff(payment.ID != nil, func() Node { return Value(*payment.ID) })),
						Label(Text("Description"),
							Input(Type("text"), Class("block w-full text-lg p-2 mt-1"), AutoComplete("off"), ID("description"), Name("description"), Required(),
								Value(payment.Description))),
						Label(Text("Account"),
							Select(Class("block w-full text-lg p-2 mt-1"), ID("account"), Name("account"),
								Map(accounts, func(account string) Node {
									return Option(Value(account), Text(account), If(account == payment.AccountID, Selected()))
								}),
								Option(Value(""), Text("New account"), If(isNewAccount, Selected())),
								hx.Get(refreshFormURL), hx.Trigger("change"), hx.Include("#paymentForm"),
								hx.Select("#paymentForm"), hx.Target("#paymentForm"), hx.Swap("outerHTML"),
							),
						),
						If(isNewAccount,
							Label(Text("New account"),
								Input(Type("text"), Class("block w-full text-lg p-2 mt-1"), AutoComplete("off"), ID("newAccount"), Name("newAccount"),
									Value(payment.AccountID)))),
						Label(Text("Date"),
							Input(Type("date"), Class("block w-full text-lg p-2 mt-1"), AutoComplete("off"), ID("date"), Name("date"), Required(),
								Value(payment.Date.Format("2006-01-02")))),
						Label(Text("Amount"),
							Input(Type("number"), Step("0.01"), Class("block w-full text-lg p-2 mt-1"), AutoComplete("off"), ID("amount"), Name("amount"), Required(),
								Value(fmt.Sprint(float64(payment.Amount)/100)))),
						Label(For("debt"), Text("Debt"),
							Input(Type("checkbox"), Class("form-check-input"), ID("debt"), Name("debt"), Value("true"), If(payment.Debt, Checked())),
						),
						Button(Class("font-bold border-none bg-blue-600 text-white text-lg py-2 hover:bg-blue-500 cursor-pointer"),
							Type("submit"), Text(createUpdate)),
						Iff(payment.ID != nil, func() Node {
							return s.Link(s.Danger.Lg().Bordered(), Text("Delete"), Href(fmt.Sprintf("/payments/%s/delete", *payment.ID)),
								Confirm("Are you sure you want to delete this payment?"))
						}),
					),
				),
			),
		),
	)
}
