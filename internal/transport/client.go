package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"rx-dispatch/internal/contracts"
	"time"
)

type AuditClienter interface {
	SendEvent(req contracts.RecordAuditEventRequest) error
}

type ResultClienter interface {
	ConsolidateReading(req interface{}) error
}

type DeliveryClienter interface {
	EnqueueResult(req interface{}) error
}

type baseClient struct {
	client   *http.Client
	endpoint string
}

func (bc *baseClient) doPost(path string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %w", err)
	}

	url := bc.endpoint + path
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := bc.client.Do(req)
	if err != nil {
		return fmt.Errorf("http request to %s failed: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("received non-2xx status code %d from %s", resp.StatusCode, url)
	}

	return nil
}

type AuditClient struct{ baseClient }

func NewAuditClient() *AuditClient {
	endpoint := os.Getenv("RX_AUDIT_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8088" // Default for local dev
	}
	return &AuditClient{
		baseClient: baseClient{
			client: &http.Client{
				Timeout: 5 * time.Second,
			},
			endpoint: endpoint,
		},
	}
}
func (c *AuditClient) SendEvent(req contracts.RecordAuditEventRequest) error {
	return c.doPost("/audit/event", req)
}

type ResultClient struct{ baseClient }

func NewResultClient() *ResultClient {
	endpoint := os.Getenv("RX_RESULT_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8086" // Default for local dev
	}
	return &ResultClient{
		baseClient: baseClient{
			client: &http.Client{
				Timeout: 5 * time.Second,
			},
			endpoint: endpoint,
		},
	}
}
func (c *ResultClient) ConsolidateReading(req interface{}) error {
	return c.doPost("/result/consolidate", req)
}

type DeliveryClient struct{ baseClient }

func NewDeliveryClient() *DeliveryClient {
	endpoint := os.Getenv("RX_DELIVERY_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8087" // Default for local dev
	}
	return &DeliveryClient{
		baseClient: baseClient{
			client: &http.Client{
				Timeout: 10 * time.Second, // Delivery might take longer
			},
			endpoint: endpoint,
		},
	}
}
func (c *DeliveryClient) EnqueueResult(req interface{}) error {
	return c.doPost("/delivery/enqueue", req)
}
