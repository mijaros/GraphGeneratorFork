package api

import (
	"bytes"
	"fmt"
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	"io"
)

func (d *DotGraph) Extension() string {
	return "dot"
}

func (d *DotGraph) Kind() string {
	return "graphviz-dot"
}

func (d *DotGraph) Convert(g generator.Graph) bool {
	if g.Properties().Weighted() {
		d.weighted = true
	}
	d.size = len(g.Edges())
	d.edges = make(map[generator.WeightedEdge]int)
	localWeights := g.Weights()
	localEdges := g.Edges()
	for k := range localEdges {
		for f, ok := range localEdges[k] {
			if !ok || f <= k {
				continue
			}
			edge := generator.WeightedEdge{
				Left:  k,
				Right: f,
			}
			if d.weighted {
				d.edges[edge] = localWeights[edge]
			} else {
				d.edges[edge] = 1
			}
		}
	}
	return true
}

func (d *DotGraph) Serialize(writer io.Writer) (io.Writer, error) {
	_, err := writer.Write([]byte("strict graph {\n"))
	if err != nil {
		return writer, err
	}
	foundVertices := make(map[int]bool)
	for k, v := range d.edges {
		if v == 0 && !d.weighted {
			continue
		}
		line := fmt.Sprintf("\t%d -- %d", k.Left, k.Right)
		foundVertices[k.Left] = true
		foundVertices[k.Right] = true
		writer.Write([]byte(line))
		if d.weighted {
			weight := fmt.Sprintf(` [label="%d"]`, v)
			writer.Write([]byte(weight))
		}
		writer.Write([]byte("\n"))
	}
	for k := 0; k < d.size; k++ {
		if ok, ex := foundVertices[k]; ok && ex {
			continue
		}
		line := fmt.Sprintf("\t%d\n", k)
		writer.Write([]byte(line))
	}
	writer.Write([]byte("}\n"))

	return writer, nil
}

func (d *DotGraph) Bytes() []byte {
	result := bytes.Buffer{}
	d.Serialize(&result)
	return result.Bytes()
}

func (d *DotGraph) ContentType() string {
	return "text/plain"
}
