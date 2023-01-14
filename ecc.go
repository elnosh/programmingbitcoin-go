package main

import (
	"fmt"
	"math"
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

type Point struct {
	x int
	y int
	a int
	b int
}

func newPoint(x, y, a, b int) *Point {
	p := &Point{x: x, y: y, a: a, b: b}

	if x == math.MinInt && y == math.MinInt {
		inf := int(math.Inf(x))
		return &Point{x: inf, y: inf, a: a, b: b}
	}

	squarey := int(math.Pow(float64(y), float64(2)))
	cubex := int(math.Pow(float64(x), float64(3)))

	if squarey != cubex+(a*x)+b {
		fmt.Printf("(%d, %d) is not in the curve\n", x, y)
		return nil
	}

	return p
}

func (p Point) eq(point Point) bool {
	return p == point
}

func (p Point) ne(point Point) bool {
	return p != point
}

func (p Point) add(point Point) *Point {
	if p.a != point.a || p.b != point.b {
		fmt.Printf("Points %v, %v are not on the same curve\n", p, point)
		return nil
	}

	if p.x == math.MinInt {
		return &point
	}

	if point.x == math.MinInt {
		return &p
	}

	if p.x == point.x && p.y != point.y {
		inf := int(math.Inf(p.x))
		return &Point{x: inf, y: inf, a: p.a, b: p.b}
	}

	if p.eq(point) && p.y == 0 {
		inf := int(math.Inf(p.x))
		return &Point{x: inf, y: inf, a: p.a, b: p.b}
	}

	if p.x != point.x {
		slope := (point.y - p.y) / (point.x - p.x)
		x := int(math.Pow(float64(slope), 2)) - p.x - point.x
		y := slope*(p.x-x) - p.y

		return &Point{x: x, y: y, a: p.a, b: p.b}
	}

	if p == point {
		slope := (3*int(math.Pow(float64(p.x), 2)) + p.a) / (2 * p.y)
		x := int(math.Pow(float64(slope), 2)) - (2 * p.x)
		y := slope*(p.x-x) - p.y

		return &Point{x: x, y: y, a: p.a, b: p.b}
	}

	return nil
}
