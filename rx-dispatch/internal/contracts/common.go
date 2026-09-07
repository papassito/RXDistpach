package contracts

import "time"

// HealthResponse is returned by the /healthz endpoint of all 9 independent services.
type HealthResponse struct {
	Service   string    `json:"service"`
	Port      int       `json:"port"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// ErrorResponse standardizes error messages across all service contracts.
type ErrorResponse struct {
	Error     string    `json:"error"`
	Code      int       `json:"code"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// SuccessResponse standardizes generic success operations.
type SuccessResponse struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
