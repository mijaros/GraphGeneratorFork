package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/soch-fit/GraphGenerator/pkg/api"
	"github.com/soch-fit/GraphGenerator/pkg/configuration"
	"github.com/soch-fit/GraphGenerator/pkg/requests"
	"net/http"
	time2 "time"
)

const (
	SessionCookieName = "RNGRSESSION"
)

var (
	ErrUnknownFailure    = errors.New("request processing failed")
	ErrSystemMaintenance = errors.New("system under maintenance")
)

func setSameSite(context *gin.Context) {
	if configuration.Default().SecureMode {
		context.SetSameSite(http.SameSiteStrictMode)
	} else {
		context.SetSameSite(http.SameSiteLaxMode)
	}
}

func SetUpCookieMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		setSameSite(context)
		cookie, err := context.Cookie(SessionCookieName)
		if err != nil {
			cookie = uuid.NewString()
			host := configuration.Default().Host
			if !configuration.Default().SecureMode {
				host = ""
			}
			context.SetCookie(SessionCookieName, cookie, 0, "", host, configuration.Default().SecureMode, true)
		}
		context.Set("identifier", cookie)
		context.Next()
	}
}

func Recover() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Middleware failed: ", r, context.Err())
				if !context.Writer.Written() {
					context.JSON(http.StatusInternalServerError, api.NewErr(ErrUnknownFailure, context.Err()))
				}
			}
		}()
		context.Next()
	}
}

func Error() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
		for _, err := range context.Errors {
			switch err.Err {
			case requests.ErrGraphNotFound:
				context.JSON(http.StatusNotFound, api.NewErr(err, nil))
				return
			case requests.ErrGraphDeleted:
				context.JSON(http.StatusNotFound, api.NewErr(err, nil))
				return
			case requests.ErrUnfinishedGraphBatch:
				context.JSON(http.StatusBadRequest, api.NewErr(err.Err, nil))
				return
			case requests.ErrGraphNotGenerated:
				context.JSON(http.StatusMethodNotAllowed, api.NewErr(err.Err, nil))
				return
			case api.ErrInvalidGraphFormat:
				context.JSON(http.StatusBadRequest, api.NewErr(err.Err, nil))
				return
			}
		}
	}
}

func Logger() gin.HandlerFunc {
	return func(context *gin.Context) {
		time := time2.Now()
		log.Debugf("Starting %s request to %s", context.Request.Method, context.FullPath())
		context.Next()

		duration := time2.Since(time)
		statusCode := context.Writer.Status()
		method := context.Request.Method
		switch {
		case statusCode >= 400 && statusCode < 500:
			log.Warningf("[REQUEST] Request %s to %s failed with %d, took %v", method, context.FullPath(), statusCode, duration)
		case statusCode >= 500:
			log.Errorf("[REQUEST] Request %s to %s failed with %d, took %v", method, context.FullPath(), statusCode, duration)
		default:
			log.Infof("[REQUEST] Request %s to %s finished with %d, took %v", method, context.FullPath(), statusCode, duration)
		}
	}
}

func SetUpRequestService(requestService requests.RequestService) gin.HandlerFunc {
	return func(context *gin.Context) {
		if requestService.CheckMaintenance() {
			log.Warningf("Got request during maintenance")
			context.JSON(http.StatusServiceUnavailable, api.NewErr(ErrSystemMaintenance, nil))
			context.Abort()
			return
		}
		context.Set("requestsService", requestService)
		context.Next()
	}
}
