package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
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

// id is 2 sha256 of tx.serialized
func (tx Tx) id() []byte {
	hash := hash256(tx.serialize())
	return reverse(hash[:])
}

func parseTx(s []byte) *Tx {
	sbuf := bytes.NewBuffer(s)
	buf := make([]byte, 4)
	// read first four bytes for version
	_, err := sbuf.Read(buf)
	if err != nil {
		fmt.Println("error reading version to little endian: ", err)
	}
	// version from buf is in little endian
	version := binary.LittleEndian.Uint32(buf)

	// get number of inputs
	numInputs, err := readVarint(sbuf)
	if err != nil {
		fmt.Println("error reading input varint: ", err)
	}

	// parse inputs and append them to input list
	var inputs []TxIn
	for i := 0; i < numInputs; i++ {
		inputs = append(inputs, *parseTxIn(sbuf))
	}

	// get number of outputs
	numOutputs, err := readVarint(sbuf)
	if err != nil {
		fmt.Println("error reading output varint: ", err)
	}

	// parse outputs and append them to output list
	var outputs []TxOut
	for i := 0; i < numOutputs; i++ {
		outputs = append(outputs, *parseTxOut(sbuf))
	}

	// read 4 bytes for locktime
	_, err = sbuf.Read(buf)
	if err != nil {
		fmt.Println("error reading locktime to little endian: ", err)
	}
	// locktime is in little endian
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

func (tx Tx) isCoinbase() bool {
	idx, err := hex.DecodeString("ffffffff")
	if err != nil {
		fmt.Println("error getting previus tx index")
		return false
	}
	prevTx := binary.BigEndian.Uint32(tx.txIns[0].prevTxId[:])
	idxInt := binary.BigEndian.Uint32(idx)
	if len(tx.txIns) == 1 && prevTx == 0 && idxInt == tx.txIns[0].prevTxIdx {
		return true
	}

	return false
}

// gets signature hash
func (tx Tx) sigHash(inputIdx uint32) *big.Int {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, tx.version)

	txInsLen, err := encodeVarint(len(tx.txIns))
	if err != nil {
		fmt.Println("error encoding length of tx inputs: ", err)
	}

	// replace scriptSigs being signed with scriptPubKey of previous tx
	var txInput []byte
	for i, txIn := range tx.txIns {
		if int(inputIdx) == i {
			modifiedTxIn := newTxIn(txIn.prevTxId, txIn.prevTxIdx, txIn.scriptPubKey(tx.testnet), txIn.sequence)
			txInput = append(txInput, modifiedTxIn.serialize()...)
		} else {
			modifiedTxIn := newTxIn(txIn.prevTxId, txIn.prevTxIdx, nil, txIn.sequence)
			txInput = append(txInput, modifiedTxIn.serialize()...)
		}
	}

	txOutsLen, err := encodeVarint(len(tx.txOuts))
	if err != nil {
		fmt.Println("error encoding length of tx outputs: ", err)
	}

	var txOutput []byte
	for _, txOut := range tx.txOuts {
		txOutput = append(txOutput, txOut.serialize()...)
	}

	locktime := make([]byte, 4)
	binary.LittleEndian.PutUint32(locktime, tx.locktime)

	hashType := make([]byte, 4)
	binary.LittleEndian.PutUint32(hashType, SIGHASH_ALL)

	modifiedTxBytes := bytes.Join([][]byte{version, txInsLen, txInput, txOutsLen, txOutput, locktime, hashType}, []byte{})

	signatureHash := hash256(modifiedTxBytes)

	return new(big.Int).SetBytes(signatureHash[:])
}

// needs index of input to sign and signs it with private key passed
func (tx *Tx) signInput(inputIdx uint32, privKey *PrivateKey) bool {
	// get the signature hash (z)
	z := tx.sigHash(inputIdx)

	// sign z with private key
	sig := privKey.sign(z).der()
	hashType := byte(SIGHASH_ALL)
	// hashType := make([]byte, 4)
	// binary.LittleEndian.PutUint32(hashType, SIGHASH_ALL)

	// signature is the der signature + hash type
	sig = append(sig, hashType)

	fmt.Println(hex.EncodeToString(sig))

	sec := privKey.point.sec(true)
	scriptSig := &Script{cmds: [][]byte{sig, sec}}

	tx.txIns[inputIdx].scriptSig = scriptSig

	// verify tx input signed is valid
	return tx.verifyInput(inputIdx)
}

func (tx Tx) verifyInput(inputIdx uint32) bool {
	txIn := tx.txIns[inputIdx]
	script := txIn.scriptSig.combine(txIn.scriptPubKey(tx.testnet))
	z := tx.sigHash(inputIdx)
	valid, err := script.evaluate(z)
	if err != nil {
		fmt.Printf("error evaluating script: %v\n", err)
	}
	return valid
}

func (tx Tx) verifyTransaction() bool {
	// this is not here but while verifying a transaction, it should also
	// check for double spends (check if the tx is in the UTXO set)

	if tx.fee(tx.testnet) < 0 {
		return false
	}
	for i := range tx.txIns {
		if !tx.verifyInput(uint32(i)) {
			return false
		}
	}
	return true
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
	scriptSig *Script  // script to unlock utxo and spend
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
	// read 32 bytes - transactionId of previous tx
	tx := make([]byte, 32)
	_, err := txHex.Read(tx)
	if err != nil {
		fmt.Println("error parsing tx input: ", err)
		return nil
	}

	// reversing because incoming prev tx hash is in little endian
	prevTx := reversePrevTxInId(tx)

	// 4 bytes for index of previous tx - utxo being spent
	txIdxbuf := make([]byte, 4)
	_, err = txHex.Read(txIdxbuf)
	if err != nil {
		fmt.Println("error parsing tx input: ", err)
		return nil
	}
	txIdx := binary.LittleEndian.Uint32(txIdxbuf)

	// parses scriptSig
	scriptSig, err := parseScript(txHex)
	if err != nil {
		fmt.Println("error parsing scriptSig: ", err)
	}

	// 4 bytes for sequence
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
		fmt.Println(err)
	}
	return t
}

// gets amount of utxo being spent
func (tx TxIn) value(testnet bool) uint64 {
	t := tx.fetchTx(testnet)
	return t.txOuts[tx.prevTxIdx].value
}

// get scriptPubKey of the previous tx being referenced in the input
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

var txCache map[string]*Tx = map[string]*Tx{}

// fetch tx
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

		decodedBody, err := hex.DecodeString(string(body))
		if err != nil {
			return nil, fmt.Errorf("error getting transaction: %v", err)
		}

		var tx *Tx
		if decodedBody[4] == 0 {
			decodedBody = bytes.Join([][]byte{decodedBody[:4], decodedBody[6:]}, []byte{})
			tx = parseTx(decodedBody)
			tx.locktime = binary.LittleEndian.Uint32(decodedBody[len(decodedBody)-4:])
		} else {
			tx = parseTx(decodedBody)
		}

		tid := hex.EncodeToString(tx.id())
		if tid != txId {
			return nil, fmt.Errorf("transaction ids do not match: %v and %v", tid, txId)
		}

		txCache[txId] = tx
	}
	txCache[txId].testnet = testnet
	return txCache[txId], nil
}
