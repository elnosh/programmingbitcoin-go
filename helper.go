package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"

	"golang.org/x/crypto/ripemd160"
)

const (
	Base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	SIGHASH_ALL    = 1
	SIGHASH_NONE   = 2
	SIGHASH_SINGLE = 3
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
	checksum := sha[:4]
	inp := bytes.Join([][]byte{input, checksum}, []byte{})
	return base58encode(inp)
}

func base58Decode(base58address string) ([]byte, error) {
	num := big.NewInt(0)

	for _, char := range base58address {
		num.Mul(num, big.NewInt(58))
		charIdx := strings.Index(Base58Alphabet, string(char))
		num.Add(num, big.NewInt(int64(charIdx)))
	}
	combined := num.Bytes()
	checksum := combined[len(combined)-4:]
	hash := hash256(combined[:len(combined)-4])

	if !bytes.Equal(hash[:4], checksum) {
		return nil, fmt.Errorf("bad address: checksum does not match '%v' '%v'", checksum, hash[:4])
	}

	return combined[1 : len(combined)-4], nil
}

func h160ToP2pkh(hash160 []byte, testnet bool) string {
	var prefix []byte
	if testnet {
		prefix = []byte{0x6f}
	} else {
		prefix = []byte{0x00}
	}
	pkhash := bytes.Join([][]byte{prefix, hash160}, []byte{})
	return base58encodeChecksum(pkhash)
}

func h160ToP2SH(hash160 []byte, testnet bool) string {
	var prefix []byte
	if testnet {
		prefix = []byte{0xc4}
	} else {
		prefix = []byte{0x05}
	}
	scriptHash := bytes.Join([][]byte{prefix, hash160}, []byte{})
	return base58encodeChecksum(scriptHash)
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

func readVarint(varint io.Reader) (int, error) {
	var numbuf []byte
	i := make([]byte, 1)
	_, err := varint.Read(i)
	if err != nil {
		return -1, err
	}

	if i[0] == 0xfd {
		numbuf = make([]byte, 2)
		_, err = varint.Read(numbuf)
		if err != nil {
			return -1, err
		}
		return int(binary.LittleEndian.Uint16(numbuf)), nil
	} else if i[0] == 0xfe {
		numbuf = make([]byte, 4)
		_, err = varint.Read(numbuf)
		if err != nil {
			return -1, err
		}
		return int(binary.LittleEndian.Uint32(numbuf)), nil
	} else if i[0] == 0xff {
		numbuf = make([]byte, 8)
		_, err = varint.Read(numbuf)
		if err != nil {
			return -1, err
		}
		return int(binary.LittleEndian.Uint64(numbuf)), nil
	}

	return int(i[0]), nil
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
		return nil, errors.New("integer too large")
	}
	return encodedRes, nil
}

func reverse(element []byte) []byte {
	reversed := make([]byte, len(element))
	counter := len(element) - 1
	for i := 0; i < len(element); i++ {
		reversed[i] = element[counter]
		counter--
	}
	return reversed
}
