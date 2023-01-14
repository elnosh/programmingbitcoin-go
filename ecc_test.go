package main

import (
	"math"
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
		exp  int
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
	a := newPoint(3, -7, 5, 7)
	b := newPoint(18, 77, 5, 7)

	c := newPoint(2, 5, 5, 7)
	d := newPoint(2, 5, 5, 7)

	cases := []struct {
		e1   Point
		e2   Point
		want bool
	}{
		{*a, *b, true},
		{*c, *d, false},
	}

	for _, test := range cases {
		ne := test.e1.ne(test.e2)
		if ne != test.want {
			t.Errorf("expected '%t' but got '%t' instead\n", test.want, ne)
		}
	}
}

func TestAddPoint(t *testing.T) {
	inf := int(math.Inf(0))
	a := newPoint(inf, inf, 5, 7)
	b := newPoint(2, 5, 5, 7)
	c := newPoint(2, -5, 5, 7)
	d := newPoint(3, 7, 5, 7)
	e := newPoint(-1, -1, 5, 7)

	cases := []struct {
		e1   Point
		e2   Point
		want Point
	}{
		{*a, *b, *b},
		{*b, *a, *b},
		{*a, *c, *c},
		{*b, *c, *a},
		{*d, *e, Point{x: 2, y: -5, a: 5, b: 7}},
		{*e, *e, Point{x: 18, y: 77, a: 5, b: 7}},
	}

	for _, test := range cases {
		result := test.e1.add(test.e2)
		if *result != test.want {
			t.Errorf("expected '%v' but got '%v' instead\n", test.want, *result)
		}
	}
}
