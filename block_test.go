package main

import (
	"encoding/hex"
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

	expectedNum = 0xe93c0118
	assert.Equal(t, expectedNum, block.bits, "bits does not match")

	expectedNum = 0xa4ffd71d
	assert.Equal(t, expectedNum, block.nonce, "nonce does not match")
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
