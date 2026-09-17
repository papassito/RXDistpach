package transport

import "rx-dispatch/internal/contracts"

type AuditClienter interface {
    SendEvent(req contracts.RecordAuditEventRequest) error
}

type ResultClienter interface {
    ConsolidateReading(req interface{}) error
}

type DeliveryClienter interface {
    EnqueueResult(req interface{}) error
}

type AuditClient struct{}

func NewAuditClient() *AuditClient { return &AuditClient{} }
func (c *AuditClient) SendEvent(req contracts.RecordAuditEventRequest) error { return nil }

type ResultClient struct{}

func NewResultClient() *ResultClient { return &ResultClient{} }
func (c *ResultClient) ConsolidateReading(req interface{}) error { return nil }

type DeliveryClient struct{}

func NewDeliveryClient() *DeliveryClient { return &DeliveryClient{} }
func (c *DeliveryClient) EnqueueResult(req interface{}) error { return nil }