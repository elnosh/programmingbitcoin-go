package main

import (
	"math/big"
	"testing"
)

func TestNeFieldElement(t *testing.T) {
	element1 := big.NewInt(2)
	element2 := big.NewInt(15)
	prime := big.NewInt(31)

	a := newFieldElement(element1, prime)
	b := newFieldElement(element1, prime)
	c := newFieldElement(element2, prime)

	cases := []struct {
		e1   FieldElement
		e2   FieldElement
		want bool
	}{
		{*a, *b, false},
		{*a, *c, true},
	}

	for _, test := range cases {
		ne := test.e1.ne(test.e2)
		if ne != test.want {
			t.Errorf("expected '%t' but got '%t' instead\n", test.want, ne)
		}
	}
}

func TestAddFieldElement(t *testing.T) {
	prime := big.NewInt(31)
	prime2 := big.NewInt(57)

	element1 := big.NewInt(2)
	element2 := big.NewInt(15)
	element3 := big.NewInt(17)
	element4 := big.NewInt(21)
	element5 := big.NewInt(44)
	element6 := big.NewInt(33)
	element7 := big.NewInt(56)
	element8 := big.NewInt(52)

	a := newFieldElement(element1, prime)
	b := newFieldElement(element2, prime)

	c := newFieldElement(element3, prime)
	d := newFieldElement(element4, prime)

	e := newFieldElement(element5, prime2)
	f := newFieldElement(element6, prime2)

	g := newFieldElement(element7, prime2)
	h := newFieldElement(element8, prime2)

	cases := []struct {
		e1   FieldElement
		e2   FieldElement
		want FieldElement
	}{
		{*a, *b, FieldElement{num: big.NewInt(17), prime: prime}},
		{*c, *d, FieldElement{num: big.NewInt(7), prime: prime}},
		{*e, *f, FieldElement{num: big.NewInt(20), prime: prime2}},
		//{*e, *f, FieldElement{num: 20, prime: prime2}},
		{*g, *h, FieldElement{num: big.NewInt(51), prime: prime2}},
		//{*g, *h, FieldElement{num: 51, prime: prime2}},
	}

	for _, test := range cases {
		result := test.e1.add(test.e2)
		if result.ne(test.want) {
			t.Errorf("expected '%v' but got '%v' instead\n", test.want.num, result.num)
		}
	}
}

func TestSubFieldElement(t *testing.T) {
	prime := big.NewInt(31)
	prime2 := big.NewInt(57)

	element1 := big.NewInt(29)
	element2 := big.NewInt(4)
	element3 := big.NewInt(15)
	element4 := big.NewInt(30)
	element5 := big.NewInt(9)
	element6 := big.NewInt(29)

	a := newFieldElement(element1, prime)
	b := newFieldElement(element2, prime)

	c := newFieldElement(element3, prime)
	d := newFieldElement(element4, prime)

	e := newFieldElement(element5, prime2)
	f := newFieldElement(element6, prime2)

	cases := []struct {
		e1   FieldElement
		e2   FieldElement
		want FieldElement
	}{
		//{*a, *b, FieldElement{num: 25, prime: 31}},
		{*a, *b, FieldElement{num: big.NewInt(25), prime: prime}},
		//{*c, *d, FieldElement{num: 16, prime: 31}},
		{*c, *d, FieldElement{num: big.NewInt(16), prime: prime}},
		//{*e, *f, FieldElement{num: 37, prime: 57}},
		{*e, *f, FieldElement{num: big.NewInt(37), prime: prime2}},
	}

	for _, test := range cases {
		result := test.e1.sub(test.e2)
		if result.ne(test.want) {
			t.Errorf("expected '%v' but got '%v' instead\n", test.want.num, result.num)
		}
	}
}

