package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"rx-dispatch/internal/contracts"
)

const (
	serviceName = "rx-delivery"
	defaultPort = "8087" // Corrected port from audit
	version     = "0.1.0-scaffold"
)

type server struct{}

func (s *server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	response := contracts.HealthResponse{
		Service:   serviceName,
		Status:    "UP",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   version,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (s *server) enqueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Protect against DoS attacks
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB limit

	// In a real implementation, we would decode the body and add it to a queue.
	// For now, just acknowledging the request is sufficient.
	// We must, however, consume and close the body to avoid resource leaks.
	defer r.Body.Close()

	w.WriteHeader(http.StatusAccepted)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	host := os.Getenv("HOST")
	// Per REQ-SEC-001, internal services must default to loopback.
	if host == "" {
		host = "127.0.0.1"
	}

	log.Printf("[%s] Starting service on port %s", serviceName, port)

	srv := &server{}
	http.HandleFunc("/healthz", srv.healthCheckHandler)
	http.HandleFunc("/delivery/enqueue", srv.enqueueHandler)

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("[%s] Listening on %s", serviceName, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[%s] Server failed: %v", serviceName, err)
	}
}
