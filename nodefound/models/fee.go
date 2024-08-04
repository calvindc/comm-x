package models

import (
	"github.com/calvindc/comm-x/nodefound/config"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
	"math/big"
)

// AccountFee 账户的缺省收费
type AccountFee struct {
	Account         string `gorm:"primary_key"`
	FeePolicy       int    //收费模式
	FeeConstantPart string //固定费用
	FeePercentPart  int64  //比例费用
}

// AccountTokenFee 某个账户针对某个Token的缺省收费
type AccountTokenFee struct {
	gorm.Model             //auto id createat updataat deleteat
	Token           string `gorm:"index"` //调用者保证Token+Account必须是唯一的
	Account         string `gorm:"index"` //调用者保证Token+Account必须是唯一的
	FeePolicy       int
	FeeConstantPart string
	FeePercentPart  int64
}

// TokenFee 针对某种token 的缺省收费,暂不启用
type TokenFee struct {
	Token           string `gorm:"primary_key"`
	FeePolicy       int
	FeeConstantPart string
	FeePercentPart  int64
}

// Fee 为了使用方便 外部定义
type Fee struct {
	FeePolicy   int      `json:"fee_policy"`
	FeeConstant *big.Int `json:"fee_constant" `
	FeePercent  int64    `json:"fee_percent"`
}

var defaultFee = &Fee{
	FeePolicy:   config.DefaultFeePolicy,
	FeeConstant: config.DefaultFeeConstantPart,
	FeePercent:  config.DefaultFeePercentPart,
}

// 设置某个账户的缺省收费,新创建的通道都会按照此缺省设置进行
func UpdataAccountDefaultFeePolicy(account common.Address, fee *Fee, db *gorm.DB) error {
	acc := &AccountFee{
		Account:         account.String(),
		FeePolicy:       fee.FeePolicy,
		FeeConstantPart: bigInt2String(fee.FeeConstant),
		FeePercentPart:  fee.FeePercent,
	}
	err := db.Where(&AccountFee{Account: account.String()}).Find(&AccountFee{}).Error
	if err == nil {
		return db.Save(acc).Error
	}

	return db.Create(acc).Error
}

// GetAccountFeePolicy 获取某个账户的缺省收费,新创建的通道都会按照此缺省设置进行
func GetAccountFeePolicy(account common.Address, db *gorm.DB) (fee *Fee) {
	a := &AccountFee{}
	err := db.Where(&AccountFee{Account: account.String()}).Find(a).Error
	if err == nil {
		return &Fee{
			FeePolicy:   a.FeePolicy,
			FeeConstant: string2BigInt(a.FeeConstantPart),
			FeePercent:  a.FeePercentPart,
		}
	}
	return &Fee{
		defaultFee.FeePolicy,
		defaultFee.FeeConstant,
		defaultFee.FeePercent,
	}
}

func GetAccountTokenFee(account, token common.Address, db *gorm.DB) (fee *Fee, err error) {
	atf := &AccountTokenFee{
		Token:   token.String(),
		Account: account.String(),
	}
	err = db.Where(atf).Find(atf).Error
	if err == nil {
		fee = &Fee{
			FeePolicy:   atf.FeePolicy,
			FeeConstant: string2BigInt(atf.FeeConstantPart),
			FeePercent:  atf.FeePercentPart,
		}
	}
	return
}
