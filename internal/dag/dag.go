package dag

import (
	"fmt"

	"github.com/dominikbraun/graph"
	"github.com/terrapilot/terrapilot/internal/descriptor"
)

// Build constructs a DAG from a list of stacks and returns them in topological order.
// Stacks are uniquely identified by their directory path.
func Build(stacks []*descriptor.Stack) ([]*descriptor.Stack, error) {
	g := graph.New(func(s *descriptor.Stack) string { return s.Dir }, graph.Directed(), graph.Acyclic())

	for _, s := range stacks {
		if err := g.AddVertex(s); err != nil {
			return nil, fmt.Errorf("duplicate stack directory %q", s.Dir)
		}
	}

	index := make(map[string]*descriptor.Stack, len(stacks))
	for _, s := range stacks {
		index[s.Dir] = s
	}

	for _, s := range stacks {
		for _, dep := range s.DependsOn {
			if _, ok := index[dep.Path]; !ok {
				return nil, fmt.Errorf("stack %q depends on %q which does not exist", s.Name, dep.Path)
			}
			if err := g.AddEdge(dep.Path, s.Dir); err != nil {
				return nil, fmt.Errorf("dependency error: %w", err)
			}
		}
	}

	order, err := graph.TopologicalSort(g)
	if err != nil {
		return nil, fmt.Errorf("cycle detected in stack dependencies: %w", err)
	}

	sorted := make([]*descriptor.Stack, 0, len(order))
	for _, dir := range order {
		sorted = append(sorted, index[dir])
	}

	return sorted, nil
}
