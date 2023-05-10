package algorithms

import (
	"errors"
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	mrand "math/rand"
)

func GenerateRandomBetween(nodes, minDegree, maxDegree int, connected bool, rand *mrand.Rand) (generator.SimpleGraph, error) {
	if nodes == 0 || minDegree > maxDegree || maxDegree >= nodes || (connected && maxDegree < 2) {
		return generator.SimpleGraph{}, errors.New("invalid ParentGraph request")
	}
	graph := make([]map[int]bool, nodes)
	filledNodes := make(map[int]bool)
	numOfEdges := 0
	if connected {
		spanDegree := minDegree
		if minDegree < 2 {
			spanDegree = 2
		}
		graph_r, err := GenerateSpanningBoruvka(nodes, spanDegree, rand)
		if err != nil {
			return generator.SimpleGraph{}, err
		}
		graph = graph_r.Edges()
		numOfEdges = nodes - 1
	} else {
		for k := range graph {
			graph[k] = make(map[int]bool)
		}
	}
	minEdges, maxEdges := (nodes*minDegree+1)/2, (nodes*maxDegree-1)/2
	expectedEdges := rand.Intn(maxEdges-minEdges) + minEdges
	currentDeg := extractNodeDegFromGraph(graph)
	pointsMin := createPoints(nodes, minDegree, currentDeg)
	for k := range currentDeg {
		if currentDeg[k] >= minDegree {
			currentDeg[k] -= minDegree
			filledNodes[k] = true
		} else {
			currentDeg[k] = 0
		}
	}
	pointsMax := createPoints(nodes, maxDegree-minDegree, currentDeg)
	counter := 0

	for pointsMin.Length() > 0 {
		leftInd, rightInd := rand.Intn(pointsMin.Length()), rand.Intn(pointsMin.Length())
		left, _ := pointsMin.GetPoint(leftInd)
		right, _ := pointsMin.GetPoint(rightInd)
		maxPoint := false

		if left == right || graph[left][right] {
			rightInd = rand.Intn(pointsMax.Length())
			right, _ = pointsMax.GetPoint(rightInd)
			if !filledNodes[right] {
				counter++
				continue
			}
			maxPoint = true
		}
		if left == right || graph[left][right] {
			counter++
			continue
		}

		graph[left][right] = true
		graph[right][left] = true
		numOfEdges++
		pointsMin.RemovePoint(left)
		if maxPoint {
			pointsMax.RemovePoint(right)
		} else {
			pointsMin.RemovePoint(right)
		}
		if len(graph[left]) >= minDegree {
			filledNodes[left] = true
		}
		if len(graph[right]) >= minDegree {
			filledNodes[right] = true
		}
	}

	counter = 0
	for numOfEdges < expectedEdges && pointsMax.Length() > 1 {

		leftInd, rightInd := rand.Intn(pointsMax.Length()), rand.Intn(pointsMax.Length())
		left, _ := pointsMax.GetPoint(leftInd)
		right, _ := pointsMax.GetPoint(rightInd)
		if left == right || graph[left][right] || len(graph[left]) >= maxDegree || len(graph[right]) >= maxDegree {
			if counter >= expectedEdges {
				break
			}
			continue
		}
		graph[left][right] = true
		graph[right][left] = true
		numOfEdges++
		pointsMax.RemovePoint(left)
		pointsMax.RemovePoint(right)
	}

	return generator.SimpleGraph{Size: nodes, EdgesMap: graph}, nil

}

func GenerateRandomAtLeast(nodes, degree int, connected bool, rand *mrand.Rand) (generator.SimpleGraph, error) {
	return GenerateRandomBetween(nodes, degree, nodes-1, connected, rand)
}
