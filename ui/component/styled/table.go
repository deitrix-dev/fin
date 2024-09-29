package styled

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

func Table(children ...g.Node) g.Node {
	return html.Table(html.Class("w-full border-collapse"), g.Group(children))
}

func Tr(children ...g.Node) g.Node {
	return html.Tr(html.Class("odd:bg-gray-200 p-4 text-center bg-gray-100 hover:bg-gray-300"), g.Group(children))
}

func Th(children ...g.Node) g.Node {
	return html.Th(html.Class("sticky top-0 p-4 bg-gray-700 text-white text-center"), g.Group(children))
}

func Td(children ...g.Node) g.Node {
	return html.Td(html.Class("p-4"), g.Group(children))
}
