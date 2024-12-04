package cyberlab

import (
	"fmt"
	"math/big"

	"github.com/janstoon/toolbox/tricks/mathx"
)

type DiffieHellmanKeyExchange struct {
	modulus *big.Int // p
	base    *big.Int // g
}

// NewDiffieHellmanKeyExchange initiates a diffie-hellman key exchange session
func NewDiffieHellmanKeyExchange(modulus, base int) DiffieHellmanKeyExchange {
	if !mathx.IsPrimitiveRoot(modulus, base) {
		panic(fmt.Errorf("expected base (%d) to be primitive root modulo of modulus (%d)", base, modulus))
	}

	return DiffieHellmanKeyExchange{
		modulus: big.NewInt(int64(modulus)),
		base:    big.NewInt(int64(base)),
	}
}

// PublicKey computes the public key using: base^secret % modulus
func (ke DiffieHellmanKeyExchange) PublicKey(secret int) int {
	return int(big.NewInt(0).Exp(ke.base, big.NewInt(int64(secret)), ke.modulus).Int64())
}

// SharedSecret computes the secret key using other party's public key and self secret key:
// counterpartPubKey^secret % modulus
func (ke DiffieHellmanKeyExchange) SharedSecret(counterpartPubKey, secret int) int {
	return int(big.NewInt(0).Exp(big.NewInt(int64(counterpartPubKey)), big.NewInt(int64(secret)), ke.modulus).Int64())
}

// DiffieHellmanPublicChannel exposes the diffie-hellman publicly-exchanged parameters
type DiffieHellmanPublicChannel struct {
	Modulus int // Key exchange session agreed modulus
	Base    int // Key exchange session agreed base
}

// PublicChannel returns public-channel of a diffie-hellman key exchange session
func (ke DiffieHellmanKeyExchange) PublicChannel() DiffieHellmanPublicChannel {
	return DiffieHellmanPublicChannel{
		Modulus: int(ke.modulus.Int64()),
		Base:    int(ke.base.Int64()),
	}
}

// GuessDiffieHellmanSharedSecret tries to guess the shared secret using publicly exchanged parameters
func (pc DiffieHellmanPublicChannel) GuessDiffieHellmanSharedSecret(pubKeyA, pubKeyB int) int {
	prr := mathx.PrimitiveRootsWithRrs(pc.Modulus)
	rrs, ok := prr[pc.Base]
	if !ok {
		return 0
	}

	for index, remainder := range rrs {
		if remainder == pubKeyB {
			sharedSecret := 1
			for i := 0; i <= index; i++ {
				sharedSecret = (sharedSecret * pubKeyA) % pc.Modulus
			}

			return sharedSecret
		}
	}

	return 0
}
