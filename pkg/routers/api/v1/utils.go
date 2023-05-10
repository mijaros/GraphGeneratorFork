package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/soch-fit/GraphGenerator/pkg/requests"
)

func getRequestsService(c *gin.Context) requests.RequestService {
	service, ok := c.Get("requestsService")
	if !ok {
		panic("No service ready")
	}
	reqService := service.(requests.RequestService)
	return reqService
}
