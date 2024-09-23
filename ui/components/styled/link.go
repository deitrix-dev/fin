package styled

import (
	. "github.com/maragudk/gomponents"
	"github.com/maragudk/gomponents/html"
)

func Link(children ...Node) Node {
	return html.A(html.Class("text-blue-700 hover:text-blue-500 no-underline"),
		Group(children),
	)
}
