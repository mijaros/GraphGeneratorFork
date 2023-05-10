package memory

import (
	"github.com/soch-fit/GraphGenerator/pkg/api"
	"github.com/soch-fit/GraphGenerator/pkg/configuration"
	"github.com/soch-fit/GraphGenerator/pkg/generator/service"
	"github.com/soch-fit/GraphGenerator/pkg/requests"
	"math/rand"
	"sync"
	"time"
)

type TypedMap[K, V any] struct {
	storage sync.Map
}

func (v *TypedMap[K, V]) Load(key K) (res V, ok bool) {
	var untypedRes any
	untypedRes, ok = v.storage.Load(key)
	if !ok {
		return
	}
	res = untypedRes.(V)
	return
}

func (v *TypedMap[K, V]) Delete(key K) {
	v.storage.Delete(key)
}

func (v *TypedMap[K, V]) Range(it func(K, V) bool) {
	v.storage.Range(func(key, value any) bool {
		mKey := key.(K)
		mVal := value.(V)
		return it(mKey, mVal)
	})
}

func (v *TypedMap[K, V]) Store(key K, val V) {
	v.storage.Store(key, val)
}

func (v *TypedMap[K, V]) LoadOrStore(key K, val V) (res V, ok bool) {
	var untypedRes any
	untypedRes, ok = v.storage.LoadOrStore(key, val)
	res = untypedRes.(V)
	return
}

type InMemoryService struct {
	rwLock    sync.RWMutex
	requests  TypedMap[uint32, api.GraphRequest]
	batches   TypedMap[uint32, api.BatchRequest]
	graphs    TypedMap[uint32, api.GraphResult]
	owners    TypedMap[string, ownerReference]
	service   *service.GenService
	stop      chan bool
	waitGroup sync.WaitGroup
}

func (i *InMemoryService) CheckMaintenance() bool {
	return false
}

func (i *InMemoryService) DeleteGraph(graphId uint32) error {
	return requests.ErrFunctionNotImplemented
}

func (i *InMemoryService) DeleteBatch(batchId uint32) error {
	return requests.ErrFunctionNotImplemented
}

type ownerReference struct {
	graphs  map[uint32]bool
	batches map[uint32]bool
}

func (i *InMemoryService) cleanUpDb() {
	i.rwLock.Lock()
	defer i.rwLock.Unlock()
	cleaningTime := time.Now()
	toDelete := make([]uint32, 100)
	batches := make([]uint32, 100)
	ownRs := make(map[string]ownerReference)
	i.batches.Range(func(key uint32, batch api.BatchRequest) bool {
		if batch.Timeout.Before(cleaningTime) {
			batches = append(batches, key)
			toDelete = append(toDelete, batch.GraphsIDs...)
			if batch.Owner == nil {
				return true
			}
			_, ok := ownRs[*batch.Owner]
			if !ok {
				ownRs[*batch.Owner] = ownerReference{
					graphs:  make(map[uint32]bool),
					batches: make(map[uint32]bool),
				}
			}
			ownRs[*batch.Owner].batches[batch.ID] = true
			for k := range batch.GraphsIDs {
				ownRs[*batch.Owner].batches[batch.GraphsIDs[k]] = true
			}
		}
		return true
	})
	for _, v := range batches {
		i.batches.Delete(v)
	}

	for _, v := range toDelete {
		i.graphs.Delete(v)
		i.requests.Delete(v)
	}
	toDelete = toDelete[:0]

	i.requests.Range(func(key uint32, graph api.GraphRequest) bool {
		if graph.Timeout.Before(cleaningTime) {
			toDelete = append(toDelete, key)
			_, ok := ownRs[*graph.Owner]
			if !ok {
				ownRs[*graph.Owner] = ownerReference{
					graphs:  make(map[uint32]bool),
					batches: make(map[uint32]bool),
				}
			}
			ownRs[*graph.Owner].graphs[graph.ID] = true
		}
		return true
	})
	for _, v := range toDelete {
		i.requests.Delete(v)
		i.graphs.Delete(v)
	}

	for k, v := range ownRs {
		owner, ok := i.owners.Load(k)
		if !ok {
			continue
		}
		for o := range v.batches {
			delete(owner.batches, o)
		}
		for o := range v.graphs {
			delete(owner.graphs, o)
		}
		if len(owner.batches) == 0 && len(owner.graphs) == 0 {
			i.owners.Delete(k)
		} else {
			i.owners.Store(k, owner)
		}
	}
}

func (i *InMemoryService) loop() {
	for {
		select {
		case <-i.stop:
			i.waitGroup.Done()
			return
		case <-time.After(configuration.Default().MaintenanceInterval):
			i.cleanUpDb()
		}
	}
}

func New(service *service.GenService) *InMemoryService {
	result := &InMemoryService{
		rwLock:  sync.RWMutex{},
		service: service,
		stop:    make(chan bool, 1),
	}
	return result
}

func getRandomId[V any](m *TypedMap[uint32, V]) (result uint32) {
	var empty V
	for {
		result = rand.Uint32()

		_, ok := m.LoadOrStore(result, empty)
		if ok {
			continue
		}
		return
	}
}

func (i *InMemoryService) storeRequestUnsafe(request api.GraphRequest, batchId *uint32) (api.GraphRequest, error) {
	requestId := getRandomId(&i.requests)
	request.ID = requestId
	if batchId != nil {
		request.BatchId = batchId
	}
	request.Timeout = time.Now().Add(configuration.Default().RequestTTL)
	request.Status = api.NotFinished
	if request.Seed == nil {
		seed := rand.Int63()
		request.Seed = &seed
	}
	i.requests.Store(request.ID, request)

	return request, nil
}

