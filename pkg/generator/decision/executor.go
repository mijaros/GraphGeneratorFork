package decision

import (
	"errors"
	"github.com/soch-fit/GraphGenerator/pkg/api"
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	"github.com/soch-fit/GraphGenerator/pkg/generator/algorithms"
	"math/rand"
)

func GenerateGraphFromRequest(request api.GraphRequest) (*api.GraphResult, error) {
	var graph generator.Graph = nil
	var err error = nil
	src := rand.NewSource(*request.Seed)
	rng := rand.New(src)
	switch request.Type {
	case api.ExactDeg:
		graph, err = algorithms.GenerateStegerWormald(request.Nodes, request.NodeDegree, request.Connected, rng)
	case api.BetweenDeg:
		graph, err = algorithms.GenerateRandomBetween(request.Nodes, request.NodeDegree, request.NodeDegreeMax, request.Connected, rng)
	case api.AtLeastDeg:
		graph, err = algorithms.GenerateRandomAtLeast(request.Nodes, request.NodeDegree, request.Connected, rng)
	case api.AverageDeg:
		graph, err = algorithms.GenerateRandomAverage(request.Nodes, request.NodeDegreeAverage, request.Connected, rng)
	case api.Complete:
		graph, err = algorithms.GenerateRandomComplete(request.Nodes)
	default:
		return nil, errors.New("invalid graph request")
	}

	if request.Weighted && err == nil {
		graph, _ = algorithms.GenerateWeights(graph, request.WeightMin, request.WeightMax, rng)
	}
	return &api.GraphResult{ID: request.ID, Generated: graph}, err
}
