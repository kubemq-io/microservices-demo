package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Log struct {
	ID       string `json:"id"`
	ClientID string `json:"client_id"`
	Metadata string `json:"metadata"`
	Body     string `json:"body"`
}

func main() {
	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGKILL)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)
	cfg, err := LoadConfig()
	if err != nil {
		log.Println("error on loading config file:")
		log.Fatal(err)
	}
	var el *Elastic
	log.Println("Wait for Elastic to be ready")
	for {
		el, err = NewElasticSearch(cfg.ElasticAddress)
		if err != nil {
			log.Printf("error connecting to elastic, error: %s retrying...\n",err.Error())
		}else {
			break
		}
	}
	el.Save(context.Background(),&History{
		Id:           uuid.New().String(),
		Source:       "history-service",
		Time:         time.Now(),
		Type:         "event",
		Method:       "init",
		Request:      "start history service",
		Response:     "",
		IsError:      false,
		ErrorMessage: "",
	})

	kube, err := NewKubeMQClient(cfg.KubeMQHost, cfg.KubeMQPort)
	if err != nil {

		log.Fatal(err)
	}
	eventsCh := make(chan *kubemq.Event, 1)
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("Wait for kubemq to be ready")
	for {
		err = kube.StartListen(ctx, cfg.Channel, cfg.Group, eventsCh, errCh)
		if err != nil {
			log.Printf("error connecting to kubemq, error: %s, retrying...\n",err.Error())
			time.Sleep(time.Second)
		}else {

			break

		}
	}

	log.Println("waiting for events from KubeMQ")
	for {
		select {
		case event := <-eventsCh:
			his:=&History{}
			err:=json.Unmarshal(event.Body,his)
			if err != nil {
				log.Println(err)
				continue
			}
			err = el.Save(ctx, his)
			if err != nil {
				log.Println(err)
			}

		case err := <-errCh:
			log.Fatal(err)
		case <-gracefulShutdown:
			kube.Close()
			return
		}
	}
}
