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
	bits          [4]byte
	nonce         [4]byte
	txHashes      [][]byte
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

	var bits [4]byte
	_, err = blockBuffer.Read(bits[:])
	if err != nil {
		fmt.Println("error getting block bits: ", err)
	}

	var nonce [4]byte
	_, err = blockBuffer.Read(nonce[:])
	if err != nil {
		fmt.Println("error getting block nonce: ", err)
	}

	return &Block{version: version, previousBlock: prevBlock, merkleRoot: merkleRoot, timestamp: timestamp, bits: bits, nonce: nonce}
}

func (b Block) serialize() []byte {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, b.version)

	prevBlock := reverseByteArr32(b.previousBlock)
	merkleRoot := reverseByteArr32(b.merkleRoot)

	timestamp := make([]byte, 4)
	binary.LittleEndian.PutUint32(timestamp, b.timestamp)

	var bits [4]byte = b.bits
	var nonce [4]byte = b.nonce

	return bytes.Join([][]byte{version, prevBlock[:], merkleRoot[:], timestamp, bits[:], nonce[:]}, []byte{})
}

// checks if the block header hash is below the target difficulty
func (b Block) checkPow() bool {
	blockHash := new(big.Int).SetBytes(b.id())
	return blockHash.Cmp(b.target()) == -1
}

// get the target number from bits field
func (b Block) target() *big.Int {
	return bitsToTarget(b.bits)
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

func (b Block) validateMerkleRoot() bool {
	merkleRoot := merkleParentRoot(b.txHashes)
	return bytes.Equal(merkleRoot, b.merkleRoot[:])
}

func bitsToTarget(bits [4]byte) *big.Int {
	// last byte in bits field is the exponent
	exponent := bits[len(bits)-1]

	// coefficient are the other 3 bytes interpreted in little endian
	coefficient := new(big.Int).SetBytes(reverse(bits[:len(bits)-1]))

	mul := new(big.Int).Exp(big.NewInt(256), big.NewInt(int64(exponent-3)), nil)

	// target = coefficient * 256^(exponent - 3)
	target := new(big.Int).Set(coefficient.Mul(coefficient, mul))
	return target
}

func targetToBits(target *big.Int) [4]byte {
	rawBytes := make([]byte, 32)
	rawBytes = target.FillBytes(rawBytes)
	rawBytes = bytes.TrimLeft(rawBytes, string(byte(0)))

	exponent := 0
	coefficient := []byte{0x00}
	if rawBytes[0] > 0x7f {
		exponent = len(rawBytes) + 1
		coefficient = append(coefficient, rawBytes[:2]...)
	} else {
		exponent = len(rawBytes)
		coefficient = rawBytes[:3]
	}

	var newBits [4]byte
	j := 3
	for i := 0; i < 3; i++ {
		newBits[i] = rawBytes[j]
		j--
	}
	newBits[3] = byte(exponent)
	return newBits
}

// time differential = (block timestamp of last block in difficulty adjustment period) - (block timestamp of first block in difficulty adjustment period)
// to calculate new target = previous target * time differential / (2 weeks)
func calculateNewBits(previousBits [4]byte, timeDifferential uint32) [4]byte {
	if timeDifferential > TWO_WEEKS*4 {
		timeDifferential = TWO_WEEKS * 4
	} else if timeDifferential < TWO_WEEKS/4 {
		timeDifferential = TWO_WEEKS / 4
	}
	previousTarget := bitsToTarget(previousBits)
	newTarget := new(big.Int).Mul(previousTarget, big.NewInt(int64(timeDifferential)))
	newTarget.Div(newTarget, big.NewInt(TWO_WEEKS))
	return targetToBits(newTarget)
}
