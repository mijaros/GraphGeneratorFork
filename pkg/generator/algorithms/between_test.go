package algorithms

import (
	"fmt"
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func getRand(seedNum int64) *rand.Rand {
	seed := rand.NewSource(seedNum)
	mr := rand.New(seed)

	return mr
}

func checkDegreeBetween(t assert.TestingT, edges []map[int]bool, min, max int) {
	for _, v := range edges {
		assert.LessOrEqual(t, len(v), max)
		assert.GreaterOrEqual(t, len(v), min)
	}
}

func TestBasicRandomBetween(t *testing.T) {
	t.Parallel()
	randGen := getRand(122535)
	graph, err := GenerateRandomBetween(20, 4, 15, true, randGen)

	assert.Nil(t, err)

	CheckGraph(t, graph.Edges())
	checkDegreeBetween(t, graph.Edges(), 4, 15)
	CheckConnectivity(t, graph.Edges())
}

func TestBasicBetween(t *testing.T) {
	nodes := 10
	minDeg := 3
	maxDeg := 5

	seed := rand.NewSource(2353)
	rand := rand.New(seed)
	graph, err := GenerateRandomBetween(nodes, minDeg, maxDeg, true, rand)
	if err != nil {
		t.Errorf("Between generation failed %s", err.Error())
		t.Fail()
	}
	CheckGraph(t, graph.Edges())
	checkDegreeBetween(t, graph.Edges(), minDeg, maxDeg)
	CheckConnectivity(t, graph.Edges())
}

func TestBasicDifferentSeeds(t *testing.T) {
	nodes := 30
	minDeg := 3
	maxDeg := 5

	seeds := []int64{80, 31, 47, 83, 99293, 123123, 81928319, 1829910291}

	for _, k := range seeds {
		t.Run(fmt.Sprintf("Seed=%d", k), func(t *testing.T) {
			seed := rand.NewSource(k)
			rand := rand.New(seed)
			graph, err := GenerateRandomBetween(nodes, minDeg, maxDeg, true, rand)
			if err != nil {
				t.Errorf("Between generation failed %s", err.Error())
			}
			checkDegreeBetween(t, graph.Edges(), minDeg, maxDeg)
			CheckConnectivity(t, graph.Edges())
		})
	}
}

func BenchmarkBasicBetween(b *testing.B) {
	minDeg, maxDeg := 43, 55
	nodes := 100
	tester := func(conn bool) func(b *testing.B) {
		return func(b *testing.B) {

			for k := 0; k < b.N; k++ {
				seed := rand.NewSource(int64(k))
				gen := rand.New(seed)

				_, err := GenerateRandomBetween(nodes, minDeg, maxDeg, conn, gen)
				assert.Nil(b, err)
			}
		}
	}
	b.Run("connected", tester(true))
	b.Run("disconnected", tester(false))
}

var betweenBench generator.SimpleGraph

func BenchmarkGenerateRandomBetweenBaseProperties(b *testing.B) {
	minDeg, maxDeg := 499, 501
	nodes := 1000
	tester := func(conn bool) func(t *testing.B) {
		return func(b *testing.B) {
			var res generator.SimpleGraph
			var err error
			for k := 0; k < b.N; k++ {
				seed := rand.NewSource(int64(3353))
				gen := rand.New(seed)

				res, err = GenerateRandomBetween(nodes, minDeg, maxDeg, conn, gen)
				assert.Nil(b, err)
			}
			betweenBench = res
		}
	}
	b.Run("connected", tester(true))
	b.Run("disconnected", tester(false))
}

func TestExactBetweenGraphForSameSeed(t *testing.T) {
	seeds := []int64{1, 3, 5, 31, 97, 123, 531, 1129239443121}

	for _, seed := range seeds {
		t.Run(fmt.Sprintf("seed=%d", seed), func(t *testing.T) {
			src := rand.NewSource(seed)
			rnd := rand.New(src)

			graph, err := GenerateRandomBetween(36, 7, 15, true, rnd)
			assert.Nil(t, err)

			src2 := rand.NewSource(seed)
			rnd2 := rand.New(src2)

			graph2, err := GenerateRandomBetween(36, 7, 15, true, rnd2)
			assert.Nil(t, err)

			checkSameGraph(t, graph, graph2)
		})
	}
}
