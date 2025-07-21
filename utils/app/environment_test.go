package app

import (
	"testing"
)

func TestGetEnv(t *testing.T) {
	key := "INVALID_KEY"
	fallback := "fallback_value"

	result := GetEnv(key, fallback)
	if result != fallback {
		t.Errorf("Expected '%s' but got '%s'", fallback, result)
	}
}
