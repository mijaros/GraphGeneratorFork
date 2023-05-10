package algorithms

import (
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	mrand "math/rand"
)

// GenerateStegerWormald implements Steger-Wormalds algorithm to quickly generate
// regular graphs.
func GenerateStegerWormald(nodes, deg int, connected bool, rand *mrand.Rand) (generator.SimpleGraph, error) {
	if nodes <= 0 || (nodes*deg)%2 != 0 || deg >= nodes || (connected && deg < 2 && nodes > 2) {
		return generator.SimpleGraph{}, generator.ErrInvalidProperties
	}
	inverted := false
	if deg > nodes/2 {
		inverted = true
		deg = (nodes - deg) - 1
	}
	points, edges := createPoints(nodes, deg, make([]int, nodes)), nodes*deg

	tuples := make([]map[int]bool, nodes)
	for k := 0; k < nodes; k++ {
		tuples[k] = make(map[int]bool)
	}

	counter := 0
	for edges != 0 {
		leftInd, rightInd := rand.Intn(points.Length()), rand.Intn(points.Length())
		left, _ := points.GetPoint(leftInd)
		right, _ := points.GetPoint(rightInd)

		if left == right || tuples[left][right] {
			counter++
			if counter > 2*edges && (left != right || points.GetRank(left) > 1) {
				success := fixStegerWormald(nodes, &tuples, left, right, &edges, points, rand)
				if !success {
					return GenerateStegerWormald(nodes, deg, connected, rand)
				}
				counter = 0
			}
			continue
		}

		counter = 0
		tuples[left][right] = true
		tuples[right][left] = true
		edges -= 2
		points.RemovePoint(left)
		points.RemovePoint(right)
	}

	if inverted {
		tuples = invertGraph(tuples)
	}
	if connected && !inverted {
		tuples = makeRegularGraphConnected(tuples, rand)
	}

	return generator.SimpleGraph{EdgesMap: tuples, Size: nodes}, nil
}

func invertGraph(b []map[int]bool) []map[int]bool {
	res := make([]map[int]bool, len(b))
	for k := range b {
		res[k] = make(map[int]bool)
	}

	for i := range b {
		for j := range res {
			if i == j {
				continue
			}
			if !b[i][j] {
				res[i][j] = true
			}
		}
	}
	return res
}

// fixStegerWormald implements fixing switching algorithm that tries to unlock
// the invalid point set currently generated.
func fixStegerWormald(nodes int, tuples *[]map[int]bool, left, right int, edges *int,
	points *Tree, rand *mrand.Rand) bool {
	candidates := make([]int, 0, len(*tuples))
	for k := 0; k < nodes; k++ {
		if !((*tuples)[left][k] || (*tuples)[right][k]) {
			candidates = append(candidates, k)
		}
	}

	candidateEdges := make([][2]int, 0)
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			fst, snd := candidates[i], candidates[j]
			if !(*tuples)[fst][snd] {
				continue
			}
			candidateEdges = append(candidateEdges, [2]int{fst, snd})
		}
	}

	if len(candidateEdges) == 0 {
		return false
	}

	edge := rand.Intn(len(candidateEdges))
	fst, snd := candidateEdges[edge][0], candidateEdges[edge][1]

	delete((*tuples)[fst], snd)
	delete((*tuples)[snd], fst)
	(*tuples)[fst][left] = true
	(*tuples)[snd][right] = true
	(*tuples)[left][fst] = true
	(*tuples)[right][snd] = true
	*edges -= 2
	points.RemovePoint(left)
	points.RemovePoint(right)

	return true
}
