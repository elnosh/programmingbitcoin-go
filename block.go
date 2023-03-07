package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
)

// block header
type Block struct {
	version       uint32
	previousBlock [32]byte
	merkleRoot    [32]byte
	timestamp     uint32
	bits          uint32
	nonce         uint32
}

func (b Block) id() []byte {
	hash := hash256(b.serialize())
	return reverse(hash[:])
}

func parseBlock(b []byte) *Block {
	blockBuffer := bytes.NewBuffer(b)

	buf := make([]byte, 4)
	_, err := blockBuffer.Read(buf)
	if err != nil {
		fmt.Println("error getting block version: ", err)
	}
	version := binary.LittleEndian.Uint32(buf)

	var prevBlock [32]byte
	_, err = blockBuffer.Read(prevBlock[:])
	if err != nil {
		fmt.Println("error getting previous block id: ", err)
	}
	prevBlock = reverseByteArr32(prevBlock)

	var merkleRoot [32]byte
	_, err = blockBuffer.Read(merkleRoot[:])
	if err != nil {
		fmt.Println("error getting merkleRoot: ", err)
	}
	merkleRoot = reverseByteArr32(merkleRoot)

	_, err = blockBuffer.Read(buf)
	if err != nil {
		fmt.Println("error getting block timestamp: ", err)
	}
	timestamp := binary.LittleEndian.Uint32(buf)

	_, err = blockBuffer.Read(buf)
	if err != nil {
		fmt.Println("error getting block bits: ", err)
	}
	bits := binary.BigEndian.Uint32(buf)

	_, err = blockBuffer.Read(buf)
	if err != nil {
		fmt.Println("error getting block nonce: ", err)
	}
	nonce := binary.BigEndian.Uint32(buf)

	return &Block{version: version, previousBlock: prevBlock, merkleRoot: merkleRoot, timestamp: timestamp, bits: bits, nonce: nonce}
}

func (b Block) serialize() []byte {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, b.version)

	prevBlock := reverseByteArr32(b.previousBlock)
	merkleRoot := reverseByteArr32(b.merkleRoot)

	timestamp := make([]byte, 4)
	binary.LittleEndian.PutUint32(timestamp, b.timestamp)

	bits := make([]byte, 4)
	binary.BigEndian.PutUint32(bits, b.bits)

	nonce := make([]byte, 4)
	binary.BigEndian.PutUint32(nonce, b.nonce)

	return bytes.Join([][]byte{version, prevBlock[:], merkleRoot[:], timestamp, bits, nonce}, []byte{})
}

func (b Block) checkPow() bool {
	blockHash := new(big.Int).SetBytes(b.id())
	return blockHash.Cmp(b.target()) == -1
}

// get the target number from bits field
func (b Block) target() *big.Int {
	bits := make([]byte, 4)
	binary.BigEndian.PutUint32(bits, b.bits)

	// last byte in bits field is the exponent
	exponent := bits[len(bits)-1]

	// coefficient are the other 3 bytes interpreted in little endian
	coefficient := new(big.Int).SetBytes(reverse(bits[:len(bits)-1]))

	mul := new(big.Int).Exp(big.NewInt(256), big.NewInt(int64(exponent-3)), nil)

	// target = coefficient * 256^(exponent - 3)
	target := new(big.Int).Set(coefficient.Mul(coefficient, mul))
	return target
}

func (b Block) difficulty() *big.Int {
	// difficulty = 0xffff * 256^(0x1d-3) / target
	exp := big.NewInt(int64(0x1d - 3))
	num := fromHex("ffff")
	mul := new(big.Int).Exp(big.NewInt(256), exp, nil)
	num.Mul(num, mul)

	return num.Div(num, b.target())

}

func (b Block) bip9() bool {
	return b.version>>29 == 1
}

func (b Block) bip91() bool {
	return b.version>>4&1 == 1
}

func (b Block) bip141() bool {
	return b.version>>1&1 == 1
}
