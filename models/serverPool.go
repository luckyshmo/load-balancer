package models

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"sync/atomic"
	"time"
)

// ServerPool пул бэкендов из конфигурации
type ServerPool struct {
	backends []*Server
	current  uint64
}

func (s *ServerPool) AddBackend(backend *Server) {
	s.backends = append(s.backends, backend)
}

// NextIndex возвращает поочереди номера бэкендов из пула
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// MarkBackendStatus для изменения статуса бэкенда в случае его недоступности или перегрузки.
func (s *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, b := range s.backends {
		if b.URL.String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}

// GetNextPeer получаем слудующий в списке живой бэкенд
func (s *ServerPool) GetNextPeer() *Server {
	// начинаем со следующего бэкенда и идем по всем поочереди, пока не встретим "живой" бэкенд
	next := s.NextIndex()
	l := len(s.backends) + next
	for i := next; i < l; i++ {
		idx := i % len(s.backends)
		if s.backends[idx].IsAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil
}

// HealthCheck ping бэкендов и вывод статуса по каждому
func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		alive := isBackendAlive(b.URL)
		b.SetAlive(alive)
		log.Println(fmt.Sprintf("%s isAlive: [%t]", b.URL, alive))
	}
}

func isBackendAlive(u *url.URL) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	_ = conn.Close()
	return true
}
