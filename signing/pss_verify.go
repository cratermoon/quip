package signing

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/hex"
)

type PSSVerifier struct {
	Cert []byte
}

func (v PSSVerifier) Verify(text, signatureHex string) error {

	signature := make([]byte, hex.DecodedLen(len(signatureHex)))
	_, err := hex.Decode(signature, []byte(signatureHex))
	if err != nil {
		return err
	}

	sum := sha512.Sum512([]byte(text))

	rsaPublicKey, err := parsePublicKeyCert(v.Cert)

	if err != nil {
		return err
	}

	return rsa.VerifyPSS(rsaPublicKey, crypto.SHA512, sum[:], signature, nil)
}
