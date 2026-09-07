package audit

import (
	"fmt"
	"sync"
	"time"

	"rx-dispatch/internal/models"
)

type AuditService struct {
	mu     sync.RWMutex
	events []models.AuditEvent
}

func NewAuditService() *AuditService {
	return &AuditService{
		events: make([]models.AuditEvent, 0),
	}
}

// Log records a compliant operational event without leaking sensitive pixel/PII content
func (s *AuditService) Log(eventType, studyID, actor, details string) models.AuditEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	ev := models.AuditEvent{
		ID:        fmt.Sprintf("AUD-%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		EventType: eventType,
		StudyID:   studyID,
		Actor:     actor,
		Details:   details,
	}
	s.events = append([]models.AuditEvent{ev}, s.events...) // Newest first
	return ev
}

func (s *AuditService) ListByStudy(studyID string) []models.AuditEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.AuditEvent
	for _, e := range s.events {
		if studyID == "" || e.StudyID == studyID {
			result = append(result, e)
		}
	}
	return result
}
