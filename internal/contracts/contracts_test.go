package contracts

import (
	"rx-dispatch/internal/models"
	"testing"
)

// TestPackageCompilation es una prueba mínima para asegurar que el paquete compila
// y para que las herramientas de testing lo registren. Este paquete solo contiene
// Data Transfer Objects (DTOs), cuya validación real ocurre en las pruebas
// de los servicios que los utilizan.
func TestPackageCompilation(t *testing.T) {
	_ = HealthResponse{}
	_ = RequestGenericReadingRequest{}
	_ = RecordAuditEventRequest{Event: models.AuditEvent{}}
}
