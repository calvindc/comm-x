package blockchainlisten

import (
	"crypto/ecdsa"
	"github.com/SmartMeshFoundation/Photon/blockchain"
	"github.com/SmartMeshFoundation/Photon/network/helper"
	"github.com/SmartMeshFoundation/Photon/network/rpc"
	"github.com/calvindc/comm-x/nodefound/models"
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

type userRequestUpdateBalanceProof struct {
	participant            common.Address
	partner                common.Address
	lockedAmount           *big.Int
	partnerBalanceProof    *models.BalanceProof
	ignoreMediatedTransfer bool
}
