package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapGraph(t *testing.T) {
	g := NewGraph[string, string]()
	// prepare graph
	g.AddVertex("A", "A")
	g.AddVertex("B", "B")
	g.AddVertex("C", "C")
	g.AddVertex("D", "D")
	g.AddVertex("E", "E")
	g.AddVertex("F", "F")
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "E")
	g.AddEdge("D", "E")
	g.AddEdge("D", "F")
	g.AddEdge("E", "F")
	// add edges to itself does nothing
	g.AddEdge("A", "A")
	g.AddEdge("B", "B")
	g.AddEdge("C", "C")
	g.AddEdge("D", "D")
	g.AddEdge("E", "E")
	g.AddEdge("F", "F")
	// add edges multiple times does nothing
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "E")
	g.AddEdge("D", "E")
	g.AddEdge("D", "F")
	g.AddEdge("E", "F")

	assert.Equal(t, 6, len(g.GetVertices()))

	assert.ElementsMatch(t, []string{"B", "C"}, g.GetVertex("A").GetNeighbors())
	assert.ElementsMatch(t, []string{"A", "D"}, g.GetVertex("B").GetNeighbors())
	assert.ElementsMatch(t, []string{"A", "E"}, g.GetVertex("C").GetNeighbors())
	assert.ElementsMatch(t, []string{"B", "E", "F"}, g.GetVertex("D").GetNeighbors())
	assert.ElementsMatch(t, []string{"C", "D", "F"}, g.GetVertex("E").GetNeighbors())
	assert.ElementsMatch(t, []string{"D", "E"}, g.GetVertex("F").GetNeighbors())

	clone := g.Clone()
	assert.Equal(t, 6, len(clone.GetVertices()))

	assert.ElementsMatch(t, []string{"B", "C"}, clone.GetVertex("A").GetNeighbors())
	assert.ElementsMatch(t, []string{"A", "D"}, clone.GetVertex("B").GetNeighbors())
	assert.ElementsMatch(t, []string{"A", "E"}, clone.GetVertex("C").GetNeighbors())
	assert.ElementsMatch(t, []string{"B", "E", "F"}, clone.GetVertex("D").GetNeighbors())
	assert.ElementsMatch(t, []string{"C", "D", "F"}, clone.GetVertex("E").GetNeighbors())
	assert.ElementsMatch(t, []string{"D", "E"}, clone.GetVertex("F").GetNeighbors())

	// remove edge
	g.RemoveEdge("A", "B")
	g.RemoveEdge("D", "E")

	assert.ElementsMatch(t, []string{"C"}, g.GetVertex("A").GetNeighbors())
	assert.ElementsMatch(t, []string{"D"}, g.GetVertex("B").GetNeighbors())
	assert.ElementsMatch(t, []string{"A", "E"}, g.GetVertex("C").GetNeighbors())
	assert.ElementsMatch(t, []string{"B", "F"}, g.GetVertex("D").GetNeighbors())
	assert.ElementsMatch(t, []string{"C", "F"}, g.GetVertex("E").GetNeighbors())
	assert.ElementsMatch(t, []string{"D", "E"}, g.GetVertex("F").GetNeighbors())

	// re-adding removed edge
	g.AddEdge("A", "B")
	g.AddEdge("D", "E")

	assert.ElementsMatch(t, []string{"B", "C"}, g.GetVertex("A").GetNeighbors())
	assert.ElementsMatch(t, []string{"A", "D"}, g.GetVertex("B").GetNeighbors())
	assert.ElementsMatch(t, []string{"A", "E"}, g.GetVertex("C").GetNeighbors())
	assert.ElementsMatch(t, []string{"B", "E", "F"}, g.GetVertex("D").GetNeighbors())
	assert.ElementsMatch(t, []string{"C", "D", "F"}, g.GetVertex("E").GetNeighbors())
	assert.ElementsMatch(t, []string{"D", "E"}, g.GetVertex("F").GetNeighbors())

	// remove vertex
	g.RemoveVertex("A")
	g.RemoveVertex("D")

	assert.ElementsMatch(t, []string{}, g.GetVertex("B").GetNeighbors())
	assert.ElementsMatch(t, []string{"E"}, g.GetVertex("C").GetNeighbors())
	assert.ElementsMatch(t, []string{"C", "F"}, g.GetVertex("E").GetNeighbors())
	assert.ElementsMatch(t, []string{"E"}, g.GetVertex("F").GetNeighbors())
}
