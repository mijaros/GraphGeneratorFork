package api

import "github.com/soch-fit/GraphGenerator/pkg/configuration"

func (g *GraphRequest) validExactDeg() bool {
	return (g.Nodes*g.NodeDegree)%2 == 0 && g.NodeDegree > 0 && g.NodeDegree < g.Nodes
}

func (g *GraphRequest) validLimits() bool {
	return g.Nodes <= configuration.Default().MaxNodes && g.Nodes > 0
}

func (g *GraphRequest) validAverage() bool {
	return int(g.NodeDegreeAverage) < g.Nodes-1 && g.NodeDegreeAverage >= 0
}

func (g *GraphRequest) validBetweenDeg() bool {
	return g.NodeDegree < g.NodeDegreeMax && g.NodeDegreeMax < g.Nodes && g.NodeDegree >= 0
}

func (g *GraphRequest) validAtLeastDeg() bool {
	return g.NodeDegree < g.Nodes && g.NodeDegree >= 0
}

func (g *GraphRequest) validWeight() bool {
	return g.WeightMin <= g.WeightMax && !(g.WeightMin == 0 && g.WeightMax == 0)
}

func (g *GraphRequest) validConnected() bool {
	switch g.Type {
	case ExactDeg:
		return g.NodeDegree >= 2
	case AverageDeg:
		return int(g.NodeDegreeAverage) >= 2 || (g.Nodes <= 2 && g.NodeDegreeAverage == 1) || (g.Nodes == 1 && g.NodeDegreeAverage == 0)
	case BetweenDeg:
		return g.NodeDegree >= 0 && (g.NodeDegreeMax >= 2 || (g.NodeDegreeMax == 1 && g.Nodes == 2))
	}
	return true
}

func (g *GraphRequest) Valid() (result bool) {
	result = g.validLimits()
	switch g.Type {
	case ExactDeg:
		result = result && g.validExactDeg()
	case AverageDeg:
		result = result && g.validAverage()
	case BetweenDeg:
		result = result && g.validBetweenDeg()
	case AtLeastDeg:
		result = result && g.validAtLeastDeg()
	}

	if g.Weighted {
		result = result && g.validWeight()
	}

	if g.Connected {
		result = result && g.validConnected()
	}
	return
}
