package main

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseScript(t *testing.T) {
	hexScriptPubKey, err := hex.DecodeString("6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937")
	if err != nil {
		t.Error("error decoding scriptPubKey")
	}
	scriptBuf := bytes.NewBuffer(hexScriptPubKey)
	scriptPubKey, err := parseScript(scriptBuf)
	if err != nil {
		t.Error("error parsing script")
	}

	testCases := []struct {
		cmd  []byte
		want string
	}{
		{scriptPubKey.cmds[0], "304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a71601"},
		{scriptPubKey.cmds[1], "035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937"},
	}

	for _, test := range testCases {
		cmdString := hex.EncodeToString(test.cmd)
		if cmdString != test.want {
			t.Errorf("expected '%v' but got '%v' instead", test.want, cmdString)
		}
	}
}

func TestSerializeScript(t *testing.T) {
	want := "6a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937"
	scriptPubKey, err := hex.DecodeString(want)
	if err != nil {
		t.Error("error decoding target script")
	}

	script, err := parseScript(bytes.NewBuffer(scriptPubKey))
	if err != nil {
		t.Error("error parsing script")
	}
	assert.Equal(t, want, hex.EncodeToString(script.serialize()), "scripts serialized do not match")
}
