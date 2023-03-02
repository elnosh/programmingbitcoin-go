package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	secret := fromHex("")
	privKey := newPrivateKey(secret)
	fmt.Println(privKey.point.address(true, true))

	var txIns []TxIn
	txId, _ := hex.DecodeString("")
	var prevTxId [32]byte
	copy(prevTxId[:], txId)
	var prevTxIdx uint32 = 1

	txIn := newTxIn(prevTxId, prevTxIdx, nil, 0xfffffffe)
	txIns = append(txIns, *txIn)

	var txOuts []TxOut

	var destAmount uint64
	destAddrHash, err := base58Decode("")
	if err != nil {
		fmt.Println(err)
	}
	destScriptPubKey := p2pkhScript(destAddrHash)
	destTxOut := TxOut{value: destAmount, scriptPubKey: destScriptPubKey}
	txOuts = append(txOuts, destTxOut)

	var changeAmount uint64
	changeAddrHash, err := base58Decode("")
	if err != nil {
		fmt.Println(err)
	}
	changeScriptPubKey := p2pkhScript(changeAddrHash)
	changeTxOut := TxOut{value: changeAmount, scriptPubKey: changeScriptPubKey}
	txOuts = append(txOuts, changeTxOut)

	tx := &Tx{version: 1, txIns: txIns, txOuts: txOuts, locktime: 0, testnet: true}
	fmt.Printf("tx hex = %v\n", hex.EncodeToString(tx.serialize()))
}
