package page

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func RecurringPaymentCreate() g.Node {
	return Layout("Create Recurring Payment",
		Div(Class("flex flex-col justify-center items-center"),
			Article(Class("flex flex-col flex-1 m-0 bg-white p-8 border-2 border-solid border-gray-300 w-[500px]"),
				Div(Class("mb-8"),
					H1(Class("text-center m-0"), g.Text("Create Recurring Payment")),
				),
				Form(Action("/create"), Method("post"),
					Div(Class("flex flex-col gap-4"),
						Label(g.Text("Name"),
							Input(Type("text"), Class("block w-full text-lg p-2 mt-1"), AutoComplete("off"), ID("name"), Name("name"), Required()),
						),
						Label(For("debt"), g.Text("Debt"),
							Input(Type("checkbox"), Class("form-check-input"), ID("debt"), Name("debt")),
						),
						Button(Type("submit"), Class("p-2"), g.Text("Create")),
					),
				),
			),
		),
	)
}
