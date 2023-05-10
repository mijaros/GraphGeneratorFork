package routers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/soch-fit/GraphGenerator/pkg/configuration"
	"github.com/soch-fit/GraphGenerator/pkg/requests"
	"github.com/soch-fit/GraphGenerator/pkg/routers/api/v1"
	"github.com/soch-fit/GraphGenerator/pkg/routers/middleware"
	"github.com/soch-fit/GraphGenerator/pkg/routers/web"
	"github.com/soch-fit/GraphGenerator/pkg/utils"
	"net/http"
	"time"
)

var (
	ErrTimeoutReached = errors.New("timeout reached")
)

type HTTPService interface {
	Start() error
	Stop(ctx context.Context) error
}

type GinHttpService struct {
	srv *http.Server
	err chan error
}

func New() *GinHttpService {
	return &GinHttpService{srv: nil}
}

func (service *GinHttpService) startServer() {
	defer utils.StopProcessOnUnhandledPanic()
	err := service.srv.ListenAndServe()
	service.err <- err
}

func (service *GinHttpService) Start(requestService requests.RequestService) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Recover())
	r.Use(middleware.Logger())
	v1.SetupREST(r, requestService)
	err := web.SetupWeb(r)
	if err != nil {
		return err
	}

	service.srv = &http.Server{
		Addr:    configuration.Default().BindAddr + ":" + configuration.Default().Port,
		Handler: r,
	}

	service.err = make(chan error)
	timer := time.After(5 * time.Second)

	log.Debug("Waiting for five seconds to start the server")
	go service.startServer()

	select {
	case err = <-service.err:
		return err
	case <-timer:
		log.Infof("Listening on %s:%s", configuration.Default().BindAddr, configuration.Default().Port)
		return nil
	}
}

func (service *GinHttpService) Stop(ctx context.Context) error {
	shutDownctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	err := service.srv.Shutdown(shutDownctx)
	if err != nil {
		return nil
	}
	select {
	case err = <-service.err:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	case <-shutDownctx.Done():
		return ErrTimeoutReached
	}
}
