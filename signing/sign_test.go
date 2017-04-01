package signing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSing(t *testing.T) {
	text := "this is a test"
	expectedSigHex := "63d89fc148877565183788fa97c69db13232481cecec8e441820f3c52ab7e7ac42a8c2398dd4824103e1c49f32fe69410e7993da319bd4a797545417a26bd5b7f0f9093d305d15cff440da2e2c09945ab96a4fedc6fa1400f793ad4b3ffa1c076e4de358651fc28018ca7e95b653d71c74886433ccd65b3890c175e50ebc6e5f"

	sigHex, err := Sign(text)
	assert.NoError(t, err, "Signing failed")
	assert.Equal(t, expectedSigHex, sigHex, "Wrong signature")
}
