package config

import "math/big"

// DefaultFeePolicy 缺省按比例收费 0固定费用
var DefaultFeePolicy = 1

// DefaultFeePercentPart 比例缺省万分之一
var DefaultFeePercentPart int64 = 10000

// DefaultFeeConstantPart 收费固定部分为0
var DefaultFeeConstantPart = big.NewInt(0)
