package utils_test

import (
	"testing"
	"time"
	"url-shortner-database/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestKeyGenerationService(t *testing.T) {
	t.Run("Test key generation service", func(t *testing.T) {
		shortUrl := utils.KeyGenerationService("https://www.google.com")
		assert.Equal(t, 7, len(shortUrl), "Key generation service failed")
	})
}

func TestGetExpirationTime(t *testing.T) {
	tests := map[string]struct {
		expiryTime time.Time
	}{
		"Expiry Time in Future": {
			expiryTime: time.Now().AddDate(0, 0, 1),
		},
		"Expiry Time in Past": {
			expiryTime: time.Now().AddDate(0, 0, -1),
		},
	}

	t.Run("Expiry Time in Future", func(t *testing.T) {
		assert.Equal(t, tests["Expiry Time in Future"].expiryTime, utils.GetExpirationTime(tests["Expiry Time in Future"].expiryTime), "GetExpirationTime failed")
	})

	t.Run("Expiry Time in Past", func(t *testing.T) {
		assert.NotEqual(t, tests["Expiry Time in Past"].expiryTime, utils.GetExpirationTime(tests["Expiry Time in Past"].expiryTime), "GetExpirationTime failed")
	})
}
