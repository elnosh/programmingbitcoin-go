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

func TestTxInputValue(t *testing.T) {
	var txHashHex [32]byte
	tx, err := hex.DecodeString("d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81")
	if err != nil {
		t.Errorf("error decoding tx hash: %v\n", err)
	}
	copy(txHashHex[:], tx)

	var idx uint32 = 0
	var want uint64 = 42505594

	txIn := newTxIn(txHashHex, idx, nil, uint32(0xfffffffe))
	assert.Equal(t, want, txIn.value(false))
}

func TestInputPubKey(t *testing.T) {
	var txHashHex [32]byte
	tx, err := hex.DecodeString("d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81")
	if err != nil {
		t.Errorf("error decoding tx hash: %v\n", err)
	}
	copy(txHashHex[:], tx)

	var idx uint32 = 0
	txIn := newTxIn(txHashHex, idx, nil, uint32(0xfffffffe))

	want, err := hex.DecodeString("1976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88ac")
	if err != nil {
		t.Errorf("error decoding expected value: %v\n", err)
	}
	assert.Equal(t, want, txIn.scriptPubKey(false).serialize(), "scriptPubKey do not match")
}

func TestFee(t *testing.T) {
	testCases := []struct {
		rawTx string
		want  uint64
	}{
		{"0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600", 40000},
		{"010000000456919960ac691763688d3d3bcea9ad6ecaf875df5339e148a1fc61c6ed7a069e010000006a47304402204585bcdef85e6b1c6af5c2669d4830ff86e42dd205c0e089bc2a821657e951c002201024a10366077f87d6bce1f7100ad8cfa8a064b39d4e8fe4ea13a7b71aa8180f012102f0da57e85eec2934a82a585ea337ce2f4998b50ae699dd79f5880e253dafafb7feffffffeb8f51f4038dc17e6313cf831d4f02281c2a468bde0fafd37f1bf882729e7fd3000000006a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937feffffff567bf40595119d1bb8a3037c356efd56170b64cbcc160fb028fa10704b45d775000000006a47304402204c7c7818424c7f7911da6cddc59655a70af1cb5eaf17c69dadbfc74ffa0b662f02207599e08bc8023693ad4e9527dc42c34210f7a7d1d1ddfc8492b654a11e7620a0012102158b46fbdff65d0172b7989aec8850aa0dae49abfb84c81ae6e5b251a58ace5cfeffffffd63a5e6c16e620f86f375925b21cabaf736c779f88fd04dcad51d26690f7f345010000006a47304402200633ea0d3314bea0d95b3cd8dadb2ef79ea8331ffe1e61f762c0f6daea0fabde022029f23b3e9c30f080446150b23852028751635dcee2be669c2a1686a4b5edf304012103ffd6f4a67e94aba353a00882e563ff2722eb4cff0ad6006e86ee20dfe7520d55feffffff0251430f00000000001976a914ab0c0b2e98b1ab6dbf67d4750b0a56244948a87988ac005a6202000000001976a9143c82d7df364eb6c75be8c80df2b3eda8db57397088ac46430600", 140500},
	}

	for _, test := range testCases {
		txHex, err := hex.DecodeString(test.rawTx)
		if err != nil {
			t.Errorf("error decoding tx hex: %v\n", err)
		}

		tx := parseTx(txHex)
		fee := tx.fee(false)
		if fee != test.want {
			t.Errorf("expected %v but got %v instead", test.want, fee)
		}
	}
}
