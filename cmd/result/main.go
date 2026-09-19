package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/models"
	"rx-dispatch/internal/transport"
)

const (
	serviceName = "rx-result"
	defaultPort = "8086"
	version     = "0.1.0-scaffold"
)

// server encapsula las dependencias para los manejadores HTTP.
type server struct {
	deliveryClient transport.DeliveryClienter
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

	srv := &server{
		deliveryClient: transport.NewDeliveryClient(),
	}

	http.HandleFunc("/healthz", srv.healthCheckHandler)
	http.HandleFunc("/result/consolidate", srv.consolidateHandler)

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("[%s] Listening on %s", serviceName, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[%s] Failed to start server: %v", serviceName, err)
	}
}

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
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[%s] ERROR: Failed to write health check response: %v", serviceName, err)
	}
}

func (s *server) consolidateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	// Inyectar protecciÃ³n DoS
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB limit

	// 2. Decodificar cuerpo de la peticiÃ³n
	var req models.GenericReading
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid or malformed JSON payload", http.StatusBadRequest)
		return
	}

	// 3. Ejecutar la llamada al cliente de entrega (Requerido por la prueba unitaria)
	if s.deliveryClient != nil {
		if err := s.deliveryClient.EnqueueResult(req); err != nil {
			log.Printf("[%s] ERROR: Failed to enqueue result for study %s: %v", serviceName, req.StudyID, err)
			http.Error(w, "Failed to enqueue result for delivery", http.StatusInternalServerError)
			return
		}
	}

	// 4. Responder con cÃ³digo 200 OK y confirmaciÃ³n
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "success"}); err != nil {
		log.Printf("[%s] ERROR: Failed to write response for study %s: %v", serviceName, req.StudyID, err)
	}
}
