package delivery

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"rx-dispatch/internal/models"
)

// QueueItemStatus defines the delivery state machine in the store-and-forward queue.
type QueueItemStatus string

const (
	StatusQueued        QueueItemStatus = "QUEUED"
	StatusRetrying      QueueItemStatus = "RETRYING"
	StatusDelivered     QueueItemStatus = "DELIVERED"
	StatusPermanentFail QueueItemStatus = "PERMANENT_FAIL"
)

// QueueItem represents a resiliently queued delivery on local disk.
type QueueItem struct {
	ID                     string          `json:"id"`
	StudyID                string          `json:"studyId"`
	PackageID              string          `json:"packageId"`
	Recipient              string          `json:"recipient"`
	Channel                string          `json:"channel"`
	TrackingToken          string          `json:"trackingToken"`
	PayloadChecksum        string          `json:"payloadChecksum"`
	Status                 QueueItemStatus `json:"status"`
	Attempts               int             `json:"attempts"`
	MaxAttempts            int             `json:"maxAttempts"`
	EnqueuedAt             time.Time       `json:"enqueuedAt"`
	NextAttemptAt          time.Time       `json:"nextAttemptAt"`
	LastAttemptAt          *time.Time      `json:"lastAttemptAt,omitempty"`
	DeliveredAt            *time.Time      `json:"deliveredAt,omitempty"`
	LastError              string          `json:"lastError,omitempty"`
	BackoffIntervalSeconds int             `json:"backoffIntervalSeconds"`
}

// StoreAndForwardQueue manages persistent local queueing for unreliable networks.
type StoreAndForwardQueue struct {
	mu                     sync.RWMutex
	filePath               string
	items                  map[string]*QueueItem
	isNetworkOnline        bool
	simulatedNetworkDrop   bool
	onSuccessfulDeliveryFn func(item *QueueItem)
	onAuditFn              func(eventType, studyID, actor, details string)
	stopChan               chan struct{}
}

// NewStoreAndForwardQueue initializes queue with disk persistence.
func NewStoreAndForwardQueue(dataDir string, onDelivery func(item *QueueItem), onAudit func(eventType, studyID, actor, details string)) *StoreAndForwardQueue {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("[STORE-FORWARD] Warning creating data dir: %v", err)
	}

	queueFilePath := filepath.Join(dataDir, "delivery_queue.json")
	q := &StoreAndForwardQueue{
		filePath:               queueFilePath,
		items:                  make(map[string]*QueueItem),
		isNetworkOnline:        true,
		simulatedNetworkDrop:   false,
		onSuccessfulDeliveryFn: onDelivery,
		onAuditFn:              onAudit,
		stopChan:               make(chan struct{}),
	}

	q.loadFromDisk()
	go q.workerLoop()
	return q
}

// Enqueue adds a package dispatch request to the resilient disk-backed queue.
func (q *StoreAndForwardQueue) Enqueue(studyID, packageID, recipient, channel, payloadChecksum string) *QueueItem {
	q.mu.Lock()
	defer q.mu.Unlock()

	randBytes := make([]byte, 6)
	_, _ = rand.Read(randBytes)
	itemID := fmt.Sprintf("QITEM-%d", time.Now().UnixNano()/1e6)
	trackingToken := fmt.Sprintf("TRK-%s", hex.EncodeToString(randBytes))

	item := &QueueItem{
		ID:                     itemID,
		StudyID:                studyID,
		PackageID:              packageID,
		Recipient:              recipient,
		Channel:                channel,
		TrackingToken:          trackingToken,
		PayloadChecksum:        payloadChecksum,
		Status:                 StatusQueued,
		Attempts:               0,
		MaxAttempts:            5,
		EnqueuedAt:             time.Now(),
		NextAttemptAt:          time.Now(), // Immediate first attempt
		BackoffIntervalSeconds: 3,          // Base 3 seconds backoff
	}

	q.items[itemID] = item
	q.saveToDiskLocked()

	if q.onAuditFn != nil {
		q.onAuditFn("delivery_enqueued_store_and_forward", studyID, "STORE_FORWARD_QUEUE",
			fmt.Sprintf("Paquete encolado en disco resiliente. Token: %s, Canal: %s a %s", trackingToken, channel, recipient))
	}

	log.Printf("[STORE-FORWARD] Enqueued delivery %s (Study: %s, Token: %s)", itemID, studyID, trackingToken)
	return item
}

// GetStatus returns queue summary statistics.
func (q *StoreAndForwardQueue) GetStatus() map[string]interface{} {
	q.mu.RLock()
	defer q.mu.RUnlock()

	queuedCount := 0
	retryingCount := 0
	deliveredCount := 0
	failedCount := 0

	var itemList []*QueueItem
	for _, item := range q.items {
		itemList = append(itemList, item)
		switch item.Status {
		case StatusQueued:
			queuedCount++
		case StatusRetrying:
			retryingCount++
		case StatusDelivered:
			deliveredCount++
		case StatusPermanentFail:
			failedCount++
		}
	}

	return map[string]interface{}{
		"totalItems":          len(q.items),
		"pendingQueued":       queuedCount,
		"retrying":            retryingCount,
		"delivered":           deliveredCount,
		"permanentFail":       failedCount,
		"isNetworkOnline":     q.isNetworkOnline && !q.simulatedNetworkDrop,
		"simulatedDropActive": q.simulatedNetworkDrop,
		"items":               itemList,
	}
}

