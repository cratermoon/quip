package signing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var pss_testKey = []byte(`-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBALgASs8C3PwCZ0Ad
FuwWKqmO11RDRvZksDGRrGoJ0tJH++ZnXJUVaCl6K3/3DhlGL1KF1l/uefhbydJy
KZZTnJ2GgxfnqYVz5sgvScT695AMA6j9xoXEc9x6eKh7wvnhDyEAILgIQ31wNyU0
y5t8ebXLHLUvKrXhDplO+bOWUA1rAgMBAAECgYBXpEv9rm9z2fE0KKbybNyFRvRp
vmHemrcR5UEqgONaJc9mP3VpzXh/ySFPIm4ku7lupTEnIIAYPCA1jQsh/1FpJ40z
Jmjq0/PE/Wkel2MXj3vpJvJ1Uo4UTZh/fpBc4iV5r1rX+ZmHJmcKUIKNV6BpbRCG
2TTJSHoe0vBleWd8gQJBAOlMd1YIis8YGkvRy4vTSvUlTt8GpHx8dzdGkl1XsvYS
UGYyYfcY7xe5+epBZWV7lj6D0F1/q8blDjgvnm/o8EUCQQDJ58+DnwAuM2D8rWlb
LugrDBM9nvnG9FJFlXCZXlL4FvX1OCZmT/IzXritw5HbTvFh3u4vyO+eh7IVe89D
bRnvAkAGCZSNBWuSMG19yNAbrjwiW/TOkL1w+0eahpkDEWHwPEkYW/VtD5ggZQ+y
oD6fgbNBqueZg9ROMV9M1O6ktsKBAkBdRufasek+buQet+qVgp0lzgVRkZFpddRQ
a1LCuA3yqYDl0hQDbmnBi8AcHt7Sh60CfyBhGR6CicQfIrzFNLEnAkBew+KSiQ9+
4VHp5CjC0sipvN9idvbnqZq+Ul7pwXv/WexSML/MQR4fZ7d47bC8CInEFeojoHqF
UU2YhLkn+FPr
-----END PRIVATE KEY-----`)

func TestPSSSign(t *testing.T) {
	text := "this is a test"

	signer := PSSSigner{pss_testKey}

	sigHex, err := signer.Sign(text)
	assert.NoError(t, err, "Signing failed")
	assert.Equal(t, 256, len(sigHex), "Wrong signature")
}
