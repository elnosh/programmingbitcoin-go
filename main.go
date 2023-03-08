package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	firstBlockHex, _ := hex.DecodeString("000000203471101bbda3fe307664b3283a9ef0e97d9a38a7eacd8800000000000000000010c8aba8479bbaa5e0848152fd3c2289ca50e1c3e58c9a4faaafbdf5803c5448ddb845597e8b0118e43a81d3")

	lastBlockHex, _ := hex.DecodeString("02000020f1472d9db4b563c35f97c428ac903f23b7fc055d1cfc26000000000000000000b3f449fcbe1bc4cfbcb8283a0d2c037f961a3fdf2b8bedc144973735eea707e1264258597e8b0118e5f00474")

	firstBlock := parseBlock(firstBlockHex)
	lastBlock := parseBlock(lastBlockHex)

	newBits := calculateNewBits(lastBlock.bits, lastBlock.timestamp-firstBlock.timestamp)
	fmt.Printf("bits = %x\n", newBits)
}
