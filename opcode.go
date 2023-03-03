package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

var opcodesFuncs map[byte]func([][]byte) (bool, [][]byte) = map[byte]func([][]byte) (bool, [][]byte){
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

	// missing arithmetic opcodes
	0x93: opcodeAdd,
	0x95: opcodeMul,

	// crypto opcodes
	0xa6: opcodeRipemd160,
	0xa7: opcodeSha1,
	0xa8: opcodeSha256,
	0xa9: opcodeHash160,
	0xaa: opcodeHash256,
}

var opcodesConditionals map[byte]func([][]byte, [][]byte) (bool, [][]byte, [][]byte) = map[byte]func([][]byte, [][]byte) (bool, [][]byte, [][]byte){
	0x63: opcodeIf,
	0x64: opcodeNotIf,
}

var opcodesAltStack map[byte]func([][]byte, [][]byte) (bool, [][]byte, [][]byte) = map[byte]func([][]byte, [][]byte) (bool, [][]byte, [][]byte){
	0x6b: opcodeToAltStack,
	0x6c: opcodeFromAltStack,
}

var opcodesSignature map[byte]func([][]byte, *big.Int) (bool, [][]byte) = map[byte]func([][]byte, *big.Int) (bool, [][]byte){
	0xac: opcodeChecksig,
	0xae: opcodeCheckMultisig,
}

var opcodesNames map[byte]string = map[byte]string{
	0x00: "OP_0",
	0x4c: "OP_PUSHDATA1",
	0x4d: "OP_PUSHDATA2",
	0x4e: "OP_PUSHDATA4",
	0x4f: "OP_1NEGATE",
	0x51: "OP_1",
	0x52: "OP_2",
	0x53: "OP_3",
	0x54: "OP_4",
	0x55: "OP_5",
	0x56: "OP_6",
	0x57: "OP_7",
	0x58: "OP_8",
	0x59: "OP_9",
	0x5a: "OP_10",
	0x5b: "OP_11",
	0x5c: "OP_12",
	0x5d: "OP_13",
	0x5e: "OP_14",
	0x5f: "OP_15",
	0x60: "OP_16",
	0x61: "OP_NOP",
	0x63: "OP_IF",
	0x64: "OP_NOTIF",
	0x69: "OP_VERIFY",
	0x6a: "OP_RETURN",

	// stack
	0x6b: "OP_TOALTSTACK",
	0x6c: "OP_FROMALTSTACK",
	0x6d: "OP_2DROP",
	0x6e: "OP_2DUP",
	0x6f: "OP_3DUP",
	0x70: "OP_2OVER",
	0x71: "OP_2ROT",
	//0x72: "OP_2SWAP", // not implemented
	0x73: "OP_IFDUP",
	0x74: "OP_DEPTH",
	0x75: "OP_DROP",
	0x76: "OP_DUP",
	0x77: "OP_NIP",
	0x78: "OP_OVER",
	0x79: "OP_PICK",
	//0x7a: "OP_ROLL", // not implemented
	//0x7b: "OP_ROT", // not implemented
	0x7c: "OP_SWAP",
	0x7d: "OP_TUCK",

	0x82: "OP_SIZE",
	0x87: "OP_EQUAL",
	0x88: "OP_EQUALVERIFY",

	// missing arithmetic opcodes
	0x93: "OP_ADD",
	0x95: "OP_MUL",

	// crypto opcodes
	0xa6: "OP_RIPEMD160",
	0xa7: "OP_SHA1",
	0xa8: "OP_SHA256",
	0xa9: "OP_HASH160",
	0xaa: "OP_HASH256",

	0xac: "OP_CHECKSIG",
	0xae: "OP_CHECKMULTISIG",
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

func pop(stack [][]byte) ([]byte, [][]byte) {
	top := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	return top, stack
}

func opcode0(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(0))
	return true, stack
}

func opcode1Negate(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(-1))
	return true, stack
}

func opcode1(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(1))
	return true, stack
}

func opcode2(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(2))
	return true, stack
}

func opcode3(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(3))
	return true, stack
}

func opcode4(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(4))
	return true, stack
}

func opcode5(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(5))
	return true, stack
}

func opcode6(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(6))
	return true, stack
}

func opcode7(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(7))
	return true, stack
}

func opcode8(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(8))
	return true, stack
}

func opcode9(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(9))
	return true, stack
}

func opcode10(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(10))
	return true, stack
}

func opcode11(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(11))
	return true, stack
}

func opcode12(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(12))
	return true, stack
}

func opcode13(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(13))
	return true, stack
}

func opcode14(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(14))
	return true, stack
}

func opcode15(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(15))
	return true, stack
}

func opcode16(stack [][]byte) (bool, [][]byte) {
	stack = append(stack, encodeNum(16))
	return true, stack
}

func opcodeNop(stack [][]byte) (bool, [][]byte) {
	return true, stack
}

