package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/models"
	"rx-dispatch/internal/transport"
)

const (
	serviceName = "rx-result"
	servicePort = 8086
	version     = "0.1.0-scaffold"
)

// server encapsula las dependencias para los manejadores HTTP.
type server struct {
	deliveryClient transport.DeliveryClienter
}

func main() {
	log.Printf("[%s] Starting service on port %d", serviceName, servicePort)

	srv := &server{
		deliveryClient: transport.NewDeliveryClient(),
	}

	http.HandleFunc("/healthz", srv.healthCheckHandler)
	http.HandleFunc("/result/consolidate", srv.consolidateHandler)

	addr := fmt.Sprintf("127.0.0.1:%d", servicePort)
	log.Printf("[%s] Listening on %s", serviceName, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("[%s] Failed to start server: %v", serviceName, err)
	}
}

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

func (s *server) consolidateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Decodificar cuerpo de la petición
	var req models.GenericReading
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 3. Ejecutar la llamada al cliente de entrega (Requerido por la prueba unitaria)
	if s.deliveryClient != nil {
		s.deliveryClient.EnqueueResult(req)
	}

	// 4. Responder con código 200 OK y confirmación
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}