package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/kubemq-io/kubemq-go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func PrettyJson(data interface{}) string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return ""
	}
	return buffer.String()
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
	var redis *Redis
	log.Println("Wait for redis to be ready")
	for {
		redis, err = NewRedisClient(cfg.RedisAddress)
		if err != nil {
			log.Printf("error connecting to redis, error: %s retrying...\n", err.Error())
		} else {
			break
		}
	}

	kube, err := NewKubeMQClient(cfg.KubeMQHost, cfg.KubeMQPort, cfg.HistoryChannel)
	if err != nil {
		log.Fatal(err)
	}
	commandsCh := make(chan *kubemq.CommandReceive, 1)
	queriesCh := make(chan *kubemq.QueryReceive, 1)
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("Wait for kubemq to be ready")
	for {
		err = kube.StartListenToCommands(ctx, cfg.Channel, cfg.Group, commandsCh, errCh)
		if err != nil {
			log.Printf("error connecting to kubemq, error: %s, retrying...\n", err.Error())
			time.Sleep(time.Second)
		} else {

			break

		}
	}
	for {
		err = kube.StartListenToQueries(ctx, cfg.Channel, cfg.Group, queriesCh, errCh)
		if err != nil {
			log.Printf("error connecting to kubemq, error: %s, retrying...\n", err.Error())
			time.Sleep(time.Second)
		} else {

			break

		}
	}
	log.Println("waiting for commands / queries from KubeMQ")
	for {
		select {
		case command := <-commandsCh:

			resp := &kubemq.Response{
				RequestId:  command.Id,
				ResponseTo: command.ResponseTo,
				Metadata:   command.Metadata,
				Body:       nil,
			}
			cm, err := NewCacheMessage(command.Body)
			if err != nil {
				log.Printf("error on parsing of cache command: %s\n", err.Error())
				continue
			}
			switch cm.Type {
			case "set":

				err = redis.Set(cm.Key, cm.Value, cm.Expiry)
			case "del":
				err = redis.Del(cm.Key)
			default:
				log.Printf("invalid cache command: %s\n", cm.Type)
				continue
			}
			log.Println(fmt.Sprintf("cache command received - Type: %s Key: %s, Value: %s", cm.Type, cm.Key, PrettyJson(cm.Value)))
			if err != nil {
				log.Printf("error on sending command to redis: %s\n", err.Error())
				resp.Err = err

			} else {
				resp.ExecutedAt = time.Now()
			}
			err = kube.SendResponse(ctx, resp)
			if err != nil {
				log.Printf("error on sending response from redis: %s\n", err.Error())

			}
			his := &History{
				Id:           uuid.New().String(),
				Source:       "cache-service",
				Time:         time.Now(),
				Type:         "command",
				Method:       cm.Type,
				Request:      PrettyJson(cm.Value),
				Response:     "",
				IsError:      false,
				ErrorMessage: "",
			}
			if err != nil {
				his.IsError = true
				his.ErrorMessage = err.Error()
			}
			go kube.SendHistory(ctx, his)
		case query := <-queriesCh:

			resp := &kubemq.Response{
				RequestId:  query.Id,
				ResponseTo: query.ResponseTo,
				Metadata:   query.Metadata,
			}
			cm, err := NewCacheMessage(query.Body)
			if err != nil {
				log.Printf("error on parsing of cache command: %s\n", err.Error())
				continue
			}
			var result []byte
			switch cm.Type {
			case "get":
				result, err = redis.Get(cm.Key)
			default:
				log.Printf("invalid cache command: %s\n", cm.Type)
				continue
			}
			log.Println(fmt.Sprintf("cache query received - Type: %s Key: %s", cm.Type, cm.Key))
			if err != nil {
				log.Printf("error on sending command to redis: %s\n", err.Error())
				resp.Err = err

			} else {
				resp.ExecutedAt = time.Now()
				resp.Body = result
			}
			err = kube.SendResponse(ctx, resp)
			if err != nil {
				log.Printf("error on sending response from redis: %s\n", err.Error())

			}
			his := &History{
				Id:           uuid.New().String(),
				Source:       "cache-service",
				Time:         time.Now(),
				Type:         "query",
				Method:       cm.Type,
				Request:      PrettyJson(cm.Key),
				Response:     PrettyJson(result),
				IsError:      false,
				ErrorMessage: "",
			}
			if err != nil {
				his.IsError = true
				his.ErrorMessage = err.Error()
			}
			go kube.SendHistory(ctx, his)
		case err := <-errCh:
			log.Fatal(err)
		case <-gracefulShutdown:
			kube.Close()
			return
		}
	}
}
