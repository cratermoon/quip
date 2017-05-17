package signing

import (
	"crypto"
	"crypto/rsa"
	"encoding/hex"
)

type PSSVerifier struct {
	Cert []byte
}

func (v PSSVerifier) Verify(message, signatureHex string) error {
	signature := make([]byte, hex.DecodedLen(len(signatureHex)))
	_, err := hex.Decode(signature, []byte(signatureHex))
	if err != nil {
		return err
	}

	h := crypto.SHA512.New()
	h.Write([]byte(message))
	hashed := h.Sum(nil)

	rsaPublicKey, err := parsePublicKey(v.Cert)
	if err != nil {
		return err
	}

	return rsa.VerifyPSS(rsaPublicKey, crypto.SHA512, hashed, signature, nil)
}
