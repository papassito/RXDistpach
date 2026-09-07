package tests

import (
	"os"
	"path/filepath"
	"testing"

	"rx-dispatch/internal/delivery"
)

func TestStoreAndForwardResilientQueue(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "rx-queue-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	var deliveredItem *delivery.QueueItem
	var auditEvents []string

	queue := delivery.NewStoreAndForwardQueue(
		filepath.Join(tempDir, "data"),
		func(item *delivery.QueueItem) {
			deliveredItem = item
		},
		func(eventType, studyID, actor, details string) {
			auditEvents = append(auditEvents, eventType)
		},
	)
	defer queue.Close()

	// 1. Simulate clinical internet drop
	queue.SetNetworkSimulation(true)

	// 2. Enqueue study delivery package
	item := queue.Enqueue(
		"STU-RESILIENCE-001",
		"PKG-001",
		"paciente.resiliente@hospital.com",
		"EMAIL",
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	)

	if item.Status != delivery.StatusQueued {
		t.Fatalf("expected initial status QUEUED, got %s", item.Status)
	}

	// 3. Force delivery attempt during simulated drop
	queue.ProcessNow()

	statusAfterDrop := queue.GetStatus()
	retryingCount := statusAfterDrop["retrying"].(int)
	if retryingCount < 1 {
		t.Fatalf("expected item to be in RETRYING state during network drop, got retrying=%d", retryingCount)
	}

	// 4. Restore internet connectivity
	queue.SetNetworkSimulation(false)

	// 5. Force immediate retry
	queue.RetryAll()
	queue.ProcessNow()

	statusRestored := queue.GetStatus()
	deliveredCount := statusRestored["delivered"].(int)
	if deliveredCount != 1 {
		t.Fatalf("expected 1 delivered item after network restored, got %d", deliveredCount)
	}

	if deliveredItem == nil || deliveredItem.Status != delivery.StatusDelivered {
		t.Fatalf("expected deliveredItem to be DELIVERED, got %+v", deliveredItem)
	}

	// 6. Verify tracking token and checksum preservation
	if deliveredItem.TrackingToken == "" {
		t.Fatalf("expected non-empty tracking token")
	}
	if deliveredItem.PayloadChecksum != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		t.Fatalf("payload checksum corrupted during queue transit")
	}
}
