package api

import (
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	"io"
)

func (j *BasicJSONGraph) Extension() string {
	return "json"
}

func (j *BasicJSONGraph) Kind() string {
	return "JSON"
}

func (j *BasicJSONGraph) Convert(g generator.Graph) bool {
	j.Nodes = make([]string, len(g.Nodes()))
	copy(j.Nodes, g.Nodes())

	j.Edges = make(map[string][]string)
	for from, to := range g.Edges() {
		fromName := g.Nodes()[from]
		namedEdges := make([]string, len(to))
		for i := range to {
			namedEdges[i] = g.Nodes()[i]
		}
		j.Edges[fromName] = namedEdges
	}

	return true
}

func (j *BasicJSONGraph) ContentType() string {
	return "application/json"
}

//func (j *BasicJSONGraph) Serialize() []byte {
//	b, _ := json.Marshal(j)
//	return b
//}

func (j *WeightedJSONGraph) ContentType() string {
	return "application/json"
}

func (j *WeightedJSONGraph) Convert(g generator.Graph) bool {
	j.Nodes = make([]string, len(g.Nodes()))
	copy(j.Nodes, g.Nodes())

	return false
}

func (j *BasicJSONGraph) Serialize(writer io.Writer) (io.Writer, error) {
	//TODO implement me
	panic("implement me")
}

func (j *BasicJSONGraph) Bytes() []byte {
	//TODO implement me
	panic("implement me")
}
