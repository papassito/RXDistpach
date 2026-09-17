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
	serviceName = "rx-image"
	defaultPort = "8084"
	version     = "0.1.0-scaffold"
)

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

	log.Printf("[%s] Starting service scaffold on port %s", serviceName, port)

	http.HandleFunc("/healthz", healthCheckHandler)

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("[%s] Listening on %s", serviceName, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[%s] Failed to start server: %v", serviceName, err)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
