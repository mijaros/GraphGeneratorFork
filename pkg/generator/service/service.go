package service

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/soch-fit/GraphGenerator/pkg/api"
	"github.com/soch-fit/GraphGenerator/pkg/configuration"
	"github.com/soch-fit/GraphGenerator/pkg/generator/decision"
	"github.com/soch-fit/GraphGenerator/pkg/utils"
	"sync"
	"sync/atomic"
)

var (
	ErrStoppedService   = errors.New("cannot stop stopped service")
	ErrCantStartService = errors.New("cannot start stopped service")
	ErrPushStoppedSvc   = errors.New("can't push to stopped service")
)

type Service interface {
	Start() error
	Pause() error
	Resume() error
	PushRequest(req api.GraphRequest) error
	GetRetriever() chan *api.GraphResult
	Stop() error
	FreeBand() int
}

type GenService struct {
	waitGroup      sync.WaitGroup
	requests       chan api.GraphRequest
	retirever      chan *api.GraphResult
	pause          chan bool
	resume         sync.WaitGroup
	workers        int
	runningWorkers atomic.Int32
	started        bool
	stopped        bool
}

func (service *GenService) GetRetriever() chan *api.GraphResult {
	return service.retirever
}

func New(workers int) *GenService {
	res := &GenService{
		waitGroup:      sync.WaitGroup{},
		requests:       make(chan api.GraphRequest, workers*configuration.Default().MaxBatchSize*2),
		retirever:      make(chan *api.GraphResult, workers*configuration.Default().MaxBatchSize*2),
		workers:        workers,
		runningWorkers: atomic.Int32{},
		started:        false,
		stopped:        false,
		pause:          make(chan bool),
		resume:         sync.WaitGroup{},
	}
	res.runningWorkers.Add(int32(workers))
	return res
}

func (service *GenService) runner() {
	defer func() {
		service.waitGroup.Done()
	}()
	defer utils.StopProcessOnUnhandledPanic()
	for {
		select {
		case <-service.pause:
			service.resume.Wait()
		default:
			select {
			case <-service.pause:
				service.resume.Wait()
			case request, cont := <-service.requests:
				if !cont {
					log.Traceln("Ending thread")
					return
				} else {
					log.Debugf("Generating %d", request.ID)
				}
				graph, err := decision.GenerateGraphFromRequest(request)
				if err != nil || graph == nil {
					log.Error(err)
					continue
				}
				service.retirever <- graph
			}
		}
	}
}
func (service *GenService) Start() error {
	if service.started || service.stopped {
		return ErrCantStartService
	}
	service.started = true
	for i := int32(0); i < service.runningWorkers.Load(); i++ {
		go service.runner()
	}
	service.waitGroup.Add(int(service.runningWorkers.Load()))
	return nil
}

func (service *GenService) Pause() error {
	service.resume.Add(1)
	for k := int32(0); k < service.runningWorkers.Load(); k++ {
		service.pause <- true
	}
	return nil
}

func (service *GenService) Resume() error {
	service.resume.Done()
	return nil
}

func (service *GenService) PushRequest(req api.GraphRequest) error {
	if !service.started {
		return ErrPushStoppedSvc
	}
	service.requests <- req
	return nil
}

func (service *GenService) Stop() error {
	if !service.started || service.stopped {

		return ErrStoppedService
	}
	service.started = false
	service.stopped = true
	close(service.requests)
	service.waitGroup.Wait()
	return nil
}

func (service *GenService) FreeBand() int {
	return cap(service.requests) - len(service.requests)
}
