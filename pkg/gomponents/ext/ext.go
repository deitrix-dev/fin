package ext

import . "github.com/maragudk/gomponents"

func MapIndex[T any](items []T, f func(int, T) Node) []Node {
	nodes := make([]Node, len(items))
	for i, item := range items {
		nodes[i] = f(i, item)
	}
	return nodes
}

func Confirm(confirmation string) Node {
	return Attr("onclick", "return window.confirm('"+confirmation+"')")
}

func IfElse(condition bool, ifTrue, ifFalse Node) Node {
	if condition {
		return ifTrue
	}
	return ifFalse
}

func IfElsef(condition bool, ifTrue, ifFalse func() Node) Node {
	if condition {
		return ifTrue()
	}
	return ifFalse()
}
