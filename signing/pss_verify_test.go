package signing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var pss_testPubKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC4AErPAtz8AmdAHRbsFiqpjtdU
Q0b2ZLAxkaxqCdLSR/vmZ1yVFWgpeit/9w4ZRi9ShdZf7nn4W8nScimWU5ydhoMX
56mFc+bIL0nE+veQDAOo/caFxHPcenioe8L54Q8hACC4CEN9cDclNMubfHm1yxy1
Lyq14Q6ZTvmzllANawIDAQAB
-----END PUBLIC KEY-----`)

var pv PSSVerifier

func TestPSSVerifyFail(t *testing.T) {
	text := "this is a test"
	expectedSigHex := "b072de1ae1a54117d2003519832c1053fbd88d3a978bc0b6088ff677991faa1947ad77a60a253585aabb5f818c9261318fc0fb14832597827ccb54ebbbd0d0142e0b9574f352d4417fb32ab4c792cfe3f571b7968c346f9070eb04cfac8fad7dfbbcc671ee8c7db0d04bdca0e8a40cc44bedcaf0bbceb8458a5218975262079d"

	assert.Error(t, v.Verify(text, expectedSigHex), "Signature verification should fail")
}

func TestPSSVerify(t *testing.T) {
	text := "this is a test"
	expectedSigHex := "4a90ac31ff9195072a0ff9302194c27e44f04128a4cd786630c36e72105c5f39fa00ac7595befbe48698e3f73922814829133366e718a66e2236601af608e6988de8379a801aefa91cbcfa12b4c1f1d0bdfc99e01e66a3da49f62433f0568edd1bc7a2d0c33af00d38cbf55c2418a76583f3100a549f1c050481351746123dce"
	resp := pv.Verify(text, expectedSigHex)
	assert.NoError(t, resp, "Signature verification failed")
}

func init() {
	pv = PSSVerifier{pss_testPubKey}
}
