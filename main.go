package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

func main() {
	script, _ := hex.DecodeString("4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73")
	scriptBuf := bytes.NewBuffer(script)

	scriptSig, _ := parseScript(scriptBuf)
	fmt.Printf("%s\n", scriptSig.cmds[2])
}
