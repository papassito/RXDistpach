package config

import (
	"os"
	"strconv"
)

// ServiceConfig holds the configuration for an individual service process.
type ServiceConfig struct {
	ServiceName string
	Port        int
	Environment string
	Host        string
}

// TopologyConfig maps the communication addresses for all 9 microservices.
type TopologyConfig struct {
	GatewayURL  string
	SecurityURL string
	StudyURL    string
	StorageURL  string
	ImageURL    string
	ReaderURL   string
	ResultURL   string
	DeliveryURL string
	AuditURL    string
}

// GetEnv helper returns value or fallback.
func GetEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// GetEnvInt helper returns integer value or fallback.
func GetEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

// LoadTopologyConfig loads standard discovery URLs for the 9 microservices.
func LoadTopologyConfig() TopologyConfig {
	return TopologyConfig{
		GatewayURL:  GetEnv("RX_GATEWAY_URL", "http://127.0.0.1:8089"),
		SecurityURL: GetEnv("RX_SECURITY_URL", "http://127.0.0.1:8081"),
		StudyURL:    GetEnv("RX_STUDY_URL", "http://127.0.0.1:8082"),
		StorageURL:  GetEnv("RX_STORAGE_URL", "http://127.0.0.1:8083"),
		ImageURL:    GetEnv("RX_IMAGE_URL", "http://127.0.0.1:8084"),
		ReaderURL:   GetEnv("RX_READER_URL", "http://127.0.0.1:8085"),
		ResultURL:   GetEnv("RX_RESULT_URL", "http://127.0.0.1:8086"),
		DeliveryURL: GetEnv("RX_DELIVERY_URL", "http://127.0.0.1:8087"),
		AuditURL:    GetEnv("RX_AUDIT_URL", "http://127.0.0.1:8088"),
	}
}
