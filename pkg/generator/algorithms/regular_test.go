package algorithms

import (
	"fmt"
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func serializeGeneratedGraph(g generator.Graph) {
	fmt.Printf("%#v", g.Edges())
}

func checkSameGraph(t assert.TestingT, first, second generator.Graph) {
	for k := range first.Edges() {
		for e, ok := range first.Edges()[k] {
			if !ok {
				continue
			}
			assert.True(t, second.Edges()[k][e])
		}
	}
}

func TestLimitValuesValid(t *testing.T) {
	t.Parallel()
	nodeDeg := []struct {
		node, deg int
		connected bool
	}{
		{1, 0, true},
		{1, 0, false},
		{2, 1, true},
		{2, 1, false},
		{3, 2, true},
		{3, 2, false},
		{4, 2, false},
		{4, 1, false},
	}

	for _, v := range nodeDeg {
		t.Run(fmt.Sprintf("n=%d,d=%d,c=%v", v.node, v.deg, v.connected), func(t *testing.T) {
			var seed int64 = int64(13)
			sd := rand.NewSource(seed)
			rnd := rand.New(sd)
			graph, err := GenerateStegerWormald(v.node, v.deg, v.connected, rnd)
			assert.Nil(t, err)
			assert.Equal(t, graph.Size, v.node)
		})
	}
}

func TestLimitValuesInvalid(t *testing.T) {
	t.Parallel()
	nodeDeg := []struct {
		node, deg int
		connected bool
	}{
		{3, 1, true},
		{3, 1, false},
		{4, 1, true},
		{6, 1, true},
	}
	for _, v := range nodeDeg {
		t.Run(fmt.Sprintf("n=%d,d=%d,c=%v", v.node, v.deg, v.connected), func(t *testing.T) {
			var seed int64 = int64(13)
			sd := rand.NewSource(seed)
			rnd := rand.New(sd)
			_, err := GenerateStegerWormald(v.node, v.deg, v.connected, rnd)
			assert.Error(t, err)
		})
	}
}

func TestKRegularBasic(t *testing.T) {
	t.Parallel()
	seeds := []int64{
		12,
		28394,
		291024,
		20394,
		421110,
		5595021111,
	}
	nodes, degs := 20, 5
	for _, seed := range seeds {
		t.Run(fmt.Sprintf("Seed=%d", seed), func(t *testing.T) {
			src := rand.NewSource(seed)
			mrand := rand.New(src)
			graph, err := GenerateStegerWormald(nodes, degs, true, mrand)
			assert.Nil(t, err)
			CheckGraph(t, graph.Edges())
			checkGraphDegrees(t, graph, nodes, degs)
			CheckConnectivity(t, graph.Edges())
		})
	}
}

func checkGraphDegrees(t *testing.T, graph generator.SimpleGraph, nodes int, degs int) {
	if len(graph.Edges()) != nodes {
		t.Error("Invalid number of nodes in graph")
	}
	for _, v := range graph.Edges() {
		assert.Equal(t, len(v), degs)
	}
}

func generatorRunner(node, degree int, seed int64, connected bool) func(t *testing.T) {
	return func(t *testing.T) {
		src := rand.NewSource(seed)
		mrand := rand.New(src)
		graph, err := GenerateStegerWormald(node, degree, connected, mrand)

		if (node*degree)%2 != 0 || degree >= node || (degree == 1 && connected) {
			if err == nil {
				t.Errorf("graph generation didn't fail even though it should for n:%d,d:%d", node, degree)

			}
			return
		}
		if err != nil {
			t.Errorf("graph generation failed %s", err.Error())
			return
		}
		CheckGraph(t, graph.Edges())
		checkGraphDegrees(t, graph, node, degree)
		if connected {
			CheckConnectivity(t, graph.Edges())
		}
	}
}

func TestKRegularLowConnectivity(t *testing.T) {
	t.Parallel()
	nodes := []int{
		10, 15, 20, 50, 33, 37, 18,
	}
	degs := []int{
		1, 2, 3, 5, 17, 23, 55, 12,
	}

	for _, node := range nodes {
		for _, degree := range degs {
			t.Run(fmt.Sprintf("Node=%d,Deg=%d", node, degree), func(t *testing.T) {
				src := rand.NewSource(1)
				mrand := rand.New(src)
				graph, err := GenerateStegerWormald(node, degree, true, mrand)

				if (node*degree)%2 != 0 || degree >= node || degree == 1 {
					assert.Error(t, err)
					return
				}
				assert.Nil(t, err)
				CheckGraph(t, graph.Edges())
				checkGraphDegrees(t, graph, node, degree)
				CheckConnectivity(t, graph.Edges())
			})
		}
	}
}

func TestBasicUsage(t *testing.T) {
	t.Parallel()
	numMin, numMax := 3, 100
	var seed int64 = 138922021

	for node := numMin; node <= numMax; node++ {
		t.Run(fmt.Sprintf("Basic %d", node), func(t *testing.T) {
			for deg := 1; deg < node; deg++ {
				generatorRunner(node, deg, seed, true)
				generatorRunner(node, deg, seed, false)
			}
		})

	}
}

var regularBench generator.SimpleGraph

func BenchmarkRegularBaseProperties(b *testing.B) {
	nodes := 1000
	degrees := 500
	tester := func(conn bool) func(b *testing.B) {
		return func(b *testing.B) {
			var err error
			var res generator.SimpleGraph
			for k := 0; k < b.N; k++ {
				src := rand.NewSource(int64(3353))
				mrand := rand.New(src)
				res, err = GenerateStegerWormald(nodes, degrees, conn, mrand)
				assert.Nil(b, err)
			}
			regularBench = res
		}
	}
	b.Run("connected", tester(true))
	b.Run("disconnected", tester(false))

}

func BenchmarkRegularLowNodes(b *testing.B) {
	nodes := 1000
	degrees := 3
	for k := 0; k < b.N; k++ {
		src := rand.NewSource(int64(k))
		mrand := rand.New(src)
		_, err := GenerateStegerWormald(nodes, degrees, true, mrand)
		assert.Nil(b, err)
	}
}

func BenchmarkRegularExpectedMax(b *testing.B) {
	nodes := 100
	degrees := 50
	for k := 0; k < b.N; k++ {
		src := rand.NewSource(int64(k))
		mrand := rand.New(src)
		_, err := GenerateStegerWormald(nodes, degrees, true, mrand)
		assert.Nil(b, err)
	}
}

func TestHighValuesKRegular(t *testing.T) {
	t.Parallel()
	nodes := []int{
		250,
		500,
		1000,
		2000,
	}
	degs := [][]int{
		{2, 3, 10, 125, 240, 247, 248},
		{2, 3, 20, 250, 480, 497, 498},
		{2, 3, 40, 500, 960, 997, 998},
		{2, 3, 80, 1000, 1920, 1997, 1998},
	}

	for k, node := range nodes {
		t.Run(fmt.Sprintf("Node=%d", node), func(t *testing.T) {
			for j, deg := range degs[k] {

				src := rand.NewSource(1337)
				mrand := rand.New(src)
				graph, err := GenerateStegerWormald(node, deg, true, mrand)
				if err != nil {
					t.Errorf("Iteration %d:%d failed with error %s", k, j, err.Error())
					return
				}
				checkGraphDegrees(t, graph, node, deg)
				CheckConnectivity(t, graph.Edges())
			}
		})
	}
}

func TestExactRegularGraphForSameSeed(t *testing.T) {
	seeds := []int64{1, 3, 5, 31, 97, 123, 531, 1129239443121}

	for _, seed := range seeds {
		t.Run(fmt.Sprintf("seed=%d", seed), func(t *testing.T) {
			src := rand.NewSource(seed)
			rnd := rand.New(src)

			graph, err := GenerateStegerWormald(36, 7, true, rnd)
			assert.Nil(t, err)

			src2 := rand.NewSource(seed)
			rnd2 := rand.New(src2)

			graph2, err := GenerateStegerWormald(36, 7, true, rnd2)
			assert.Nil(t, err)

			checkSameGraph(t, graph, graph2)
		})
	}
}
