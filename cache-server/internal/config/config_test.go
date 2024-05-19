package config_test

import (
	"cache-server/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGet(t *testing.T) {
	mockLogger := zap.NewNop().Sugar()
	config, err := config.NewConfig(mockLogger)

	if err != nil {
		t.Fatalf("Error creating config: %v", err)
	}

	var tests = map[string]struct {
		name     string
		key      string
		expected string
	}{
		"NON_EXISTENT_KEY": {
			key:      "ABC",
			expected: "",
		},
		"EXISTENT_KEY": {
			key:      "DATABASE_SERVICE_BASE_URL",
			expected: "http://localhost:8081",
		},
	}

	for tc, test := range tests {
		t.Run(tc, func(t *testing.T) {
			actual := config.Get(test.key)
			assert.Equal(t, test.expected, actual, "Expected and actual values do not match")
		})
	}
}
