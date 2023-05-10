package algorithms

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestGenerateRandomComplete(t *testing.T) {
	if testing.Short() {
		t.Skip("Long test")
	}
	nodes := []int{100, 1000, 2000, 531, 243, 98}

	for _, node := range nodes {
		t.Run(fmt.Sprintf("n=%d", node), func(t *testing.T) {
			graph, err := GenerateRandomComplete(node)
			assert.Nil(t, err)

			checkGraphDegrees(t, graph, node, node-1)

			var buff bytes.Buffer
			enc := gob.NewEncoder(&buff)
			enc.Encode(graph)

			fmt.Println(node, buff.Len())

			seed := rand.NewSource(2353)
			rnd := rand.New(seed)

			graphWeighted, err := GenerateWeights(graph, -10, 10, rnd)
			assert.Nil(t, err)
			checkWeights(t, graphWeighted, -10, 10)

			var buff2 bytes.Buffer

			enc = gob.NewEncoder(&buff2)
			enc.Encode(graphWeighted)

			fmt.Println(node, buff2.Len())
		})
	}
}

func BenchmarkGenerateRandomComplete(b *testing.B) {

}

var completeBench generator.SimpleGraph

func BenchmarkGenerateCompleteBaseProperties(b *testing.B) {
	nodes := 1000
	tester := func(t *testing.B) {
		var res generator.SimpleGraph
		var err error
		for k := 0; k < b.N; k++ {
			res, err = GenerateRandomComplete(nodes)
			assert.Nil(b, err)
		}
		completeBench = res
	}

	b.Run("connected", tester)
}
