package depgraph

import (
	"errors"
	"fmt"
)

// Our type consists of only a map of Nodes, indexed by strings. The graph's
// edge data is stored in the Nodes themselves, which makes cycle detection a
// bit easier.
type DependencyGraph struct {
	NodeMap map[string]*Node
}

// A Node of a directed graph, with incoming and outgoing edges.
type Node struct {
	Value    Keyer
	EdgesOut map[string]*Node
	EdgesIn  map[string]*Node
}

func (n *Node) Key() string {
	return n.Value.Key()
}

type Keyer interface {
	Key() string
}

type StringNode string

func (s StringNode) Key() string {
	return string(s)
}

// Add an incoming edge.  This needs to be paired with addEdgeOut.
func (n *Node) addEdgeIn(edgeNode *Node) {
	n.EdgesIn[edgeNode.Key()] = edgeNode
}

// Add an outgoing edge.  This needs to be paired with addEdgeIn.
func (n *Node) addEdgeOut(edgeNode *Node) {
	n.EdgesOut[edgeNode.Key()] = edgeNode
}

// Remove an incoming edge.  This needs to be paired with removeEdgeOut.
func (n *Node) removeEdgeIn(edgeNode *Node) {
	delete(n.EdgesIn, edgeNode.Key())
}

/*
// Remove an outgoing edge.  This needs to be paired with removeEdgeIn.
func (n *Node) removeEdgeOut(edgeNode *Node) {
	delete(n.EdgesOut, edgeNode.Key())
}
*/

// TopSort creates a topological sort of the Nodes of a Graph.  If there is a
// cycle, an error is returned, otherwise the topological sort is returned as a
// list of node names.
func (dg *DependencyGraph) TopSort() ([]string, error) {
	sorted := make([]string, 0)
	copy := dg.copy()

	// Initially, add all nodes without dependencies
	empty := make([]*Node, 0)
	for _, node := range copy.NodeMap {
		if len(node.EdgesIn) == 0 {
			empty = append(empty, node)
		}
	}

	for len(empty) > 0 {
		node := empty[0]
		sorted = append(sorted, node.Key())
		empty = empty[1:]
		for _, outgoing := range node.EdgesOut {
			// delete the edge from node -> outgoing
			outgoing.removeEdgeIn(node)
			if len(outgoing.EdgesIn) == 0 {
				empty = append(empty, outgoing)
			}
		}
		node.EdgesOut = nil
	}

	// if there are any edges left, we have a cycle
	for _, n := range copy.NodeMap {
		if len(n.EdgesIn) > 0 || len(n.EdgesOut) > 0 {
			return nil, errors.New("Cycle!")
		}
	}
	return sorted, nil
}

// Copy an existing graph into an independent structure
// (i.e. new nodes/edges are created - pointers aren't copied)
func (dg *DependencyGraph) copy() *DependencyGraph {
	// Copy nodes
	nodes := make([]Keyer, 0, len(dg.NodeMap))
	for _, node := range dg.NodeMap {
		nodes = append(nodes, node.Value)
	}
	other := New(nodes)

	// Copy edges
	for fromId, node := range dg.NodeMap {
		for toId := range node.EdgesOut {
			other.addEdge(other.NodeMap[fromId], other.NodeMap[toId])
		}
	}

	return other
}

// Create a new graph consisting of a set of nodes with no edges.
func New(nodes []Keyer) *DependencyGraph {
	dg := &DependencyGraph{}
	dg.NodeMap = make(map[string]*Node)

	for _, s := range nodes {
		dg.AddNode(s)
	}
	return dg
}

func (dg *DependencyGraph) AddNode(node Keyer) {
	dg.NodeMap[node.Key()] = &Node{
		Value:    node,
		EdgesIn:  make(map[string]*Node),
		EdgesOut: make(map[string]*Node),
	}
}

func (dg *DependencyGraph) addEdge(from *Node, to *Node) {
	from.addEdgeOut(to)
	to.addEdgeIn(from)
}

// Add an edge from <from> to <to>
func (g *DependencyGraph) AddEdge(from Keyer, to Keyer) error {
	if from == to {
		return errors.New("From node cannot be the same as To node")
	}

	fromNode, ok1 := g.NodeMap[from.Key()]
	toNode, ok2 := g.NodeMap[to.Key()]

	if !ok1 {
		return errors.New("from node not found")
	} else if !ok2 {
		return errors.New("to node not found")
	} else {
		g.addEdge(fromNode, toNode)
		return nil
	}
}

// A Dependency is a Dependent node that is DependentOn on another node.
type Dependency struct {
	Dependent   string
	DependentOn string
}

func (dg *DependencyGraph) AddDependency(dep *Dependency) error {
	from, ok := dg.NodeMap[dep.Dependent]
	if !ok {
		return fmt.Errorf("Cannot find dependent node '%v'", dep.Dependent)
	}

	to, ok := dg.NodeMap[dep.DependentOn]
	if !ok {
		return fmt.Errorf("Cannot find dependency node '%v'", dep.DependentOn)
	}

	return dg.AddEdge(from.Value, to.Value)
}

func (dg *DependencyGraph) AddDependencies(deps []*Dependency) error {
	for _, dep := range deps {
		if err := dg.AddDependency(dep); err != nil {
			return err
		}
	}
	return nil
}

func (dg *DependencyGraph) AddDependenciesForNode(dependentKey string, dependencyKeys []string) error {
	for _, dep := range dependencyKeys {
		if err := dg.AddDependency(&Dependency{dependentKey, dep}); err != nil {
			return err
		}
	}
	return nil
}
