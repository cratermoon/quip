package signing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const pss_testKey = `
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDNu4fdJJ8FDlpdad/I78GqCrtncAW6IDOHhxLicZJAbiyUsg9Q
y+1p3+0KVFGSWo3ouBgIGH7wJaln+AwqYtIIAJOmNsbYs3OiZsd57ix5pyNgD8We
CmATG76juSNIr+EZ6b/ZCkpMwVfaYF7rgBbo6w17mk7KxZ5/08uwkN7OiQIDAQAB
AoGAMfAnVoKhJvUI8kvUyk2IUOIyAzWp7jgKygb7ioPO4Fzd7WobVZ1qj5yPSUaW
VlQsxiSJkM2BYoGB0X7gVRmPqXotDtPSu5A+CMPWor144dwQdtFG662wvdwcQ6xl
Q7OLiaL53U6QKbi28MZ/suIE7xzAGu0xvqxjCC66qau7UZECQQDpHVTne3iw2PDC
ZlgwWGYZ/Gmp7P6S0JtQ/jgAJot9q8ckAM7P8oCd/UYe359AVVbEEZ60Yv9NeWzU
kknId0fbAkEA4e4Hp/irVvfBxNTkZSHpHcElnhPkdkQpgRP6CMqHXqLMQ0apy2jZ
1bZ5xszXdDiCOcrh5ekHbY+IgivxE2EyawJBAJOav/0GwGfyJZhiJ2sNPjEsE1fG
OXeK4R2KqrjlryN21lRkso8XNPtUuMapv/ODVbo2kfAUUyWiQhfjPRbS+EsCQDJW
oKoSQ8rKxQegD4tg9NnGUSVZdUvMgBrcYpdW2LaDO1O6CNbjc7WkRJnAxjiE5q8N
vytEsnz8wAOQ2tPgkiUCQQDl5MbmhGCKt+2HiSd/Bvi0BhRL0/g051S6mNWMhjg/
Hs7kq+w14RMwFC6F3HBf/rYKsmiU7epy6IAcY6G56b1q
-----END RSA PRIVATE KEY-----`

func TestPSSSign(t *testing.T) {
	text := "this is a test"

	signer := PSSSigner{[]byte(pss_testKey)}

	sigHex, err := signer.Sign(text)
	t.Logf("%s\n", sigHex)
	assert.NoError(t, err, "Signing failed")
	assert.Equal(t, 256, len(sigHex), "Wrong signature")
}
