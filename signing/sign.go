package signing

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
)

func Sign(text string) (string, error) {
	hashed := sha256.Sum256([]byte(text))

	rsaPrivateKey, err := readKey()

	if err != nil {
		return "", err
	}

	sig, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}

	dst := make([]byte, hex.EncodedLen(len(sig)))
	hex.Encode(dst, sig)
	buf := bytes.NewBuffer(dst)
	return buf.String(), nil

}
