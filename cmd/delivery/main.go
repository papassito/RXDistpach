package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "rx-dispatch/internal/contracts"
)

const (
    serviceName = "rx-delivery"
    servicePort = 8085
    version     = "0.1.0-scaffold"
)

type server struct{}

func (s *server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    response := contracts.HealthResponse{
        Service:   serviceName,
        Port:      servicePort,
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
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }
    w.WriteHeader(http.StatusAccepted)
}

func main() {
    log.Printf("[%s] Starting service on port %d", serviceName, servicePort)

    srv := &server{}
    http.HandleFunc("/healthz", srv.healthCheckHandler)
    http.HandleFunc("/delivery/enqueue", srv.enqueueHandler)

    addr := fmt.Sprintf("127.0.0.1:%d", servicePort)
    log.Printf("[%s] Listening on %s", serviceName, addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("[%s] Server failed: %v", serviceName, err)
    }
}