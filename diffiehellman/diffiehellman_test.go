package diffiehellman_test

import (
	"math/rand"
	"testing"

	"github.com/janstoon/toolbox/tricks/mathx"
	"github.com/stretchr/testify/assert"

	"github.com/pouyanh/cyberlab/diffiehellman"
)

func TestDiscoverSharedSecret(t *testing.T) {
	sessions := []struct {
		modulus int
		base    int
	}{
		{11, 6},
		{5, 0},
		{5, 2},
		{23, 0},
		{23, 5},
		{107, 63},
	}

	privateKeyPairs := []struct {
		alice int
		bob   int
	}{
		{5, 11}, {5, 12}, {5, 23}, {5, 10},
		{6, 11}, {6, 6},
		{9, 6}, {9, 5},
		{7, 5}, {7, 6}, {7, 8},
		{11, 5},
		{110, 49}, {110, 50}, {110, 51},
	}

	for _, session := range sessions {
		if session.base == 0 {
			prr := mathx.PrimitiveRoots(session.modulus)
			session.base = prr[rand.Intn(len(prr))]
		}
		ke := diffiehellman.NewKeyExchange(session.modulus, session.base)

		pc := ke.PublicChannel()
		for _, privateKeys := range privateKeyPairs {
			AlicePublicKey := ke.PublicKey(privateKeys.alice)
			BobPublicKey := ke.PublicKey(privateKeys.bob)

			aliceSharedSecret := ke.SharedSecret(BobPublicKey, privateKeys.alice)
			bobSharedSecret := ke.SharedSecret(AlicePublicKey, privateKeys.bob)
			assert.Equal(t, aliceSharedSecret, bobSharedSecret)

			assert.Equal(t, aliceSharedSecret, pc.DiscoverSharedSecret(AlicePublicKey, BobPublicKey))
			assert.Equal(t, bobSharedSecret, pc.DiscoverSharedSecret(BobPublicKey, AlicePublicKey))
			assert.Equal(t, bobSharedSecret, pc.WithPublicKeySelector(diffiehellman.NearestPublicKeySelector).
				DiscoverSharedSecret(BobPublicKey, AlicePublicKey))
		}
	}
}
