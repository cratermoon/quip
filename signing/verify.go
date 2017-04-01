package signing

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
)

func Verify(text, signatureHex string) error {

	signature := make([]byte, hex.DecodedLen(len(signatureHex)))
	_, err := hex.Decode(signature, []byte(signatureHex))
	if err != nil {
		return err
	}

	sum := sha256.Sum256([]byte(text))

	rsaPrivateKey, err := readKey()

	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(&rsaPrivateKey.PublicKey, crypto.SHA256, sum[:], signature)
}
