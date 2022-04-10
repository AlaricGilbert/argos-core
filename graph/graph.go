package graph

// Vertex is a vertex in a MapGraph.
type Vertex[Key comparable, Value any] struct {
	key       Key
	value     Value
	neighbors map[Key]struct{}
}

// Graph is a graph implemented as a map[Key, Vertex[Key, Value]] of vertices.
type Graph[Key comparable, Value any] struct {
	vertices map[Key]*Vertex[Key, Value]
}

// NewGraph returns a new graph.
func NewGraph[Key comparable, Value any]() *Graph[Key, Value] {
	return &Graph[Key, Value]{
		vertices: make(map[Key]*Vertex[Key, Value]),
	}
}

// AddVertex adds a vertex to the graph.
func (g *Graph[Key, Value]) AddVertex(key Key, value Value) {
	if _, ok := g.vertices[key]; !ok {
		g.vertices[key] = &Vertex[Key, Value]{
			key:       key,
			value:     value,
			neighbors: make(map[Key]struct{}, 0),
		}
	}
}

// AddVertex adds a vertex to the graph.
func (g *Graph[Key, Value]) AddVertexWithFactory(key Key, factory func() Value) {
	if _, ok := g.vertices[key]; !ok {
		g.vertices[key] = &Vertex[Key, Value]{
			key:       key,
			value:     factory(),
			neighbors: make(map[Key]struct{}, 0),
		}
	}
}

func (g *Graph[Key, Value]) getTwoVertices(k1, k2 Key) (v1, v2 *Vertex[Key, Value], ok bool) {
	if v1, ok = g.vertices[k1]; !ok {
		return
	}
	v2, ok = g.vertices[k2]
	return
}

// AddEdge adds an edge to the graph.
func (g *Graph[Key, Value]) AddEdge(from Key, to Key) {
	if from == to {
		return
	}
	if fromVertex, toVertex, ok := g.getTwoVertices(from, to); ok {
		fromVertex.neighbors[to] = struct{}{}
		toVertex.neighbors[from] = struct{}{}
	}
}

// RemoveVertex removes a vertex from the graph.
func (g *Graph[Key, Value]) RemoveVertex(key Key) {
	if _, ok := g.vertices[key]; ok {
		for neighborKey := range g.vertices[key].neighbors {
			if neighborVertex, ok := g.vertices[neighborKey]; ok {
				delete(neighborVertex.neighbors, key)
			}
		}
		delete(g.vertices, key)
	}
}

// RemoveEdge removes an edge from the graph.
func (g *Graph[Key, Value]) RemoveEdge(from Key, to Key) {
	if fromVertex, toVertex, ok := g.getTwoVertices(from, to); ok {
		delete(fromVertex.neighbors, to)
		delete(toVertex.neighbors, from)
	}
}

// GetVertex returns a vertex from the graph.
func (g *Graph[Key, Value]) GetVertex(key Key) *Vertex[Key, Value] {
	if v, ok := g.vertices[key]; ok {
		return v
	}
	return nil
}

// GetVertices returns all vertices from the graph.
func (g *Graph[Key, Value]) GetVertices() []*Vertex[Key, Value] {
	keys := make([]*Vertex[Key, Value], 0, len(g.vertices))
	for key := range g.vertices {
		keys = append(keys, g.vertices[key])
	}
	return keys
}

// Clone returns a copy of the graph.
func (g *Graph[Key, Value]) Clone() *Graph[Key, Value] {
	clone := NewGraph[Key, Value]()
	for key := range g.vertices {
		clone.AddVertex(key, g.vertices[key].value)
	}
	for key := range g.vertices {
		for neighbor := range g.vertices[key].neighbors {
			clone.AddEdge(key, neighbor)
		}
	}
	return clone
}

func (v *Vertex[Key, Value]) GetKey() Key {
	return v.key
}

// GetValue returns the value of a vertex.
func (v *Vertex[Key, Value]) GetValue() Value {
	return v.value
}

// GetNeighbors returns all neighbors of a vertex.
func (v *Vertex[Key, Value]) GetNeighbors() []Key {
	neighbors := make([]Key, 0)
	for neighbor := range v.neighbors {
		neighbors = append(neighbors, neighbor)
	}
	return neighbors
}
