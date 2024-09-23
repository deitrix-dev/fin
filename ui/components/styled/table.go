package styled

import (
	. "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

func Table(children ...Node) Node {
	return html.Table(html.Class("w-full border-collapse"), Group(children))
}

func Tr(children ...Node) Node {
	return html.Tr(html.Class("odd:bg-gray-200 p-4 text-center bg-gray-100 hover:bg-gray-300"), Group(children))
}

func Th(children ...Node) Node {
	return html.Th(html.Class("p-4 bg-gray-700 text-white text-center"), Group(children))
}

func Td(children ...Node) Node {
	return html.Td(html.Class("p-4"), Group(children))
}