// SetNetworkSimulation toggles simulated network disruption for resilience tests.
func (q *StoreAndForwardQueue) SetNetworkSimulation(dropActive bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.simulatedNetworkDrop = dropActive
	statusText := "ONLINE"
	if dropActive {
		statusText = "OFFLINE (SIMULATED DROP)"
	}
	log.Printf("[STORE-FORWARD] Network simulation set to: %s", statusText)
	if q.onAuditFn != nil {
		q.onAuditFn("network_simulation_toggled", "GLOBAL", "ADMIN_TOOL",
			fmt.Sprintf("Estado de red simulado cambiado a: %s", statusText))
	}
}

// RetryAll forces immediate re-attempt of all pending/retrying items.
func (q *StoreAndForwardQueue) RetryAll() int {
	q.mu.Lock()
	count := 0
	now := time.Now()
	for _, item := range q.items {
		if item.Status == StatusQueued || item.Status == StatusRetrying {
			item.NextAttemptAt = now
			count++
		}
	}
	q.saveToDiskLocked()
	q.mu.Unlock()

	log.Printf("[STORE-FORWARD] Forced retry on %d pending queue items", count)
	if count > 0 {
		go q.processPendingItems()
	}
	return count
}

// ProcessNow immediately evaluates and dispatches any eligible items.
func (q *StoreAndForwardQueue) ProcessNow() {
	q.processPendingItems()
}

func (q *StoreAndForwardQueue) workerLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-q.stopChan:
			return
		case <-ticker.C:
			q.processPendingItems()
		}
	}
}

func (q *StoreAndForwardQueue) processPendingItems() {
	q.mu.Lock()
	now := time.Now()
	var toProcess []*QueueItem

	for _, item := range q.items {
		if (item.Status == StatusQueued || item.Status == StatusRetrying) && now.After(item.NextAttemptAt) {
			toProcess = append(toProcess, item)
		}
	}
	q.mu.Unlock()

	for _, item := range toProcess {
		q.attemptDelivery(item)
	}
}

func (q *StoreAndForwardQueue) attemptDelivery(item *QueueItem) {
	q.mu.Lock()
	now := time.Now()
	item.Attempts++
	item.LastAttemptAt = &now

	// Check if network is simulated dropped
	if q.simulatedNetworkDrop || !q.isNetworkOnline {
		item.Status = StatusRetrying
		item.LastError = "Conexión a internet no disponible: servidor remoto inalcanzable (simulación de corte de red)."
		// Exponential backoff: 3s, 6s, 12s, 24s...
		backoff := item.BackoffIntervalSeconds * (1 << (item.Attempts - 1))
		if backoff > 60 {
			backoff = 60
		}
		item.NextAttemptAt = now.Add(time.Duration(backoff) * time.Second)

		if item.Attempts >= item.MaxAttempts {
			item.Status = StatusPermanentFail
			item.LastError = fmt.Sprintf("Máximo número de reintentos alcanzado (%d). Retenido en disco para intervención técnica.", item.MaxAttempts)
			log.Printf("[STORE-FORWARD] Item %s reached max attempts (%d). PERMANENT_FAIL.", item.ID, item.MaxAttempts)
		} else {
			log.Printf("[STORE-FORWARD] Network drop detected. Item %s set to RETRYING in %ds (Attempt %d/%d)",
				item.ID, backoff, item.Attempts, item.MaxAttempts)
		}

		q.saveToDiskLocked()
		q.mu.Unlock()
		return
	}

	// Successful delivery dispatch simulation
	item.Status = StatusDelivered
	item.DeliveredAt = &now
	item.LastError = ""
	q.saveToDiskLocked()
	q.mu.Unlock()

	log.Printf("[STORE-FORWARD] Delivery %s successfully dispatched to %s via %s (Token: %s)",
		item.ID, item.Recipient, item.Channel, item.TrackingToken)

	if q.onSuccessfulDeliveryFn != nil {
		q.onSuccessfulDeliveryFn(item)
	}

	if q.onAuditFn != nil {
		q.onAuditFn("delivery_dispatched_via_store_and_forward", item.StudyID, "STORE_FORWARD_DISPATCHER",
			fmt.Sprintf("Resultado entregado exitosamente tras %d intentos a %s (%s). Token: %s",
				item.Attempts, item.Recipient, item.Channel, item.TrackingToken))
	}
}

func (q *StoreAndForwardQueue) saveToDiskLocked() {
	data, err := json.MarshalIndent(q.items, "", "  ")
	if err != nil {
		log.Printf("[STORE-FORWARD] Error marshaling queue: %v", err)
		return
	}
	_ = os.WriteFile(q.filePath, data, 0644)
}

func (q *StoreAndForwardQueue) loadFromDisk() {
	data, err := os.ReadFile(q.filePath)
	if err != nil {
		return // File doesn't exist yet, empty queue
	}
	var loaded map[string]*QueueItem
	if err := json.Unmarshal(data, &loaded); err == nil {
		q.items = loaded
		log.Printf("[STORE-FORWARD] Loaded %d queue items from disk (%s)", len(loaded), q.filePath)
	}
}

// Close stops the background worker loop.
func (q *StoreAndForwardQueue) Close() {
	close(q.stopChan)
}

// ConvertToDeliveryRecord maps a QueueItem to models.DeliveryRecord
func (item *QueueItem) ToDeliveryRecord() models.DeliveryRecord {
	status := models.DeliverySent
	if item.Status == StatusPermanentFail {
		status = models.DeliveryFailed
	}
	return models.DeliveryRecord{
		ID:            item.ID,
		PackageID:     item.PackageID,
		StudyID:       item.StudyID,
		Recipient:     item.Recipient,
		Channel:       item.Channel,
		Status:        status,
		SentAt:        item.EnqueuedAt,
		TrackingToken: item.TrackingToken,
	}
}
