package algorithms

import (
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

var (
	testingGraph = generator.SimpleGraph{
		Size: 5,
		EdgesMap: []map[int]bool{
			{1: true, 2: true, 3: true, 4: true},
			{0: true, 2: true, 3: true, 4: true},
			{0: true, 1: true, 3: true, 4: true},
			{0: true, 1: true, 2: true, 4: true},
			{0: true, 1: true, 2: true, 3: true},
		},
	}
)

func checkWeights(t assert.TestingT, g generator.Graph, min, max int) {
	assert.True(t, g.Properties().Weighted())
	for i := range g.Edges() {
		for j := range g.Edges()[i] {
			if !g.Edges()[i][j] {
				continue
			}
			edge := generator.CreateEdge(i, j)
			v, ok := g.Weights()[edge]
			assert.True(t, ok)
			assert.NotEqual(t, 0, v)
			assert.LessOrEqual(t, min, v)
			assert.GreaterOrEqual(t, max, v)
		}
	}
}

func TestInvalidWeights(t *testing.T) {
	rnd := rand.New(rand.NewSource(25))
	_, err := GenerateWeights(testingGraph, 1, 0, rnd)
	assert.Error(t, err)
	_, err = GenerateWeights(testingGraph, 0, 0, rnd)
	assert.Error(t, err)
	_, err = GenerateWeights(testingGraph, -3, -5, rnd)
	assert.Error(t, err)
	_, err = GenerateWeights(testingGraph, 10, -20, rnd)
	assert.Error(t, err)
	_, err = GenerateWeights(testingGraph, 10, 20, nil)
	assert.Error(t, err)

}

func TestValidWeights(t *testing.T) {
	rnd := rand.New(rand.NewSource(25))
	res, err := GenerateWeights(testingGraph, 1, 1, rnd)
	assert.Nil(t, err)
	checkWeights(t, res, 1, 1)
}
