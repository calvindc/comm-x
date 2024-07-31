package handshake

import (
	"bytes"
	"fmt"

	"crypto/rand"
	"crypto/sha256"

	"errors"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/nacl/auth"
	"golang.org/x/crypto/nacl/box"
)

// Ed25519KeyPair is a keypair implements the Ed25519 signature algorithm. See https://ed25519.cr.yp.to
type Ed25519KeyPair struct {
	Secret ed25519.PrivateKey
	Public ed25519.PublicKey
}

// Curve25519KeyPair
type Curve25519KeyPair struct {
	Secret [32]byte //different whit ed25519
	Public [32]byte
}

// PeerState define the state each peer holds during the handshark
type PeerState struct {
	networkIdentifier        [32]byte
	secretHash               []byte
	localAppMac              [32]byte
	remoteAppMac             []byte
	localExchange            Curve25519KeyPair
	local                    Ed25519KeyPair
	remoteExchange           Curve25519KeyPair
	remotePublic             ed25519.PublicKey //long-term
	secret, secret2, secret3 [32]byte
	sayhello                 []byte
	bAlice, aBob             [32]byte
}

func NewEd25519KeyPair(secret, public []byte) (*Ed25519KeyPair, error) {
	var kp Ed25519KeyPair

	if err := checkEd25519KeySize(secret, public); err != nil {
		return nil, err
	}
	if hasSmallOrder(public) {
		return nil, fmt.Errorf("invalid public key(check ed25519 black list)")
	}

	kp.Secret = secret
	kp.Public = public

	return &kp, nil
}

func checkEd25519KeySize(secretKey, publicKey []byte) error {
	if n := len(secretKey); n != ed25519.PrivateKeySize {
		return fmt.Errorf("invalid size of ed25519 private key, expect %v, but got %v", ed25519.PrivateKeySize, n)
	}
	if n := len(publicKey); n != ed25519.PublicKeySize {
		return fmt.Errorf("invalid size of ed25519 public key, expect %v, but got %v", ed25519.PublicKeySize, n)
	}
	return nil
}

// newPeerState initializes the state needed by both client and server
func newPeerState(networkIdentifier []byte, local Ed25519KeyPair) (*PeerState, error) {
	// peer local generate a random to share secret derivation
	pubKey, secKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	pstate := PeerState{
		remotePublic: make([]byte, ed25519.PublicKeySize),
	}
	copy(pstate.networkIdentifier[:], networkIdentifier)
	copy(pstate.localExchange.Secret[:], secKey[:])
	copy(pstate.localExchange.Public[:], pubKey[:])

	pstate.local = local

	if err := checkEd25519KeySize(pstate.local.Secret, pstate.local.Public); err != nil {
		return nil, err
	}

	return &pstate, nil
}

// NewClientState initializes the state for the client side
// the client must know the server’s public key before connecting.
func NewClientPeerState(networkIdentifier []byte, local Ed25519KeyPair, remotePublic ed25519.PublicKey) (*PeerState, error) {
	state, err := newPeerState(networkIdentifier, local)
	if err != nil {
		return state, err //but my data is correct
	}
	state.remotePublic = remotePublic

	var noRemoteSecret [ed25519.PrivateKeySize]byte
	if err = checkEd25519KeySize(noRemoteSecret[:], state.remotePublic); err != nil {
		return nil, fmt.Errorf("remote publick key err = %s", err) //but remote pk is wrong
	}
	return state, err
}

// NewClientState initializes the state for the client side
// the server learns the client’s public key during the handshake.
func NewServerPeerState(networkIdentifier []byte, local Ed25519KeyPair) (*PeerState, error) {
	state, err := newPeerState(networkIdentifier, local)
	return state, err
}

// createChallenge returns a buffer with a challenge
func (prs *PeerState) createChallenge() []byte {
	mac := auth.Sum(prs.localExchange.Public[:], &prs.networkIdentifier) //returns 32-byte digest
	copy(prs.localAppMac[:], mac[:])

	challengeData := append(prs.localAppMac[:], prs.localExchange.Public[:]...)
	return challengeData //32(net-mac)+32(local pk)
}

// verifyChallenge returns whether the passed buffer is valid
func (prs *PeerState) verifyChallenge(ch []byte) bool {
	mac := ch[:32]
	//for me ,i know the challenge creator in [32:], ephemeral
	remoteEphPubKey := ch[32:]

	// Verify checks that digest is a valid authenticator of message m under the
	// given secret key. Verify does not leak timing information.
	vret := auth.Verify(mac, remoteEphPubKey, &prs.networkIdentifier)

	copy(prs.remoteExchange.Public[:], remoteEphPubKey)
	prs.remoteAppMac = mac
	var sec [32]byte
	curve25519.ScalarMult(&sec, &prs.localExchange.Secret, &prs.remoteExchange.Public)
	copy(prs.secret[:], sec[:])
	secHasher := sha256.New()
	secHasher.Write(prs.secret[:])
	prs.secretHash = secHasher.Sum(nil)

	return vret
}

