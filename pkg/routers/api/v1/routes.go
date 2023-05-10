package v1

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/soch-fit/GraphGenerator/pkg/api"
	"github.com/soch-fit/GraphGenerator/pkg/configuration"
	"github.com/soch-fit/GraphGenerator/pkg/requests"
	"github.com/soch-fit/GraphGenerator/pkg/routers/middleware"
	"io"
	"net/http"
	"strconv"
)

var (
	ErrInvalidRequest    = errors.New("invalid request body")
	ErrInvalidAttributes = errors.New("attributes of request are invalid")
)

func handleGraphList(r *gin.Context) {
	c, ok := r.Get("identifier")
	if !ok {
		r.JSON(200, gin.H{"graphs": []int32{}})
		return
	}

	service := getRequestsService(r)
	list, err := service.ListRequests(c.(string))
	if err != nil {
		r.JSON(400, gin.H{"error": err.Error()})
		return
	}

	r.JSON(200, gin.H{"graphs": list})
}

func handleGraphGet(r *gin.Context) {
	var graphId uint32

	grId, err := strconv.Atoi(r.Param("graphId"))
	if err != nil || grId < 0 {

		return
	}
	graphId = uint32(grId)

	service, ok := r.Get("requestsService")
	if !ok {
		panic("No service ready")
	}
	reqService := service.(requests.RequestService)
	res, err := reqService.GetGraphRequest(graphId)
	if err != nil {
		log.Warning("error: ", err)
		r.JSON(http.StatusNotFound, api.NewErr(err, nil))
		return
	}
	s := res.Type.String()
	log.Debug(s)
	r.JSON(200, res)

}

func handleGraphCreate(r *gin.Context) {
	var req api.GraphRequest
	if err := r.ShouldBind(&req); err != nil {
		r.JSON(http.StatusBadRequest, api.NewErr(ErrInvalidRequest, err))
		return
	}
	if !req.Valid() {
		r.JSON(http.StatusBadRequest, api.NewErr(ErrInvalidAttributes, nil))
		return
	}

	own, ok := r.Get("identifier")

	if ok {
		owner := own.(string)
		req.Owner = &owner
	}

	service, ok := r.Get("requestsService")
	if !ok {
		panic("No service ready")
	}
	reqService := service.(requests.RequestService)
	newReq, err := reqService.StoreNewRequest(req)
	if err != nil {
		panic("coldn't store")
	}

	r.JSON(201, newReq)
}

func handleGraphDownload(r *gin.Context) {
	var graphId uint32

	grId, err := strconv.Atoi(r.Param("graphId"))
	if err != nil || grId < 0 {
		r.JSON(http.StatusBadRequest, api.NewErr(ErrInvalidRequest, err))
		return
	}
	graphId = uint32(grId)

	requestService := getRequestsService(r)

	translator := getGraphFormat(r)

	v, err := requestService.GetGraph(graphId)
	if err != nil {
		r.Error(err)
		return
	}

	translator.Convert(v.Generated)
	data := translator.Bytes()
	reader := bytes.NewReader(data)
	attachment := fmt.Sprintf(`attachment; filename="rngr-%d.%s"`, graphId, translator.Extension())
	r.DataFromReader(200, reader.Size(), translator.ContentType(), reader, map[string]string{"Content-Disposition": attachment})

}

func handleBatchList(r *gin.Context) {
	c, ok := r.Get("identifier")
	if !ok {
		r.JSON(200, gin.H{"batches": []api.BatchRequest{}})
		return
	}

	srv := getRequestsService(r)

	batches, err := srv.ListBatches(c.(string))
	if err != nil {
		r.Error(err)
		return
	}

	r.JSON(200, gin.H{"batches": batches})
}

func handleBatchCreate(r *gin.Context) {
	var request api.BatchRequest

	if err := r.ShouldBind(&request); err != nil {
		//r.Error(ErrInvalidRequest)
		r.JSON(http.StatusBadRequest, api.NewErr(ErrInvalidRequest, err))
		return
	}

	if request.Number < 0 || request.Number > configuration.Default().MaxBatchSize {
		r.JSON(http.StatusBadRequest, gin.H{"error": "batch size outside of configuration"})
		return
	}
	if !request.BaseGraph.Valid() {
		r.JSON(http.StatusBadRequest, api.NewErr(ErrInvalidRequest, ErrInvalidAttributes))
		return
	}

	own, ok := r.Get("identifier")

	if ok {
		owner := own.(string)
		request.Owner = &owner
	}

	service := getRequestsService(r)
	res, err := service.StoreNewBatch(request)
	if err != nil {
		r.JSON(http.StatusInternalServerError, gin.H{"error": "could not store batch"})
		return
	}
	r.JSON(201, res)
}

