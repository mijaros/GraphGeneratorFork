package algorithms

import (
	"errors"
	"github.com/soch-fit/GraphGenerator/pkg/generator"
)

func GenerateRandomComplete(nodes int) (generator.SimpleGraph, error) {
	if nodes <= 0 {
		return generator.SimpleGraph{}, errors.New("invalid ParentGraph request")
	}
	graph := generator.SimpleGraph{
		Size:     nodes,
		EdgesMap: nil,
	}

	edges := make([]map[int]bool, nodes)

	for k := range edges {
		edges[k] = make(map[int]bool)
		for j := range edges {
			if k == j {
				continue
			}
			edges[k][j] = true
		}
	}
	graph.EdgesMap = edges
	return graph, nil
}
