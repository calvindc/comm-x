package handshake

import (
	"reflect"
	"testing"

	"golang.org/x/crypto/ed25519"
)

func TestNewEd25519KeyPair(t *testing.T) {
	type args struct {
		secret []byte
		public []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Ed25519KeyPair
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEd25519KeyPair(tt.args.secret, tt.args.public)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEd25519KeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEd25519KeyPair() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkEd25519KeySize(t *testing.T) {
	type args struct {
		secretKey []byte
		publicKey []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkEd25519KeySize(tt.args.secretKey, tt.args.publicKey); (err != nil) != tt.wantErr {
				t.Errorf("checkEd25519KeySize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newPeerState(t *testing.T) {
	type args struct {
		networkIdentifier []byte
		local             Ed25519KeyPair
	}
	tests := []struct {
		name    string
		args    args
		want    *PeerState
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newPeerState(tt.args.networkIdentifier, tt.args.local)
			if (err != nil) != tt.wantErr {
				t.Errorf("newPeerState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPeerState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClientPeerState(t *testing.T) {
	type args struct {
		networkIdentifier []byte
		local             Ed25519KeyPair
		remotePublic      ed25519.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    *PeerState
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClientPeerState(tt.args.networkIdentifier, tt.args.local, tt.args.remotePublic)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientPeerState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientPeerState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewServerPeerState(t *testing.T) {
	type args struct {
		networkIdentifier []byte
		local             Ed25519KeyPair
	}
	tests := []struct {
		name    string
		args    args
		want    *PeerState
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServerPeerState(tt.args.networkIdentifier, tt.args.local)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewServerPeerState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServerPeerState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerState_createChallenge(t *testing.T) {
	tests := []struct {
		name string
		prs  *PeerState
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.prs.createChallenge(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerState.createChallenge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerState_verifyChallenge(t *testing.T) {
	type args struct {
		ch []byte
	}
	tests := []struct {
		name string
		prs  *PeerState
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.prs.verifyChallenge(tt.args.ch); got != tt.want {
				t.Errorf("PeerState.verifyChallenge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerState_createClientAuth(t *testing.T) {
	tests := []struct {
		name    string
		prs     *PeerState
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.prs.createClientAuth()
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerState.createClientAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerState.createClientAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerState_verifyClientAuth(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		prs  *PeerState
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.prs.verifyClientAuth(tt.args.data); got != tt.want {
				t.Errorf("PeerState.verifyClientAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerState_createServerAccept(t *testing.T) {
	tests := []struct {
		name    string
		prs     *PeerState
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.prs.createServerAccept()
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerState.createServerAccept() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerState.createServerAccept() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerState_verifyServerAccept(t *testing.T) {
	type args struct {
		boxedOkay []byte
	}
	tests := []struct {
		name string
		prs  *PeerState
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.prs.verifyServerAccept(tt.args.boxedOkay); got != tt.want {
				t.Errorf("PeerState.verifyServerAccept() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerState_overwriteSecrets(t *testing.T) {
	tests := []struct {
		name string
		prs  *PeerState
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prs.overwriteSecrets()
		})
	}
}

func TestPeerState_Remote(t *testing.T) {
	tests := []struct {
		name string
		prs  *PeerState
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.prs.Remote(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerState.Remote() = %v, want %v", got, tt.want)
			}
		})
	}
}
