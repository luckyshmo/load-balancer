package service

import (
	"log"
	"time"

	"github.com/luckyshmo/load-balancer/models"
)

// healthCheck проверяем доступность бэкендов каждые 10 секунд
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