func (i *InMemoryService) StoreNewRequest(request api.GraphRequest) (api.GraphRequest, error) {
	var err error
	i.rwLock.RLock()
	defer i.rwLock.RUnlock()
	request, err = i.storeRequestUnsafe(request, nil)

	if err != nil {
		return request, err
	}
	err = i.service.PushRequest(request)

	i.addToOwner(request.Owner, request.ID, true)

	return request, err
}

func (i *InMemoryService) addToOwner(owner *string, id uint32, graph bool) {
	if owner != nil {
		owners, ok := i.owners.Load(*owner)
		if !ok {
			owners = ownerReference{
				graphs:  make(map[uint32]bool),
				batches: make(map[uint32]bool),
			}
		}
		if graph {
			owners.graphs[id] = true
		} else {
			owners.batches[id] = true
		}
		i.owners.Store(*owner, owners)
	}
}

func (i *InMemoryService) StoreNewBatch(request api.BatchRequest) (api.BatchRequest, error) {
	i.rwLock.RLock()
	defer i.rwLock.RUnlock()
	ids := make([]uint32, request.Number)
	request.Timeout = time.Now().Add(configuration.Default().RequestTTL)
	request.Status = api.NotFinished
	batchId := request.ID
	for j := range ids {
		r, err := i.storeRequestUnsafe(request.BaseGraph, &batchId)
		if err != nil {
			panic("Kalm")
		}
		ids[j] = r.ID
		i.service.PushRequest(r)
	}

	request.ID = getRandomId(&i.batches)
	rCopy := request
	rCopy.GraphsIDs = append([]uint32{}, ids...)
	i.batches.Store(rCopy.ID, rCopy)
	i.addToOwner(request.Owner, request.ID, false)

	request.GraphsIDs = ids
	return request, nil

}

func (i *InMemoryService) ListRequests(sessionId string) ([]uint32, error) {
	i.rwLock.Lock()
	defer i.rwLock.Unlock()

	owners, ok := i.owners.Load(sessionId)
	if !ok {
		return []uint32{}, nil
	}
	result := make([]uint32, 0, len(owners.graphs))
	for k, v := range owners.graphs {
		if v {
			result = append(result, k)
		}
	}

	return result, nil
}

func (i *InMemoryService) ListBatches(sessionId string) ([]uint32, error) {
	i.rwLock.RLock()
	defer i.rwLock.RUnlock()

	owners, ok := i.owners.Load(sessionId)
	if !ok {
		return []uint32{}, nil
	}

	result := make([]uint32, 0, len(owners.batches))

	for k, v := range owners.batches {
		if !v {
			continue
		}

		result = append(result, k)
	}

	return result, nil
}

func (i *InMemoryService) GetGraphRequest(graphId uint32) (api.GraphRequest, error) {
	i.rwLock.RLock()
	defer i.rwLock.RUnlock()
	res, ok := i.requests.Load(graphId)
	if !ok {
		return api.GraphRequest{}, requests.ErrGraphNotFound
	}
	return res, nil
}

func (i *InMemoryService) GetBatch(batchId uint32) (api.BatchRequest, error) {
	i.rwLock.RLock()
	defer i.rwLock.RUnlock()

	res, ok := i.batches.Load(batchId)
	if !ok {
		return api.BatchRequest{}, requests.ErrBatchNotFound
	}

	status := res.Status
	if status != api.Finished {
		actualState := api.Finished
		for _, k := range res.GraphsIDs {
			if _, ok := i.graphs.Load(k); !ok {
				actualState = api.NotFinished
				break
			}
		}
		if actualState != status {
			res.Status = actualState
			i.batches.Store(res.ID, res)
		}
	}
	return res, nil
}

func (i *InMemoryService) GetGraph(graphId uint32) (api.GraphResult, error) {
	i.rwLock.RLock()
	defer i.rwLock.RUnlock()

	result, ok := i.graphs.Load(graphId)
	if !ok {
		return api.GraphResult{}, requests.ErrGraphNotFound
	}
	return result, nil
}

func (i *InMemoryService) StoreGraph(graph *api.GraphResult) error {
	i.rwLock.RLock()
	defer i.rwLock.RUnlock()
	val, ok := i.requests.Load(graph.ID)
	if !ok {
		return requests.ErrGraphDeleted
	}
	val.Status = api.Finished
	i.graphs.Store(graph.ID, *graph)
	i.requests.Store(graph.ID, val)

	return nil
}

func (i *InMemoryService) Start() error {
	i.waitGroup.Add(1)
	go i.loop()
	return nil
}

func (i *InMemoryService) Stop() error {
	i.stop <- true
	i.waitGroup.Wait()
	return nil
}

func (i *InMemoryService) GetBatchResult(batchId uint32) ([]api.GraphResult, error) {
	i.rwLock.RLock()
	defer i.rwLock.RUnlock()

	batch, ok := i.batches.Load(batchId)
	if !ok {
		return []api.GraphResult{}, requests.ErrBatchNotFound
	}
	result := make([]api.GraphResult, len(batch.GraphsIDs))

	for k, v := range batch.GraphsIDs {
		graph, ok := i.requests.Load(v)
		if !ok {
			return []api.GraphResult{}, requests.ErrGraphNotFound
		}
		if graph.Status != api.Finished {
			return []api.GraphResult{}, requests.ErrUnfinishedGraphBatch
		}
		res, ok := i.graphs.Load(v)
		if !ok {
			panic("This shouldn't have happened")
		}
		result[k] = res
	}
	return result, nil
}
