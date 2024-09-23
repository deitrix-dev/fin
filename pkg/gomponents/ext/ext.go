package ext

import . "github.com/maragudk/gomponents"

func MapIndex[T any](items []T, f func(int, T) Node) []Node {
	nodes := make([]Node, len(items))
	for i, item := range items {
		nodes[i] = f(i, item)
	}
	return nodes
}
