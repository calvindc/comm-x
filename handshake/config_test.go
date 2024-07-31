package handshake

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"
)

func Test_hasSmallOrder(t *testing.T) {
	//var elementKey interface{}
	//elementKey = 10 //Ed25519KeyPair{}
	quit := make(chan os.Signal, 1)
	go func() {
		signal.Notify(quit, os.Interrupt, os.Kill)
		<-quit
		os.Exit(0)
	}()

	wg := sync.WaitGroup{}

	starttimt := time.Now()
	for i := 0; i < 1000; i++ {
		go func() {
			wg.Add(1)
			keySrv, err := GenerateEd25519KeyPair(nil) //random(0x11)
			if err != nil {
				t.Error(fmt.Sprintf("generate key error = %s", err.Error()))
			}
			pubkey := keySrv.Public
			t.Log("public  key=", keySrv.Public)
			paddr := PeerAddr{keySrv.Public}
			t.Log("base64  id=", paddr.String())
			if len(pubkey) != 32 {
				t.Error("key length invalid")
			}
			if hasSmallOrder(pubkey) {
				t.Error("check passed group element is  low order ,not passed")
			}

			wg.Done()

		}()
	}

	wg.Wait()

	t.Logf("test cost %s", time.Since(starttimt))

}
