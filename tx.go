package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

// transaction
type Tx struct {
	version  uint32
	txIns    []TxIn
	txOuts   []TxOut
	locktime uint32
	testnet  bool
}

// binary hash of legacy serialization
func (tx Tx) id() []byte {
	hash := hash256(tx.serialize())
	return reverse(hash[:])
}

func parseTx(s []byte) *Tx {
	sbuf := bytes.NewBuffer(s)
	buf := make([]byte, 4)
	_, err := sbuf.Read(buf)
	if err != nil {
		fmt.Println("error reading version to little endian: ", err)
	}
	version := binary.LittleEndian.Uint32(buf)

	numInputs, err := readVarint(sbuf)
	if err != nil {
		fmt.Println("error reading input varint: ", err)
	}

	var inputs []TxIn
	for i := 0; i < numInputs; i++ {
		inputs = append(inputs, *parseTxIn(sbuf))
	}

	numOutputs, err := readVarint(sbuf)
	if err != nil {
		fmt.Println("error reading output varint: ", err)
	}

	var outputs []TxOut
	for i := 0; i < numOutputs; i++ {
		outputs = append(outputs, *parseTxOut(sbuf))
	}

	_, err = sbuf.Read(buf)
	if err != nil {
		fmt.Println("error reading locktime to little endian: ", err)
	}
	locktime := binary.LittleEndian.Uint32(buf)

	return &Tx{version: version, txIns: inputs, txOuts: outputs, locktime: locktime}
}

func (tx Tx) serialize() []byte {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, tx.version)

	txInsLen, err := encodeVarint(len(tx.txIns))
	if err != nil {
		fmt.Println("error encoding length of tx inputs: ", err)
	}

	var txIns []byte
	for _, txIn := range tx.txIns {
		txIns = append(txIns, txIn.serialize()...)
	}

	txOutsLen, err := encodeVarint(len(tx.txOuts))
	if err != nil {
		fmt.Println("error encoding length of tx outputs: ", err)
	}

	var txOuts []byte
	for _, txOut := range tx.txOuts {
		txOuts = append(txOuts, txOut.serialize()...)
	}

	locktime := make([]byte, 4)
	binary.LittleEndian.PutUint32(locktime, tx.locktime)

	return bytes.Join([][]byte{version, txInsLen, txIns, txOutsLen, txOuts, locktime}, []byte{})
}

func sigHash() {
}

// fee = sum(inputs) - sum(outputs)
func (tx Tx) fee(testnet bool) uint64 {
	var inputSum, outputSum uint64

	for _, input := range tx.txIns {
		inputSum += input.value(testnet)
	}

	for _, output := range tx.txOuts {
		outputSum += output.value
	}
	return inputSum - outputSum
}

// transaction input
type TxIn struct {
	prevTxId  [32]byte // hash of previous referenced transaction
	prevTxIdx uint32   // index of output from referenced transaction
	scriptSig *Script
	sequence  uint32
}

func newTxIn(prevTx [32]byte, prevTxIdx uint32, scriptSig *Script, sequence uint32) *TxIn {
	var script *Script
	if scriptSig == nil {
		script = &Script{}
	} else {
		script = scriptSig
	}

	return &TxIn{prevTxId: prevTx, prevTxIdx: prevTxIdx, scriptSig: script, sequence: sequence}
}

func parseTxIn(txHex io.Reader) *TxIn {
	tx := make([]byte, 32)
	_, err := txHex.Read(tx)
	if err != nil {
		fmt.Println("error parsing tx input: ", err)
		return nil
	}

	// reversing because incoming prev tx hash is in little endian
	prevTx := reversePrevTxInId(tx)

	txIdxbuf := make([]byte, 4)
	_, err = txHex.Read(txIdxbuf)
	if err != nil {
		fmt.Println("error parsing tx input: ", err)
		return nil
	}
	txIdx := binary.LittleEndian.Uint32(txIdxbuf)

	scriptSig, err := parseScript(txHex)
	if err != nil {
		fmt.Println("error parsing scriptSig: ", err)
	}

	sequencebuf := make([]byte, 4)
	_, err = txHex.Read(sequencebuf)
	if err != nil {
		fmt.Println("error parsing tx input: ", err)
		return nil
	}
	sequence := binary.LittleEndian.Uint32(sequencebuf)

	return &TxIn{prevTxId: prevTx, prevTxIdx: txIdx, scriptSig: scriptSig, sequence: sequence}
}

