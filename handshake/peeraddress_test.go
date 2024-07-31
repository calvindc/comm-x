package handshake

import "testing"

func TestPeerAddr_NetworkID(t *testing.T) {
	tests := []struct {
		name string
		pa   PeerAddr
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pa.NetworkID(); got != tt.want {
				t.Errorf("PeerAddr.NetworkID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerAddr_String(t *testing.T) {
	tests := []struct {
		name string
		pa   PeerAddr
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pa.String(); got != tt.want {
				t.Errorf("PeerAddr.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
