package algorithms

import (
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	mrand "math/rand"
)

func valueGen(min, max int, rand *mrand.Rand) int {
	if min == max {
		return min
	}
	var value int = 0
	for value == 0 {
		value = rand.Intn(max-min) + min
	}
	return value
}

func GenerateWeights(graph generator.Graph, min, max int, rand *mrand.Rand) (generator.Graph, error) {
	if max < min {
		return graph, generator.ErrInvalidWeight
	}
	if rand == nil {
		return graph, generator.ErrMissingRand
	}
	if min == 0 && max == 0 {
		return graph, generator.ErrInvalidWeight
	}

	if graph.Properties().Weighted() {
		return graph, nil
	}
	result := generator.WeightedGraph{ParentGraph: graph, WeightsMap: make(map[generator.WeightedEdge]int)}

	dimensions := len(graph.Edges())

	for i := 0; i < dimensions; i++ {
		for j := i + 1; j < dimensions; j++ {
			if !result.Edges()[i][j] {
				continue
			}
			value := valueGen(min, max, rand)
			edge := generator.WeightedEdge{
				Left:  i,
				Right: j,
			}
			result.WeightsMap[edge] = value
		}
	}
	return result, nil
}
