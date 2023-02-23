package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"
)

type Script struct {
	cmds [][]byte
}

// combine scripts (scriptPubKey + scriptSig) for evaluation
func (sc Script) combine(script *Script) *Script {
	scriptBytes := make([][]byte, len(sc.cmds)+len(script.cmds))
	count := 0
	for i := len(script.cmds) - 1; i >= 0; i-- {
		scriptBytes[count] = script.cmds[i]
		count++
	}
	for i := len(sc.cmds) - 1; i >= 0; i-- {
		scriptBytes[count] = sc.cmds[i]
		count++
	}
	return &Script{cmds: scriptBytes}
}

func parseScript(script io.Reader) (*Script, error) {
	var cmds [][]byte
	count := 0
	scriptLength, err := readVarint(script)
	if err != nil {
		return nil, err
	}

	for count < scriptLength {
		curbuf := make([]byte, 1)
		readNextBytes(script, curbuf)
		count += 1

		cur := curbuf[0]
		if cur >= 1 && cur <= 75 {
			n := int(cur)
			element := make([]byte, n)
			readNextBytes(script, element)
			cmds = append(cmds, element)
			count += n
		} else if cur == 76 { // 76 == opcode for OP_PUSHDATA1
			// read next byte, curbuf will be length of next element to read
			readNextBytes(script, curbuf)

			// create element buffer with length read from previous readNextBytes
			element := make([]byte, curbuf[0])
			readNextBytes(script, element)
			// append element read
			cmds = append(cmds, element)
			//count += int(curbuf[0])
			count = count + int(curbuf[0]) + 1
		} else if cur == 77 { // 77 == opcode for OP_PUSHDATA2
			lengthArr := make([]byte, 2)
			// read next byte, curbuf will be length of next element to read
			readNextBytes(script, lengthArr)
			// convert length from little endian byte slice to uint16
			length := binary.LittleEndian.Uint16(lengthArr)

			// create element buffer with length read from previous readNextBytes
			element := make([]byte, length)
			readNextBytes(script, element)
			// append element read
			cmds = append(cmds, element)
			count = count + int(length) + 2
		} else { // else next byte is an opcode
			opcode := curbuf
			cmds = append(cmds, opcode)
		}
	}

	if count != scriptLength {
		return nil, fmt.Errorf("parsing script failed")
	}

	return &Script{cmds: cmds}, nil
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

// z - signature
func (sc Script) evaluate(z *big.Int) (bool, error) {
	cmds := make([][]byte, len(sc.cmds))
	copy(cmds, sc.cmds)

	var eval bool
	stack := [][]byte{}
	altStack := [][]byte{}
	items := [][]byte{}

	var cmd []byte

	for len(cmds) > 0 {
		cmd, cmds = pop(cmds)
		opcodeType, ok := isOpcode(cmd)
		if ok {
			cmdByte := byte(cmd[0])
			if opcodeType == "opcode" {
				instruction := opcodesFuncs[cmdByte]
				eval, stack = instruction(stack)
				if !eval {
					return false, fmt.Errorf("bad op: %v", opcodesNames[cmdByte])
				}
			} else if opcodeType == "opConditional" {
				instruction := opcodesConditionals[cmdByte]
				eval, stack, items = instruction(stack, items)
				if !eval {
					return false, fmt.Errorf("bad op: %v", opcodesNames[cmdByte])
				}
			} else if opcodeType == "opStack" {
				instruction := opcodesAltStack[cmdByte]
				eval, stack, altStack = instruction(stack, altStack)
				if !eval {
					return false, fmt.Errorf("bad op: %v", opcodesNames[cmdByte])
				}
			} else { // if not any of previous, then is op signature
				instruction := opcodesSignature[cmdByte]
				eval, stack = instruction(stack, z)
				if !eval {
					return false, fmt.Errorf("bad op: %v", opcodesNames[cmdByte])
				}
			}
		} else { // if cmd is not an opcode, then append data to stack
			stack = append(stack, cmd)
		}
	}

	if len(stack) == 0 || stack[0] == nil {
		return false, errors.New("invalid signature")
	}

	return true, nil
}

func isOpcode(cmd []byte) (string, bool) {
	if len(cmd) == 1 {
		cmdByte := byte(cmd[0])
		_, ok := opcodesFuncs[cmdByte]
		if ok {
			return "opcode", true
		}
		_, ok = opcodesConditionals[cmdByte]
		if ok {
			return "opConditional", true
		}
		_, ok = opcodesAltStack[cmdByte]
		if ok {
			return "opStack", true
		}
		_, ok = opcodesSignature[cmdByte]
		if ok {
			return "opSignature", true
		}
	}
	return "", false
}

func readNextBytes(rd io.Reader, buf []byte) {
	_, err := rd.Read(buf)
	if err != nil {
		fmt.Println("error reading script: ", err)
		return
	}
}
