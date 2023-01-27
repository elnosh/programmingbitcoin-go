package main

import (
	"fmt"
	"math/big"
)

type FieldElement struct {
	num   *big.Int // single finite field element
	prime *big.Int // field
}

func newFieldElement(num, prime *big.Int) *FieldElement {
	// if num < 0 || num >= prime
	if num.Sign() == -1 || num.Cmp(prime) == 1 {
		fmt.Printf("num %d not in field range 0 to %d\n", num, prime.Sub(prime, big.NewInt(1)))
		return nil
	}
	return &FieldElement{num: num, prime: prime}
}

func (e FieldElement) eq(element FieldElement) bool {
	if e.num.Cmp(element.num) == 0 && e.prime.Cmp(element.prime) == 0 {
		return true
	}
	return false
	//return e.num == element.num && e.prime == element.prime
}

func (e FieldElement) ne(element FieldElement) bool {
	if e.num.Cmp(element.num) != 0 || e.prime.Cmp(element.prime) != 0 {
		return true
	}
	return false
	//return e.num != element.num || e.prime != element.prime
}

func (e FieldElement) add(element FieldElement) *FieldElement {
	if e.prime.Cmp(element.prime) != 0 {
		fmt.Println("cannot add two numbers in different fields")
		return nil
	}

	//num := mod((e.num + element.num), e.prime)
	ec := newFieldElement(new(big.Int).Set(e.num), new(big.Int).Set(e.prime))
	elc := newFieldElement(new(big.Int).Set(element.num), new(big.Int).Set(element.prime))
	sum := ec.num.Add(ec.num, elc.num)
	num := sum.Mod(sum, ec.prime)
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) sub(element FieldElement) *FieldElement {
	if e.prime.Cmp(element.prime) != 0 {
		fmt.Println("cannot subtract two numbers in different fields")
		return nil
	}

	//num := mod((e.num - element.num), e.prime)
	ec := newFieldElement(new(big.Int).Set(e.num), new(big.Int).Set(e.prime))
	elc := newFieldElement(new(big.Int).Set(element.num), new(big.Int).Set(element.prime))
	sub := ec.num.Sub(ec.num, elc.num)
	num := sub.Mod(sub, ec.prime)
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) mul(element FieldElement) *FieldElement {
	if e.prime.Cmp(element.prime) != 0 {
		fmt.Println("cannot multiply two numbers in different fields")
		return nil
	}

	//num := mod((e.num * element.num), e.prime)
	ec := newFieldElement(new(big.Int).Set(e.num), new(big.Int).Set(e.prime))
	elc := newFieldElement(new(big.Int).Set(element.num), new(big.Int).Set(element.prime))
	mul := ec.num.Mul(ec.num, elc.num)
	num := mul.Mod(mul, ec.prime)
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) pow(exponent *big.Int) *FieldElement {
	num := new(big.Int).Exp(e.num, exponent, e.prime)

	//num := mod(int(powres.Int64()), e.prime)
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) div(divisor FieldElement) *FieldElement {
	if e.prime.Cmp(divisor.prime) != 0 {
		fmt.Println("cannot divide two numbers in different fields")
		return nil
	}

	// divpow := divisor.pow(e.prime - 2)
	// num := mod((e.mul(*divpow).num), e.prime)

	temp := new(big.Int).Set(e.prime)
	divpow := divisor.pow(temp.Sub(e.prime, big.NewInt(2)))
	divres := e.mul(*divpow)
	num := divpow.num.Mod(divres.num, e.prime)

	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) repr() {
	fmt.Printf("FieldElement_%d (%d)\n", e.prime, e.num)
}

func isInf(e FieldElement) bool {
	if e.num == nil && e.prime == nil {
		return true
	}
	return false
}

type Point struct {
	x FieldElement
	y FieldElement
	a FieldElement
	b FieldElement
}

