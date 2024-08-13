package blockchainlisten

import (
	"github.com/calvindc/comm-x/nodefound/models"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"sync"
)

type channel struct {
	Participant1        common.Address
	Participant2        common.Address
	Participant1Balance *big.Int
	Participant2Balance *big.Int
	Participant1Fee     *models.Fee
	Participant2Fee     *models.Fee
	Token               common.Address
}

// TokenNetwork token network view
type TokenNetwork struct {
	TokensNetworkAddress common.Address
	channelViews         map[common.Address][]*channel //token to channels
	channels             map[common.Hash]*channel      //channel id to chann
	token2TokenNetwork   map[common.Address]common.Address
	decimals             map[common.Address]int
	viewlock             sync.RWMutex
	participantStatus    map[common.Address]nodeStatus
	nodeLock             sync.Mutex
	transport            Transporter
}

type nodeStatus struct {
	isMobile               bool
	isOnline               bool
	ignoreMediatedTransfer bool
}
