package main

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
