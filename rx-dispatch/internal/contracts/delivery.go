package contracts

import (
	"rx-dispatch/internal/models"
)

// DispatchDeliveryRequest instructs RX Delivery to transmit a result package.
type DispatchDeliveryRequest struct {
	PackageID models.ResultPackage `json:"package"`
	Recipient string               `json:"recipient"`
	Channel   string               `json:"channel"` // EMAIL, SMS, PORTAL, SECURE_LINK
}

// DispatchDeliveryResponse reports dispatch initiation and tracking token.
type DispatchDeliveryResponse struct {
	DeliveryID    string                `json:"deliveryId"`
	TrackingToken string                `json:"trackingToken"`
	Status        models.DeliveryStatus `json:"status"`
	Recipient     string                `json:"recipient"`
	Channel       string                `json:"channel"`
}

// QueryDeliveriesResponse returns dispatch history for a study or package.
type QueryDeliveriesResponse struct {
	Deliveries []models.DeliveryRecord `json:"deliveries"`
	Count      int                     `json:"count"`
}
