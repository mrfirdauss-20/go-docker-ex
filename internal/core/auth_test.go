package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateAPIKeyMatch(t *testing.T) {
	exampleAPIKey := "c4211664-47dc-4887-a2fe-9e694fbaf55a"

	// initialize new auth
	auth, err := NewAuth(AuthConfig{
		APIKey: "c4211664-47dc-4887-a2fe-9e694fbaf55a",
	})
	require.NoError(t, err)
	// test the validate api key
	err = auth.ValidateAPIKey(context.Background(), exampleAPIKey)
	// api key is validated when there is no error
	assert.NoError(t, err, "api keys are mismatch")
}

func TestValidateAPIKeyMismatch(t *testing.T) {
	exampleAPIKey := "Ini Contoh API Key"

	// initialize new auth
	auth, err := NewAuth(AuthConfig{
		APIKey: "c4211664-47dc-4887-a2fe-9e694fbaf55a",
	})
	require.NoError(t, err)
	// test the validate api key
	err = auth.ValidateAPIKey(context.Background(), exampleAPIKey)
	// it should return error when the key sent and auth key are mismatch
	assert.Error(t, err, "api keys are match")
	assert.Equal(t, ErrInvalidAPIKey, err, "mismatch error type")
}
