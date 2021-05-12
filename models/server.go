package models

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

// Server данные одного сервера
type Server struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex //записей будет больше, чем чтения.
	ReverseProxy *httputil.ReverseProxy
}

func (b *Server) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

func (b *Server) IsAlive() (alive bool) {
	b.mux.RLock()
	alive = b.Alive
	b.mux.RUnlock()
	return
}
