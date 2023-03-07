package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {
	b, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")

	block := parseBlock(b)
	fmt.Printf("hash = %x\n", block.id())
	hashProof := new(big.Int).SetBytes(block.id())
	fmt.Printf("targ = %x\n", block.target())
	if hashProof.Cmp(block.target()) == -1 {
		fmt.Println("hash is below target")
	}

	fmt.Println(block.difficulty())
}
