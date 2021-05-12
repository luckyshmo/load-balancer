package service

import (
	"log"
	"time"

	"github.com/luckyshmo/load-balancer/models"
)

// healthCheck runs a routine for check status of the backends every 2 mins
func HealthCheck(serverPool models.ServerPool) {
	t := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-t.C:
			log.Println("Starting health check...")
			serverPool.HealthCheck()
			log.Println("Health check completed")
		}
	}
}
