package generator

import (
	"encoding/gob"
	"errors"
)

// WeightedEdge associates nodes into one edge
type WeightedEdge struct {
	Left, Right int
}

// Graph Basic interface for representation of a graph
type Graph interface {
	Nodes() []string
	Edges() []map[int]bool
	Weights() map[WeightedEdge]int
	Properties() GraphProperties
}

func init() {
	gob.Register(SimpleGraph{})
	gob.Register(WeightedGraph{})
	gob.Register(NamedGraph{})
}

var (
	ErrInvalidWeight     = errors.New("invalid weights")
	ErrInvalidProperties = errors.New("invalid graph properties")
	ErrMissingRand       = errors.New("rand must be passed")
)

// GraphProperties represents properties of ParentGraph
type GraphProperties uint8

const (
	NONE GraphProperties = 1 << iota
	NAMED
	WEIGHTED
)

func (g GraphProperties) Weighted() bool {
	return g&WEIGHTED != 0
}

func (g GraphProperties) Named() bool {
	return g&NAMED != 0
}

type SimpleGraph struct {
	Size     int
	EdgesMap []map[int]bool
}

func (b SimpleGraph) Nodes() []string {
	return []string{}
}

func (b SimpleGraph) Properties() GraphProperties {
	return NONE
}

func (b SimpleGraph) Edges() []map[int]bool {
	return b.EdgesMap
}

func (b SimpleGraph) Weights() map[WeightedEdge]int {
	return map[WeightedEdge]int{}
}

type NamedGraph struct {
	ParentGraph Graph
	VertexNames []string
}

func (n NamedGraph) Nodes() []string {
	return n.VertexNames
}

func (n NamedGraph) Edges() []map[int]bool {
	return n.ParentGraph.Edges()
}

func (n NamedGraph) Weights() map[WeightedEdge]int {
	return n.ParentGraph.Weights()
}

func (n NamedGraph) Properties() GraphProperties {
	return n.ParentGraph.Properties() | NAMED
}

type WeightedGraph struct {
	ParentGraph Graph
	WeightsMap  map[WeightedEdge]int
}

func (w WeightedGraph) Nodes() []string {
	return w.ParentGraph.Nodes()
}

func (w WeightedGraph) Edges() []map[int]bool {
	return w.ParentGraph.Edges()
}

func (w WeightedGraph) Weights() map[WeightedEdge]int {
	return w.WeightsMap
}

func (w WeightedGraph) Properties() GraphProperties {
	return w.ParentGraph.Properties() | WEIGHTED
}

func CreateEdge(u, v int) WeightedEdge {
	if v < u {
		u, v = v, u
	}
	return WeightedEdge{
		Left:  u,
		Right: v,
	}
}
