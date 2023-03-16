package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	networkMessage, _ := hex.DecodeString("f9beb4d976657261636b000000000000000000005df6e0e2")

	netenvelope := parseNetworkEnvelope(networkMessage, false)
	fmt.Println(string(netenvelope.command[:]))
	fmt.Println(string(netenvelope.payload))
}
