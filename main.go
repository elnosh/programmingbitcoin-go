package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	// encint := 18005558675309
	// varintbyte, _ := encodeVarint(encint)
	// fmt.Println(hex.EncodeToString(varintbyte))

	// varintbyte := []byte{0xfe, 0x7f, 0x11, 0x01, 0x00}
	// varint, _ := readVarint(varintbyte)
	// fmt.Println(varint)

	txBytes, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")

	tx := parseTx(txBytes)
	fmt.Printf("tx version = %d\n", tx.version)

	fmt.Printf("number of inputs = %d\n", len(tx.txIns))
	fmt.Printf("prev tx id = %x\n", tx.txIns[0].prevTxId[:])
	fmt.Printf("prev tx index = %v\n", tx.txIns[0].prevTxIdx)
}

// func main() {
// 	// hexval := "c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6"
// 	// hexbytes, _ := hex.DecodeString(hexval)

// 	// base58str := base58encode(hexbytes)
// 	// fmt.Println(base58str)

// 	privKey := newPrivateKey(fromHex("54321deadbeef"))
// 	fmt.Println(privKey.wif(true, false))
// }

// func main() {
// 	privKey := newPrivateKey(big.NewInt(12345))

// 	z := hash256([]byte("some message"))
// 	zint := new(big.Int).SetBytes(z[:])
// 	signature := privKey.sign(zint)

// 	point := g.rmulS256(privKey.secret)
// 	point.repr()
// 	fmt.Println(zint.Text(16))
// 	fmt.Println(signature.r.Text(16))
// 	fmt.Println(signature.s.Text(16))
// }
