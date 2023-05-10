package api

import (
	"bytes"
	"errors"
	"github.com/goccy/go-json"
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	"io"
	"time"
)

var (
	ErrInvalidGraphFormat   = errors.New("invalid string for graphFormat")
	ErrInvalidGraphType     = errors.New("invalid graph type passed")
	ErrInvalidRequestStatus = errors.New("invalid graph status passed")
)

type GraphTranslator interface {
	Convert(g generator.Graph) bool
	Serialize(writer io.Writer) (io.Writer, error)
	Bytes() []byte
	ContentType() string
	Extension() string
	Kind() string
}

type GraphResult struct {
	ID        uint32
	Generated generator.Graph
}

type GraphFormat uint8

const (
	JSON GraphFormat = iota
	Matrix
	Dot
)

var graphFormatString = map[GraphFormat]string{
	JSON:   "JSON",
	Dot:    "dot",
	Matrix: "matrix"}

var stringGraphFormat = map[string]GraphFormat{
	"JSON":   JSON,
	"dot":    Dot,
	"matrix": Matrix}

func (p GraphFormat) String() string {
	return graphFormatString[p]
}

func (p GraphFormat) MarshalJSON() ([]byte, error) {
	builder := bytes.Buffer{}
	builder.WriteByte('"')
	builder.WriteString(p.String())
	builder.WriteByte('"')
	return builder.Bytes(), nil
}

func (p *GraphFormat) UnmarshalJSON(in []byte) error {
	var str string
	if err := json.Unmarshal(in, &str); err != nil {
		return err
	}

	v, ok := stringGraphFormat[str]
	if !ok {
		return ErrInvalidGraphFormat
	}
	*p = v
	return nil
}

func (p GraphFormat) GetGraphRepre() GraphTranslator {
	switch p {
	case Matrix:
		return &MatrixGraph{}
	case JSON:
		return &BasicJSONGraph{}
	}
	return nil
}

type GraphType uint8

const (
	ExactDeg GraphType = iota
	AtLeastDeg
	BetweenDeg
	AverageDeg
	Complete
)

var graphToString = map[GraphType]string{
	ExactDeg:   "exact-degree",
	AtLeastDeg: "at-least-degree",
	BetweenDeg: "between-degree",
	AverageDeg: "average-degree",
	Complete:   "complete"}

var stringToGraph = map[string]GraphType{
	"exact-degree":    ExactDeg,
	"at-least-degree": AtLeastDeg,
	"between-degree":  BetweenDeg,
	"average-degree":  AverageDeg,
	"complete":        Complete}

func (g GraphType) String() string {
	return graphToString[g]
}

func (g GraphType) MarshalJSON() ([]byte, error) {
	buffer := bytes.Buffer{}
	buffer.WriteByte('"')
	buffer.WriteString(g.String())
	buffer.WriteByte('"')
	return buffer.Bytes(), nil
}

func (g *GraphType) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	var val GraphType
	var ok bool

	if val, ok = stringToGraph[str]; !ok {
		return ErrInvalidGraphType
	}

	*g = val
	return nil
}

type RequestStatus int

const (
	Undefined RequestStatus = iota
	NotFinished
	Finished
)

var statusToString = map[RequestStatus]string{
	Undefined:   "undefined",
	NotFinished: "not-finished",
	Finished:    "finished",
}

var stringToStatus = map[string]RequestStatus{
	"undefined":    Undefined,
	"not-finished": NotFinished,
	"finished":     Finished,
}

func (s RequestStatus) String() string {
	return statusToString[s]
}

func (s RequestStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.Buffer{}
	buffer.WriteByte('"')
	buffer.WriteString(s.String())
	buffer.WriteByte('"')
	return buffer.Bytes(), nil
}

func (s *RequestStatus) UnmarshalJSON(i []byte) error {
	var str string
	err := json.Unmarshal(i, &str)
	if err != nil {
		return err
	}

	var val RequestStatus
	var ok bool

	if val, ok = stringToStatus[str]; !ok {
		return ErrInvalidRequestStatus
	}

	*s = val
	return nil
}

type GraphRequest struct {
	Type              GraphType     `json:"type"`
	Weighted          bool          `json:"weighted"`
	Nodes             int           `json:"nodes"`
	NodeDegree        int           `json:"node_degree"`
	Status            RequestStatus `json:"status"`
	Seed              *int64        `json:"seed,omitempty"`
	Timeout           time.Time     `json:"deleted,omitempty"`
	NodeDegreeMax     int           `json:"node_degree_max,omitempty"`
	NodeDegreeAverage float32       `json:"node_degree_average,omitempty"`
	WeightMin         int           `json:"weight_min"`
	WeightMax         int           `json:"weight_max"`
	Connected         bool          `json:"connected"`
	ID                uint32        `json:"id"`
	Owner             *string       `json:"-"`
	BatchId           *uint32       `json:"-"`
}

type GraphBatchStatus struct {
	GraphId uint32        `json:"graph_id"`
	Status  RequestStatus `json:"status"`
}

type BatchRequest struct {
	BaseGraph GraphRequest  `json:"base"`
	Number    int           `json:"number"`
	ID        uint32        `json:"id,omitempty"`
	Timeout   time.Time     `json:"deleted,omitempty"`
	GraphsIDs []uint32      `json:"graph_ids,omitempty"`
	Status    RequestStatus `json:"status,omitempty"`
	Owner     *string       `json:"-"`
	Finished  int           `json:"-"`
}

type ErrorResponse struct {
	Error string  `json:"error"`
	Cause *string `json:"cause"`
}

type LimitsResponse struct {
	MaxNodes     int `json:"max_nodes"`
	MaxBatchSize int `json:"max_batch_size"`
}

type JSONGraph interface {
}

type BasicJSONGraph struct {
	Nodes []string            `json:"nodes"`
	Edges map[string][]string `json:"edges"`
}

type WeightedJSONGraph struct {
	Nodes []string `json:"nodes"`
	Edges map[string]map[string]int
}

type MatrixGraph struct {
	edges [][]int
}

type DotGraph struct {
	weighted bool
	size     int
	edges    map[generator.WeightedEdge]int
}
