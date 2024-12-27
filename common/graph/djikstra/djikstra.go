package djikstra

import (
	hp "container/heap"
	"github.com/mbordner/aoc2018/common/graph"
	"math"
)

type NodeValue struct {
	graph.NodeValue
	visited bool
}

// ShortestPaths will hold all the nodes, visited and unvisited
type ShortestPaths map[interface{}]*NodeValue

func (sps ShortestPaths) GetShortestPathWithEdges(n *graph.Node) ([]*graph.Node, []*graph.Edge, float64) {
	if _, ok := sps[n.GetID()]; !ok {
		return nil, nil, float64(0)
	}

	current := sps[n.GetID()].Node
	value := float64(0)
	if sps[current.GetID()].PreviousNode != nil {
		value = sps[current.GetID()].Value
	}

	nodes := make([]*graph.Node, 0, 50)
	edges := make([]*graph.Edge, 0, 50)

	for sps[current.GetID()].PreviousNode != nil {
		if current.IsTraversable() == false {
			return []*graph.Node{}, []*graph.Edge{}, float64(0)
		}
		nodes = append(nodes, current)
		edges = append(edges, sps[current.GetID()].EdgeTaken)
		current = sps[current.GetID()].PreviousNode
	}

	// reverse the array
	for i, j, h := 0, len(nodes)-1, len(nodes)/2; i < h; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
		edges[i], edges[j] = edges[j], edges[i]
	}

	return nodes, edges, value
}

func (sps ShortestPaths) GetShortestPath(n *graph.Node) ([]*graph.Node, float64) {
	nodes, _, value := sps.GetShortestPathWithEdges(n)
	return nodes, value
}

// this will hold all the unvisited node values sorted with minimum values at the top
type nodeValues []*NodeValue

func (h nodeValues) Len() int {
	return len(h)
}
func (h nodeValues) Less(i, j int) bool {
	return h[i].Value < h[j].Value
}
func (h nodeValues) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *nodeValues) Push(nv interface{}) {
	*h = append(*h, nv.(*NodeValue))
}

func (h *nodeValues) Pop() interface{} {
	nv := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return nv
}

type heap struct {
	values nodeValues
}

func newNodeValueHeap(capacity int) *heap {
	h := new(heap)
	h.values = make(nodeValues, 0, capacity)
	return h
}

func (h *heap) index(nv *NodeValue) int {
	for i := range h.values {
		if nv == h.values[i] {
			return i
		}
	}
	return -1
}

func (h *heap) remove(nv *NodeValue) {
	i := h.index(nv)
	if i != -1 {
		hp.Remove(&h.values, i)
	}
}

func (h *heap) fix(nv *NodeValue) {
	i := h.index(nv)
	if i != -1 {
		hp.Fix(&h.values, i)
	}
}

func (h *heap) push(nv *NodeValue) {
	hp.Push(&h.values, nv)
}

func (h *heap) pop() *NodeValue {
	i := hp.Pop(&h.values)
	return i.(*NodeValue)
}

func GenerateShortestPaths(g *graph.Graph, source *graph.Node) ShortestPaths {
	// shortest paths from n to all other nodes
	sps := make(ShortestPaths)

	// node value heap used to sort current distances through nodes
	nvh := newNodeValueHeap(g.Len())

	for _, node := range g.GetTraversableNodes() {
		nv := &NodeValue{NodeValue: graph.NodeValue{Node: node, Value: math.MaxFloat64, PreviousNode: nil, PreviousNodeValue: nil, EdgeTaken: nil}, visited: false}
		if node == source {
			// this is our source node, and we need to treat it differently
			nv.Value = float64(0)
		}
		// add to the heap
		nvh.push(nv)
		sps[nv.Node.GetID()] = nv
	}

	// at this point, all unvisited nodes are in the heap

	for nvh.values.Len() > 0 {
		current := nvh.pop()

		for _, e := range current.Node.GetTraversableEdges() {

			if _, ok := sps[e.GetDestination().GetID()]; ok {
				// we don't want to explore edge destinations that already have been visited, i.e. removed from nvh
				if sps[e.GetDestination().GetID()].visited == false {

					// get the edge node value (env) for this edge's destination node
					env := sps[e.GetDestination().GetID()]

					// this value is the cost up to current node + cost to destination from current
					value := current.Value + e.GetNodeValue(current.NodeValue)

					// check if this new value is less than anything we found before
					if value < env.Value {
						// we found a shorter path from source -> e.destination through current

						env.Value = value
						env.PreviousNode = current.Node
						env.PreviousNodeValue = &current.NodeValue
						env.EdgeTaken = e
						// need to reorder the heap after this change
						nvh.fix(env)
					}

				}
			}

		}

		// at this point, current is marked as visited, and will remain removed
		// we are marking it visited, so we don't explore it again from another node's edges.
		sps[current.Node.GetID()].visited = true
	}

	return sps
}
