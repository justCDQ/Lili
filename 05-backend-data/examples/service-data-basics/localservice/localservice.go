package localservice

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func Start() (baseURL string, stop func(context.Context) error, err error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, fmt.Errorf("listen: %w", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok\n")
	})
	server := &http.Server{Handler: mux, ReadHeaderTimeout: 2 * time.Second, IdleTimeout: 10 * time.Second}
	go func() { _ = server.Serve(listener) }()
	return "http://" + listener.Addr().String(), server.Shutdown, nil
}
