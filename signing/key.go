package signing

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func parseKey(pemKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemKey)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func parsePublicKeyCert(pemCert []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemCert)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert.PublicKey.(*rsa.PublicKey), nil
}

func parseRSAKey(pemKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemKey)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch key.(type) {
	default:
		return nil, errors.New("unknown key type")
	case *rsa.PrivateKey:
		return key.(*rsa.PrivateKey), nil
	}
}

func parsePublicKey(pemPubKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemPubKey)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch pubKey.(type) {
	default:
		return nil, errors.New("unknown key type")
	case *rsa.PublicKey:
		return pubKey.(*rsa.PublicKey), nil
	}
}
