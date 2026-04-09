package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TimeResponse struct {
	Time string `json:"time"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TimeResponse{Time: time.Now().Format(time.RFC3339)})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}

func selfCheck() {
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil || resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
	os.Exit(0)
}

func newServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/time", timeHandler)
	mux.HandleFunc("/health", healthHandler)
	return &http.Server{Addr: ":8080", Handler: mux}
}

func runServer(srv *http.Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	}
}

func main() {
	healthCheck := flag.Bool("health", false, "perform health check and exit")
	flag.Parse()

	if *healthCheck {
		selfCheck()
	}

	srv := newServer()
	if err := runServer(srv); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
