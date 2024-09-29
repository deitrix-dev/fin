package component

import (
	s "github.com/deitrix/fin/ui/component/styled"
	. "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func HeaderUser(email string) Node {
	if email != "" {
		return Div(Class("flex gap-4"),
			Span(Text(email)),
			s.Link(s.Primary.Text(), Href("/cdn-cgi/access/logout"), Text("Logout")),
		)
	}
	return nil
}
