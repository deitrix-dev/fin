package components

import (
	s "github.com/deitrix/fin/ui/components/styled"
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
			Div(Class("flex flex-col h-[100vh] p-4 gap-2 m-auto"),
				Div(Class("flex justify-between items-center px-2"),
					s.Link(Href("/"), H1(Class("m-0"), Text("Fin"))),
					Div(hx.Get("/api/header-user"), hx.Trigger("load"), hx.Swap("outerHTML")),
				),
				Group(components),
			),
			Script(
				Src("https://unpkg.com/htmx.org@2.0.2/dist/htmx.js"),
				Integrity("sha384-yZq+5izaUBKcRgFbxgkRYwpHhHHCpp5nseXp0MEQ1A4MTWVMnqkmcuFez8x5qfxr"),
				CrossOrigin("anonymous"),
			),
		},
	})
}
