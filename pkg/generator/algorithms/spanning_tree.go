package algorithms

import (
	"errors"
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	mrand "math/rand"
)

func GenerateSpanningBoruvka(nodes, maxDegree int, rand *mrand.Rand) (generator.SimpleGraph, error) {
	if nodes == 0 || (nodes > 2 && maxDegree < 2) {
		return generator.SimpleGraph{}, errors.New("invalid ParentGraph request")
	}
	components := make([][]int, nodes)
	edges := make([]map[int]bool, nodes)
	for k := range components {
		components[k] = make([]int, 1)
		components[k][0] = k
		edges[k] = make(map[int]bool)
	}

	for len(components) > 1 {
		leftInd := rand.Intn(len(components))
		rightInd := rand.Intn(len(components))
		if leftInd == rightInd {
			continue
		}
		if leftInd > rightInd {
			leftInd, rightInd = rightInd, leftInd
		}

		leftComp, rightComp := components[leftInd], components[rightInd]
		components = removeIndex(rightInd, components)
		leftWidth, rightWidth := len(leftComp), len(rightComp)
		firstNode := leftComp[rand.Intn(leftWidth)]
		secondNode := rightComp[rand.Intn(rightWidth)]
		for len(edges[firstNode]) >= maxDegree || len(edges[secondNode]) >= maxDegree {
			firstNode = leftComp[rand.Intn(leftWidth)]
			secondNode = rightComp[rand.Intn(rightWidth)]
		}
		edges[firstNode][secondNode] = true
		edges[secondNode][firstNode] = true
		components[leftInd] = append(leftComp, rightComp...)
	}

	return generator.SimpleGraph{Size: nodes, EdgesMap: edges}, nil
}
