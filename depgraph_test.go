package depgraph

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGraphCycles(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	nodes := []Keyer{
		StringNode("shell"),
		StringNode("boot"),
		StringNode("badcall2"),
		StringNode("randcall"),
		StringNode("badcall"),
		StringNode("shll"),
		StringNode("badcall3"),
	}

	graph := New(nodes)
	assert.Equal(len(nodes), len(graph.NodeMap))
	if len(nodes) == len(graph.NodeMap) {
		for _, s := range nodes {
			var sn StringNode = s.(StringNode)
			key := string(sn)
			assert.NotNil(graph.NodeMap[key])
		}
	}

	graph.AddEdge(StringNode("shell"), StringNode("boot"))
	graph.AddEdge(StringNode("badcall"), StringNode("shell"))
	graph.AddEdge(StringNode("randcall"), StringNode("shell"))
	graph.AddEdge(StringNode("badcall2"), StringNode("badcall"))
	graph.AddEdge(StringNode("shll"), StringNode("boot"))
	graph.AddEdge(StringNode("badcall3"), StringNode("badcall2"))

	sorted, err := graph.TopSort()
	assert.Nil(err)
	t.Log(sorted)

	graph.AddEdge(StringNode("shell"), StringNode("badcall3"))
	_, err = graph.TopSort()
	assert.NotNil(err)
	t.Log(err)
}

func TestGraphForest(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	nodes := []Keyer{
		StringNode("shell"),
		StringNode("boot"),
		StringNode("badcall2"),
		StringNode("randcall"),
		StringNode("badcall"),
		StringNode("shll"),
		StringNode("badcall3"),
		StringNode("boot2"),
		StringNode("shell2"),
	}

	graph := New(nodes)
	graph.AddEdge(StringNode("shell"), StringNode("boot"))
	graph.AddEdge(StringNode("badcall"), StringNode("shell"))
	graph.AddEdge(StringNode("randcall"), StringNode("shell"))
	graph.AddEdge(StringNode("badcall2"), StringNode("badcall"))
	graph.AddEdge(StringNode("shll"), StringNode("boot"))
	graph.AddEdge(StringNode("badcall3"), StringNode("badcall2"))
	graph.AddEdge(StringNode("shell2"), StringNode("boot2"))

	sorted, err := graph.TopSort()
	assert.Nil(err)
	t.Log(sorted)

	graph.AddEdge(StringNode("boot2"), StringNode("shell2"))
	_, err = graph.TopSort()
	assert.NotNil(err)
	t.Log(err)
}

func TestAddDependencies(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	nodes := []Keyer{
		StringNode("A"),
		StringNode("B"),
		StringNode("C"),
		StringNode("D"),
		StringNode("E"),
		StringNode("F"),
	}

	dependencies := []*Dependency{
		&Dependency{"F", "A"},
		&Dependency{"F", "B"},
		&Dependency{"F", "C"},
		&Dependency{"F", "D"},
		&Dependency{"F", "E"},

		&Dependency{"E", "A"},
		&Dependency{"E", "B"},
		&Dependency{"E", "C"},
		&Dependency{"E", "D"},

		&Dependency{"D", "A"},
		&Dependency{"D", "B"},
		&Dependency{"D", "C"},

		&Dependency{"C", "A"},
		&Dependency{"C", "B"},

		&Dependency{"B", "A"},
	}

	graph := New(nodes)
	graph.AddDependencies(dependencies)

	sorted, err := graph.TopSort()
	assert.Nil(err)
	t.Log(sorted)

	assert.Equal([]string{"F", "E", "D", "C", "B", "A"}, sorted)

}

func TestAddErrors(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var err error

	nodes := []Keyer{
		StringNode("A"),
		StringNode("B"),
	}

	dg := New(nodes)

	err = dg.AddDependency(&Dependency{"A", "B"})
	assert.Nil(err)

	tests := [][]string{
		[]string{"A", "C"},
		[]string{"C", "A"},
		[]string{"A", "A"},
	}

	for _, test := range tests {
		err = dg.AddDependency(&Dependency{test[0], test[1]})
		assert.NotNil(err)

		err = dg.AddEdge(StringNode(test[0]), StringNode(test[1]))
		assert.NotNil(err)
	}

	err = dg.AddDependencies([]*Dependency{
		&Dependency{"A", "A"},
	})
	assert.NotNil(err)
}

func TestAddDependencyLists(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	nodes := []Keyer{
		StringNode("A"),
		StringNode("B"),
		StringNode("C"),
		StringNode("D"),
		StringNode("E"),
		StringNode("F"),
	}

	graph := New(nodes)
	var err error

	err = graph.AddDependenciesForNode("F", []string{"A", "B", "C", "D", "E"})
	assert.Nil(err)

	err = graph.AddDependenciesForNode("E", []string{"A", "B", "C", "D"})
	assert.Nil(err)

	err = graph.AddDependenciesForNode("D", []string{"A", "B", "C"})
	assert.Nil(err)

	err = graph.AddDependenciesForNode("C", []string{"A", "B"})
	assert.Nil(err)

	err = graph.AddDependenciesForNode("B", []string{"A"})
	assert.Nil(err)

	sorted, err := graph.TopSort()
	assert.Nil(err)
	t.Log(sorted)

	assert.Equal([]string{"F", "E", "D", "C", "B", "A"}, sorted)
}

func TestAddDependencyListsErr(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	nodes := []Keyer{
		StringNode("A"),
		StringNode("B"),
		StringNode("C"),
		StringNode("D"),
		StringNode("E"),
		StringNode("F"),
	}

	graph := New(nodes)
	err := graph.AddDependenciesForNode("F", []string{"A", "B", "C", "D", "E", "G"})
	assert.NotNil(err)
}
