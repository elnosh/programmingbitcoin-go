package main

func main() {
	x := fromHex("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798")
	y := fromHex("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")

	g := newS256Point(x, y)

	ng := g.rmulS256()
	ng.repr()
}

//func main() {
//	twopow256 := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))
//	twopow32 := new(big.Int).Exp(big.NewInt(2), big.NewInt(32), big.NewInt(0))

//	sub := twopow256.Sub(twopow256, twopow32)
//	p := sub.Sub(sub, big.NewInt(977))

//	a := newFieldElement(big.NewInt(0), p)
//	b := newFieldElement(big.NewInt(7), p)

//	// gxs := "0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
//	// gys := "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"

//	gx := fromHex("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798")
//	gy := fromHex("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
//	n := fromHex("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141")

//	x := newFieldElement(gx, p)
//	y := newFieldElement(gy, p)

//	g := newPoint(*x, *y, *a, *b)
//	if g != nil {
//		fmt.Println("g is on the curve")
//		//fmt.Printf("p.x = %v, p.y = %v", g.x.num, g.y.num)
//	}

//	ng := g.rmul(n)
//	ng.repr()
//}
