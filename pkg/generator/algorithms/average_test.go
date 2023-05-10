package algorithms

import (
	"fmt"
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"testing"
)

func checkAverageDegree(t *testing.T, averageDeg float64, edges []map[int]bool) {
	sum := 0
	size := len(edges)
	for _, k := range edges {
		sum += len(k)
	}
	actualDeg := float64(sum) / float64(size)
	numberOfEdges := int((averageDeg * float64(size)) / 2.0)
	allowedDiff := math.Abs(averageDeg*float64(size) - (2.0 * float64(numberOfEdges)))

	diff := math.Abs(averageDeg - actualDeg)
	assert.LessOrEqual(t, diff, allowedDiff)
}

func TestGenerateRandomAverage(t *testing.T) {
	nodes := 50
	avg := 6.8

	src := rand.NewSource(12)
	mrand := rand.New(src)
	graph, err := GenerateRandomAverage(nodes, float32(avg), true, mrand)
	assert.Nil(t, err)
	CheckGraph(t, graph.Edges())
	checkAverageDegree(t, avg, graph.Edges())
	CheckConnectivity(t, graph.Edges())
	graph, err = GenerateRandomAverage(nodes, float32(avg), false, mrand)
	assert.Nil(t, err)
	CheckGraph(t, graph.Edges())
	checkAverageDegree(t, avg, graph.Edges())
}

func TestGenerateRandomAverageBadInputs(t *testing.T) {
	inputs := []struct {
		node      int
		deg       float64
		connected bool
	}{
		{12, 11.8, false},
		{1, 1.1, true},
		{3, 1.1, true},
	}

	src := rand.NewSource(1)
	mrand := rand.New(src)

	for _, k := range inputs {
		_, err := GenerateRandomAverage(k.node, float32(k.deg), k.connected, mrand)
		assert.Error(t, err)
	}
}

func TestExactAverageGraphForSameSeed(t *testing.T) {
	seeds := []int64{1, 3, 5, 31, 97, 123, 531, 1129239443121}

	for _, seed := range seeds {
		t.Run(fmt.Sprintf("seed=%d", seed), func(t *testing.T) {
			src := rand.NewSource(seed)
			rnd := rand.New(src)

			graph, err := GenerateRandomAverage(36, 12.3, true, rnd)
			assert.Nil(t, err)

			src2 := rand.NewSource(seed)
			rnd2 := rand.New(src2)

			graph2, err := GenerateRandomAverage(36, 12.3, true, rnd2)
			assert.Nil(t, err)

			checkSameGraph(t, graph, graph2)
		})
	}
}

var averageBench generator.SimpleGraph

func BenchmarkGenerateRandomAverageBaseProperties(b *testing.B) {
	average := float32(500.0)
	nodes := 1000
	tester := func(conn bool) func(t *testing.B) {
		return func(b *testing.B) {
			var res generator.SimpleGraph
			var err error
			for k := 0; k < b.N; k++ {
				seed := rand.NewSource(int64(3353))
				gen := rand.New(seed)

				res, err = GenerateRandomAverage(nodes, average, conn, gen)
				assert.Nil(b, err)
			}
			averageBench = res
		}
	}
	b.Run("connected", tester(true))
	b.Run("disconnected", tester(false))
}
