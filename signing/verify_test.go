package signing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testCert = 
`-----BEGIN CERTIFICATE-----
MIIDFzCCAsGgAwIBAgIJAKMvcSar/SX5MA0GCSqGSIb3DQEBCwUAMIGRMQswCQYD
VQQGEwJVUzEPMA0GA1UECBMGT3JlZ29uMREwDwYDVQQHEwhQb3J0bGFuZDEgMB4G
A1UEChMXQ3JhdGVyIE1vb24gRGV2ZWxvcG1lbnQxEjAQBgNVBAMTCWNtZGV2LmNv
bTEoMCYGCSqGSIb3DQEJARYZc3RldmVuLmUubmV3dG9uQGdtYWlsLmNvbTAeFw0x
NzA0MDMxNzQ0MTdaFw0xNzA1MDMxNzQ0MTdaMIGRMQswCQYDVQQGEwJVUzEPMA0G
A1UECBMGT3JlZ29uMREwDwYDVQQHEwhQb3J0bGFuZDEgMB4GA1UEChMXQ3JhdGVy
IE1vb24gRGV2ZWxvcG1lbnQxEjAQBgNVBAMTCWNtZGV2LmNvbTEoMCYGCSqGSIb3
DQEJARYZc3RldmVuLmUubmV3dG9uQGdtYWlsLmNvbTBcMA0GCSqGSIb3DQEBAQUA
A0sAMEgCQQC1fFM5hH/djcUrdygY5eTrb4+vc8+d4pUvT8BLrF8Dhh+at9uPwiNV
mahn1qtPzysuqHisgQfcCt0PXUk1VeATAgMBAAGjgfkwgfYwHQYDVR0OBBYEFPVh
CQ4T/szs+C+IAkGUd7oVMVapMIHGBgNVHSMEgb4wgbuAFPVhCQ4T/szs+C+IAkGU
d7oVMVapoYGXpIGUMIGRMQswCQYDVQQGEwJVUzEPMA0GA1UECBMGT3JlZ29uMREw
DwYDVQQHEwhQb3J0bGFuZDEgMB4GA1UEChMXQ3JhdGVyIE1vb24gRGV2ZWxvcG1l
bnQxEjAQBgNVBAMTCWNtZGV2LmNvbTEoMCYGCSqGSIb3DQEJARYZc3RldmVuLmUu
bmV3dG9uQGdtYWlsLmNvbYIJAKMvcSar/SX5MAwGA1UdEwQFMAMBAf8wDQYJKoZI
hvcNAQELBQADQQB6iV6fDh5R0Bv3tFlWrkO6gxybrGSHkOv8ZoA6pcpmeXwcuZRp
eBfl+Ra8ix8icuPDEtJmCA4bJeqbq5e9tCbD
-----END CERTIFICATE-----`

var v Verifier

func TestVerifyFail(t *testing.T) {
	text := "this is a test"
	expectedSigHex := "38b0dac7c1308cb876476b666758b724b1531ed2f5c565254e7f35d4b54dedca16ff7f6caadb32c6ae5c51cb9b120e11d3da3abd9e7eb23c107a1f984110c13e"

	assert.Error(t, v.Verify(text, expectedSigHex), "Signature verification should fail")
}

func TestVerify(t *testing.T) {
	text := "this is a test"
	expectedSigHex := "38b0dac7c1308cb876476b666758b724b1531ed2f5c565254e7f35d4b54dedca16ff7f6caadb32c6ae5c51cb9b120e11d3da3abd9e7eb23c107a1f984110c13d"

	assert.NoError(t, v.Verify(text, expectedSigHex), "Signature verification failed")
}

func init() {
	v = Verifier{[]byte(testCert)}
}
