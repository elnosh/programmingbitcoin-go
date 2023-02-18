package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/ripemd160"
)

//var opcodeFuncs map[byte]interface{}

var opcodeFuncs map[byte]func([][]byte) bool = map[byte]func([][]byte) bool{
	0x00: opcode0,
	0x4f: opcode1Negate,
	0x51: opcode1,
	0x52: opcode2,
	0x53: opcode3,
	0x54: opcode4,
	0x55: opcode5,
	0x56: opcode6,
	0x57: opcode7,
	0x58: opcode8,
	0x59: opcode9,
	0x5a: opcode10,
	0x5b: opcode11,
	0x5c: opcode12,
	0x5d: opcode13,
	0x5e: opcode14,
	0x5f: opcode15,
	0x60: opcode16,
	0x61: opcodeNop,
	0x69: opcodeVerify,
	0x6a: opcodeReturn,
	0x6d: opcode2Drop,
	0x6e: opcode2Dup,
	0x6f: opcode3Dup,
	0x70: opcode2Over,
	0x71: opcode2Rot,
	//0x72: opcode2Swap, // not implemented
	0x73: opcodeIfDup,
	0x74: opcodeDepth,
	0x75: opcodeDrop,
	0x76: opcodeDup,
	0x77: opcodeNip,
	0x78: opcodeOver,
	0x79: opcodePick,
	//0x7a: opcodeRoll, // not implemented
	//0x7b: opcodeRot, // not implemented
	0x7c: opcodeSwap,
	0x7d: opcodeTuck,
	0x82: opcodeSize,
	0x87: opcodeEqual,
	0x88: opcodeEqualVerify,
	// missing arithmetic opcoded

	// crypto opcodes
	0xa6: opcodeRipemd160,
	0xa7: opcodeSha1,
	0xa8: opcodeSha256,
	0xa9: opcodeHash160,
	0xaa: opcodeHash256,
}

var opcodesConditionals map[byte]func([][]byte, []byte) bool = map[byte]func([][]byte, []byte) bool{
	0x63: opcodeIf,
	0x64: opcodeNotIf,
}

var opcodesAltStack map[byte]func([][]byte, [][]byte) bool = map[byte]func([][]byte, [][]byte) bool{
	0x6b: opcodeToAltStack,
	0x6c: opcodeFromAltStack,
}

