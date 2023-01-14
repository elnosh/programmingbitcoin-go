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
	x FieldElement
	y FieldElement
	a FieldElement
	b FieldElement
}

func newPoint(x, y, a, b FieldElement) *Point {
	p := &Point{x: x, y: y, a: a, b: b}

	if x.num == math.MinInt && y.num == math.MinInt {
		inf := int(math.Inf(x.num))
		infelement := newFieldElement(inf, x.prime)
		return &Point{x: *infelement, y: *infelement, a: a, b: b}
	}

	squarey := y.pow(2)
	cubex := x.pow(3)

	if *squarey != *cubex.add(*a.mul(x)).add(b) {
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

	if p.x.num == math.MinInt {
		return &point
	}

	if point.x.num == math.MinInt {
		return &p
	}

	if p.x == point.x && p.y != point.y {
		inf := int(math.Inf(p.x.num))
		infelement := newFieldElement(inf, p.a.prime)
		return &Point{x: *infelement, y: *infelement, a: p.a, b: p.b}
	}

	if p.eq(point) && p.y.num == 0 {
		inf := int(math.Inf(p.x.num))
		infelement := newFieldElement(inf, p.a.prime)
		return &Point{x: *infelement, y: *infelement, a: p.a, b: p.b}
	}

	if p.x != point.x {
		//slope := (point.y.num - p.y.num) / (point.x.num - p.x.num)
		// x := int(math.Pow(float64(slope), 2)) - p.x.num - point.x.num
		// y := slope*(p.x.num-x) - p.y.num

		slope := point.y.sub(p.y).div(*point.x.sub(p.x))
		x := slope.pow(2).sub(p.x).sub(point.x)
		y := slope.mul(*p.x.sub(*x)).sub(p.y)

		// xelement := newFieldElement(x, p.x.prime)
		// yelement := newFieldElement(y, p.y.prime)

		return &Point{x: *x, y: *y, a: p.a, b: p.b}
	}

	if p == point {
		// slope := (3*int(math.Pow(float64(p.x.num), 2)) + p.a.num) / (2 * p.y.num)
		// x := int(math.Pow(float64(slope), 2)) - (2 * p.x.num)
		// y := slope*(p.x.num-x) - p.y.num
		// xelement := newFieldElement(x, p.x.prime)
		// yelement := newFieldElement(y, p.y.prime)

		three := newFieldElement(3, p.x.prime)
		two := newFieldElement(2, p.x.prime)

		slope := p.x.pow(2).mul(*three).add(p.a).div(*two.mul(p.y))
		x := slope.pow(2).sub(p.x).sub(point.x)
		y := slope.mul(*p.x.sub(*x)).sub(p.y)

		return &Point{x: *x, y: *y, a: p.a, b: p.b}
	}

	return nil
}

func (p Point) repr() {
	fmt.Printf("Point(%d, %d)_%d_%d FieldElement(%d)\n", p.x.num, p.y.num, p.a.num, p.b.num, p.x.prime)
}
