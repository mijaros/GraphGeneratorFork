package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/soch-fit/GraphGenerator/pkg/api"
	"testing"
)

type RequestServiceMock struct {
}

func (r RequestServiceMock) StoreNewRequest(request api.GraphRequest) (api.GraphRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (r RequestServiceMock) StoreNewBatch(request api.BatchRequest) (api.BatchRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (r RequestServiceMock) StoreGraph(graph *api.GraphResult) error {
	//TODO implement me
	panic("implement me")
}

func (r RequestServiceMock) ListRequests(sessionId string) ([]uint32, error) {
	//TODO implement me
	panic("implement me")
}

func (r RequestServiceMock) ListBatches(sessionId string) ([]uint32, error) {
	//TODO implement me
	panic("implement me")
}

func (r RequestServiceMock) GetGraphRequest(graphId uint32) (api.GraphRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (r RequestServiceMock) GetBatch(batchId uint32) (api.BatchRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (r RequestServiceMock) GetGraph(graphId uint32) (api.GraphResult, error) {
	//TODO implement me
	panic("implement me")
}

func (r RequestServiceMock) GetBatchResult(batchId uint32) ([]api.GraphResult, error) {
	//TODO implement me
	panic("implement me")
}

func (r RequestServiceMock) DeleteGraph(graphId uint32) error {
	//TODO implement me
	panic("implement me")
}

func (r RequestServiceMock) DeleteBatch(batchId uint32) error {
	//TODO implement me
	panic("implement me")
}

func (r RequestServiceMock) Start() error {
	//TODO implement me
	panic("implement me")
}

func (r RequestServiceMock) Stop() error {
	//TODO implement me
	panic("implement me")
}

func (r RequestServiceMock) CheckMaintenance() bool {
	//TODO implement me
	panic("implement me")
}

func TestBasicInit(t *testing.T) {
	test := gin.New()
	mockRequests := RequestServiceMock{}
	SetupREST(test, &mockRequests)

}
