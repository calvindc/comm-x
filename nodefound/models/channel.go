package models

import (
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
	"math/big"
)

type BalanceProof struct {
	Nonce           uint64      `json:"nonce"`
	TransferAmount  *big.Int    `json:"transfer_amount"`
	LocksRoot       common.Hash `json:"locks_root"`
	ChannelID       common.Hash `json:"channel_identifier"`
	OpenBlockNumber int64       `json:"open_block_number"`
	AdditionalHash  common.Hash `json:"addition_hash"`
	Signature       []byte      `json:"signature"`
}

type ChannelParticipantInfo struct {
	ID                     int
	ChannelID              string `gorm:"index"`
	Participant            string
	Nonce                  uint64
	Balance                string
	Deposit                string
	LockedAmount           string
	TransferedAmount       string
	IgnoreMediatedTransfer bool
}

func (c *ChannelParticipantInfo) BalanceValue() *big.Int {
	return string2BigInt(c.Balance)
}

func (c *ChannelParticipantInfo) Fee(token common.Address, db *gorm.DB) *Fee {
	return GetChannelFeeRate(common.HexToHash(c.ChannelID), common.HexToAddress(c.Participant), token, db)
}

// GetChannelFeeRate get channel's fee rate
func GetChannelFeeRate(channelIdentifier common.Hash, participant, token common.Address, db *gorm.DB) (fee *Fee) {
	cf, err := getDirectChannelFee(channelIdentifier, participant, db)
	if err == nil {
		fee = &Fee{
			FeePolicy:   cf.FeePolicy,
			FeeConstant: string2BigInt(cf.FeeConstantPart),
			FeePercent:  cf.FeePercentPart,
		}
		return
	}
	//从来没有针对通道设置过
	fee, err = GetAccountTokenFee(participant, token, db)
	if err == nil {
		return
	}
	fee = GetAccountFeePolicy(participant, db)
	return
}

func getDirectChannelFee(channelIdentifier common.Hash, participant common.Address, db *gorm.DB) (cf *ChannelParticipantFee, err error) {
	cf = &ChannelParticipantFee{
		ChannelID:   channelIdentifier.String(),
		Participant: participant.String(),
	}
	err = db.Where(cf).Find(cf).Error
	return
}
