package main

import (
	"crypto/sha256"
	"math/big"
)

// do two rounds of sha256
func hash256(s []byte) *big.Int {
	sum := sha256.Sum256(s)
	sum2 := sha256.Sum256([]byte(sum[:]))
	return new(big.Int).SetBytes(sum2[:])
}

func fromHex(s string) *big.Int {
	if s == "" {
		return big.NewInt(0)
	}
	r, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("invalid hex: " + s)
	}
	return r
}
