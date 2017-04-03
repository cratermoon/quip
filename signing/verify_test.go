package signing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var v = Verifier{"quip.crt"}

func TestVerifyFail(t *testing.T) {
	text := "this is a test"
	signatureHex := "7c87a24a7afb97529c5e6f43211cf6430e9503ea6d14fb2652be8ac2ba8996f2b8b4e4a10745a458afdde43968b344ff915e6a883369a33bb7ea1aa1cbc13be4fca052757725a6c72b5f1541f8c0ed269f5772327e470a74fd6c51bca5f90fd3ddc894a1dcebeb5e7e14d3c287c8073632d81c0fe326ed93ea1937507a1cc318"

	assert.Error(t, v.Verify(text, signatureHex), "Signature verification should fail")
}

func TestVerify(t *testing.T) {
	text := "this is a test"
	signatureHex := "63d89fc148877565183788fa97c69db13232481cecec8e441820f3c52ab7e7ac42a8c2398dd4824103e1c49f32fe69410e7993da319bd4a797545417a26bd5b7f0f9093d305d15cff440da2e2c09945ab96a4fedc6fa1400f793ad4b3ffa1c076e4de358651fc28018ca7e95b653d71c74886433ccd65b3890c175e50ebc6e5f"

	assert.NoError(t, v.Verify(text, signatureHex), "Signature verification failed")
}