func reversePrevTxInId(prevTx []byte) [32]byte {
	var reversed [32]byte
	counter := 31
	for i := 0; i < 32; i++ {
		reversed[i] = prevTx[counter]
		counter--
	}
	return reversed
}

func (tx TxIn) serialize() []byte {
	prevTxId := reversePrevTxInId(tx.prevTxId[:])

	prevTxIdx := make([]byte, 4)
	binary.LittleEndian.PutUint32(prevTxIdx, tx.prevTxIdx)

	scriptSig := tx.scriptSig.serialize()

	sequence := make([]byte, 4)
	binary.LittleEndian.PutUint32(sequence, tx.sequence)

	return bytes.Join([][]byte{prevTxId[:], prevTxIdx, scriptSig, sequence}, []byte{})
}

func (tx TxIn) fetchTx(testnet bool) *Tx {
	t, err := fetch(hex.EncodeToString(tx.prevTxId[:]), testnet)
	if err != nil {
		fmt.Println("error fetching transaction: ", err)
	}
	return t
}

func (tx TxIn) value(testnet bool) uint64 {
	t := tx.fetchTx(testnet)
	return t.txOuts[tx.prevTxIdx].value
}

func (tx TxIn) scriptPubKey(testnet bool) *Script {
	t := tx.fetchTx(testnet)
	return t.txOuts[tx.prevTxIdx].scriptPubKey
}

type TxOut struct {
	value        uint64  // amount in satoshis being transferred
	scriptPubKey *Script // locking script
}

func parseTxOut(txHex io.Reader) *TxOut {
	// parse amount (# is in satoshis) - amount is in little endian stored in 8 bytes
	amountbuf := make([]byte, 8)
	_, err := txHex.Read(amountbuf)
	if err != nil {
		fmt.Println("error parsing tx output: ", err)
		return nil
	}
	amount := binary.LittleEndian.Uint64(amountbuf)

	scriptPubKey, err := parseScript(txHex)
	if err != nil {
		fmt.Println("error parsing scriptPubKey: ", err)
	}

	return &TxOut{value: amount, scriptPubKey: scriptPubKey}
}

func (tx TxOut) serialize() []byte {
	amount := make([]byte, 8)
	binary.LittleEndian.PutUint64(amount, tx.value)

	script := tx.scriptPubKey.serialize()
	return bytes.Join([][]byte{amount, script}, []byte{})
}

// type TxFetcher struct {
// 	cache map[string]*Tx
// }

var txCache map[string]*Tx

func fetch(txId string, testnet bool) (*Tx, error) {
	// get correct url
	url := "https://blockstream.info/api/"
	if testnet {
		url = "https://blockstream.info/testnet/api/"
	}

	// if tx is not in cache, fetch it
	_, ok := txCache[txId]
	if !ok {
		url += "tx/" + txId + "/hex"
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var tx *Tx
		if body[4] == 0 {
			body = append(body[:4], body[6:]...)
			tx = parseTx(body)
			binary.LittleEndian.PutUint32(body[len(body)-4:], tx.locktime)
		} else {
			tx = parseTx(body)
		}

		tid := string(tx.id())
		if tid != txId {
			return nil, fmt.Errorf("transaction ids do not match: %v and %v\n", tid, txId)
		}

		txCache[txId] = tx
	}
	txCache[txId].testnet = testnet
	return txCache[txId], nil
}