func encodeNum(num int) []byte {
	if num == 0 {
		return nil
	}
	absNum := abs(num)
	negative := num < 0
	result := []byte{}

	for absNum > 0 {
		result = append(result, byte(absNum&0xff))
		absNum >>= 8
	}

	if result[len(result)-1]&0x80 != 0 {
		if negative {
			result = append(result, 0x80)
		} else {
			result = append(result, 0)
		}
	} else if negative {
		result[len(result)-1] |= 0x80
	}
	return result
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func decodeNum(element []byte) int {
	if element == nil {
		return 0
	}
	var result int

	negative := true
	bigEndian := reverse(element)
	if bigEndian[0]&0x80 != 0 {
		result = int(bigEndian[0] & 0x7f)
	} else {
		negative = false
		result = int(bigEndian[0])
	}

	for _, val := range bigEndian[1:] {
		result <<= 8
		result += int(val)
	}

	if negative {
		return -result
	}
	return result
}

func reverse(element []byte) []byte {
	reversed := make([]byte, len(element))
	counter := len(element) - 1
	for i := 0; i < len(element); i++ {
		reversed[i] = element[counter]
		counter--
	}
	return reversed
}

func pop(stack [][]byte) []byte {
	top := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	return top
}

func opcode0(stack [][]byte) bool {
	stack = append(stack, encodeNum(0))
	return true
}

func opcode1Negate(stack [][]byte) bool {
	stack = append(stack, encodeNum(-1))
	return true
}

func opcode1(stack [][]byte) bool {
	stack = append(stack, encodeNum(1))
	return true
}

func opcode2(stack [][]byte) bool {
	stack = append(stack, encodeNum(2))
	return true
}

func opcode3(stack [][]byte) bool {
	stack = append(stack, encodeNum(3))
	return true
}

func opcode4(stack [][]byte) bool {
	stack = append(stack, encodeNum(4))
	return true
}

func opcode5(stack [][]byte) bool {
	stack = append(stack, encodeNum(5))
	return true
}

func opcode6(stack [][]byte) bool {
	stack = append(stack, encodeNum(6))
	return true
}

func opcode7(stack [][]byte) bool {
	stack = append(stack, encodeNum(7))
	return true
}

func opcode8(stack [][]byte) bool {
	stack = append(stack, encodeNum(8))
	return true
}

func opcode9(stack [][]byte) bool {
	stack = append(stack, encodeNum(9))
	return true
}

func opcode10(stack [][]byte) bool {
	stack = append(stack, encodeNum(10))
	return true
}

func opcode11(stack [][]byte) bool {
	stack = append(stack, encodeNum(11))
	return true
}

func opcode12(stack [][]byte) bool {
	stack = append(stack, encodeNum(12))
	return true
}

func opcode13(stack [][]byte) bool {
	stack = append(stack, encodeNum(13))
	return true
}

func opcode14(stack [][]byte) bool {
	stack = append(stack, encodeNum(14))
	return true
}

func opcode15(stack [][]byte) bool {
	stack = append(stack, encodeNum(15))
	return true
}

func opcode16(stack [][]byte) bool {
	stack = append(stack, encodeNum(16))
	return true
}

func opcodeNop(stack [][]byte) bool {
	return true
}

func opcodeIf(stack [][]byte, items []byte) bool {
	if len(stack) < 1 {
		return false
	}
	trueItems, falseItems := []byte{}, []byte{}
	currentArr := trueItems
	found := false
	endifsNeeded := 1

	for len(items) > 0 {
		item := items[0]
		items = append(items[:0], items[1:]...)

		if item == 99 || item == 100 {
			endifsNeeded++
			currentArr = append(currentArr, item)
		} else if endifsNeeded == 1 && item == 103 {
			currentArr = falseItems
		} else if item == 104 {
			if endifsNeeded == 1 {
				found = true
				break
			} else {
				endifsNeeded--
				currentArr = append(currentArr, item)
			}
		} else {
			currentArr = append(currentArr, item)
		}
	}
	if !found {
		return false
	}

	// element := stack[len(stack)-1]
	// stack = stack[:len(stack)-1]
	// if decodeNum(element) == 0 {
	// 	items[:0] = falseItems
	// } else {
	// 	items[:0] = trueItems
	// }

	return true
}

func opcodeNotIf(stack [][]byte, items []byte) bool {
	if len(stack) < 1 {
		return false
	}
	trueItems, falseItems := []byte{}, []byte{}
	currentArr := trueItems
	found := false
	endifsNeeded := 1

	for len(items) > 0 {
		item := items[0]
		items = append(items[:0], items[1:]...)

		if item == 99 || item == 100 {
			endifsNeeded++
			currentArr = append(currentArr, item)
		} else if endifsNeeded == 1 && item == 103 {
			currentArr = falseItems
		} else if item == 104 {
			if endifsNeeded == 1 {
				found = true
				break
			} else {
				endifsNeeded--
				currentArr = append(currentArr, item)
			}
		} else {
			currentArr = append(currentArr, item)
		}
	}
	if !found {
		return false
	}

	// element := stack[len(stack)-1]
	// stack = stack[:len(stack)-1]
	// if decodeNum(element) == 0 {
	// 	items[:0] = trueItems
	// } else {
	// 	items[:0] = falseItems
	// }

	return true
}

func opcodeVerify(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	item := pop(stack)
	if decodeNum(item) == 0 {
		return false
	}
	return true
}

func opcodeReturn(stack [][]byte) bool {
	return false
}

func opcodeToAltStack(stack, altStack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	item := pop(stack)
	altStack = append(altStack, item)
	return true
}

func opcodeFromAltStack(stack, altStack [][]byte) bool {
	if len(altStack) < 1 {
		return false
	}
	item := pop(altStack)
	stack = append(stack, item)
	return true
}

func opcode2Drop(stack [][]byte) bool {
	if len(stack) < 2 {
		return false
	}
	stack = stack[:len(stack)-2]
	return true
}

func opcode2Dup(stack [][]byte) bool {
	if len(stack) < 2 {
		return false
	}
	stack = append(stack, stack[len(stack)-2:]...)
	return true
}

func opcode3Dup(stack [][]byte) bool {
	if len(stack) < 3 {
		return false
	}
	stack = append(stack, stack[len(stack)-3:]...)
	return true
}

func opcode2Over(stack [][]byte) bool {
	if len(stack) < 4 {
		return false
	}
	stacklen := len(stack)
	stack = append(stack, stack[stacklen-4:stacklen-2]...)
	return true
}

func opcode2Rot(stack [][]byte) bool {
	if len(stack) < 6 {
		return false
	}
	stacklen := len(stack)
	stack = append(stack, stack[stacklen-6:stacklen-4]...)
	return true
}

// func opcode2Swap(stack [][]byte) bool {
// 	if len(stack) < 4 {
// 		return false
// 	}
// 	stacklen := len(stack)
// 	last2 := stack[stacklen-2:]
// 	concat := bytes.Join([][]byte{stack[stacklen-2:], stack[stacklen-4 : stacklen-2]}, []byte{})
// 	stack[len(stack)-4:] = concat

// 	return true
// }

func opcodeIfDup(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	top := stack[len(stack)-1]
	if decodeNum(top) != 0 {
		stack = append(stack, stack[len(stack)-1])
	}
	return true
}

func opcodeDepth(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	stack = append(stack, encodeNum(len(stack)))
	return true
}

func opcodeDrop(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	stack = stack[:len(stack)-1]
	return true
}

func opcodeDup(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	stack = append(stack, stack[len(stack)-1])
	return true
}

func opcodeNip(stack [][]byte) bool {
	if len(stack) < 2 {
		return false
	}
	stack[len(stack)-2] = stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	return true
}

func opcodeOver(stack [][]byte) bool {
	if len(stack) < 2 {
		return false
	}
	stack = append(stack, stack[len(stack)-2])
	return true
}

func opcodePick(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	n := decodeNum(stack[len(stack)-1])
	if len(stack) < n+1 {
		return false
	}
	stack = append(stack, stack[len(stack)-n-1])
	return true
}

// TODO: func opcodeRoll(stack [][]byte) bool {}

// TODO: func opcodeRot(stack [][]byte) bool {}

func opcodeSwap(stack [][]byte) bool {
	if len(stack) < 2 {
		return false
	}
	stack[len(stack)-1], stack[len(stack)-2] = stack[len(stack)-2], stack[len(stack)-1]
	return true
}

func opcodeTuck(stack [][]byte) bool {
	if len(stack) < 2 {
		return false
	}
	top := stack[len(stack)-1]
	stack = append(stack[:len(stack)-2], stack[len(stack)-2])
	stack[len(stack)-2] = top
	return true
}

func opcodeSize(stack [][]byte) bool {
	if len(stack) < 1 {
		return true
	}
	stack = append(stack, encodeNum(len(stack[len(stack)-1])))
	return true
}

func opcodeEqual(stack [][]byte) bool {
	if len(stack) < 2 {
		return false
	}

	item1 := pop(stack)
	item2 := pop(stack)
	if decodeNum(item1) == decodeNum(item2) {
		return true
	}
	return false
}

func opcodeEqualVerify(stack [][]byte) bool {
	return opcodeEqual(stack) && opcodeVerify(stack)
}

// TODO - all arithmetic opcodes

func opcodeRipemd160(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	item := pop(stack)
	r160 := ripemd160.New()
	_, err := r160.Write(item)
	if err != nil {
		fmt.Printf("error ripemd160 hash: %v\n", err)
	}
	stack = append(stack, r160.Sum(nil))
	return true
}

func opcodeSha1(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	item := pop(stack)
	hash := sha1.Sum(item)
	stack = append(stack, hash[:])
	return true
}

func opcodeSha256(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	item := pop(stack)
	hash := sha256.Sum256(item)
	stack = append(stack, hash[:])
	return true
}

func opcodeHash160(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	item := pop(stack)
	hash := hash160(item)
	stack = append(stack, hash[:])
	return true
}

func opcodeHash256(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	item := pop(stack)
	hash := hash160(item)
	stack = append(stack, hash[:])
	return true
}
