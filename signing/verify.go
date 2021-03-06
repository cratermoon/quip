package signing

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
)

type Verifier struct {
	Cert []byte
}

func (v Verifier) Verify(text, signatureHex string) error {

	signature := make([]byte, hex.DecodedLen(len(signatureHex)))
	_, err := hex.Decode(signature, []byte(signatureHex))
	if err != nil {
		return err
	}

	sum := sha256.Sum256([]byte(text))

	rsaPublicKey, err := parsePublicKeyCert(v.Cert)

	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, sum[:], signature)
}
