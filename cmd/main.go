package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/luckyshmo/load-balancer/models"
	"github.com/luckyshmo/load-balancer/pkg/handler"
	"github.com/luckyshmo/load-balancer/pkg/service"
)

var serverPool models.ServerPool

func main() {
	var serverList string
	var port int
	// flag.StringVar(&serverList, "backends", "", "Load balanced backends, use commas to separate")
	// flag.IntVar(&port, "port", 3030, "Port to serve")
	// flag.Parse()

	port = 9090
	serverList = "http://localhost:8081,http://localhost:8082,http://localhost:8083,http://localhost:8084,http://localhost:8080"

	if len(serverList) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}

	serverHandler := handler.NewServerHandler(&serverPool)

	// parse servers
	tokens := strings.Split(serverList, ",")
	for _, tok := range tokens {
		serverUrl, err := url.Parse(tok)
		if err != nil {
			log.Fatal(err)
		}

		serverPool.AddBackend(&models.Server{
			URL:          serverUrl,
			Alive:        true,
			ReverseProxy: serverHandler.GetProxy(serverUrl),
		})
		log.Printf("Configured server: %s\n", serverUrl)
	}

	// create http server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(serverHandler.LoadBalancingProxy),
	}

	go service.HealthCheck(serverPool)

	log.Printf("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	//if app get SIGTERM it will exit
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("App Shutting Down")

}
