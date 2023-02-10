package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// transaction
type Tx struct {
	version int
	txIns   []TxIn
	//txOuts
	locktime int
	testnet  bool
}

// hex of tx hash
func (tx Tx) id() {}

// binary hash of legacy serialization
func (tx Tx) hash() {}

func parse(s []byte) {
	sbuf := bytes.NewBuffer(s)
	versionbuf := make([]byte, 4)
	_, err := sbuf.Read(versionbuf)
	if err != nil {
		fmt.Println(err.Error())
	}
	version := binary.LittleEndian.Uint32(versionbuf)
	fmt.Println(version)
	// var version int32
	// vbuf := bytes.NewBuffer(versionbuf)
	// err = binary.Read(vbuf, binary.LittleEndian, &version)
	// if err != nil {
	// 	fmt.Println("error reading version to little endian: ", err)
	// }
	// fmt.Println(version)
}

// transaction input
type TxIn struct {
	prevTxId  [32]byte // hash of previous referenced transaction
	prevTxIdx [4]byte  // index of output from referenced transaction
	scriptSig []byte
	sequence  [4]byte
}

func newTxIn(prevTx [32]byte, prevTxIdx [4]byte, scriptSig []byte, sequence [4]byte) *TxIn {
	var script []byte
	if scriptSig == nil {
		script = []byte{}
	} else {
		script = scriptSig
	}

	return &TxIn{prevTxId: prevTx, prevTxIdx: prevTxIdx, scriptSig: script, sequence: sequence}
}

type TxOut struct {
	value        int // amount in satoshis being transferred
	scriptPubKey []byte
}
