package main

import (
	"crypto/sha256"
	"math/big"
)

const (
	Base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

// do two rounds of sha256
func hash256(s []byte) *big.Int {
	sum := sha256.Sum256(s)
	sum2 := sha256.Sum256([]byte(sum[:]))
	return new(big.Int).SetBytes(sum2[:])
}

func base58encode(input []byte) string {
	prefix := ""
	for _, inbyte := range input {
		if inbyte == 0 {
			prefix += "1"
		} else {
			break
		}
	}

	num := big.NewInt(0).SetBytes(input)
	result := ""
	for num.Sign() > 0 {
		mod := new(big.Int)
		num, mod = num.DivMod(num, big.NewInt(58), mod)
		result = string(Base58Alphabet[mod.Int64()]) + result
	}
	return prefix + result
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
