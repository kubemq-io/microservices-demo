package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const envFile = `
(function (window) {
  window.__env = window.__env || {};
  window.__env.apiUrl = '%s';
  window.__env.enableDebug = true;
}(this));
`

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
	log.Println("staring web server....")
	log.Println(fmt.Sprintf("Api Address: %s", cfg.ApiAddress))
	newEnvFile := fmt.Sprintf(envFile, cfg.ApiAddress)
	err = ioutil.WriteFile(`./web/users/dist/users/env.js`, []byte(newEnvFile), 0644)
	if err != nil {
		log.Fatalln(err.Error())
	}
	http.Handle("/", http.FileServer(http.Dir("./web/users/dist/users")))
	go http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), nil)
	<-gracefulShutdown
	log.Println("shutdown web server....")

}