func opcodeIf(stack [][]byte, items [][]byte) (bool, [][]byte, [][]byte) {
	if len(stack) < 1 {
		return false, stack, items
	}
	trueItems, falseItems := [][]byte{}, [][]byte{}
	currentArr := trueItems
	found := false
	endifsNeeded := 1

	for len(items) > 0 {
		item := items[0]
		items = append(items[:0], items[1:]...)

		ilen := len(item)
		if ilen == 1 && item[0] == 0x63 || item[0] == 0x64 {
			endifsNeeded++
			currentArr = append(currentArr, item)
		} else if ilen == 1 && endifsNeeded == 1 && item[0] == 0x67 { // 0x67 = OP_ELSE
			currentArr = falseItems
		} else if ilen == 1 && item[0] == 0x68 { // 0x68 = OP_ENDIF
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
		return false, stack, items
	}

	item, stack := pop(stack)
	if decodeNum(item) == 0 {
		copy(items[:0], falseItems)
	} else {
		copy(items[:0], trueItems)
	}

	return true, stack, items
}

func opcodeNotIf(stack [][]byte, items [][]byte) (bool, [][]byte, [][]byte) {
	if len(stack) < 1 {
		return false, stack, items
	}
	trueItems, falseItems := [][]byte{}, [][]byte{}
	currentArr := trueItems
	found := false
	endifsNeeded := 1

	for len(items) > 0 {
		item := items[0]
		items = append(items[:0], items[1:]...)

		ilen := len(item)
		if ilen == 1 && item[0] == 0x63 || item[0] == 0x64 {
			endifsNeeded++
			currentArr = append(currentArr, item)
		} else if ilen == 1 && endifsNeeded == 1 && item[0] == 0x67 { // 0x67 = OP_ELSE
			currentArr = falseItems
		} else if ilen == 1 && item[0] == 0x68 { // 0x68 = OP_ENDIF
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
		return false, stack, items
	}

	item, stack := pop(stack)
	if decodeNum(item) == 0 {
		copy(items[:0], trueItems)
	} else {
		copy(items[:0], falseItems)
	}

	return true, stack, items
}

func opcodeVerify(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 1 {
		return false, stack
	}
	item, stack := pop(stack)
	if decodeNum(item) == 0 {
		return false, stack
	}
	return true, stack
}

func opcodeReturn(stack [][]byte) (bool, [][]byte) {
	return false, stack
}

func opcodeToAltStack(stack, altStack [][]byte) (bool, [][]byte, [][]byte) {
	if len(stack) < 1 {
		return false, stack, altStack
	}
	item, stack := pop(stack)
	altStack = append(altStack, item)
	return true, stack, altStack
}

func opcodeFromAltStack(stack, altStack [][]byte) (bool, [][]byte, [][]byte) {
	if len(altStack) < 1 {
		return false, stack, altStack
	}
	item, stack := pop(altStack)
	stack = append(stack, item)
	return true, stack, altStack
}

func opcode2Drop(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 2 {
		return false, stack
	}
	stack = stack[:len(stack)-2]
	return true, stack
}

func opcode2Dup(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 2 {
		return false, stack
	}
	stack = append(stack, stack[len(stack)-2:]...)
	return true, stack
}

func opcode3Dup(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 3 {
		return false, stack
	}
	stack = append(stack, stack[len(stack)-3:]...)
	return true, stack
}

func opcode2Over(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 4 {
		return false, stack
	}
	stacklen := len(stack)
	stack = append(stack, stack[stacklen-4:stacklen-2]...)
	return true, stack
}

func opcode2Rot(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 6 {
		return false, stack
	}
	stacklen := len(stack)
	stack = append(stack, stack[stacklen-6:stacklen-4]...)
	return true, stack
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

func opcodeIfDup(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 1 {
		return false, stack
	}
	top := stack[len(stack)-1]
	if decodeNum(top) != 0 {
		stack = append(stack, stack[len(stack)-1])
	}
	return true, stack
}

func opcodeDepth(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 1 {
		return false, stack
	}
	stack = append(stack, encodeNum(len(stack)))
	return true, stack
}

func opcodeDrop(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 1 {
		return false, stack
	}
	stack = stack[:len(stack)-1]
	return true, stack
}

func opcodeDup(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 1 {
		return false, stack
	}
	stack = append(stack, stack[len(stack)-1])
	return true, stack
}

func opcodeNip(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 2 {
		return false, stack
	}
	stack[len(stack)-2] = stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	return true, stack
}

func opcodeOver(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 2 {
		return false, stack
	}
	stack = append(stack, stack[len(stack)-2])
	return true, stack
}

func opcodePick(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 1 {
		return false, stack
	}
	n := decodeNum(stack[len(stack)-1])
	if len(stack) < n+1 {
		return false, stack
	}
	stack = append(stack, stack[len(stack)-n-1])
	return true, stack
}

// TODO: func opcodeRoll(stack [][]byte) bool {}

// TODO: func opcodeRot(stack [][]byte) bool {}

func opcodeSwap(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 2 {
		return false, stack
	}
	stack[len(stack)-1], stack[len(stack)-2] = stack[len(stack)-2], stack[len(stack)-1]
	return true, stack
}

func opcodeTuck(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 2 {
		return false, stack
	}
	top := stack[len(stack)-1]
	stack = append(stack[:len(stack)-2], stack[len(stack)-2])
	stack[len(stack)-2] = top
	return true, stack
}

func opcodeSize(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 1 {
		return true, stack
	}
	stack = append(stack, encodeNum(len(stack[len(stack)-1])))
	return true, stack
}

func opcodeEqual(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 2 {
		return false, stack
	}

	item1, stack := pop(stack)
	item2, stack := pop(stack)
	if decodeNum(item1) == decodeNum(item2) {
		stack = append(stack, encodeNum(1))
	} else {
		stack = append(stack, encodeNum(0))
	}
	return true, stack
}

func opcodeEqualVerify(stack [][]byte) (bool, [][]byte) {
	equal, stack := opcodeEqual(stack)
	verify, stack := opcodeVerify(stack)
	return (equal && verify), stack
}

func opcodeAdd(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 2 {
		return false, stack
	}

	item1, stack := pop(stack)
	item2, stack := pop(stack)
	i1, i2 := decodeNum(item1), decodeNum(item2)
	stack = append(stack, encodeNum(i1+i2))
	return true, stack
}

func opcodeMul(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 2 {
		return false, stack
	}

	item1, stack := pop(stack)
	item2, stack := pop(stack)
	i1, i2 := decodeNum(item1), decodeNum(item2)
	stack = append(stack, encodeNum(i1*i2))
	return true, stack
}

// TODO - all arithmetic opcodes

func opcodeRipemd160(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 1 {
		return false, stack
	}
	item, stack := pop(stack)
	r160 := ripemd160.New()
	_, err := r160.Write(item)
	if err != nil {
		fmt.Printf("error ripemd160 hash: %v\n", err)
	}
	stack = append(stack, r160.Sum(nil))
	return true, stack
}

func opcodeSha1(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 1 {
		return false, stack
	}
	item, stack := pop(stack)
	hash := sha1.Sum(item)
	stack = append(stack, hash[:])
	return true, stack
}

func opcodeSha256(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 1 {
		return false, stack
	}
	item, stack := pop(stack)
	hash := sha256.Sum256(item)
	stack = append(stack, hash[:])
	return true, stack
}

func opcodeHash160(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 1 {
		return false, stack
	}
	item, stack := pop(stack)
	hash := hash160(item)
	stack = append(stack, hash[:])
	return true, stack
}

func opcodeHash256(stack [][]byte) (bool, [][]byte) {
	if len(stack) < 1 {
		return false, stack
	}
	item, stack := pop(stack)
	hash := hash160(item)
	stack = append(stack, hash[:])
	return true, stack
}

func opcodeChecksig(stack [][]byte, z *big.Int) (bool, [][]byte) {
	if len(stack) < 2 {
		return false, stack
	}
	pubKey, stack := pop(stack)
	pubKeyPoint := parsePubKey(pubKey)

	signature, stack := pop(stack)
	sig, err := parseSignature(signature)
	if err != nil {
		fmt.Printf("invalid signature: %v\n", err)
		return false, stack
	}

	if !pubKeyPoint.verifySignature(*sig, z) {
		stack = append(stack, encodeNum(0))
		return false, stack
	}

	stack = append(stack, encodeNum(1))
	return true, stack
}

func opcodeCheckMultisig(stack [][]byte, z *big.Int) (bool, [][]byte) {
	if len(stack) < 1 {
		return false, stack
	}
	// m-of-n multisig
	nbyte, stack := pop(stack)
	n := decodeNum(nbyte)
	if len(stack) < n+1 {
		return false, stack
	}
	pubKeys := make([]*Point, n)
	var pubKey []byte
	for i := n; i > 0; i-- {
		pubKey, stack = pop(stack)
		pubKeyPoint := parsePubKey(pubKey)
		pubKeys = append(pubKeys, pubKeyPoint)
	}

	mbyte, stack := pop(stack)
	m := decodeNum(mbyte)
	if len(stack) < m+1 {
		return false, stack
	}
	sigs := make([]*Signature, m)
	var sigByte []byte
	for i := 0; i < m; i++ {
		sigByte, stack = pop(stack)
		sig, err := parseSignature(sigByte)
		if err != nil {
			fmt.Printf("error parsing signature in multisig - '%v'\n", err)
		}
		sigs = append(sigs, sig)
	}
	// pop for bug of extra value unused
	_, stack = pop(stack)

	for _, sig := range sigs {
		if len(pubKeys) < 1 {
			stack = append(stack, encodeNum(0))
			return false, stack
		}

		var pubKey *Point
		for len(pubKeys) > 0 {
			pubKey, pubKeys = popKey(pubKeys)
			if pubKey.verifySignature(*sig, z) {
				break
			}
		}
	}

	stack = append(stack, encodeNum(1))
	return true, stack
}

func popKey(arr []*Point) (*Point, []*Point) {
	key := arr[len(arr)-1]
	arr = arr[:len(arr)-1]
	return key, arr
}
