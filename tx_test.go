package main

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVersion(t *testing.T) {
	txHex, err := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	if err != nil {
		t.Errorf("error decoding tx hex: %v\n", err)
	}
	tx := parseTx(txHex)
	assert.Equal(t, uint32(1), tx.version)
}

func TestParseTxInput(t *testing.T) {
	txHex, err := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	if err != nil {
		t.Errorf("error decoding tx hex: %v\n", err)
	}

	tx := parseTx(txHex)

	assert.Equal(t, 1, len(tx.txIns), "unexpected TxIns length")

	want, _ := hex.DecodeString("d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81")
	assert.Equal(t, want, tx.txIns[0].prevTxId[:], "tx in bytes do not match")
	assert.Equal(t, uint32(0), tx.txIns[0].prevTxIdx, "tx idxs do not match")

	want, _ = hex.DecodeString("6b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a")
	assert.Equal(t, want, tx.txIns[0].scriptSig.serialize(), "scriptSig does not match")
	assert.Equal(t, uint32(0xfffffffe), tx.txIns[0].sequence, "scriptSig does not match")
}

func TestParseTxOutput(t *testing.T) {
	txHex, err := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	if err != nil {
		t.Errorf("error decoding tx hex: %v\n", err)
	}

	tx := parseTx(txHex)

	assert.Equal(t, 2, len(tx.txOuts), "unexpected TxOuts length")
	var want uint64 = 32454049
	assert.Equal(t, want, tx.txOuts[0].value, "txOut amount does not match")

	pubKeyWant, _ := hex.DecodeString("1976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac")
	assert.Equal(t, pubKeyWant, tx.txOuts[0].scriptPubKey.serialize(), "public key do not match")

	want = 10011545
	assert.Equal(t, want, tx.txOuts[1].value, "txOut amount does not match")

	pubKeyWant, _ = hex.DecodeString("1976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac")
	assert.Equal(t, pubKeyWant, tx.txOuts[1].scriptPubKey.serialize(), "public key do not match")
}

func TestParseLocktime(t *testing.T) {
	txHex, err := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	if err != nil {
		t.Errorf("error decoding tx hex: %v\n", err)
	}
	tx := parseTx(txHex)
	assert.Equal(t, uint32(410393), tx.locktime)
}

func TestSerialize(t *testing.T) {
	txHex, err := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	if err != nil {
		t.Errorf("error decoding tx hex: %v\n", err)
	}

	tx := parseTx(txHex)
	assert.Equal(t, txHex, tx.serialize(), "hex value of serialize does not match")
}

// func TestTxInputValue(t *testing.T) {
// 	var txHashHex [32]byte
// 	tx, err := hex.DecodeString("d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81")
// 	if err != nil {
// 		t.Errorf("error decoding tx hash: %v\n", err)
// 	}
// 	copy(txHashHex[:], tx)

// 	var idx uint32 = 0
// 	want := 42505594

// 	txIn := newTxIn(txHashHex, idx, nil, uint32(0xfffffffe))
// 	assert.Equal(t, want, txIn.value(false))
// }
