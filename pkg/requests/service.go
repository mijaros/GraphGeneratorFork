package requests

import (
	"errors"
	"github.com/soch-fit/GraphGenerator/pkg/api"
)

var (
	ErrGraphNotFound          = errors.New("graph not found")
	ErrGraphNotGenerated      = errors.New("graph not generated yet")
	ErrBatchNotGenerated      = errors.New("batch not finished")
	ErrUnfinishedGraphBatch   = errors.New("unfinished graph in batch")
	ErrBatchNotFound          = errors.New("batch not found")
	ErrGraphDeleted           = errors.New("graph was deleted")
	ErrBatchDeleted           = errors.New("batch was deleted")
	ErrServiceMaintenance     = errors.New("service is under maintenance")
	ErrFunctionNotImplemented = errors.New("requested action is not implemented")
)

// RequestService is responsible for data storage and retrieval
// it is also responsible for requests processing and results retrieval
type RequestService interface {

	// StoreNewRequest stores new GraphRequest into the persistence and sends it
	// for processing, generated graph should be afterwards retrieved by GetGraph method.
	// It sets the Timeout and assigns new Id to the request object.
	// Returns updated request or error on failure in persistence layer.
	StoreNewRequest(request api.GraphRequest) (api.GraphRequest, error)

	// StoreNewBatch stores new BatchRequest, creates GraphRequests, sets ids
	StoreNewBatch(request api.BatchRequest) (api.BatchRequest, error)

	// StoreGraph method stores result of graph generation and updates
	StoreGraph(graph *api.GraphResult) error

	ListRequests(sessionId string) ([]uint32, error)

	ListBatches(sessionId string) ([]uint32, error)

	GetGraphRequest(graphId uint32) (api.GraphRequest, error)

	GetBatch(batchId uint32) (api.BatchRequest, error)

	GetGraph(graphId uint32) (api.GraphResult, error)

	GetBatchResult(batchId uint32) ([]api.GraphResult, error)

	DeleteGraph(graphId uint32) error

	DeleteBatch(batchId uint32) error

	Start() error

	Stop() error

	CheckMaintenance() bool
}
