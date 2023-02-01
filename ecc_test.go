package main

import (
	"crypto/rand"
	"encoding/hex"
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

func TestRmul(t *testing.T) {
	prime := big.NewInt(223)
	a := newFieldElement(big.NewInt(0), prime)
	b := newFieldElement(big.NewInt(7), prime)

	coef := big.NewInt(2)
	coef2 := big.NewInt(4)
	coef3 := big.NewInt(8)
	coef4 := big.NewInt(21)

	x1 := big.NewInt(192)
	y1 := big.NewInt(105)
	x2 := big.NewInt(49)
	y2 := big.NewInt(71)

	x3 := big.NewInt(143)
	y3 := big.NewInt(98)
	x4 := big.NewInt(64)
	y4 := big.NewInt(168)

	x5 := big.NewInt(47)
	y5 := big.NewInt(71)
	x6 := big.NewInt(36)
	y6 := big.NewInt(111)

	x7 := big.NewInt(194)
	y7 := big.NewInt(51)

	x8 := big.NewInt(116)
	y8 := big.NewInt(55)

	cases := [][5]*big.Int{
		{coef, x1, y1, x2, y2},
		{coef, x3, y3, x4, y4},
		{coef, x5, y5, x6, y6},
		{coef2, x5, y5, x7, y7},
		{coef3, x5, y5, x8, y8},
		{coef4, x5, y5, nil, nil},
	}

	for _, test := range cases {
		x1 := newFieldElement(test[1], prime)
		y1 := newFieldElement(test[2], prime)
		p1 := newPoint(*x1, *y1, *a, *b)

		var x2 *FieldElement
		var y2 *FieldElement
		var p2 *Point

		if test[3] == nil {
			var infelement FieldElement
			p2 = newPoint(infelement, infelement, *a, *b)
		} else {
			x2 = newFieldElement(test[3], prime)
			y2 = newFieldElement(test[4], prime)
			p2 = newPoint(*x2, *y2, *a, *b)
		}

		mul := p1.rmul(test[0])
		if p2.ne(*mul) {
			t.Errorf("expected '%v' but got '%v' instead\n", p2.x.num, *p1.rmul(test[0]).x.num)
		}
	}
}

func TestOrder(t *testing.T) {
	if g.rmulS256(n).x.num != nil {
		t.Errorf("expected '%v' but got '%v' instead\n", nil, g.rmulS256(n).x.num)
	}
}

func TestPublicPoint(t *testing.T) {
	secret := big.NewInt(7)
	secret2 := big.NewInt(1485)
	secret3 := new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil)

	exp := new(big.Int).Exp(big.NewInt(2), big.NewInt(240), nil)
	exp2 := new(big.Int).Exp(big.NewInt(2), big.NewInt(31), nil)
	secret4 := exp.Add(exp, exp2)

	points := [][3]*big.Int{
		{secret, fromHex("5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc"),
			fromHex("6aebca40ba255960a3178d6d861a54dba813d0b813fde7b5a5082628087264da")},
		{secret2, fromHex("c982196a7466fbbbb0e27a940b6af926c1a74d5ad07128c82824a11b5398afda"),
			fromHex("7a91f9eae64438afb9ce6448a1c133db2d8fb9254e4546b6f001637d50901f55")},
		{secret3, fromHex("8f68b9d2f63b5f339239c1ad981f162ee88c5678723ea3351b7b444c9ec4c0da"),
			fromHex("662a9f2dba063986de1d90c2b6be215dbbea2cfe95510bfdf23cbf79501fff82")},
		{secret4, fromHex("9577ff57c8234558f293df502ca4f09cbc65a6572c842b39b366f21717945116"),
			fromHex("10b49c67fa9365ad7b90dab070be339a1daf9052373ec30ffae4f72d5e66d053")},
	}

	for _, test := range points {
		point := newS256Point(test[1], test[2])
		pubPoint := g.rmulS256(test[0])

		if pubPoint.ne(*point) {
			t.Errorf("expected '%v' but got '%v' instead\n", pubPoint.repr(), point.repr())
		}
	}
}

func TestVerifySignature(t *testing.T) {
	x1 := fromHex("887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c")
	y1 := fromHex("61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")
	point := newS256Point(x1, y1)

	cases := []struct {
		z    *big.Int
		r    *big.Int
		s    *big.Int
		want bool
	}{
		{fromHex("ec208baa0fc1c19f708a9ca96fdeff3ac3f230bb4a7ba4aede4942ad003c0f60"),
			fromHex("ac8d1c87e51d0d441be8b3dd5b05c8795b48875dffe00b7ffcfac23010d3a395"),
			fromHex("68342ceff8935ededd102dd876ffd6ba72d6a427a3edb13d26eb0781cb423c4"), true},
		{fromHex("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d"),
			fromHex("eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c"),
			fromHex("c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6"), true},
	}

	for _, test := range cases {
		verified := point.verifySignature(Signature{r: test.r, s: test.s}, test.z)
		if verified != test.want {
			t.Errorf("expected '%v' but got '%v' instead\n", test.want, verified)
		}
	}
}

func TestSign(t *testing.T) {
	randk, _ := rand.Int(rand.Reader, n)
	privKey := newPrivateKey(randk)

	randz := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	z, _ := rand.Int(rand.Reader, randz)

	signature := privKey.sign(z)

	verified := privKey.point.verifySignature(*signature, z)
	if verified != true {
		t.Errorf("expected '%v' but got '%v' instead\n", true, verified)
	}
}

func TestSec(t *testing.T) {
	cases := []struct {
		coefficient      *big.Int
		wantUncompressed string
		wantCompressed   string
	}{
		{big.NewInt(997002999),
			"049d5ca49670cbe4c3bfa84c96a8c87df086c6ea6a24ba6b809c9de234496808d56fa15cc7f3d38cda98dee2419f415b7513dde1301f8643cd9245aea7f3f911f9",
			"039d5ca49670cbe4c3bfa84c96a8c87df086c6ea6a24ba6b809c9de234496808d5"},
		{big.NewInt(123),
			"04a598a8030da6d86c6bc7f2f5144ea549d28211ea58faa70ebf4c1e665c1fe9b5204b5d6f84822c307e4b4a7140737aec23fc63b65b35f86a10026dbd2d864e6b",
			"03a598a8030da6d86c6bc7f2f5144ea549d28211ea58faa70ebf4c1e665c1fe9b5"},
		{big.NewInt(42424242),
			"04aee2e7d843f7430097859e2bc603abcc3274ff8169c1a469fee0f20614066f8e21ec53f40efac47ac1c5211b2123527e0e9b57ede790c4da1e72c91fb7da54a3",
			"03aee2e7d843f7430097859e2bc603abcc3274ff8169c1a469fee0f20614066f8e"},
	}

	for _, test := range cases {
		point := g.rmulS256(test.coefficient)
		uncompressed := hex.EncodeToString(point.sec(false))
		if test.wantUncompressed != uncompressed {
			t.Errorf("expected '%v' but got '%v' instead\n", test.wantUncompressed, uncompressed)
		}

		compressed := hex.EncodeToString(point.sec(true))
		if test.wantCompressed != compressed {
			t.Errorf("expected '%v' but got '%v' instead\n", test.wantCompressed, compressed)
		}
	}
}
