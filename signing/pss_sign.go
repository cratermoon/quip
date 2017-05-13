package signing

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
)

type PSSSigner struct {
	Key []byte
}

func (s PSSSigner) Sign(text string) (string, error) {
	rsaPrivateKey, err := parseKey(s.Key)
	if err != nil {
		return "", err
	}

	h := crypto.SHA512.New()
	h.Write([]byte(text))
	hashed := h.Sum(nil)

	sig, err := rsa.SignPSS(rand.Reader, rsaPrivateKey, crypto.SHA512, hashed[:], nil)
	if err != nil {
		return "", err
	}

	dst := make([]byte, hex.EncodedLen(len(sig)))
	hex.Encode(dst, sig)
	buf := bytes.NewBuffer(dst)
	return buf.String(), nil

}
