package blockchainlisten

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ChainEvents struct {
	client            *helper.SafeEthClient
	be                *blockchain.Events
	bcs               *rpc.BlockChainService
	key               *ecdsa.PrivateKey
	quitChan          chan struct{}
	updateBalanceChan chan *userRequestUpdateBalanceProof
	stopped           bool
	TokenNetwork      *TokenNetwork
}
