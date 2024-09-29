package styled

import (
	g "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

// so that tailwind build tool can pick up the classes
var _ = `
text-blue-600 hover:text-blue-500 border-blue-600 hover:bg-blue-100 bg-blue-600 hover:bg-blue-500 hover:border-blue-500
text-gray-600 hover:text-gray-500 border-gray-600 hover:bg-gray-100 bg-gray-600 hover:bg-gray-500 hover:border-gray-500
text-red-600 hover:text-red-500 border-red-600 hover:bg-red-100 bg-red-600 hover:bg-red-500 hover:border-red-500
text-yellow-600 hover:text-yellow-500 border-yellow-600 hover:bg-yellow-100 bg-yellow-600 hover:bg-yellow-500 hover:border-yellow-500
text-green-600 hover:text-green-500 border-green-600 hover:bg-green-100 bg-green-600 hover:bg-green-500 hover:border-green-500
`

func Link(opt Options, children ...g.Node) g.Node {
	var color string
	switch {
	case opt.Has(Primary):
		color = "blue"
	case opt.Has(Secondary):
		color = "gray"
	case opt.Has(Danger):
		color = "red"
	case opt.Has(Warn):
		color = "yellow"
	case opt.Has(Success):
		color = "green"
	}

	// text color
	class := "text-center no-underline"
	switch {
	case opt.Has(Text) || opt.Has(Bordered):
		class += " text-" + color + "-600 hover:text-" + color + "-500"
	default:
		class += " text-white"
	}

	// background color
	if opt.Has(Bordered) {
		class += " border-" + color + "-600 border-2 border-solid hover:bg-" + color + "-100"
	} else if !opt.Has(Text) {
		class += " border-" + color + "-600 border-2 border-solid bg-" + color + "-600 hover:bg-" + color + "-500" + " hover:border-" + color + "-500"
	}

	// text size
	switch {
	case opt.Has(Sm):
		class += " text-sm"
	case opt.Has(Lg):
		class += " text-lg"
	}

	// padding
	if !opt.Has(Text) {
		switch {
		case opt.Has(Sm):
			class += " px-2 py-0.5"
		case opt.Has(Lg):
			class += " px-4 py-2"
		default:
			class += " px-4 py-2"
		}
	}

	return html.A(html.Class(class), g.Group(children))
}
