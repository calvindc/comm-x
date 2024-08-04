package models

import "math/big"

func bigInt2String(b *big.Int) string {
	if b == nil {
		return "0"
	}
	return b.String()
}

func string2BigInt(s string) *big.Int {
	bi, b := new(big.Int).SetString(s, 10)
	if !b {
		bi = new(big.Int)
	}
	return bi
}
