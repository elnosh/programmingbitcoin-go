package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	modifiedTx, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000001976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88acfeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac1943060001000000")

	hash := hash256(modifiedTx)
	z := new(big.Int).SetBytes(hash[:])

	sec, _ := hex.DecodeString("0349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a")
	der, _ := hex.DecodeString("3045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed")

	pk := parsePubKey(sec)

	signature, err := parseSignature(der)
	if err != nil {
		fmt.Println(err)
	}

	if pk.verifySignature(*signature, z) {
		fmt.Println("valid signature")
	} else {
		fmt.Println("invalid signature")
	}
}
