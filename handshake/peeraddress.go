package handshake

import "encoding/base64"

// PeerAddr wrapps peer's publick key and NetworkID
type PeerAddr struct {
	PublicKey []byte
}

func (pa PeerAddr) NetworkID() string {
	return NetworkString
}

func (pa PeerAddr) String() string {
	return "@" + base64.StdEncoding.EncodeToString(pa.PublicKey) + ".ed25519"
}
