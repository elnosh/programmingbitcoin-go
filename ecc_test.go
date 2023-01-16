package main

import (
	"testing"
)

func TestNeFieldElement(t *testing.T) {
	a := newFieldElement(2, 31)
	b := newFieldElement(2, 31)
	c := newFieldElement(15, 31)

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
	a := newFieldElement(2, 31)
	b := newFieldElement(15, 31)

	c := newFieldElement(17, 31)
	d := newFieldElement(21, 31)

	e := newFieldElement(44, 57)
	f := newFieldElement(33, 57)

	g := newFieldElement(56, 57)
	h := newFieldElement(52, 57)

	cases := []struct {
		e1   FieldElement
		e2   FieldElement
		want FieldElement
	}{
		{*a, *b, FieldElement{num: 17, prime: 31}},
		{*c, *d, FieldElement{num: 7, prime: 31}},
		{*e, *f, FieldElement{num: 20, prime: 57}},
		{*g, *h, FieldElement{num: 51, prime: 57}},
	}

	for _, test := range cases {
		result := test.e1.add(test.e2)
		if *result != test.want {
			t.Errorf("expected '%v' but got '%v' instead\n", test.want, *result)
		}
	}
}

func TestSubFieldElement(t *testing.T) {
	a := newFieldElement(29, 31)
	b := newFieldElement(4, 31)

	c := newFieldElement(15, 31)
	d := newFieldElement(30, 31)

	e := newFieldElement(9, 57)
	f := newFieldElement(29, 57)

	cases := []struct {
		e1   FieldElement
		e2   FieldElement
		want FieldElement
	}{
		{*a, *b, FieldElement{num: 25, prime: 31}},
		{*c, *d, FieldElement{num: 16, prime: 31}},
		{*e, *f, FieldElement{num: 37, prime: 57}},
	}

	for _, test := range cases {
		result := test.e1.sub(test.e2)
		if *result != test.want {
			t.Errorf("expected '%v' but got '%v' instead\n", test.want, *result)
		}
	}
}

func TestMulFieldElement(t *testing.T) {
	a := newFieldElement(24, 31)
	b := newFieldElement(19, 31)

	c := newFieldElement(95, 97)
	d := newFieldElement(45, 97)

	e := newFieldElement(c.mul(*d).num, 97)
	f := newFieldElement(31, 97)

	g := newFieldElement(5, 31)
	powresult := g.pow(5)
	h := newFieldElement(18, 31)

	cases := []struct {
		e1   FieldElement
		e2   FieldElement
		want FieldElement
	}{
		{*a, *b, FieldElement{num: 22, prime: 31}},
		{*c, *d, FieldElement{num: 7, prime: 97}},
		{*e, *f, FieldElement{num: 23, prime: 97}},
		{*powresult, *h, FieldElement{num: 16, prime: 31}},
	}

	for _, test := range cases {
		result := test.e1.mul(test.e2)
		if *result != test.want {
			t.Errorf("expected '%v' but got '%v' instead\n", test.want, *result)
		}
	}
}

func TestPowFieldElement(t *testing.T) {
	a := newFieldElement(17, 31)
	b := newFieldElement(5, 31)

	cases := []struct {
		e1   FieldElement
		exp  float64
		want FieldElement
	}{
		{*a, 3, FieldElement{num: 15, prime: 31}},
		{*b, 5, FieldElement{num: 25, prime: 31}},
		{*a, -3, FieldElement{num: 29, prime: 31}},
	}

	for _, test := range cases {
		result := test.e1.pow(test.exp)
		if *result != test.want {
			t.Errorf("expected '%v' but got '%v' instead\n", test.want, *result)
		}
	}
}

func TestDivFieldElement(t *testing.T) {
	a := newFieldElement(3, 31)
	b := newFieldElement(24, 31)

	cases := []struct {
		e1   FieldElement
		e2   FieldElement
		want FieldElement
	}{
		{*a, *b, FieldElement{num: 4, prime: 31}},
	}

	for _, test := range cases {
		result := test.e1.div(test.e2)
		if *result != test.want {
			t.Errorf("expected '%v' but got '%v' instead\n", test.want, *result)
		}
	}
}

func TestNePoint(t *testing.T) {
	var prime float64 = 98
	a := newFieldElement(5, prime)
	b := newFieldElement(7, prime)

	ap := newPoint(*newFieldElement(3, prime), *newFieldElement(7, prime), *a, *b)
	bp := newPoint(*newFieldElement(18, prime), *newFieldElement(77, prime), *a, *b)

	cp := newPoint(*newFieldElement(2, prime), *newFieldElement(5, prime), *a, *b)
	dp := newPoint(*newFieldElement(2, prime), *newFieldElement(5, prime), *a, *b)

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
	var prime float64 = 223
	a := newFieldElement(0, prime)
	b := newFieldElement(7, prime)

	cases := [][6]float64{
		{192, 105, 17, 56, 170, 142},
		{170, 142, 60, 139, 220, 181},
		{47, 71, 17, 56, 215, 68},
		{143, 98, 76, 66, 47, 71},
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

		if *p1.add(*p2) != *p3 {
			t.Errorf("expected '%v' but got '%v' instead\n", p3, *p1.add(*p2))
		}
	}
}

func TestOnCurve(t *testing.T) {
	var prime float64 = 223
	a := newFieldElement(0, prime)
	b := newFieldElement(7, prime)

	validPoints := [][2]float64{
		{192, 105},
		{17, 56},
		{1, 193},
	}

	invalidPoints := [][2]float64{
		{200, 119},
		{42, 99},
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
