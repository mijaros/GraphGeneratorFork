package service

import (
	"github.com/soch-fit/GraphGenerator/pkg/api"
	"testing"
	"time"
)

func TestBasicGeneration(t *testing.T) {
	genService := New(1)
	results := genService.GetRetriever()
	checker := make(chan error)

	err := genService.Start()
	if err != nil {
		t.Error("GenService start ended", err)
	}

	seed := int64(255)
	request := api.GraphRequest{
		Type:  api.Complete,
		Nodes: 3,
		Seed:  &seed,
		ID:    2233,
	}
	timeout := time.NewTimer(1 * time.Second)
	go func() {
		err := genService.PushRequest(request)
		checker <- err
	}()

	select {
	case err = <-checker:
		if err != nil {
			t.Error("Push failed", err.Error())
		}
	case <-timeout.C:
		t.Error("Push didn't finish in time.")
	}
	if !timeout.Stop() {
		<-timeout.C
	}

	timeout.Reset(2 * time.Second)
	select {
	case dat := <-results:
		if dat.ID != 2233 {
			t.Error("Invalid result poped")
		}
		if len(dat.Generated.Edges()) != 3 {
			t.Error("Invalid result obtained")
		}
	case <-timeout.C:
		t.Error("Graph was not generated in time")
	}

}

func TestPauseAndResume(t *testing.T) {
	counter := 0
	//results := make(chan *api.GraphResult, 5)
	genService := New(5)
	results := genService.GetRetriever()
	checker := make(chan error)

	err := genService.Start()
	if err != nil {
		t.Error("GenService start ended", err)
	}

	seed := int64(255)
	request := api.GraphRequest{
		Type:  api.Complete,
		Nodes: 3,
		Seed:  &seed,
		ID:    2233,
	}
	go func() {
		for i := 0; i < 10; i++ {
			err = genService.PushRequest(request)
			if err != nil {
				checker <- err
				return
			}
		}
		checker <- nil
	}()
	err = genService.Pause()
	if err != nil {
		t.Error("GenService was not paused", err.Error())
	}

	timeout := time.NewTimer(2 * time.Second)
	select {
	case err = <-checker:
		if err != nil {
			t.Error("Pushing failed", err.Error())
		}
	case <-timeout.C:
		t.Error("Timeout reached for push!")
	}
	timeout = time.NewTimer(2 * time.Second)

wait:
	for {
		select {
		case d := <-results:
			counter++
			if d.ID != 2233 {
				t.Error("Invalid graph obtained")
			} else {
				t.Log("Emptying after pause", d.ID)
			}
		case <-timeout.C:
			break wait
		}
	}

	timeout = time.NewTimer(2 * time.Second)
	select {
	case d := <-results:
		counter++
		t.Error("Result delivered after service was paused", d.ID)
	case <-timeout.C:
		break
	}

	genService.Resume()

	timeout = time.NewTimer(1 * time.Second)
wait2:
	for {
		select {
		case d := <-results:
			counter++
			if d.ID != 2233 {
				t.Error("Invalid graph obtained")
			}
			if counter == 10 {
				break wait2
			}
		case <-timeout.C:
			t.Error("Timeout reached for popping from service", counter)
			break wait2
		}
	}
	genService.Stop()

}
