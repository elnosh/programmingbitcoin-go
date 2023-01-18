package main

import (
	"fmt"
	"math"
	"math/big"
	"math/bits"
)

type FieldElement struct {
	num   float64 // single finite field element
	prime float64 // field
}

func newFieldElement(num, prime float64) *FieldElement {
	if math.IsInf(num, 0) {
		inf := math.Inf(int(num))
		return &FieldElement{num: inf, prime: prime}
	}
	if num < 0 || num >= prime {
		fmt.Printf("num %f not in field range 0 to %f\n", num, prime-1)
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

	num := float64(mod((int(e.num + element.num)), int(e.prime)))
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) sub(element FieldElement) *FieldElement {
	if e.prime != element.prime {
		fmt.Println("cannot subtract two numbers in different fields")
		return nil
	}

	num := float64(mod((int(e.num - element.num)), int(e.prime)))
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) mul(element FieldElement) *FieldElement {
	if e.prime != element.prime {
		fmt.Println("cannot multiply two numbers in different fields")
		return nil
	}

	num := float64(mod((int(e.num * element.num)), int(e.prime)))
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) pow(exponent float64) *FieldElement {
	powres := new(big.Int).Exp(big.NewInt(int64(e.num)), big.NewInt(int64(exponent)), big.NewInt(int64(e.prime)))

	num := float64(mod(int(powres.Int64()), int(e.prime)))
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) div(divisor FieldElement) *FieldElement {
	if e.prime != divisor.prime {
		fmt.Println("cannot divide two numbers in different fields")
		return nil
	}

	divpow := divisor.pow(e.prime - 2)
	num := float64(mod((int(e.mul(*divpow).num)), int(e.prime)))
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) repr() {
	fmt.Printf("FieldElement_%f (%f)\n", e.prime, e.num)
}

type Point struct {
	x FieldElement
	y FieldElement
	a FieldElement
	b FieldElement
}

func newPoint(x, y, a, b FieldElement) *Point {
	p := &Point{x: x, y: y, a: a, b: b}

	if math.IsInf(x.num, 0) && math.IsInf(y.num, 0) {
		inf := math.Inf(int(x.num))
		infelement := newFieldElement(inf, x.prime)
		return &Point{x: *infelement, y: *infelement, a: a, b: b}
	}

	squarey := y.pow(2)
	cubex := x.pow(3)

	if *squarey != *cubex.add(*a.mul(x)).add(b) {
		fmt.Printf("(%f, %f) is not in the curve\n", x, y)
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

	if math.IsInf(p.x.num, 0) {
		return &point
	}

	if math.IsInf(point.x.num, 0) {
		return &p
	}

	if p.x == point.x && p.y != point.y {
		inf := math.Inf(int(p.x.num))
		infelement := newFieldElement(inf, p.a.prime)
		return &Point{x: *infelement, y: *infelement, a: p.a, b: p.b}
	}

	if p.eq(point) && p.y.num == 0 {
		inf := math.Inf(int(p.x.num))
		infelement := newFieldElement(inf, p.a.prime)
		return &Point{x: *infelement, y: *infelement, a: p.a, b: p.b}
	}

	if p.x != point.x {
		// (y2 - y1) / (x2 - x1)
		slope := point.y.sub(p.y).div(*point.x.sub(p.x))

		// x3 = slope^2 - x1 - x2
		x := slope.pow(2).sub(p.x).sub(point.x)

		// y3 = slope(x1 - x3) - y1
		y := slope.mul(*p.x.sub(*x)).sub(p.y)

		return &Point{x: *x, y: *y, a: p.a, b: p.b}
	}

	if p.eq(point) {
		cube := newFieldElement(3, p.x.prime)
		t := newFieldElement(2, p.x.prime)

		// (3x1^2 + a) / (2y1)
		slope := p.x.pow(2).mul(*cube).add(p.a).div(*t.mul(p.y))

		// slope^2 - 2x1
		x := slope.pow(2).sub(p.x).sub(point.x)

		// slope(x1 - x3) - y1
		y := slope.mul(*p.x.sub(*x)).sub(p.y)

		return &Point{x: *x, y: *y, a: p.a, b: p.b}
	}

	return nil
}

func (p Point) rmul(num int) *Point {
	current := &p

	coef := num
	inf := math.Inf(int(0))
	infelement := newFieldElement(inf, p.a.prime)
	result := newPoint(*infelement, *infelement, p.a, p.b)

	for i := 0; i < bits.Len(uint(num)); i++ {
		if (coef & 1) != 0 {
			result = result.add(*current)
		}
		current = current.add(*current)
		coef = coef >> 1
	}

	return result
}

func (p Point) repr() {
	fmt.Printf("Point(%f, %f)_%f_%f FieldElement(%f)\n", p.x.num, p.y.num, p.a.num, p.b.num, p.x.prime)
}
