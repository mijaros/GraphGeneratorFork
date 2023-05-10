package persistent

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v4"
	log "github.com/sirupsen/logrus"
	"github.com/soch-fit/GraphGenerator/pkg/api"
	"github.com/soch-fit/GraphGenerator/pkg/configuration"
	"github.com/soch-fit/GraphGenerator/pkg/generator/service"
	"github.com/soch-fit/GraphGenerator/pkg/requests"
	"github.com/soch-fit/GraphGenerator/pkg/utils"
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type KeyProvider interface {
	GetKey() []byte
	GetPrefix() []byte
}

type OwnerReference struct {
	Id   uint32
	Span time.Time
}

type OwnerReferenceList []OwnerReference

func (o OwnerReferenceList) Len() int {
	return len(o)
}

func (o OwnerReferenceList) Less(i, j int) bool {
	return o[i].Span.Before(o[j].Span)
}

func (o OwnerReferenceList) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

type DbGraphRequest struct {
	ID uint32
}

func (g *DbGraphRequest) GetKey() []byte {
	return []byte(fmt.Sprintf("request-%d", g.ID))
}

func (g *DbGraphRequest) GetPrefix() []byte {
	return []byte("request-")
}

type DbGraphResult struct {
	ID uint32
}

func (d DbGraphResult) GetKey() []byte {
	return []byte(fmt.Sprintf("result-%d", d.ID))
}

func (d DbGraphResult) GetPrefix() []byte {
	return []byte("result-")
}

type DbOwnerGraphId struct {
	Session string
	GraphId uint32
}

func (d DbOwnerGraphId) GetKey() []byte {
	return []byte(fmt.Sprintf("graphowner-%s-%d", d.Session, d.GraphId))
}

func (d DbOwnerGraphId) GetPrefix() []byte {
	return []byte(fmt.Sprintf("graphowner-%s-", d.Session))
}

type DbOwnerBatch struct {
	Session string
	BatchId uint32
}

func (d DbOwnerBatch) GetKey() []byte {
	return []byte(fmt.Sprintf("batchowner-%s-%d", d.Session, d.BatchId))
}

func (d DbOwnerBatch) GetPrefix() []byte {
	return []byte(fmt.Sprintf("batchowner-%s-", d.Session))
}

type DbBatch struct {
	ID uint32
}

type DbBatchGraph struct {
	BatchId uint32
	GraphId uint32
}

func (d DbBatchGraph) GetKey() []byte {
	return []byte(fmt.Sprintf("batch-graph-%d-%d", d.BatchId, d.GraphId))
}

func (d DbBatchGraph) GetPrefix() []byte {
	return []byte(fmt.Sprintf("batch-graph-%d-", d.BatchId))
}

func (d DbBatch) GetKey() []byte {
	return []byte(fmt.Sprintf("batch-%d", d.ID))
}

func (d DbBatch) GetPrefix() []byte {
	return []byte("batch-")
}

type DbHandleFunc func(txn *badger.Txn) error

func marshall[K any](val K) []byte {
	var buff bytes.Buffer
	translator := gob.NewEncoder(&buff)
	translator.Encode(val)
	return buff.Bytes()
}

func unMarshall[K any](val []byte) K {
	var result K
	reader := bytes.NewReader(val)
	translator := gob.NewDecoder(reader)
	translator.Decode(&result)

	return result
}

