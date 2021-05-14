package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/luckyshmo/load-balancer/models"
)

type ServerPoolHandler struct {
	ServerPool *models.ServerPool
}

func NewServerHandler(serverPool *models.ServerPool) *ServerPoolHandler {
	return &ServerPoolHandler{
		ServerPool: serverPool,
	}
}

const (
	Attempts int = iota
	Retry
)

func (h *ServerPoolHandler) getAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

func (h *ServerPoolHandler) getRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}

// loadBalancingProxy основная функция баланисровки пакетов
func (h *ServerPoolHandler) LoadBalancingProxy(w http.ResponseWriter, r *http.Request) {
	attempts := h.getAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	backend := h.ServerPool.GetNextPeer()
	if backend != nil {
		fmt.Print(backend.URL)
		backend.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}
