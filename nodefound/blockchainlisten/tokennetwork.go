package blockchainlisten

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"sort"
	"sync"
	"time"

	"errors"

	"github.com/ethereum/go-ethereum/common"
)

type channel struct {
	Participant1        common.Address
	Participant2        common.Address
	Participant1Balance *big.Int
	Participant2Balance *big.Int
	Participant1Fee     *model.Fee
	Participant2Fee     *model.Fee
	Token               common.Address
}

type channel struct {
	Participant1        common.Address
	Participant2        common.Address
	Participant1Balance *big.Int
	Participant2Balance *big.Int
	Participant1Fee     *model.Fee
	Participant2Fee     *model.Fee
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
