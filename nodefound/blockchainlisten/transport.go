package blockchainlisten

import "github.com/ethereum/go-ethereum/common"

type Transporter interface {
	SubscribeNeighbors(addrs []common.Address) error
	Unsubscribe(addr common.Address) error

	Stop() //stop transport service
}
