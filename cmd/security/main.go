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
    serviceName = "rx-security"
    servicePort = 8081
    version     = "0.1.0-scaffold"
)

func main() {
    log.Printf("[%s] Starting service scaffold on port %d", serviceName, servicePort)

    http.HandleFunc("/healthz", healthCheckHandler)

    addr := fmt.Sprintf("127.0.0.1:%d", servicePort)
    log.Printf("[%s] Listening on %s", serviceName, addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("[%s] Failed to start server: %v", serviceName, err)
    }
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    response := contracts.HealthResponse{
        Service:   serviceName,
        Port:      servicePort,
        Status:    "UP",
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        Version:   version,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(response)
}