func TestMulFieldElement(t *testing.T) {
	prime := big.NewInt(31)
	prime2 := big.NewInt(97)

	element1 := big.NewInt(24)
	element2 := big.NewInt(19)
	element3 := big.NewInt(95)
	element4 := big.NewInt(45)
	element5 := big.NewInt(7)
	element6 := big.NewInt(31)
	element7 := big.NewInt(5)
	element8 := big.NewInt(18)

	a := newFieldElement(element1, prime)
	b := newFieldElement(element2, prime)

	c := newFieldElement(element3, prime2)
	d := newFieldElement(element4, prime2)

	e := newFieldElement(element5, prime2)
	f := newFieldElement(element6, prime2)

	g := newFieldElement(element7, prime)
	powresult := g.pow(big.NewInt(5))
	h := newFieldElement(element8, prime)

	cases := []struct {
		e1   FieldElement
		e2   FieldElement
		want FieldElement
	}{
		//{*a, *b, FieldElement{num: 22, prime: 31}},
		{*a, *b, FieldElement{num: big.NewInt(22), prime: prime}},
		//{*c, *d, FieldElement{num: 7, prime: 97}},
		{*c, *d, FieldElement{num: big.NewInt(7), prime: prime2}},
		//{*e, *f, FieldElement{num: 23, prime: 97}},
		{*e, *f, FieldElement{num: big.NewInt(23), prime: prime2}},
		//{*powresult, *h, FieldElement{num: 16, prime: 31}},
		{*powresult, *h, FieldElement{num: big.NewInt(16), prime: prime}},
	}

	for _, test := range cases {
		result := test.e1.mul(test.e2)
		if result.ne(test.want) {
			t.Errorf("expected '%v' but got '%v' instead\n", test.want.num, result.num)
		}
	}
}

func TestPowFieldElement(t *testing.T) {
	prime := big.NewInt(31)

	element1 := big.NewInt(17)
	element2 := big.NewInt(5)

	a := newFieldElement(element1, prime)
	b := newFieldElement(element2, prime)

	cases := []struct {
		e1   FieldElement
		exp  *big.Int
		want FieldElement
	}{
		{*a, big.NewInt(3), FieldElement{num: big.NewInt(15), prime: prime}},
		{*b, big.NewInt(5), FieldElement{num: big.NewInt(25), prime: prime}},
		{*a, big.NewInt(-3), FieldElement{num: big.NewInt(29), prime: prime}},
	}

	for _, test := range cases {
		result := test.e1.pow(test.exp)
		if result.ne(test.want) {
			t.Errorf("expected '%v' but got '%v' instead\n", test.want.num, result.num)
		}
	}
}

func TestDivFieldElement(t *testing.T) {
	prime := big.NewInt(31)

	element1 := big.NewInt(3)
	element2 := big.NewInt(24)

	a := newFieldElement(element1, prime)
	b := newFieldElement(element2, prime)

	cases := []struct {
		e1   FieldElement
		e2   FieldElement
		want FieldElement
	}{
		//{*a, *b, FieldElement{num: 4, prime: 31}},
		{*a, *b, FieldElement{num: big.NewInt(4), prime: prime}},
	}

	for _, test := range cases {
		result := test.e1.div(test.e2)
		if result.ne(test.want) {
			t.Errorf("expected '%v' but got '%v' instead\n", test.want.num, result.num)
		}
	}
}

func TestNePoint(t *testing.T) {
	prime := big.NewInt(98)
	a := newFieldElement(big.NewInt(5), prime)
	b := newFieldElement(big.NewInt(7), prime)

	x1 := big.NewInt(3)
	y1 := big.NewInt(7)
	x2 := big.NewInt(18)
	y2 := big.NewInt(77)
	x3 := big.NewInt(2)
	y3 := big.NewInt(5)

	ap := newPoint(*newFieldElement(x1, prime), *newFieldElement(y1, prime), *a, *b)
	bp := newPoint(*newFieldElement(x2, prime), *newFieldElement(y2, prime), *a, *b)

	cp := newPoint(*newFieldElement(x3, prime), *newFieldElement(y3, prime), *a, *b)
	dp := newPoint(*newFieldElement(x3, prime), *newFieldElement(y3, prime), *a, *b)

	cases := []struct {
		e1   Point
		e2   Point
		want bool
	}{
		{*ap, *bp, true},
		{*cp, *dp, false},
	}

	for _, test := range cases {
		ne := test.e1.ne(test.e2)
		if ne != test.want {
			t.Errorf("expected '%t' but got '%t' instead\n", test.want, ne)
		}
	}
}

