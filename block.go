package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Block struct {
	version       uint32
	previousBlock [32]byte
	merkleRoot    [32]byte
	timestamp     uint32
	bits          uint32
	nonce         uint32
}

func (b Block) id() []byte {
	hash := hash256(b.serialize())
	return reverse(hash[:])
}

func parseBlock(b []byte) *Block {
	blockBuffer := bytes.NewBuffer(b)

	buf := make([]byte, 4)
	_, err := blockBuffer.Read(buf)
	if err != nil {
		fmt.Println("error getting block version: ", err)
	}
	version := binary.LittleEndian.Uint32(buf)

	var prevBlock [32]byte
	_, err = blockBuffer.Read(prevBlock[:])
	if err != nil {
		fmt.Println("error getting previous block id: ", err)
	}
	prevBlock = reverseByteArr32(prevBlock)

	var merkleRoot [32]byte
	_, err = blockBuffer.Read(merkleRoot[:])
	if err != nil {
		fmt.Println("error getting merkleRoot: ", err)
	}
	merkleRoot = reverseByteArr32(merkleRoot)

	_, err = blockBuffer.Read(buf)
	if err != nil {
		fmt.Println("error getting block timestamp: ", err)
	}
	timestamp := binary.LittleEndian.Uint32(buf)

	_, err = blockBuffer.Read(buf)
	if err != nil {
		fmt.Println("error getting block bits: ", err)
	}
	bits := binary.BigEndian.Uint32(buf)

	_, err = blockBuffer.Read(buf)
	if err != nil {
		fmt.Println("error getting block nonce: ", err)
	}
	nonce := binary.BigEndian.Uint32(buf)

	return &Block{version: version, previousBlock: prevBlock, merkleRoot: merkleRoot, timestamp: timestamp, bits: bits, nonce: nonce}
}

func (b Block) serialize() []byte {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, b.version)

	prevBlock := reverseByteArr32(b.previousBlock)
	merkleRoot := reverseByteArr32(b.merkleRoot)

	timestamp := make([]byte, 4)
	binary.LittleEndian.PutUint32(timestamp, b.timestamp)

	bits := make([]byte, 4)
	binary.BigEndian.PutUint32(bits, b.bits)

	nonce := make([]byte, 4)
	binary.BigEndian.PutUint32(nonce, b.nonce)

	return bytes.Join([][]byte{version, prevBlock[:], merkleRoot[:], timestamp, bits, nonce}, []byte{})
}
