package signing

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func parseKey(pemKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemKey)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func ReadKey(keyFile string) (*rsa.PrivateKey, error) {

	prvKeyPem, err := ioutil.ReadFile(keyFile)

	if err != nil {
		return nil, err
	}
	return parseKey(prvKeyPem)
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

func readPublicKeyCert(certFile string) (*rsa.PublicKey, error) {

	publicKeyCert, err := ioutil.ReadFile(certFile)

	if err != nil {
		return nil, err
	}

	return parsePublicKeyCert(publicKeyCert)
}
