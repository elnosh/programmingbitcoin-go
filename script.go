package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Script struct {
	cmds [][]byte
}

func parse(script []byte) *Script {
	scriptbuf := bytes.NewBuffer(script)

	var cmds [][]byte
	count := 0
	scriptLength, err := readVarint(scriptbuf)
	if err != nil {
		fmt.Println("error reading script length: ", err)
	}

	for count < scriptLength {
		curbuf := make([]byte, 1)
		readNextBytes(scriptbuf, curbuf)

		cur := curbuf[0]
		if cur >= 1 && cur <= 75 {
			n := cur
			element := make([]byte, n)
			readNextBytes(scriptbuf, element)
			cmds = append(cmds, element)
			count += int(n)
		} else if cur == 76 { // 76 == opcode for OP_PUSHDATA1
			// read next byte, curbuf will be length of next element to read
			readNextBytes(scriptbuf, curbuf)

			// create element buffer with length read from previous readNextBytes
			element := make([]byte, curbuf[0])
			readNextBytes(scriptbuf, element)
			// append element read
			cmds = append(cmds, element)
			count += int(curbuf[0])
		} else if cur == 77 { // 77 == opcode for OP_PUSHDATA2
			lengthArr := make([]byte, 2)
			// read next byte, curbuf will be length of next element to read
			readNextBytes(scriptbuf, lengthArr)
			// convert length from little endian byte slice to uint16
			length := binary.LittleEndian.Uint16(lengthArr)

			// create element buffer with length read from previous readNextBytes
			element := make([]byte, length)
			readNextBytes(scriptbuf, element)
			// append element read
			cmds = append(cmds, element)
			count += int(length)
		} else { // else next byte is an opcode
			opcode := curbuf
			cmds = append(cmds, opcode)
		}
	}
	if count != scriptLength {
		fmt.Println("error parsing script")
	}

	return &Script{cmds: cmds}
}

func (sc Script) rawSerialize() []byte {
	var result []byte

	for _, cmd := range sc.cmds {
		length := len(cmd)
		if length == 1 {
			result = append(result, cmd[0])
		} else {
			if length <= 75 {
				// append the length of cmd
				result = append(result, byte(length))
			} else if length > 75 && length < 0x100 { // 0x100 == 256
				// append 76 == opcode for OP_PUSHDATA1
				result = append(result, byte(76))
				result = append(result, byte(length))
			} else if length >= 0x100 && length < 520 { // 0x100 == 256
				// append 77 == opcode for OP_PUSHDATA2
				result = append(result, byte(77))
				lenLittleEndian := make([]byte, 2)
				binary.LittleEndian.PutUint16(lenLittleEndian, uint16(length))
				result = append(result, lenLittleEndian...)
			} else {
				fmt.Println("cmd too long")
			}
			// append element after appending opcode and length
			result = append(result, cmd...)
		}
	}
	return result
}

func (sc Script) serialize() []byte {
	result := sc.rawSerialize()
	resultLen := len(result)
	encodedLen, err := encodeVarint(resultLen)
	if err != nil {
		fmt.Println(err)
	}
	return bytes.Join([][]byte{encodedLen, result}, []byte{})
}

func readNextBytes(rd io.Reader, buf []byte) {
	_, err := rd.Read(buf)
	if err != nil {
		fmt.Println("error reading script: ", err)
	}
}
