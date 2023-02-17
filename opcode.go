package main

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
		item := items[len(items)-1]
		items = items[:len(items)-1]

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

	element := stack[len(stack)-1]
	stack = stack[:len(stack)-1]

	return false
}

func opcodeDup(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	stack = append(stack, stack[len(stack)-1])
	return true
}

func opcodeHash256(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	hash256 := hash256(stack[len(stack)-1])
	stack[len(stack)-1] = hash256[:]
	return true
}

func opcodeHash160(stack [][]byte) bool {
	if len(stack) < 1 {
		return false
	}
	hash160 := hash160(stack[len(stack)-1])
	stack[len(stack)-1] = hash160
	return true
}
