package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
)

var (
	twopow256 *big.Int = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))
	twopow32  *big.Int = new(big.Int).Exp(big.NewInt(2), big.NewInt(32), big.NewInt(0))
	sub       *big.Int = twopow256.Sub(twopow256, twopow32)
	prime256  *big.Int = sub.Sub(sub, big.NewInt(977))
	n         *big.Int = fromHex("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141")
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

func newS256FieldElement(num *big.Int) *FieldElement {
	return newFieldElement(num, prime256)
}

func (e FieldElement) eq(element FieldElement) bool {
	if e.num.Cmp(element.num) == 0 && e.prime.Cmp(element.prime) == 0 {
		return true
	}
	return false
}

func (e FieldElement) ne(element FieldElement) bool {
	if e.num.Cmp(element.num) != 0 || e.prime.Cmp(element.prime) != 0 {
		return true
	}
	return false
}

func (e FieldElement) add(element FieldElement) *FieldElement {
	if e.prime.Cmp(element.prime) != 0 {
		fmt.Println("cannot add two numbers in different fields")
		return nil
	}

	// (e.num + element.num) % e.prime
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

	// (e.num - element.num) % e.prime
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

	// (e.num * element.num) % e.prime
	ec := newFieldElement(new(big.Int).Set(e.num), new(big.Int).Set(e.prime))
	elc := newFieldElement(new(big.Int).Set(element.num), new(big.Int).Set(element.prime))
	mul := ec.num.Mul(ec.num, elc.num)
	num := mul.Mod(mul, ec.prime)
	return &FieldElement{num: num, prime: e.prime}
}

