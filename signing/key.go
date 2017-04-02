package signing

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func readKey() (*rsa.PrivateKey, error) {

	keyFile := "quip.key"

	prvKeyPem, err := ioutil.ReadFile(keyFile)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(prvKeyPem)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func readPublicKeyCert() (*rsa.PublicKey, error) {

	certFile := "quip.crt"

	publicKeyCert, err := ioutil.ReadFile(certFile)

	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(publicKeyCert)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert.PublicKey.(*rsa.PublicKey), nil
}
