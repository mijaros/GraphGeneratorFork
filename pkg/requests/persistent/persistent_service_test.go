package persistent

import (
	"github.com/soch-fit/GraphGenerator/pkg/api"
	"github.com/soch-fit/GraphGenerator/pkg/configuration"
	"github.com/soch-fit/GraphGenerator/pkg/generator"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	configuration.SetupTestingEnv()
}

type GeneratorSvcMock struct {
	pushedGraphRequests []api.GraphRequest
	pushedBatchRequests []api.BatchRequest
	commChan            chan *api.GraphResult
	pauseCounter        int
	band                int
}

func buildGenSvcMock() GeneratorSvcMock {
	return GeneratorSvcMock{
		pushedGraphRequests: make([]api.GraphRequest, 0),
		pushedBatchRequests: make([]api.BatchRequest, 0),
		commChan:            make(chan *api.GraphResult, 0),
		pauseCounter:        0,
		band:                10,
	}
}

func (g *GeneratorSvcMock) Start() error {
	//TODO implement me
	panic("implement me")
}

func (g *GeneratorSvcMock) Pause() error {
	//TODO implement me
	panic("implement me")
}

func (g *GeneratorSvcMock) Resume() error {
	//TODO implement me
	panic("implement me")
}

func (g *GeneratorSvcMock) PushRequest(req api.GraphRequest) error {
	g.pushedGraphRequests = append(g.pushedGraphRequests, req)
	return nil
}

func (g *GeneratorSvcMock) GetRetriever() chan *api.GraphResult {
	return g.commChan
}

func (g *GeneratorSvcMock) Stop() error {
	//TODO implement me
	panic("implement me")
}

func (g *GeneratorSvcMock) FreeBand() int {
	return g.band
}

func TestBasicStore(t *testing.T) {
	if testing.Short() {
		t.Skip("Long test skipping")
	}
	dbRoot := t.TempDir()
	configuration.SetTestingDBRoot(dbRoot)
	genService := buildGenSvcMock()
	ps, err := New(&genService)
	assert.Nil(t, err)
	err = ps.Start()
	assert.Nil(t, err)
	request := api.GraphRequest{
		Type:              api.AverageDeg,
		Weighted:          true,
		Nodes:             10,
		NodeDegreeAverage: 2.7,
		WeightMin:         24,
		WeightMax:         39,
		Connected:         false,
	}
	resultGraph := api.GraphResult{
		ID: 0,
		Generated: generator.SimpleGraph{
			Size: 25,
			EdgesMap: []map[int]bool{
				{2: true, 4: true}}}}
	req, err := ps.StoreNewRequest(request)
	assert.Nil(t, err)
	assert.Len(t, genService.pushedGraphRequests, 1)
	assert.Equal(t, req, genService.pushedGraphRequests[0])
	assert.NotNil(t, req.Seed)
	reqRes, err := ps.GetGraphRequest(req.ID)
	assert.Nil(t, err)
	assert.Equal(t, req.ID, reqRes.ID)
	assert.Equal(t, req.Status, reqRes.Status)
	assert.Equal(t, *req.Seed, *reqRes.Seed)
	resultGraph.ID = req.ID
	genService.commChan <- &resultGraph
	time.Sleep(1 * time.Second)
	reqRes, err = ps.GetGraphRequest(req.ID)
	assert.Nil(t, err)
	assert.Equal(t, api.Finished, reqRes.Status)
	genGraph, err := ps.GetGraph(req.ID)
	assert.Nil(t, err)
	assert.Equal(t, genGraph, resultGraph)
	err = ps.Stop()
	assert.Nil(t, err)
}

func TestBatchStore(t *testing.T) {
	request := api.BatchRequest{
		BaseGraph: api.GraphRequest{
			Type:              api.AverageDeg,
			Weighted:          true,
			Nodes:             10,
			NodeDegreeAverage: 2.7,
			WeightMin:         24,
			WeightMax:         39,
			Connected:         false,
		},
		Number: 5,
	}

	if testing.Short() {
		t.Skip("Long test skipping")
	}
	dbRoot := t.TempDir()
	configuration.SetTestingDBRoot(dbRoot)
	genService := buildGenSvcMock()
	ps, err := New(&genService)
	assert.Nil(t, err)
	err = ps.Start()
	assert.Nil(t, err)

	bat, err := ps.StoreNewBatch(request)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, bat.ID)

	//assert.Len(t, genService.pushedBatchRequests, 1)
	assert.Len(t, genService.pushedGraphRequests, request.Number)

}
