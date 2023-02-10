package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
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

func readVarint(hexnum []byte) int {
	i := hexnum[0]

	numbuf := hexnum[1:]
	if i == 0xfd {
		return int(binary.LittleEndian.Uint16(numbuf))
	} else if i == 0xfe {
		return int(binary.LittleEndian.Uint32(numbuf))
	} else if i == 0xff {
		return int(binary.LittleEndian.Uint64(numbuf))
	}

	return int(i)
}

func encodeVarint(num int) ([]byte, error) {
	cmpInt := []byte{0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0}
	var varintbuf, prefix, encodedRes []byte

	if num < 0xfd {
		return []byte{byte(num)}, nil
	} else if num < 0x10000 {
		varintbuf = make([]byte, 2)
		prefix = []byte{0xfd}
		binary.LittleEndian.PutUint16(varintbuf, uint16(num))
		encodedRes = bytes.Join([][]byte{prefix, varintbuf}, []byte{})
	} else if num < 0x100000000 {
		varintbuf = make([]byte, 4)
		prefix = []byte{0xfe}
		binary.LittleEndian.PutUint32(varintbuf, uint32(num))
		encodedRes = bytes.Join([][]byte{prefix, varintbuf}, []byte{})
	} else if big.NewInt(int64(num)).Cmp(new(big.Int).SetBytes(cmpInt)) == -1 {
		varintbuf = make([]byte, 8)
		prefix = []byte{0xff}
		binary.LittleEndian.PutUint64(varintbuf, uint64(num))
		encodedRes = bytes.Join([][]byte{prefix, varintbuf}, []byte{})
	} else {
		// err value too large
		return nil, errors.New("error encoding varint: integer too large")
	}
	return encodedRes, nil
}
