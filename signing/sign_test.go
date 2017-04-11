package signing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIBPQIBAAJBALV8UzmEf92NxSt3KBjl5Otvj69zz53ilS9PwEusXwOGH5q324/C
I1WZqGfWq0/PKy6oeKyBB9wK3Q9dSTVV4BMCAwEAAQJBAIu3A4cTJzDUFIeKuxa3
3U0W3KHw4VOl/L2ogtx+/cKCcZW4QUNqwULT7ZBEms4hGfJ3yeDkXHZaKpBagSLj
dykCIQDhr5mExTTm2AxJV6BbP6ic5GkXDhBym0B1u24Hk1tYXQIhAM3c3H5Rz1zY
NX4Dj87f5Fl4uNJGLQKn2zBeUNKI7NMvAiEAw/ehm3tOK2DQkmLnWDSXqdxgMGfC
+nE68MAWk7dtqvUCIQCzbTodO27yLFxLTg1ssUYVGZx1YcbfVrA7syjcp41K7wIh
AI6Ncz4fq5BBAIugQ/JhUYutJLdB64nwZT/f+cyGmJY2
-----END RSA PRIVATE KEY-----`

func TestSing(t *testing.T) {
	text := "this is a test"
	expectedSigHex := "38b0dac7c1308cb876476b666758b724b1531ed2f5c565254e7f35d4b54dedca16ff7f6caadb32c6ae5c51cb9b120e11d3da3abd9e7eb23c107a1f984110c13d"

	signer := Signer{[]byte(testKey)}

	sigHex, err := signer.Sign(text)
	assert.NoError(t, err, "Signing failed")
	assert.Equal(t, expectedSigHex, sigHex, "Wrong signature")
}