// createClientAuth returns a buffer containing a clientAuth message
func (prs *PeerState) createClientAuth() ([]byte, error) {
	var curveRemotePubKey [32]byte
	if !PublicKeyToCurve25519(&curveRemotePubKey, prs.remotePublic[:]) {
		return nil, errors.New("could not convert remote pubkey to Curve25519")
	}
	var aBob [32]byte
	curve25519.ScalarMult(&aBob, &prs.localExchange.Secret, &curveRemotePubKey)
	copy(prs.aBob[:], aBob[:])

	secHasher := sha256.New()
	secHasher.Write(prs.networkIdentifier[:])
	secHasher.Write(prs.secret[:])
	secHasher.Write(prs.aBob[:])
	copy(prs.secret2[:], secHasher.Sum(nil))

	/*Client computes
	detached_signature_A = nacl_sign_detached(
		  msg: concat(
			network_identifier,
			server_longterm_pk,
			sha256(shared_secret_ab)
		  ),
		  key: client_longterm_sk
		)
	*/
	var sigMsg bytes.Buffer
	sigMsg.Write(prs.networkIdentifier[:])
	sigMsg.Write(prs.remotePublic[:])
	sigMsg.Write(prs.secretHash)
	signDetached := ed25519.Sign(prs.local.Secret, sigMsg.Bytes())

	/*Client sends (112 bytes)
	nacl_secret_box(
	  msg: concat(
		detached_signature_A,
		client_longterm_pk
	  ),
	  nonce: 24_bytes_of_zeros,
	  key: sha256(
		concat(
		  network_identifier,
		  shared_secret_ab,
		  shared_secret_aB
		)
	  )
	)
	*/
	var helloBuf bytes.Buffer
	helloBuf.Write(signDetached[:])
	helloBuf.Write(prs.local.Public[:])
	prs.sayhello = helloBuf.Bytes()

	out := make([]byte, 0, len(prs.sayhello)-box.Overhead)
	var n [24]byte
	out = box.SealAfterPrecomputation(out, prs.sayhello, &n, &prs.secret2)

	return out, nil
}

// plainHello assert(length(msg3_plaintext) == 96)
var plainHello [ed25519.SignatureSize + ed25519.PublicKeySize]byte

// verifyClientAuth returns whether a buffer contains a valid clientAuth message
func (prs *PeerState) verifyClientAuth(data []byte) bool {
	var cvSec, aBob [32]byte
	PrivateKeyToCurve25519(&cvSec, prs.local.Secret)
	curve25519.ScalarMult(&aBob, &cvSec, &prs.remoteExchange.Public)
	copy(prs.aBob[:], aBob[:])

	/*Server verifies
	msg3_plaintext = assert_nacl_secretbox_open(
				ciphertext: msg3,
				nonce: 24_bytes_of_zeros,
				key: sha256(
				concat(
					network_identifier,
					shared_secret_ab,
					shared_secret_aB
			)
		)
		)
	*/
	secHasher := sha256.New()
	secHasher.Write(prs.networkIdentifier[:])
	secHasher.Write(prs.secret[:])
	secHasher.Write(prs.aBob[:])
	copy(prs.secret2[:], secHasher.Sum(nil))

	prs.sayhello = make([]byte, 0, len(data)-16)

	var nonce [24]byte
	var openOk bool
	prs.sayhello, openOk = box.OpenAfterPrecomputation(prs.sayhello, data, &nonce, &prs.secret2)

	/*
		detached_signature_A = first_64_bytes(msg3_plaintext)
		client_longterm_pk = last_32_bytes(msg3_plaintext)
	*/
	var sig = make([]byte, ed25519.SignatureSize)
	var public = make([]byte, ed25519.PublicKeySize)

	if openOk {
		copy(sig, prs.sayhello[:ed25519.SignatureSize])
		copy(public[:], prs.sayhello[ed25519.SignatureSize:])

	} else {
		copy(sig, plainHello[:ed25519.SignatureSize])
		copy(public[:], plainHello[ed25519.SignatureSize:])
	}

	if hasSmallOrder(sig[:32]) {
		openOk = false
	}
	/*Server verifies
	assert_nacl_sign_verify_detached(
	  sig: detached_signature_A,
	  msg: concat(
		network_identifier,
		server_longterm_pk,
		sha256(shared_secret_ab)
	  ),
	  key: client_longterm_pk
	)
	*/
	var sigMsg bytes.Buffer
	sigMsg.Write(prs.networkIdentifier[:])
	sigMsg.Write(prs.local.Public[:])
	sigMsg.Write(prs.secretHash)
	verifyOk := ed25519.Verify(public, sigMsg.Bytes(), sig)
	copy(prs.remotePublic, public)

	return openOk && verifyOk
}

