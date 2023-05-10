package web

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/soch-fit/GraphGenerator/pkg/configuration"
	"github.com/soch-fit/GraphGenerator/pkg/routers/middleware"
	"path/filepath"
)

var (
	ErrNoUiPresent = errors.New("no ui present to setup")
)

func handleNoRoute(ctx *gin.Context) {
	origUri := ctx.Request.RequestURI
	dir, file := filepath.Split(origUri)
	ext := filepath.Ext(file)
	log.Debugf("Got request to %s %s", file, ext)
	if dir == "/" && (file == "" || ext == "") {
		ctx.File(fmt.Sprintf("%s/index.html", *configuration.Default().UiLocation))
		return
	}
}

func SetupWeb(engine *gin.Engine) error {
	if configuration.Default().UiLocation == nil {
		return ErrNoUiPresent
	}
	log.Debugf("Serving data from %s", *configuration.Default().UiLocation)
	engine.NoRoute(handleNoRoute)
	storage := static.LocalFile(*configuration.Default().UiLocation, false)
	engine.Use(middleware.SetUpCookieMiddleware())
	engine.Use(static.Serve("/", storage))
	return nil
}