func (e FieldElement) pow(exponent *big.Int) *FieldElement {
	// (e.num ** exponent) % e.prime
	num := new(big.Int).Exp(e.num, exponent, e.prime)
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
	fmt.Printf("FieldElement_%f (%f)\n", e.prime, e.num)
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

func newS256Point(x, y *big.Int) *Point {
	a := newS256FieldElement(big.NewInt(0))
	b := newS256FieldElement(big.NewInt(7))
	xp := newS256FieldElement(x)
	yp := newS256FieldElement(y)
	return newPoint(*xp, *yp, *a, *b)
}

var (
	gx *big.Int = fromHex("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798")
	gy *big.Int = fromHex("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
	g  *Point   = newS256Point(gx, gy)
)

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

func (p Point) rmul(coefficient *big.Int) *Point {
	current := &p
	coef := new(big.Int).Set(coefficient)
	var infelement FieldElement
	result := newPoint(infelement, infelement, p.a, p.b)

	numlen := coefficient.BitLen()
	for i := 0; i < numlen; i++ {
		temp := new(big.Int).Set(coef)
		coefand1 := temp.And(coef, big.NewInt(1))
		// if (coef & 1) != 0 {
		if coefand1.Sign() != 0 {
			result = result.add(*current)
		}
		current = current.add(*current)
		//coef = coef >> 1
		coef.Rsh(coef, 1)
	}

	return result
}

func (p Point) rmulS256(coefficient *big.Int) *Point {
	coefc := new(big.Int).Set(coefficient)
	coefc.Mod(coefc, n)
	return p.rmul(coefc)
}

func (p Point) verifySignature(s Signature, z *big.Int) bool {
	nc := new(big.Int).Set(n)
	s_inv := new(big.Int).Exp(s.s, nc.Sub(nc, big.NewInt(2)), n)

	zc := new(big.Int).Set(z)
	rc := new(big.Int).Set(s.r)

	umul := zc.Mul(zc, s_inv)
	u := umul.Mod(umul, n)

	vmul := rc.Mul(rc, s_inv)
	v := vmul.Mod(vmul, n)

	uG := g.rmulS256(u)
	vP := p.rmulS256(v)

	sum := uG.add(*vP)
	return s.r.Cmp(sum.x.num) == 0
}

func (p Point) sec(compressed bool) []byte {
	prefixbuf := make([]byte, 1)
	xbuf := make([]byte, 32)

	xbuf = p.x.num.FillBytes(xbuf)
	if compressed {
		yc := new(big.Int).Set(p.y.num)
		yc.Mod(yc, big.NewInt(2))
		// if y is even - prefix 02. Else prefix 03
		if yc.Sign() == 0 {
			prefixbuf = big.NewInt(2).FillBytes(prefixbuf)
		} else {
			prefixbuf = big.NewInt(3).FillBytes(prefixbuf)
		}
	} else {
		prefixbuf = big.NewInt(4).FillBytes(prefixbuf)
		ybuf := make([]byte, 32)
		ybuf = p.y.num.FillBytes(ybuf)
		return bytes.Join([][]byte{prefixbuf, xbuf, ybuf}, []byte{})
	}

	return bytes.Join([][]byte{prefixbuf, xbuf}, []byte{})
}

func (p Point) parse(sec_arr []byte) *Point {
	prefix := int(sec_arr[0])
	if prefix == 4 {
		x := new(big.Int).SetBytes(sec_arr[1:33])
		y := new(big.Int).SetBytes(sec_arr[33:])
		return newS256Point(x, y)
	}

	x := new(big.Int).SetBytes(sec_arr[1:])
	isEven := prefix == 2

	// y^2 = x^3 + 7
	powr := new(big.Int).Set(x).Exp(x, big.NewInt(3), nil)
	right := powr.Add(powr, big.NewInt(7))
	left := sqrt(right)

	var even_left, odd_left *big.Int
	if new(big.Int).Set(left).Mod(left, big.NewInt(2)).Sign() == 0 {
		even_left = left
		odd_left = new(big.Int).Set(prime256).Sub(prime256, left)
	} else {
		even_left = new(big.Int).Set(prime256).Sub(prime256, left)
		odd_left = left
	}

	if isEven {
		return newS256Point(x, even_left)
	} else {
		return newS256Point(x, odd_left)
	}
}

func sqrt(num *big.Int) *big.Int {
	exp := new(big.Int).Set(prime256).Add(prime256, big.NewInt(1))
	exp.Div(exp, big.NewInt(4))

	result := new(big.Int).Set(num).Exp(num, exp, nil)
	return result
}

func (p Point) repr() string {
	if isInf(p.x) && isInf(p.y) {
		return fmt.Sprintf("Point(infinity, infinity)_%d_%d FieldElement(%d)\n", p.a.num, p.b.num, p.a.prime)
	}
	return fmt.Sprintf("Point(%x, %x)_%d_%d FieldElement(%d)\n", p.x.num, p.y.num, p.a.num, p.b.num, p.a.prime)
}

type Signature struct {
	r *big.Int
	s *big.Int
}

func (s Signature) repr() {
	fmt.Printf("Signature(%d, %d)\n", s.r, s.s)
}

// der encoding
func (s Signature) der() []byte {
	prepfix := []byte{0x00}
	marker := []byte{0x02}

	rbytes := new(big.Int).Set(s.r).Bytes()
	if rbytes[0] >= 0x80 {
		rbytes = bytes.Join([][]byte{prepfix, rbytes}, []byte{})
	}

	rlen := []byte{byte(len(rbytes))}
	result := bytes.Join([][]byte{marker, rlen, rbytes}, []byte{})

	sbytes := new(big.Int).Set(s.s).Bytes()
	if sbytes[0] >= 0x80 {
		sbytes = bytes.Join([][]byte{prepfix, sbytes}, []byte{})
	}
	slen := []byte{byte(len(sbytes))}
	result = bytes.Join([][]byte{result, marker, slen, sbytes}, []byte{})
	marker = []byte{0x30}
	reslen := []byte{byte(len(result))}
	return bytes.Join([][]byte{marker, reslen, result}, []byte{})
}

type PrivateKey struct {
	secret *big.Int
	point  Point // public key
}

func newPrivateKey(secret *big.Int) *PrivateKey {
	publicKey := g.rmulS256(secret)
	return &PrivateKey{secret: secret, point: *publicKey}
}

func (pp PrivateKey) sign(z *big.Int) *Signature {
	zc := new(big.Int).Set(z)

	k, _ := rand.Int(rand.Reader, n)
	r := g.rmulS256(k).x.num
	rc := new(big.Int).Set(r)

	nc := new(big.Int).Set(n)
	k_inv := new(big.Int).Exp(k, nc.Sub(nc, big.NewInt(2)), n)

	re := rc.Mul(rc, pp.secret)
	zre := zc.Add(zc, re)
	zrek := zre.Mul(zre, k_inv)
	s := zrek.Mod(zrek, n)

	return &Signature{r: r, s: s}
}
