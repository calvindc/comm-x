package handshake

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestPeerToPeerHandShake(t *testing.T) {
	log.SetOutput(os.Stdout)
	/* ERROR io r&w
	rwServer := new(io.ReadWriter)
	rwClient := new(io.ReadWriter)*/

	networkId := []byte("netwrkidnetwrkidnetwrkidnetwrkid1")
	/*networkId := make([]byte, 32)
	io.ReadFull(random(0xed), networkId)*/
	t.Log("networkId= " + base64.StdEncoding.EncodeToString(networkId))

	keySrv, err := GenerateEd25519KeyPair(random(0x11))
	if err != nil {
		t.Fatal(err)
	}
	serverState, err := NewServerPeerState(networkId, *keySrv)
	if err != nil {
		t.Error("NewServerPeerState err:", err)
	}
	paddrserver := PeerAddr{keySrv.Public}
	t.Log("server public key=", keySrv.Public)
	t.Log("server id: " + paddrserver.String())

	keyCli, err := GenerateEd25519KeyPair(random(0x22))
	if err != nil {
		t.Fatal(err)
	}
	clientState, err := NewClientPeerState(networkId, *keyCli, keySrv.Public) //client know remote pubkey
	if err != nil {
		t.Error("NewClientPeerState err:", err)
	}
	paddrclient := PeerAddr{keyCli.Public}
	t.Log("client public key=", keyCli.Public)
	t.Log("client id: " + paddrclient.String())
	// create server and client's r&w, io.Pipe()'s r&w will lock io's byte when r&w not coming,
	// conforming to communication scenarios
	rServer, wClient := io.Pipe()
	rClient, wServer := io.Pipe()
	rwServer := rw{rServer, wServer}
	rwClient := rw{rClient, wClient}

	ch := make(chan error, 2)
	go func() {
		err := ServerHandShake(serverState, rwServer)
		ch <- err
		//w EOF
		wServer.Close()
	}()
	time.Sleep(2 * time.Second) //simulate timeout scenarios in a peer to peer's secret handshake
	go func() {
		err := ClientHandShake(clientState, rwClient)
		ch <- err
		//w EOF
		wClient.Close()
	}()

	if err = <-ch; err != nil {
		t.Errorf("ch-1 read: %v", err)
	}
	time.Sleep(2 * time.Millisecond)
	if err = <-ch; err != nil {
		t.Errorf("ch-2 read: %v", err)
	}
	t.Logf("server secret len=%v,data=%v", len(serverState.secret), serverState.secret)
	t.Logf("client secret len=%v,data=%v", len(clientState.secret), clientState.secret)
	if reflect.DeepEqual(serverState.secret, clientState.secret) == false {
		t.Error("2 peer secret not equal")
	}
}

type rw struct { //implents io.ReadWriter
	io.Reader
	io.Writer
}

type random byte

// a io.Reader
func (r random) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = byte(r)
	}
	n = len(p)
	return
}

func TestGenerateEd25519KeyPair(t *testing.T) {
	quitSig := make(chan os.Signal)
	go func() {
		q := <-quitSig
		t.Logf("generate 25519 key interrupted,sig=%v", q)
		os.Exit(0)
	}()

	for turn := 10000; turn >= 0; turn-- {
		r := rand.Reader
		keypair, err := GenerateEd25519KeyPair(r)
		if err != nil {
			t.Error(err)
			break
		}
		t.Logf("pubkey=%v", keypair.Public)
		t.Logf("prvkey=%v", keypair.Secret)
		if len(keypair.Public) != 32 || len(keypair.Secret) != 64 {
			t.Error("key len error")
			break
		}
	}
}

func TestClientHandShake(t *testing.T) {
	//todo add test
}

func TestServerHandShake(t *testing.T) {
	//todo add test
}
