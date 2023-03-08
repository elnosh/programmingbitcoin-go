package main

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBlock(t *testing.T) {
	rawBlock, err := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if err != nil {
		t.Error("error decoding block")
	}

	block := parseBlock(rawBlock)

	var expectedNum uint32 = 0x20000002
	assert.Equal(t, expectedNum, block.version, "version does not match")

	var want [32]byte
	prevBlockHex, _ := hex.DecodeString("000000000000000000fd0c220a0a8c3bc5a7b487e8c8de0dfa2373b12894c38e")
	copy(want[:], prevBlockHex)
	assert.Equal(t, want, block.previousBlock)

	merkleRootHex, _ := hex.DecodeString("be258bfd38db61f957315c3f9e9c5e15216857398d50402d5089a8e0fc50075b")
	copy(want[:], merkleRootHex)
	assert.Equal(t, want, block.merkleRoot)

	expectedNum = 0x59a7771e
	assert.Equal(t, expectedNum, block.timestamp, "timestamp does not match")

	expectedBytes := [4]byte{0xe9, 0x3c, 0x01, 0x18}
	assert.Equal(t, expectedBytes, block.bits, "bits does not match")

	expectedBytes = [4]byte{0xa4, 0xff, 0xd7, 0x1d}
	assert.Equal(t, expectedBytes, block.nonce, "nonce does not match")
}

func TestSerializeBlock(t *testing.T) {
	rawBlock, err := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if err != nil {
		t.Error("error decoding block")
	}
	block := parseBlock(rawBlock)
	assert.Equal(t, rawBlock, block.serialize(), "blocks serialized do not match")
}

func TestTarget(t *testing.T) {
	rawBlock, err := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if err != nil {
		t.Error("error decoding block")
	}

	block := parseBlock(rawBlock)
	want := fromHex("13ce9000000000000000000000000000000000000000000")
	assert.Equal(t, want, block.target(), "targets do not match")
}

func TestDifficulty(t *testing.T) {
	rawBlock, err := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if err != nil {
		t.Error("error decoding block")
	}

	block := parseBlock(rawBlock)
	assert.Equal(t, big.NewInt(888171856257), block.difficulty(), "difficulty does not match")
}

func TestCheckPow(t *testing.T) {
	rawBlock, err := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	if err != nil {
		t.Error("error decoding block")
	}

	block := parseBlock(rawBlock)
	assert.Equal(t, true, block.checkPow())
}

//func TestCalculateNewBits(t *testing.T) {
//	prevBits := [4]byte{0x54, 0xd8, 0x01, 0x18}
//	var timeDifferential uint32 = 302400
//	want := [4]byte{0x00, 0x00, 0x15, 0x17}
//	//want := [4]byte{0x00, 0x15, 0x76, 0x17}

//	fmt.Printf("calculated new bits = %x\n", calculateNewBits(prevBits, timeDifferential))
//	assert.Equal(t, want, calculateNewBits(prevBits, timeDifferential), "bits do not match")
//}
