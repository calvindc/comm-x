package handshake

import (
	"io"

	"crypto/rand"
	"errors"

	"golang.org/x/crypto/ed25519"
)

// ErrProcessing is returned if I/O fails during the handshake
type ErrProcessing struct {
	where string
	cause error
}

// Error() overwrite error.Error()
func (e ErrProcessing) Error() string {
	errStr := "secrethandshake: failed during data transfer of [" + e.where
	errStr += "],occur err = : [" + e.cause.Error() + "]"
	return errStr
}

// ClientHandShake client role use the cryptographic identity when hand shake
func ClientHandShake(peerstate *PeerState, conn io.ReadWriter) error {
	// 1. send challenge
	_, err := conn.Write(peerstate.createChallenge())
	if err != nil {
		return ErrProcessing{where: "sending challenge", cause: err}
	}

	// 2. recv challenge
	chalResp := make([]byte, ChallengeLength)
	_, err = io.ReadFull(conn, chalResp)
	if err != nil {
		return ErrProcessing{where: "receiving challenge", cause: err}
	}

	// 3. verify challenge
	if !peerstate.verifyChallenge(chalResp) {
		return ErrProcessing{where: "verify challenge", cause: errors.New("wrong Network Identifier or wrong protocol")}
	}

	// 4. send authentication vector
	clientAuthBytes, err := peerstate.createClientAuth()
	if err != nil {
		return ErrProcessing{where: "sending client hello(create client auth first)", cause: err}
	}
	_, err = conn.Write(clientAuthBytes)
	if err != nil {
		return ErrProcessing{where: "sending client hello", cause: err}
	}

	// 5. recv authentication vector
	boxedSig := make([]byte, ServerAuthLength)
	_, err = io.ReadFull(conn, boxedSig)
	if err != nil {
		return ErrProcessing{where: "receiving server auth", cause: err}
	}

	// 6. authenticate remote
	if !peerstate.verifyServerAccept(boxedSig) {
		return ErrProcessing{where: "authenticate remote", cause: errors.New("other side not authenticated")}
	}

	// 7. clean all temp secret this turn
	peerstate.overwriteSecrets()

	return nil
}

// Server shakes hands using the cryptographic identity specified in s using conn in the server role
func ServerHandShake(peerstate *PeerState, conn io.ReadWriter) (err error) {
	// 1. recv challenge
	challenge := make([]byte, ChallengeLength)
	_, err = io.ReadFull(conn, challenge)
	if err != nil {
		return ErrProcessing{where: "receiving challenge", cause: err}
	}

	// 2. verify challenge
	if !peerstate.verifyChallenge(challenge) {
		return ErrProcessing{where: "verify challenge", cause: errors.New("wrong Network Identifier or wrong protocol")}
	}

	// 3. send challenge
	_, err = conn.Write(peerstate.createChallenge())
	if err != nil {
		return ErrProcessing{where: "sending challenge", cause: err}
	}

	// 4. recv authentication vector
	hello := make([]byte, ClientAuthLength)
	_, err = io.ReadFull(conn, hello)
	if err != nil {
		return ErrProcessing{where: "receiving client hello", cause: err}
	}

	// 5. authenticate remote
	if !peerstate.verifyClientAuth(hello) {
		return ErrProcessing{where: "authenticate remote", cause: errors.New("other side not authenticated")}
	}

	// 6.accept
	serverAcceptBytes, err := peerstate.createServerAccept()
	if err != nil {
		return ErrProcessing{where: "sending server accept(create server accept first)", cause: err}
	}
	_, err = conn.Write(serverAcceptBytes)
	if err != nil {
		return ErrProcessing{where: "sending server accept", cause: err}
	}

	//7.  clean all temp secret this turn
	peerstate.overwriteSecrets()

	return nil
}

// GenerateEd25519KeyPair generates a ed25519 keyPair using the passed reader
func GenerateEd25519KeyPair(r io.Reader) (*Ed25519KeyPair, error) {
	if r == nil {
		r = rand.Reader
	}
	pubkey, sectkey, err := ed25519.GenerateKey(r)
	if err != nil {
		return nil, err
	}
	if hasSmallOrder(pubkey) {
		//try again
		pubkey, sectkey, err = ed25519.GenerateKey(r)
	}

	return &Ed25519KeyPair{
		Secret: sectkey,
		Public: pubkey,
	}, nil
}
