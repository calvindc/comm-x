package handshake

import (
	"crypto/sha512"

	"filippo.io/edwards25519"
	"golang.org/x/crypto/ed25519"
)

// PublicKeyToCurve25519 converts an Ed25519 public key into the curve25519
// public key that would be generated from the same private key.
func PublicKeyToCurve25519(curveBytes *[32]byte, edBytes ed25519.PublicKey) bool {
	if hasSmallOrder(edBytes) {
		return false
	}

	edPoint, err := new(edwards25519.Point).SetBytes(edBytes)
	if err != nil {
		return false
	}

	copy(curveBytes[:], edPoint.BytesMontgomery())
	return true
}

// PrivateKeyToCurve25519 converts an ed25519 private key into a corresponding curve25519 private key
// calculates a private key from a seed. This function is provided for interoperabilitywith RFC 8032.
// RFC 8032's private keys correspond to seeds in this package.
func PrivateKeyToCurve25519(curve25519PrivateKey *[32]byte, ed25519PrivateKey ed25519.PrivateKey) {
	h := sha512.New()
	h.Write(ed25519PrivateKey[:32])
	digest := h.Sum(nil)
	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	copy(curve25519PrivateKey[:], digest)
}
