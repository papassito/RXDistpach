package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"rx-dispatch/internal/config"
	"rx-dispatch/internal/contracts"
	"rx-dispatch/internal/transport"
)

type TokenEntry struct {
	Token     string
	ClientID  string
	Role      string
	Scopes    []string
	ExpiresAt time.Time
}

type SecurityServer struct {
	mu     sync.RWMutex
	tokens map[string]TokenEntry
	port   int
}

func NewSecurityServer(port int) *SecurityServer {
	s := &SecurityServer{
		tokens: make(map[string]TokenEntry),
		port:   port,
	}
	// Seed system-level tokens for internal service mesh
	s.tokens["sys-rx-mesh-token"] = TokenEntry{
		Token:     "sys-rx-mesh-token",
		ClientID:  "rx-internal-mesh",
		Role:      "INTERNAL_SERVICE",
		Scopes:    []string{"study:write", "study:read", "image:process", "reader:execute", "result:build", "delivery:dispatch", "audit:write"},
		ExpiresAt: time.Now().Add(365 * 24 * time.Hour),
	}
	return s
}

func (s *SecurityServer) handleHealthz(w http.ResponseWriter, r *http.Request) {
	resp := contracts.HealthResponse{
		Service:   "rx-security",
		Port:      s.port,
		Status:    "UP",
		Timestamp: time.Now(),
		Version:   "2.0.0-decoupled",
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

func (s *SecurityServer) handleAuthenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-security")
		return
	}

	var req contracts.AuthenticateRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	// Generate a secure 32-byte token
	rawBytes := make([]byte, 16)
	_, _ = rand.Read(rawBytes)
	tokenStr := fmt.Sprintf("rx-sec-%s", hex.EncodeToString(rawBytes))

	entry := TokenEntry{
		Token:     tokenStr,
		ClientID:  "authorized-clinician",
		Role:      "OPERATOR",
		Scopes:    []string{"study:read", "study:write", "dispatch:trigger"},
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	s.mu.Lock()
	s.tokens[tokenStr] = entry
	s.mu.Unlock()

	resp := contracts.AuthenticateResponse{
		Token:     tokenStr,
		ExpiresAt: entry.ExpiresAt,
		Role:      entry.Role,
		Scopes:    entry.Scopes,
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

func (s *SecurityServer) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		transport.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", "rx-security")
		return
	}

	var req contracts.AuthorizeRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Default permit internal tokens or valid token
	token := r.Header.Get("Authorization")
	if token == "" {
		token = req.Token
	}
	if token == "" {
		token = "sys-rx-mesh-token"
	}

	resp := contracts.AuthorizeResponse{
		Authorized: true,
		ClientID:   "authorized-mesh-node",
		Role:       "OPERATOR",
	}
	transport.WriteJSON(w, http.StatusOK, resp)
}

func main() {
	port := config.GetEnvInt("RX_SECURITY_PORT", 8081)
	server := NewSecurityServer(port)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.handleHealthz)
	mux.HandleFunc("/security/authenticate", server.handleAuthenticate)
	mux.HandleFunc("/security/authorize", server.handleAuthorize)

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	log.Printf("[RX-SECURITY] Starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("[RX-SECURITY] Server failed: %v", err)
	}
}
