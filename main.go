package main

import (
	"strings"
)

func main() {
	//var prime float64 = 223
	// twopow256 := new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(10)).Int64()
	// twopow32 := new(big.Int).Exp(big.NewInt(2), big.NewInt(32), big.NewInt(0)).Int64()

	// fmt.Printf("2^256 = %v\n", twopow256)
	// fmt.Printf("2^32 = %v\n", twopow32)
	// twopow256 := math.Pow(2, 256)
	// twopow32 := math.Pow(2, 32)

	// p := int(twopow256) - int(twopow32) - 977

	// a := newFieldElement(0, p)
	// b := newFieldElement(7, p)

	// gxs := "0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
	// gys := "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"

	// gx, _ := strconv.ParseUint(hexToInt(gxs), 16, 64)
	// gy, _ := strconv.ParseUint(hexToInt(gys), 16, 64)

	// // x1 := newFieldElement(55066263022277343669578718895168534326250603453777594175500187360389116729240, p)
	// // y1 := newFieldElement(32670510020758816978083085130507043184471273380659243275938904335757337482424, p)

	// x1 := newFieldElement(int(gx), p)
	// y1 := newFieldElement(int(gy), p)

	// p1 := newPoint(*x1, *y1, *a, *b)

	// fmt.Println(p1)
}

func hexToInt(hexString string) string {
	numberStr := strings.Replace(hexString, "0x", "", -1)
	return numberStr
}
