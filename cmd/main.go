package main

import (
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

	serverList := os.Getenv("SERVER_LIST")
	port := os.Getenv("APP_PORT")

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
		Addr:    ":" + port,
		Handler: http.HandlerFunc(serverHandler.LoadBalancingProxy),
	}

	go service.HealthCheck(serverPool)

	log.Printf("Load Balancer started at :%s\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	quit := make(chan os.Signal, 1)
	//if app get SIGTERM it will exit
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("App Shutting Down")

}
