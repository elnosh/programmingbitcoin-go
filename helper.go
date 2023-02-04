package main

import (
	"bytes"
	"crypto/sha256"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

const (
	Base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

// do two rounds of sha256
func hash256(input []byte) [32]byte {
	sum := sha256.Sum256(input)
	return sha256.Sum256([]byte(sum[:]))
	//return new(big.Int).SetBytes(sum2[:])
}

// sha256 + ripemd160
func hash160(input []byte) []byte {
	h256 := sha256.Sum256(input)
	h := ripemd160.New()
	h.Write(h256[:])
	return h.Sum(nil)
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

func base58encodeChecksum(input []byte) string {
	sha := hash256(input)
	firstFour := sha[:4]
	inp := bytes.Join([][]byte{input, firstFour}, []byte{})
	return base58encode(inp)
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
