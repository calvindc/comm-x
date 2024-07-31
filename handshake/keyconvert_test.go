package handshake

import (
	"testing"

	"golang.org/x/crypto/ed25519"
)

func TestPublicKeyToCurve25519(t *testing.T) {
	type args struct {
		curveBytes *[32]byte
		edBytes    ed25519.PublicKey
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PublicKeyToCurve25519(tt.args.curveBytes, tt.args.edBytes); got != tt.want {
				t.Errorf("PublicKeyToCurve25519() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKeyToCurve25519(t *testing.T) {
	type args struct {
		curve25519PrivateKey *[32]byte
		ed25519PrivateKey    ed25519.PrivateKey
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrivateKeyToCurve25519(tt.args.curve25519PrivateKey, tt.args.ed25519PrivateKey)
		})
	}
}
