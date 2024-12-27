package graph

import "fmt"

type NodeValue struct {
	Node              *Node
	Value             float64
	PreviousNode      *Node
	PreviousNodeValue *NodeValue
	EdgeTaken         *Edge
}

type TraversableNodeFunction func(n *Node) bool
type TraversableEdgeFunction func(e *Edge) bool

type EdgeNodeValueFunction func(e *Edge, nv NodeValue) float64

type VisitedNodes []*Node

type PosNodeMap[V comparable] map[V]*Node

func (vn VisitedNodes) Contains(n *Node) bool {
	for _, node := range vn {
		if node == n {
			return true
		}
	}
	return false
}

type Edge struct {
	source          *Node
	destination     *Node
	value           float64
	properties      map[string]interface{}
	traversable     bool
	traversableFunc *TraversableEdgeFunction
	nodeValueFunc   *EdgeNodeValueFunction
}

func (e *Edge) IsTraversable() bool {
	traversable := e.traversable
	if e.traversableFunc != nil {
		f := *(e.traversableFunc)
		traversable = f(e)
	}
	return traversable && e.destination != nil && e.destination.IsTraversable()
}

func (e *Edge) SetTraversable(b bool) *Edge {
	e.traversable = b
	return e
}

func (e *Edge) GetSource() *Node {
	return e.source
}

func (e *Edge) GetDestination() *Node {
	return e.destination
}

func (e *Edge) SetDestination(o *Node) *Edge {
	e.destination = o
	return e
}

func (e *Edge) GetValue() float64 {
	return e.value
}

func (e *Edge) GetNodeValue(nv NodeValue) float64 {
	if e.nodeValueFunc != nil {
		f := *(e.nodeValueFunc)
		return f(e, nv)
	}
	return e.GetValue()
}

func (e *Edge) SetNodeValueFunction(nvf *EdgeNodeValueFunction) *Edge {
	e.nodeValueFunc = nvf
	return e
}

func (e *Edge) AddProperty(id string, value interface{}) *Edge {
	e.properties[id] = value
	return e
}

func (e *Edge) GetProperty(id string) interface{} {
	if v, ok := e.properties[id]; ok {
		return v
	}
	return nil
}

type Node struct {
	id              interface{}
	edges           []*Edge
	properties      map[string]interface{}
	traversable     bool
	traversableFunc *TraversableNodeFunction
}

func (n Node) String() string {
	return fmt.Sprintf("%v, %v", n.id, n.properties)
}

func (n *Node) GetID() interface{} {
	return n.id
}

func (n *Node) GetEdges() []*Edge {
	edges := make([]*Edge, len(n.edges), len(n.edges))
	for i := range n.edges {
		edges[i] = n.edges[i]
	}
	return edges
}

func (n *Node) GetTraversableEdges() []*Edge {
	edges := make([]*Edge, 0, len(n.edges))
	for i := range n.edges {
		if n.edges[i].IsTraversable() {
			edges = append(edges, n.edges[i])
		}
	}
	return edges
}

func (n *Node) IsTraversable() bool {
	if n.traversableFunc != nil {
		f := *(n.traversableFunc)
		return f(n)
	}
	return n.traversable
}

func (n *Node) SetTraversable(b bool) *Node {
	n.traversable = b
	return n
}

func (n *Node) SetTraversableFunction(f TraversableNodeFunction) *Node {
	n.traversableFunc = &f
	return n
}

func (n *Node) AddProperty(id string, value interface{}) *Node {
	n.properties[id] = value
	return n
}

func (n *Node) GetProperty(id string) interface{} {
	if v, ok := n.properties[id]; ok {
		return v
	}
	return nil
}

func (n *Node) AddEdge(o *Node, w float64) *Edge {
	e := Edge{source: n, destination: o, value: w, traversable: true}
	e.properties = make(map[string]interface{})
	if n.edges == nil {
		n.edges = make([]*Edge, 0, 8)
	}
	n.edges = append(n.edges, &e)
	return &e
}

type Graph struct {
	nodes map[interface{}]*Node
}

func NewGraph() *Graph {
	g := new(Graph)
	g.nodes = make(map[interface{}]*Node)
	return g
}

func (g *Graph) Len() int {
	return len(g.nodes)
}

func (g *Graph) CreateNode(id interface{}) *Node {
	n := new(Node)
	n.id = id
	n.properties = make(map[string]interface{})
	n.traversable = true
	g.nodes[n.id] = n
	return n
}

func (g *Graph) GetNode(id interface{}) *Node {
	if n, ok := g.nodes[id]; ok {
		return n
	}
	return nil
}

func (g *Graph) GetNodes() []*Node {
	ns := make([]*Node, len(g.nodes), len(g.nodes))
	i := 0
	for _, n := range g.nodes {
		ns[i] = n
		i++
	}
	return ns
}

func (g *Graph) GetTraversableNodes() []*Node {
	ns := make([]*Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		if n.IsTraversable() {
			ns = append(ns, n)
		}
	}
	return ns
}

func (g *Graph) GetNonTraversableNodes() []*Node {
	ns := make([]*Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		if !n.IsTraversable() {
			ns = append(ns, n)
		}
	}
	return ns
}

func (g *Graph) Merge(og *Graph) {
	for id, n := range og.nodes {
		g.nodes[id] = n
	}
}

func (g *Graph) GetNodeCount() int {
	return len(g.nodes)
}
