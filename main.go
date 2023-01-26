package main

import (
	"fmt"
	"math/big"
)

func main() {
	twopow256 := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))
	twopow32 := new(big.Int).Exp(big.NewInt(2), big.NewInt(32), big.NewInt(0))

	sub := twopow256.Sub(twopow256, twopow32)
	p := sub.Sub(sub, big.NewInt(977))

	a := newFieldElement(big.NewInt(0), p)
	b := newFieldElement(big.NewInt(7), p)

	// gxs := "0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
	// gys := "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"

	gx := fromHex("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798")
	gy := fromHex("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")

	x1 := newFieldElement(gx, p)
	y1 := newFieldElement(gy, p)

	p1 := newPoint(*x1, *y1, *a, *b)
	if p1 != nil {
		fmt.Printf("p.x = %v, p.y = %v", p1.x.num, p1.y.num)
	}
}

func fromHex(s string) *big.Int {
	if s == "" {
		return big.NewInt(0)
	}
	r, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("invalid hex: " + s)
	}
	return r
}
