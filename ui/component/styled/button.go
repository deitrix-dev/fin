package styled

import (
	"strings"

	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

func Button(opt string, children ...g.Node) g.Node {
	var variant string
	var size string
	for _, o := range strings.Split(opt, ",") {
		switch o {
		case "text", "primary", "secondary", "danger", "success":
			variant = o
		case "sm", "lg":
			size = o
		}
	}

	class := "no-underline"
	switch variant {
	case "text":
		class += " text-blue-600 hover:text-blue-500"
	case "primary":
		class += " border-none bg-blue-600 text-white text-lg py-2 hover:bg-blue-500 cursor-pointer"
	case "secondary":
		class += " text-blue-600 hover:bg-blue-100 border-blue-600 border-2 border-solid"
	case "danger":
		class += " bg-red-600 text-white hover:bg-red-500 border-red-600 border-2 border-solid"
	case "success":
		class += " bg-green-600 text-white hover:bg-green-500 border-green-600 border-2 border-solid"
	}

	switch size {
	case "sm":
		class += " px-2 py-0.5 text-sm"
	case "lg":
		class += " px-4 py-2 text-lg"
	default:
		class += " px-4 py-2"
	}

	return html.A(g.If(class != "", html.Class(class)), g.Group(children))
}
