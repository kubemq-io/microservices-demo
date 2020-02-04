package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

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

	kube, err := NewKubeMQClient(cfg.KubeMQHost, cfg.KubeMQPort, cfg.HistoryChannel)
	if err != nil {
		log.Fatal(err)
	}
	_, err = NewServer(kube, cfg)
	if err != nil {
		log.Fatal(err)
	}

	<-gracefulShutdown
	kube.Close()

}
