package models

import "testing"

func TestHashPayload(t *testing.T) {
	t.Run("should return a consistent SHA-256 hash", func(t *testing.T) {
		data := []byte("hello world")
		expectedHash := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"

		hash1 := HashPayload(data)
		hash2 := HashPayload(data)

		if hash1 != expectedHash {
			t.Errorf("HashPayload() got = %v, want %v", hash1, expectedHash)
		}
		if hash1 != hash2 {
			t.Errorf("HashPayload() is not deterministic, got %v and %v", hash1, hash2)
		}
	})
}