func newPoint(x, y, a, b FieldElement) *Point {
	p := &Point{x: x, y: y, a: a, b: b}

	if isInf(x) && isInf(y) {
		var infelement FieldElement
		return &Point{x: infelement, y: infelement, a: a, b: b}
	}

	squarey := y.pow(big.NewInt(2))
	cubex := x.pow(big.NewInt(3))
	rights := cubex.add(*a.mul(x)).add(b)

	if squarey.ne(*rights) {
		fmt.Printf("(%d, %d) is not in the curve\n", x.num, y.num)
		return nil
	}

	return p
}

func (p Point) eq(point Point) bool {
	if p.x.eq(point.x) && p.y.eq(point.y) && p.a.eq(point.a) && p.b.eq(point.b) {
		return true
	}
	return false
}

func (p Point) ne(point Point) bool {
	if p.x.ne(point.x) || p.y.ne(point.y) || p.a.ne(point.a) || p.b.ne(point.b) {
		return true
	}
	return false
}

func (p Point) add(point Point) *Point {
	if p.a.ne(point.a) || p.b.ne(point.b) {
		fmt.Printf("Points %v, %v are not on the same curve\n", p, point)
		return nil
	}

	if isInf(p.x) {
		return &point
	}
	if isInf(point.x) {
		return &p
	}

	if p.x.eq(point.x) && p.y.ne(point.y) {
		var infelement FieldElement
		return newPoint(infelement, infelement, p.a, p.b)
	}

	if p.eq(point) && p.y.num.Sign() == 0 {
		var infelement FieldElement
		return newPoint(infelement, infelement, p.a, p.b)
	}

	if p.x.ne(point.x) {
		// (y2 - y1) / (x2 - x1)
		slope := point.y.sub(p.y).div(*point.x.sub(p.x))

		// x3 = slope^2 - x1 - x2
		x := slope.pow(big.NewInt(2)).sub(p.x).sub(point.x)

		// y3 = slope(x1 - x3) - y1
		y := slope.mul(*p.x.sub(*x)).sub(p.y)

		return &Point{x: *x, y: *y, a: p.a, b: p.b}
	}

	if p.eq(point) {
		three := newFieldElement(big.NewInt(3), p.x.prime)
		two := newFieldElement(big.NewInt(2), p.x.prime)

		// (3x1^2 + a) / (2y1)
		slope := p.x.pow(big.NewInt(2)).mul(*three).add(p.a).div(*two.mul(p.y))

		// slope^2 - 2x1
		x := slope.pow(big.NewInt(2)).sub(p.x).sub(point.x)

		// slope(x1 - x3) - y1
		y := slope.mul(*p.x.sub(*x)).sub(p.y)

		return &Point{x: *x, y: *y, a: p.a, b: p.b}
	}

	return nil
}

func (p Point) rmul(num *big.Int) *Point {
	current := &p
	coef := num
	var infelement FieldElement
	result := newPoint(infelement, infelement, p.a, p.b)

	numlen := num.BitLen()
	for i := 0; i < numlen; i++ {
		temp := new(big.Int).Set(coef)
		coefand1 := temp.And(coef, big.NewInt(1))
		if coefand1.Sign() != 0 {
			result = result.add(*current)
		}
		current = current.add(*current)
		coef.Rsh(coef, 1)
		//fmt.Printf("current x, y = (%v, %v)\n", current.x.num, current.y.num)

		// if (coef & 1) != 0 {
		// 	result = result.add(*current)
		// }
		// current = current.add(*current)
		//coef = coef >> 1
	}

	return result
}

func (p Point) repr() {
	if isInf(p.x) && isInf(p.y) {
		fmt.Printf("Point(infinity, infinity)_%d_%d FieldElement(%d)\n", p.a.num, p.b.num, p.a.prime)
		return
	}
	fmt.Printf("Point(%d, %d)_%d_%d FieldElement(%d)\n", p.x.num, p.y.num, p.a.num, p.b.num, p.a.prime)
}
