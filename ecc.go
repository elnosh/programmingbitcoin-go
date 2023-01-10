package main

import (
	"fmt"
	"math/big"
)

type FieldElement struct {
	num   int // single finite field element
	prime int // field
}

func newFieldElement(num, prime int) *FieldElement {
	if num < 0 || num >= prime {
		fmt.Printf("num %d not in field range 0 to %d\n", num, prime-1)
		return nil
	}
	return &FieldElement{num: num, prime: prime}
}

func (e FieldElement) eq(element FieldElement) bool {
	return e.num == element.num && e.prime == element.prime
}

func (e FieldElement) ne(element FieldElement) bool {
	return e.num != element.num || e.prime != element.prime
}

func (e FieldElement) add(element FieldElement) *FieldElement {
	if e.prime != element.prime {
		fmt.Println("cannot add two numbers in different fields")
		return nil
	}

	num := mod((e.num + element.num), e.prime)
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) sub(element FieldElement) *FieldElement {
	if e.prime != element.prime {
		fmt.Println("cannot subtract two numbers in different fields")
		return nil
	}

	num := mod((e.num - element.num), e.prime)
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) mul(element FieldElement) *FieldElement {
	if e.prime != element.prime {
		fmt.Println("cannot multiply two numbers in different fields")
		return nil
	}

	num := mod((e.num * element.num), e.prime)
	return &FieldElement{num: num, prime: e.prime}
}

// func (e FieldElement) pow(exponent int) *FieldElement {
// 	powint := int(math.Pow(float64(e.num), float64(exponent)))

// 	num := mod(powint, e.prime)
// 	return &FieldElement{num: num, prime: e.prime}
// }

func (e FieldElement) pow(exponent int) *FieldElement {
	powres := new(big.Int).Exp(big.NewInt(int64(e.num)), big.NewInt(int64(exponent)), big.NewInt(int64(e.prime)))

	num := mod(int(powres.Int64()), e.prime)
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) div(divisor FieldElement) *FieldElement {
	if e.prime != divisor.prime {
		fmt.Println("cannot divide two numbers in different fields")
		return nil
	}

	divpow := divisor.pow(e.prime - 2)
	num := mod((e.mul(*divpow).num), e.prime)
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) repr() {
	fmt.Printf("FieldElement_%d (%d)\n", e.prime, e.num)
}
