package handler

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func (h *ServerPoolHandler) GetProxy(serverUrl *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(serverUrl)
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
		log.Printf("[%s] %s\n", serverUrl.Host, e.Error())
		retries := h.getRetryFromContext(request)
		if retries < 3 {
			select {
			case <-time.After(10 * time.Millisecond):
				ctx := context.WithValue(request.Context(), Retry, retries+1)
				proxy.ServeHTTP(writer, request.WithContext(ctx))
			}
			return
		}

		// считаем, что бэкенд мертв в случае 3 неудачных попыток
		h.ServerPool.MarkBackendStatus(serverUrl, false)

		attempts := h.getAttemptsFromContext(request)
		ctx := context.WithValue(request.Context(), Attempts, attempts+1)
		h.LoadBalancingProxy(writer, request.WithContext(ctx))
	}

	return proxy
}
