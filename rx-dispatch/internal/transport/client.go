package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"rx-dispatch/internal/contracts"
)

// DefaultHTTPClient provides a standard resilient HTTP client for inter-service RPC.
var DefaultHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

// WriteJSON sends a JSON response with status code.
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// WriteError sends a standard error response envelope.
func WriteError(w http.ResponseWriter, status int, message string, serviceName string) {
	resp := contracts.ErrorResponse{
		Error:     message,
		Code:      status,
		Timestamp: time.Now(),
		Service:   serviceName,
	}
	WriteJSON(w, status, resp)
}

// PostJSON performs an HTTP POST sending `payload` as JSON and unmarshaling into `result`.
func PostJSON(url string, payload any, result any) error {
	var bodyReader io.Reader
	if payload != nil {
		buf, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal payload error: %w", err)
		}
		bodyReader = bytes.NewReader(buf)
	}

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "RX-Dispatch-Mesh/2.0")

	resp, err := DefaultHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("http call to %s failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("remote service returned %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode response error: %w", err)
		}
	}
	return nil
}

// GetJSON performs an HTTP GET unmarshaling JSON response into `result`.
func GetJSON(url string, result any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "RX-Dispatch-Mesh/2.0")

	resp, err := DefaultHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("http get to %s failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("remote service returned %d: %s", resp.StatusCode, string(body))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode response error: %w", err)
		}
	}
	return nil
}

// CheckHealth queries the /healthz endpoint of a service.
func CheckHealth(serviceURL string) (contracts.HealthResponse, error) {
	var health contracts.HealthResponse
	url := fmt.Sprintf("%s/healthz", serviceURL)
	err := GetJSON(url, &health)
	return health, err
}
