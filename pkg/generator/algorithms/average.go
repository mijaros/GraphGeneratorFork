package algorithms

import (
	"errors"
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	mrand "math/rand"
)

func GenerateRandomAverage(nodes int, degree float32, connected bool, rand *mrand.Rand) (generator.SimpleGraph, error) {
	if nodes == 0 || (nodes > 2 && degree < float32(2) || (int(degree) >= (nodes - 1))) {
		return generator.SimpleGraph{}, errors.New("invalid ParentGraph request")
	}

	targetEdges := int((float32(nodes) * degree) / 2.0)
	numOfEdges := 0
	var edges []map[int]bool
	if connected {
		graph, err := GenerateSpanningBoruvka(nodes, nodes-1, rand)
		if err != nil {
			return generator.SimpleGraph{}, err
		}
		edges = graph.Edges()
		numOfEdges = nodes - 1
	} else {
		edges = make([]map[int]bool, nodes)
		for k := 0; k < nodes; k++ {
			edges[k] = make(map[int]bool)
		}
	}

	for numOfEdges < targetEdges {
		leftInd := rand.Intn(nodes)
		rightInd := rand.Intn(nodes)
		if leftInd == rightInd {
			continue
		}
		if edges[leftInd][rightInd] {
			continue
		}
		numOfEdges += 1
		edges[leftInd][rightInd] = true
		edges[rightInd][leftInd] = true
	}

	result := generator.SimpleGraph{Size: nodes, EdgesMap: edges}
	return result, nil

}