func handleBatchGet(r *gin.Context) {
	var batchId uint32

	grId, err := strconv.Atoi(r.Param("batchId"))
	if err != nil || grId < 0 {
		r.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	batchId = uint32(grId)

	service := getRequestsService(r)
	batch, err := service.GetBatch(batchId)
	if err != nil {
		r.JSON(404, gin.H{"error": "couldn't find batch", "reason": err.Error()})
		return
	}

	r.JSON(200, batch)

}

func handleGraphDelete(r *gin.Context) {
	var graphId uint32

	graphRaw, err := strconv.Atoi(r.Param("graphId"))
	if err != nil {
		r.JSON(400, gin.H{"error": "invalid graph id", "cause": err.Error()})
		return
	}

	graphId = uint32(graphRaw)
	service := getRequestsService(r)
	err = service.DeleteGraph(graphId)
	if err != nil {
		r.Error(err)
		return
	}
	r.Status(http.StatusNoContent)
}

func handleBatchDelete(r *gin.Context) {

}

func getGraphFormat(r *gin.Context) api.GraphTranslator {
	rep := r.DefaultQuery("graphKind", "matrix")
	var translator api.GraphTranslator

	switch rep {
	case "matrix":
		translator = &api.MatrixGraph{}
	case "JSON":
		translator = &api.BasicJSONGraph{}
	case "dot":
		translator = &api.DotGraph{}
	default:
		translator = &api.MatrixGraph{}
	}
	return translator
}

func handleBatchDownload(r *gin.Context) {
	var batchId uint32

	grId, err := strconv.Atoi(r.Param("batchId"))
	if err != nil || grId < 0 {
		r.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	batchId = uint32(grId)

	service := getRequestsService(r)
	graphs, err := service.GetBatchResult(batchId)
	//batch, err := service.GetBatch(batchId)
	if err != nil {
		r.JSON(404, gin.H{"error": "batch couldn't be obtained", "reason": err.Error()})
		return
	}

	translator := getGraphFormat(r)
	attachment := fmt.Sprintf(`attachment; filename=rnrg-%d-%s.zip`, batchId, translator.Kind())
	r.Writer.Header().Add("Content-Disposition", attachment)

	r.Stream(func(w io.Writer) bool {
		zp := zip.NewWriter(w)
		for k := range graphs {
			fileName := fmt.Sprintf("rngr-%d.%s", graphs[k].ID, translator.Extension())
			translator.Convert(graphs[k].Generated)
			f, err := zp.Create(fileName)
			if err != nil {
				panic("something terrible happened")
			}
			_, err = translator.Serialize(f)
			if err != nil {
				panic("something else terrible had happened")
			}
		}
		err = zp.Close()
		if err != nil {
			panic("could not finalize zip")
		}
		return false
	})
}

func handleLimitsRequest(r *gin.Context) {
	response := api.LimitsResponse{
		MaxNodes:     configuration.Default().MaxNodes,
		MaxBatchSize: configuration.Default().MaxBatchSize,
	}
	r.JSON(http.StatusOK, response)
}

func SetupREST(engine *gin.Engine, svc requests.RequestService) *gin.RouterGroup {
	if engine == nil {
		engine = gin.Default()
	}
	r := engine.Group("/api/v1")
	gin.EnableJsonDecoderDisallowUnknownFields()
	r.Use(middleware.SetUpCookieMiddleware())
	r.Use(middleware.SetUpRequestService(svc))
	r.Use(middleware.Error())
	r.GET("/limits", handleLimitsRequest)
	r.GET("graph", handleGraphList)
	r.GET("graph/:graphId", handleGraphGet)
	r.DELETE("graph/:graphId", handleGraphDelete)
	r.POST("graph", handleGraphCreate)
	r.GET("graph/:graphId/download", handleGraphDownload)

	r.GET("batch", handleBatchList)
	r.POST("batch", handleBatchCreate)
	r.GET("batch/:batchId", handleBatchGet)
	r.DELETE("batch/:batchId", handleBatchDelete)
	r.GET("batch/:batchId/download", handleBatchDownload)

	return r

}