// createServerAccept returns a buffer containing a serverAccept message
func (prs *PeerState) createServerAccept() ([]byte, error) {
	var curveRemotePubKey [32]byte
	if !PublicKeyToCurve25519(&curveRemotePubKey, prs.remotePublic) {
		return nil, errors.New("could not convert remote pubkey to Curve25519")
	}
	var bAlice [32]byte
	curve25519.ScalarMult(&bAlice, &prs.localExchange.Secret, &curveRemotePubKey)
	copy(prs.bAlice[:], bAlice[:])

	/*Client verifies
	detached_signature_B = assert_nacl_secretbox_open(
	  ciphertext: msg4,
	  nonce: 24_bytes_of_zeros,
	  key: sha256(
		concat(
		  network_identifier,
		  shared_secret_ab,
		  shared_secret_aB,
		  shared_secret_Ab
		)
	  )
	)
	*/
	secHasher := sha256.New()
	secHasher.Write(prs.networkIdentifier[:])
	secHasher.Write(prs.secret[:])
	secHasher.Write(prs.aBob[:])
	secHasher.Write(prs.bAlice[:])
	copy(prs.secret3[:], secHasher.Sum(nil))

	/*
			assert_nacl_sign_verify_detached(
			  sig: detached_signature_B,
			  msg: concat(
				network_identifier,
				detached_signature_A,
				client_longterm_pk,
				sha256(shared_secret_ab)
			  ),
			  key: server_longterm_pk
		)
	*/
	var sigMsg bytes.Buffer
	sigMsg.Write(prs.networkIdentifier[:])
	sigMsg.Write(prs.sayhello[:])
	sigMsg.Write(prs.secretHash)

	sig := ed25519.Sign(prs.local.Secret, sigMsg.Bytes())

	var out = make([]byte, 0, len(sig)+16)
	var nonce [24]byte //nonce: 24_bytes_of_zeros, when handshake
	out = box.SealAfterPrecomputation(out, sig[:], &nonce, &prs.secret3)

	return out, nil
}

// verifyServerAccept returns whether the passed buffer contains a valid serverAccept message
func (prs *PeerState) verifyServerAccept(boxedOkay []byte) bool {
	var curveLocalSec [32]byte
	PrivateKeyToCurve25519(&curveLocalSec, prs.local.Secret)
	var bAlice [32]byte
	curve25519.ScalarMult(&bAlice, &curveLocalSec, &prs.remoteExchange.Public)
	copy(prs.bAlice[:], bAlice[:])

	/*Server sends (80 bytes)
	nacl_secret_box(
		  msg: detached_signature_B,
		  nonce: 24_bytes_of_zeros,
		  key: sha256(
			concat(
			  network_identifier,
			  shared_secret_ab,
			  shared_secret_aB,
			  shared_secret_Ab
			)
		  )
		)
	*/
	secHasher := sha256.New()
	secHasher.Write(prs.networkIdentifier[:])
	secHasher.Write(prs.secret[:])
	secHasher.Write(prs.aBob[:])
	secHasher.Write(prs.bAlice[:])
	copy(prs.secret3[:], secHasher.Sum(nil))

	/*Server computes
	detached_signature_B = nacl_sign_detached(
	  msg: concat(
		network_identifier,
		detached_signature_A,
		client_longterm_pk,
		sha256(shared_secret_ab)
	  ),
	  key: server_longterm_sk
	)
	*/
	var nonce [24]byte
	sig := make([]byte, 0, len(boxedOkay)-16)
	sig, openOk := box.OpenAfterPrecomputation(nil, boxedOkay, &nonce, &prs.secret3)

	var sigMsg bytes.Buffer
	sigMsg.Write(prs.networkIdentifier[:])
	sigMsg.Write(prs.sayhello[:])
	sigMsg.Write(prs.secretHash)

	verifyOk := ed25519.Verify(prs.remotePublic, sigMsg.Bytes(), sig)
	return verifyOk && openOk
}

// overwriteSecrets overwrites all intermediate secrets
func (prs *PeerState) overwriteSecrets() {
	var zeros [64]byte

	copy(prs.secretHash, zeros[:])
	copy(prs.secret[:], zeros[:])
	copy(prs.aBob[:], zeros[:])
	copy(prs.bAlice[:], zeros[:])

	h := sha256.New()
	h.Write(prs.secret3[:])
	copy(prs.secret[:], h.Sum(nil))
	copy(prs.secret2[:], zeros[:])
	copy(prs.secret3[:], zeros[:])
	copy(prs.localExchange.Secret[:], zeros[:])
}

// Remote returns the public key of the remote peer
func (prs *PeerState) Remote() []byte {
	return prs.remotePublic[:]
}