func getUserGraphs(owner string, res *OwnerReferenceList) DbHandleFunc {
	or := DbOwnerGraphId{Session: owner}
	pref := or.GetPrefix()
	return func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		opt.Prefix = pref
		it := txn.NewIterator(opt)
		defer it.Close()
		for it.Seek(pref); it.ValidForPrefix(pref); it.Next() {
			i := it.Item()
			err := i.Value(func(val []byte) error {
				*res = append(*res, unMarshall[OwnerReference](val))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func getUserBatches(owner string, res *OwnerReferenceList) DbHandleFunc {
	or := DbOwnerBatch{
		Session: owner,
		BatchId: 0,
	}.GetPrefix()
	return func(txn *badger.Txn) error {
		opt := badger.DefaultIteratorOptions
		it := txn.NewIterator(opt)
		defer it.Close()

		for it.Seek(or); it.ValidForPrefix(or); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				*res = append(*res, unMarshall[OwnerReference](val))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func getGraphRequest(graphId uint32, result *api.GraphRequest) DbHandleFunc {
	return func(txn *badger.Txn) error {
		key := DbGraphRequest{graphId}
		res, err := txn.Get(key.GetKey())
		if err != nil {
			return err
		}
		err = res.Value(func(val []byte) error {
			*result = unMarshall[api.GraphRequest](val)
			return nil
		})
		return err
	}
}

func storeGraphRequest(request *api.GraphRequest, newId bool) DbHandleFunc {
	return func(txn *badger.Txn) error {
		var err error
		id := DbGraphRequest{request.ID}
		for newId {
			id = DbGraphRequest{rand.Uint32()}
			_, err = txn.Get(id.GetKey())
			if err != nil && err == badger.ErrKeyNotFound {
				break
			}
		}
		request.ID = id.ID
		if newId && err != badger.ErrKeyNotFound {
			return err
		}
		log.Debugln("Store request", request.ID, time.Until(request.Timeout))

		entry := badger.NewEntry(id.GetKey(), marshall(request)).WithTTL(time.Until(request.Timeout))
		return txn.SetEntry(entry)
	}
}

func updateGraphOwner(graphId uint32, owner string, ttl time.Time) DbHandleFunc {
	ref := DbOwnerGraphId{
		Session: owner,
		GraphId: graphId,
	}.GetKey()
	ownRef := OwnerReference{
		Id:   graphId,
		Span: ttl,
	}
	log.Debugln("Store owner", owner, graphId, ttl)
	return func(txn *badger.Txn) error {
		entry := badger.NewEntry(ref, marshall(ownRef)).WithTTL(time.Until(ttl))
		return txn.SetEntry(entry)
	}
}

func updateBatchOwner(batchId uint32, owner string, ttl time.Time) DbHandleFunc {
	ref := DbOwnerBatch{
		Session: owner,
		BatchId: batchId,
	}.GetKey()
	ownRef := OwnerReference{
		Id:   batchId,
		Span: ttl}
	log.Debugln("Store owner for batch", owner, batchId, ttl)
	return func(txn *badger.Txn) error {
		entry := badger.NewEntry(ref, marshall(ownRef)).WithTTL(time.Until(ttl))
		return txn.SetEntry(entry)
	}
}

type PersistentService struct {
	dbHandle        *badger.DB
	generator       service.Service
	closer          chan bool
	wg              sync.WaitGroup
	maintenance     atomic.Bool
	maintenanceStop sync.WaitGroup
	lastMaintenance time.Time
	nextMaintenance time.Time
}

func (p *PersistentService) CheckMaintenance() bool {
	return p.maintenance.Load()
}

func handleDeleteGraph(graphId uint32) DbHandleFunc {
	return func(txn *badger.Txn) error {
		var graph api.GraphRequest
		e := getGraphRequest(graphId, &graph)(txn)
		if e != nil {
			return e
		}
		if graph.Status != api.Finished {
			return requests.ErrGraphNotGenerated
		}
		grId, resId := DbGraphRequest{ID: graphId}, DbGraphResult{ID: graphId}
		e = txn.Delete(grId.GetKey())
		if e != nil {
			return e
		}
		return txn.Delete(resId.GetKey())
	}
}

func (p *PersistentService) DeleteGraph(graphId uint32) error {
	if p.CheckMaintenance() {
		return requests.ErrServiceMaintenance
	}
	err := p.dbHandle.Update(handleDeleteGraph(graphId))
	if err == badger.ErrKeyNotFound {
		return requests.ErrGraphNotFound
	}
	return err
}

func (p *PersistentService) DeleteBatch(batchId uint32) error {
	if p.CheckMaintenance() {
		return requests.ErrServiceMaintenance
	}
	err := p.dbHandle.Update(func(txn *badger.Txn) error {
		var batch api.BatchRequest
		e := getBatch(batchId, &batch)(txn)
		if e != nil {
			return e
		}
		if batch.Status != api.Finished {
			return requests.ErrBatchNotGenerated
		}
		for _, v := range batch.GraphsIDs {
			e = handleDeleteGraph(v)(txn)
			if e == badger.ErrKeyNotFound {
				return requests.ErrGraphNotFound
			}
			if e != nil {
				return e
			}
		}
		bId := DbBatch{ID: batchId}
		return txn.Delete(bId.GetKey())
	})
	if err == badger.ErrKeyNotFound {
		return requests.ErrBatchNotFound
	}
	return err
}

func getNextMaintenance(lastMaintenance time.Time, interval time.Duration, baseHour int) time.Time {
	nextMaint := lastMaintenance.Add(interval)
	if interval.Hours() == 24 {
		nextMaint.Truncate(time.Hour)
	}
	return nextMaint
}

func getFirstMaintenance(interval time.Duration, baseHour int) time.Time {

	start := time.Now()
	y, m, d := start.Year(), start.Month(), start.Day()
	start = time.Date(y, m, d, baseHour, 0, 0, 0, time.Local)
	for start.Before(time.Now()) {
		start = start.Add(interval)
	}
	log.Debugf("First maintenance planned for %s", start.String())
	return start
}

func New(svc service.Service) (*PersistentService, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}

	return &PersistentService{dbHandle: db,
		generator:       svc,
		closer:          make(chan bool),
		wg:              sync.WaitGroup{},
		maintenance:     atomic.Bool{},
		maintenanceStop: sync.WaitGroup{},
		lastMaintenance: time.Now(),
		nextMaintenance: getFirstMaintenance(configuration.Default().MaintenanceInterval, configuration.Default().MaintenanceHour.Value()),
	}, nil
}

func connectDB() (*badger.DB, error) {
	opts := badger.DefaultOptions(configuration.Default().DbRoot).
		WithLogger(log.StandardLogger()).
		WithNumVersionsToKeep(1).
		WithCompactL0OnClose(true).WithVLogPercentile(0.0).WithNumCompactors(10)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	err = nil
	for err == nil {
		err = db.RunValueLogGC(0.1)
	}
	return db, nil
}

func (p *PersistentService) doMaintenance() {
	var err error
	defer utils.StopProcessOnUnhandledPanic()
	log.Warning("Starting maintenance now!")
	p.generator.Pause()

	p.maintenance.Store(true)

	p.dbHandle.Close()
	p.dbHandle, err = connectDB()
	if err != nil {
		log.Panic("Restart of server failed ", err)
	}
	p.generator.Resume()
	time.Sleep(5 * time.Second)
	p.maintenance.Store(false)
	log.Info("Maintenance ended successfully")

}

func (p *PersistentService) cleanLoop() {
	defer utils.StopProcessOnUnhandledPanic()
	timer := time.NewTicker(1 * time.Minute)
	for {
		select {
		case currTime := <-timer.C:
			if currTime.After(p.nextMaintenance) {
				p.doMaintenance()
				p.lastMaintenance = time.Now()
				p.nextMaintenance = getNextMaintenance(p.lastMaintenance, configuration.Default().MaintenanceInterval,
					configuration.Default().MaintenanceHour.Value())
				log.Debugf("Next maintenance planned for %s", p.nextMaintenance)
				continue
			}
		case <-p.closer:
			timer.Stop()
			p.wg.Done()
			return
		}
	}
}

func (p *PersistentService) StoreNewRequest(request api.GraphRequest) (api.GraphRequest, error) {
	if p.CheckMaintenance() {
		return api.GraphRequest{}, requests.ErrServiceMaintenance
	}
	var err error
	if request.Seed == nil {
		seed := rand.Int63()
		request.Seed = &seed
	}
	request.Timeout = time.Now().Add(configuration.Default().RequestTTL)
	request.Status = api.NotFinished
	log.Debugln("request duration ", request.Timeout)
	err = p.dbHandle.Update(func(txn *badger.Txn) error {
		e := storeGraphRequest(&request, true)(txn)
		if e != nil {
			return e
		}
		if request.Owner != nil {
			e = updateGraphOwner(request.ID, *request.Owner, request.Timeout)(txn)
		}
		return e
	})
	err = p.generator.PushRequest(request)
	return request, err
}

func (p *PersistentService) StoreNewBatch(request api.BatchRequest) (api.BatchRequest, error) {
	if p.CheckMaintenance() {
		return api.BatchRequest{}, requests.ErrServiceMaintenance
	}
	if p.generator.FreeBand() < request.Number {
		return request, errors.New("insufficient bandwidth")
	}
	baseGraph := request.BaseGraph
	request.Status = api.NotFinished
	request.Timeout = time.Now().Add(configuration.Default().RequestTTL)
	request.Finished = 0
	request.GraphsIDs = make([]uint32, request.Number)
	baseGraph.Timeout = request.Timeout

	txn := p.dbHandle.NewTransaction(true)

	id := DbBatch{ID: rand.Uint32()}

	for {
		_, e := txn.Get(id.GetKey())
		if e == badger.ErrKeyNotFound {
			break
		}
		id.ID = rand.Uint32()
	}
	request.ID = id.ID

	graphRequests := make([]api.GraphRequest, request.Number)
	for i := 0; i < request.Number; i++ {
		ownId := id.ID
		newGraph := baseGraph
		newGraph.BatchId = &ownId
		seed := rand.Int63()
		newGraph.Seed = &seed

		err := storeGraphRequest(&newGraph, true)(txn)
		if err != nil {
			txn.Discard()
			return api.BatchRequest{}, err
		}
		request.GraphsIDs[i] = newGraph.ID
		graphRequests[i] = newGraph
	}
	entry := badger.NewEntry(id.GetKey(), marshall(request)).WithTTL(time.Until(request.Timeout))
	err := txn.SetEntry(entry)
	if err != nil {
		txn.Discard()
		return request, err
	}
	if request.Owner != nil {
		err = updateBatchOwner(request.ID, *request.Owner, request.Timeout)(txn)
		if err != nil {
			txn.Discard()
			return request, err
		}
	}
	err = txn.Commit()
	for _, v := range graphRequests {
		p.generator.PushRequest(v)
	}
	return request, err

}

func updateBatch(batchId, graphId uint32) DbHandleFunc {
	return func(txn *badger.Txn) error {
		recId := DbBatchGraph{BatchId: batchId, GraphId: graphId}

		ref := DbBatch{ID: batchId}.GetKey()
		batch, err := txn.Get(ref)
		if err != nil {
			return err
		}

		var batchRequest api.BatchRequest
		err = batch.Value(func(val []byte) error {
			batchRequest = unMarshall[api.BatchRequest](val)
			return nil
		})
		if err != nil {
			return err
		}

		entry := badger.NewEntry(recId.GetKey(), marshall(graphId)).WithTTL(time.Until(batchRequest.Timeout))
		return txn.SetEntry(entry)
	}
}

func (p *PersistentService) retrieveGraphs() {
	defer utils.StopProcessOnUnhandledPanic()
	for {
		select {
		case graph := <-p.generator.GetRetriever():
			if graph == nil {
				continue
			}
			if err := p.StoreGraph(graph); err != nil {
				log.Errorf("Storage of graph %d failed %s", graph.ID, err.Error())
			}
		case <-p.closer:
			p.wg.Done()
			return
		}
	}
}

func (p *PersistentService) StoreGraph(graph *api.GraphResult) error {
	if p.CheckMaintenance() {
		p.maintenanceStop.Wait()
	}
	return p.dbHandle.Update(func(txn *badger.Txn) error {
		id := DbGraphResult{graph.ID}
		var graphRequest api.GraphRequest
		err := getGraphRequest(graph.ID, &graphRequest)(txn)
		if err != nil {
			return requests.ErrGraphDeleted
		}
		gr := marshall(*graph)
		dur := time.Until(graphRequest.Timeout)
		entry := badger.NewEntry(id.GetKey(), gr).WithTTL(dur)
		err = txn.SetEntry(entry)
		if err != nil {
			return err
		}
		graphRequest.Status = api.Finished
		err = storeGraphRequest(&graphRequest, false)(txn)
		if err != nil {
			return err
		}
		if graphRequest.BatchId != nil {
			return updateBatch(*graphRequest.BatchId, graph.ID)(txn)
		}
		return nil
	})
}

func (p *PersistentService) ListRequests(sessionId string) (result []uint32, e error) {
	if p.CheckMaintenance() {
		return []uint32{}, requests.ErrServiceMaintenance
	}
	tmpRes := make(OwnerReferenceList, 0)
	e = p.dbHandle.View(getUserGraphs(sessionId, &tmpRes))
	if e != nil {
		return
	}
	sort.Sort(tmpRes)
	result = make([]uint32, tmpRes.Len())
	for i, v := range tmpRes {
		result[i] = v.Id
	}
	return
}

func (p *PersistentService) ListBatches(sessionId string) (result []uint32, e error) {
	if p.CheckMaintenance() {
		return []uint32{}, requests.ErrServiceMaintenance
	}
	tmpRes := make(OwnerReferenceList, 0)
	e = p.dbHandle.View(getUserBatches(sessionId, &tmpRes))
	if e != nil {
		return
	}
	sort.Sort(tmpRes)
	result = make([]uint32, tmpRes.Len())
	for i, v := range tmpRes {
		result[i] = v.Id
	}
	return
}

func (p *PersistentService) GetGraphRequest(graphId uint32) (result api.GraphRequest, e error) {
	if p.CheckMaintenance() {
		return api.GraphRequest{}, requests.ErrServiceMaintenance
	}
	e = p.dbHandle.View(getGraphRequest(graphId, &result))
	return
}

func batchFinished(batchId uint32, txn *badger.Txn) int {
	it := badger.DefaultIteratorOptions
	it.PrefetchValues = false
	pref := DbBatchGraph{BatchId: batchId}.GetPrefix()

	iter := txn.NewIterator(it)
	defer iter.Close()
	count := 0
	for iter.Seek(pref); iter.ValidForPrefix(pref); iter.Next() {
		count++
	}
	return count
}

func getBatch(batchId uint32, result *api.BatchRequest) DbHandleFunc {
	return func(txn *badger.Txn) error {
		id := DbBatch{ID: batchId}

		v, e := txn.Get(id.GetKey())
		if e != nil {
			return e
		}
		e = v.Value(func(val []byte) error {
			*result = unMarshall[api.BatchRequest](val)
			return nil
		})

		return e
	}
}

func (p *PersistentService) GetBatch(batchId uint32) (result api.BatchRequest, err error) {
	if p.CheckMaintenance() {
		return api.BatchRequest{}, requests.ErrServiceMaintenance
	}
	err = p.dbHandle.View(getBatch(batchId, &result))
	if err == nil && result.Status == api.NotFinished {
		err = p.dbHandle.Update(func(txn *badger.Txn) error {
			count := batchFinished(result.ID, txn)
			if count == result.Number {
				batchKey := DbBatch{ID: result.ID}
				result.Status = api.Finished
				entry := badger.NewEntry(batchKey.GetKey(), marshall(result)).WithTTL(time.Until(result.Timeout))
				return txn.SetEntry(entry)
			}
			return nil
		})
	}
	return
}
func getGraphResult(graphId uint32, result *api.GraphResult) DbHandleFunc {
	return func(txn *badger.Txn) error {
		id := DbGraphResult{graphId}
		it, e := txn.Get(id.GetKey())
		if e != nil {
			return e
		}
		return it.Value(func(val []byte) error {
			*result = unMarshall[api.GraphResult](val)
			return nil
		})
	}
}
func (p *PersistentService) GetGraph(graphId uint32) (res api.GraphResult, err error) {
	if p.CheckMaintenance() {
		return api.GraphResult{}, requests.ErrServiceMaintenance
	}
	err = p.dbHandle.View(getGraphResult(graphId, &res))
	if err == badger.ErrKeyNotFound {
		err = requests.ErrGraphNotFound
	}
	return
}

func (p *PersistentService) GetBatchResult(batchId uint32) (result []api.GraphResult, err error) {
	if p.CheckMaintenance() {
		return []api.GraphResult{}, requests.ErrServiceMaintenance
	}
	err = p.dbHandle.View(func(txn *badger.Txn) error {
		var graphRequest api.BatchRequest
		e := getBatch(batchId, &graphRequest)(txn)
		if e != nil {
			return e
		}
		result = make([]api.GraphResult, graphRequest.Number)
		for k, v := range graphRequest.GraphsIDs {
			e = getGraphResult(v, &result[k])(txn)
			if e != nil {
				return requests.ErrUnfinishedGraphBatch
			}
		}
		return nil
	})
	if err == badger.ErrKeyNotFound {
		err = requests.ErrUnfinishedGraphBatch
	}
	return
}

func (p *PersistentService) Start() error {
	p.wg.Add(2)
	go p.cleanLoop()
	go p.retrieveGraphs()
	return nil
}

func (p *PersistentService) Stop() error {
	p.closer <- true
	p.closer <- true
	p.wg.Wait()
	return p.dbHandle.Close()
}
