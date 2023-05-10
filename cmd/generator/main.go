package main

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/soch-fit/GraphGenerator/pkg/configuration"
	"github.com/soch-fit/GraphGenerator/pkg/generator/service"
	"github.com/soch-fit/GraphGenerator/pkg/requests/persistent"
	"github.com/soch-fit/GraphGenerator/pkg/routers"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var processEndingWg = sync.WaitGroup{}

func processEndingTimeout() {
	timer := time.NewTimer(30 * time.Second)
	stopped := make(chan bool)
	go func() {
		processEndingWg.Wait()
		stopped <- true
	}()

	select {
	case <-stopped:
		return
	case <-timer.C:
		log.Panic("Server didn't stop in time, killing now")
	}
}

func main() {
	processEndingWg.Add(1)
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Unhandled panic ending program %v", r)
		} else {
			log.Info("Server shut down correctly, bye!")
		}
		processEndingWg.Done()
	}()
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	configuration.ParseFlags(flags)
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Panicln("Couldn't parse flags")
	}
	log.SetLevel(configuration.Default().LogLevel)
	log.Info("Starting the application")

	grService := service.New(configuration.Default().Workers)
	if err = grService.Start(); err != nil {
		log.Panicf("Couldn't start the graph generating service %s", err.Error())
	}
	defer grService.Stop()
	requestsService, err := persistent.New(grService)
	if err != nil {
		log.Panicf("Couldn't open service: %v", err)
	}

	if err = requestsService.Start(); err != nil {
		log.Panicf("Couldn't start the database service")
	}
	defer requestsService.Stop()

	httpService := routers.New()
	if err = httpService.Start(requestsService); err != nil {
		log.Panicf("Server didn't start up properly %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		log.Info("Stopping the web server")
		err = httpService.Stop(ctx)
		if err != nil {
			log.Error("Service stopping failed with ", err.Error())
		}
		select {
		case <-ctx.Done():
			log.Info("Web Server shut down.")
		}
	}()

	log.Info("Startup finished, blocking for signals to appear")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Warn("Terminating signal caught! Shutting down!")
	go processEndingTimeout()
}
