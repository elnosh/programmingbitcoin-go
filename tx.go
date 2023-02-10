package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// transaction
type Tx struct {
	version uint32
	txIns   []TxIn
	//txOuts
	locktime int
	testnet  bool
}

// hex of tx hash
func (tx Tx) id() {}

// binary hash of legacy serialization
func (tx Tx) hash() {}

func parseTx(s []byte) *Tx {
	sbuf := bytes.NewBuffer(s)
	versionbuf := make([]byte, 4)
	_, err := sbuf.Read(versionbuf)
	if err != nil {
		fmt.Println("error reading version to little endian: ", err)
	}
	version := binary.LittleEndian.Uint32(versionbuf)

	numInputs, numBytes, err := readVarint(s)
	if err != nil {
		fmt.Println("error reading input varint: ", err)
	}
	sbuf.Next(numBytes)

	var inputs []TxIn
	for i := 0; i < numInputs; i++ {
		inputs = append(inputs, *parseTxIn(sbuf))
	}

	return &Tx{version: version, txIns: inputs}
}

// transaction input
type TxIn struct {
	prevTxId  [32]byte // hash of previous referenced transaction
	prevTxIdx uint32   // index of output from referenced transaction
	scriptSig []byte
	sequence  uint32
}

func newTxIn(prevTx [32]byte, prevTxIdx uint32, scriptSig []byte, sequence uint32) *TxIn {
	var script []byte
	if scriptSig == nil {
		script = []byte{}
	} else {
		script = scriptSig
	}

	return &TxIn{prevTxId: prevTx, prevTxIdx: prevTxIdx, scriptSig: script, sequence: sequence}
}

func parseTxIn(txHex io.Reader) *TxIn {
	var tx [32]byte
	_, err := txHex.Read(tx[:])
	if err != nil {
		fmt.Println("error parsing tx input: ", err)
		return nil
	}
	// reversing because incoming prev tx hash is in little endian
	prevTx := reversePrevTxInId(tx)

	var txIdxbuf [4]byte
	_, err = txHex.Read(txIdxbuf[:])
	if err != nil {
		fmt.Println("error parsing tx input: ", err)
		return nil
	}
	txIdx := binary.LittleEndian.Uint32(txIdxbuf[:])

	// next: parse scriptSig

	var sequencebuf [32]byte
	_, err = txHex.Read(sequencebuf[:])
	if err != nil {
		fmt.Println("error parsing tx input: ", err)
		return nil
	}
	sequence := binary.LittleEndian.Uint32(sequencebuf[:])

	return &TxIn{prevTxId: prevTx, prevTxIdx: txIdx, scriptSig: nil, sequence: sequence}
}

func reversePrevTxInId(prevTx [32]byte) [32]byte {
	var reversed [32]byte
	counter := 31
	for i := 0; i < 32; i++ {
		reversed[i] = prevTx[counter]
		counter--
	}
	return reversed
}

type TxOut struct {
	value        int // amount in satoshis being transferred
	scriptPubKey []byte
}
