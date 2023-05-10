package api

import (
	"bytes"
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	"io"
	"strconv"
)

func (m *MatrixGraph) Extension() string {
	return "txt"
}

func (m *MatrixGraph) Kind() string {
	return "matrix"
}
func (m *MatrixGraph) Convert(g generator.Graph) bool {
	m.edges = make([][]int, len(g.Edges()))
	for k := range m.edges {
		m.edges[k] = make([]int, len(g.Edges()))
		for j, v := range g.Edges()[k] {
			if !v {
				continue
			}
			if g.Properties().Weighted() {
				l, r := k, j
				if l > r {
					l, r = r, l
				}
				edge := generator.WeightedEdge{
					Left:  l,
					Right: r,
				}
				m.edges[k][j] = g.Weights()[edge]
			} else {
				m.edges[k][j] = 1
			}
		}
	}
	return true
}

func (m *MatrixGraph) Serialize(writer io.Writer) (io.Writer, error) {
	for k := range m.edges {
		for j, v := range m.edges[k] {
			if j != 0 {
				_, err := writer.Write([]byte{' '})
				if err != nil {
					return writer, err
				}
			}
			_, err := writer.Write([]byte(strconv.Itoa(v)))
			if err != nil {
				return writer, err
			}
		}
		_, err := writer.Write([]byte{'\n'})
		if err != nil {
			return writer, err
		}
	}
	return writer, nil
}

func (m *MatrixGraph) ContentType() string {
	return "text/plain"
}

func (m *MatrixGraph) Bytes() []byte {
	buffer := bytes.Buffer{}

	m.Serialize(&buffer)
	return buffer.Bytes()

}
