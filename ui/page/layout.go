package page

import (
	s "github.com/deitrix/fin/ui/component/styled"
	. "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func Layout(title string, components ...Node) Node {
	return HTML5(HTML5Props{
		Title: title,
		Head: []Node{
			Link(Type("text/css"), Rel("stylesheet"), Href("/assets/style.css")),
			Meta(Charset("utf-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
		},
		Body: []Node{Class("m-0 p-0 bg-gray-100"),
			Div(Class("flex flex-col h-[100vh] p-4 gap-2 max-w-[1690px] m-auto"),
				Div(Class("flex justify-between items-center px-2"),
					Div(Class("flex flex-1"), s.Link(s.Primary.Text(), Href("/"), H1(Class("m-0"), Text("Fin")))),
					Div(Class("flex gap-10 flex-1 justify-center"),
						s.Link(s.Primary.Text(), Href("#"), Text("Accounts")),
						s.Link(s.Primary.Text(), Href("#"), Text("Income")),
						s.Link(s.Primary.Text(), Href("#"), Text("Spending")),
					),
					Div(Class("flex flex-1 justify-end"),
						Div(hx.Get("/api/header-user"), hx.Trigger("load"), hx.Swap("outerHTML"))),
				),
				Group(components),
			),
			Script(Src("/assets/index.js")),
		},
	})
}
