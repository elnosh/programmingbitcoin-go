package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type NetworkEnvelope struct {
	magic   [4]byte
	command [12]byte
	payload []byte
}

var (
	MAINNET_NETWORK_MAGIC = [4]byte{0xf9, 0xbe, 0xb4, 0xd9}
	TESTNET_NETWORK_MAGIC = [4]byte{0x0b, 0x11, 0x09, 0x07}
)

func newNetworkEnvelope(command [12]byte, payload []byte, testnet bool) *NetworkEnvelope {
	if testnet {
		return &NetworkEnvelope{magic: TESTNET_NETWORK_MAGIC, command: command, payload: payload}
	}
	return &NetworkEnvelope{magic: MAINNET_NETWORK_MAGIC, command: command, payload: payload}
}

func parseNetworkEnvelope(b []byte, testnet bool) *NetworkEnvelope {
	netbuffer := bytes.NewBuffer(b)

	var buf [4]byte
	_, err := netbuffer.Read(buf[:])
	if err != nil {
		fmt.Println("error getting network magic: ", err)
	}
	if testnet {
		if buf != TESTNET_NETWORK_MAGIC {
			fmt.Println("Testnet network bytes do not match")
			return nil
		}
	} else {
		if buf != MAINNET_NETWORK_MAGIC {
			fmt.Println("Network bytes do not match")
			return nil
		}
	}

	var command [12]byte
	_, err = netbuffer.Read(command[:])
	if err != nil {
		fmt.Println("error getting command: ", err)
	}

	// next 4 bytes to read payload length
	_, err = netbuffer.Read(buf[:])
	if err != nil {
		fmt.Println("error getting payload length: ", err)
	}
	payloadLength := binary.LittleEndian.Uint32(buf[:])

	// next 4 bytes are payload checksum
	_, err = netbuffer.Read(buf[:])
	if err != nil {
		fmt.Println("error getting payload checksum: ", err)
	}

	payload := make([]byte, payloadLength)
	_, err = netbuffer.Read(payload)
	if err != nil {
		fmt.Println("error getting payload: ", err)
	}

	payloadHash := hash256(payload)
	if bytes.Compare(payloadHash[:4], buf[:]) != 0 {
		fmt.Println("payload checksum does not match")
		return nil
	}

	return newNetworkEnvelope(command, payload, false)
}

func (n NetworkEnvelope) serialize() []byte {
	payloadLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(payloadLength, uint32(len(n.payload)))

	// first 4 bytes of hash is the checksum
	hash := hash256(n.payload)

	return bytes.Join([][]byte{n.magic[:], n.command[:], payloadLength, hash[:4], n.payload}, []byte{})
}