func TestAddPointFiniteField(t *testing.T) {
	prime := big.NewInt(223)
	a := newFieldElement(big.NewInt(0), prime)
	b := newFieldElement(big.NewInt(7), prime)

	x1 := big.NewInt(192)
	y1 := big.NewInt(105)
	x2 := big.NewInt(17)
	y2 := big.NewInt(56)
	x3 := big.NewInt(170)
	y3 := big.NewInt(142)

	x4 := big.NewInt(170)
	y4 := big.NewInt(142)
	x5 := big.NewInt(60)
	y5 := big.NewInt(139)
	x6 := big.NewInt(220)
	y6 := big.NewInt(181)

	x7 := big.NewInt(47)
	y7 := big.NewInt(71)
	x8 := big.NewInt(17)
	y8 := big.NewInt(56)
	x9 := big.NewInt(215)
	y9 := big.NewInt(68)

	x10 := big.NewInt(143)
	y10 := big.NewInt(98)
	x11 := big.NewInt(76)
	y11 := big.NewInt(66)
	x12 := big.NewInt(47)
	y12 := big.NewInt(71)

	cases := [][6]*big.Int{
		{x1, y1, x2, y2, x3, y3},
		{x4, y4, x5, y5, x6, y6},
		{x7, y7, x8, y8, x9, y9},
		{x10, y10, x11, y11, x12, y12},
	}

	for _, test := range cases {
		x1 := newFieldElement(test[0], prime)
		y1 := newFieldElement(test[1], prime)
		p1 := newPoint(*x1, *y1, *a, *b)

		x2 := newFieldElement(test[2], prime)
		y2 := newFieldElement(test[3], prime)
		p2 := newPoint(*x2, *y2, *a, *b)

		x3 := newFieldElement(test[4], prime)
		y3 := newFieldElement(test[5], prime)
		p3 := newPoint(*x3, *y3, *a, *b)

		sum := p1.add(*p2)
		if p3.ne(*sum) {
			t.Errorf("expected '%v' but got '%v' instead\n", p3.x.num, *p1.add(*p2).x.num)
		}
	}
}

func TestOnCurve(t *testing.T) {
	prime := big.NewInt(223)
	a := newFieldElement(big.NewInt(0), prime)
	b := newFieldElement(big.NewInt(7), prime)

	x1 := big.NewInt(192)
	y1 := big.NewInt(105)
	x2 := big.NewInt(17)
	y2 := big.NewInt(56)
	x3 := big.NewInt(1)
	y3 := big.NewInt(193)

	x4 := big.NewInt(200)
	y4 := big.NewInt(119)
	x5 := big.NewInt(42)
	y5 := big.NewInt(99)

	validPoints := [][2]*big.Int{
		{x1, y1},
		{x2, y2},
		{x3, y3},
	}

	invalidPoints := [][2]*big.Int{
		{x4, y4},
		{x5, y5},
	}

	for _, point := range validPoints {
		x := newFieldElement(point[0], prime)
		y := newFieldElement(point[1], prime)
		npoint := newPoint(*x, *y, *a, *b)
		if npoint == nil {
			t.Errorf("point %v should not be nil\n", point)
		}
	}

	for _, point := range invalidPoints {
		x := newFieldElement(point[0], prime)
		y := newFieldElement(point[1], prime)
		npoint := newPoint(*x, *y, *a, *b)
		if npoint != nil {
			t.Errorf("point should be nil for invalid point %v\n", point)
		}
	}

}

//func TestAddPoint(t *testing.T) {
//	prime := 98
//	a := newFieldElement(5, prime)
//	b := newFieldElement(7, prime)

//	// infelement := newFieldElement(int(math.Inf(0)), prime)
//	// ap := newPoint(*infelement, *infelement, *a, *b)
//	// bp := newPoint(*newFieldElement(2, prime), *newFieldElement(5, prime), *a, *b)
//	// cp := newPoint(*newFieldElement(2, prime), *newFieldElement(-5, prime), *a, *b)
//	dp := newPoint(*newFieldElement(3, prime), *newFieldElement(7, prime), *a, *b)
//	ep := newPoint(*newFieldElement(-1, prime), *newFieldElement(-1, prime), *a, *b)

//	cases := []struct {
//		e1   Point
//		e2   Point
//		want Point
//	}{
//		//{*ap, *bp, *bp},
//		//{*bp, *ap, *bp},
//		//{*ap, *cp, *cp},
//		//{*bp, *cp, *ap},
//		{*dp, *ep, Point{x: *newFieldElement(2, prime), y: *newFieldElement(-5, prime), a: *a, b: *b}},
//		{*ep, *ep, Point{x: *newFieldElement(18, prime), y: *newFieldElement(77, prime), a: *a, b: *b}},
//	}

//	for _, test := range cases {
//		result := test.e1.add(test.e2)
//		if *result != test.want {
//			t.Errorf("expected '%v' but got '%v' instead\n", test.want, *result)
//		}
//	}
//}
