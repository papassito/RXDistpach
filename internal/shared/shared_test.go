package shared

import "testing"

// TestConstants es una prueba mínima para asegurar que las constantes
// normativas del sistema están definidas y no están vacías.
func TestConstants(t *testing.T) {
	if MandatoryLegalDisclaimer == "" {
		t.Error("La constante MandatoryLegalDisclaimer no puede estar vacía.")
	}
}
