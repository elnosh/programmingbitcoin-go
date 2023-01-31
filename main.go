package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	r := fromHex("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6")
	s := fromHex("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec")

	sig := Signature{r: r, s: s}
	der := sig.der()
	fmt.Println(hex.EncodeToString(der))
}

//func main() {
//	pk := fromHex("deadbeef54321")
//	//pk := big.NewInt(5001)
//	//pk := new(big.Int).Exp(big.NewInt(2019), big.NewInt(5), nil)
//	privKey := newPrivateKey(pk)

//	uncompressed := privKey.point.sec(true)
//	uncompressedStr := hex.EncodeToString(uncompressed)
//	fmt.Println(uncompressedStr)
//}

// func main() {
// 	privKey := newPrivateKey(big.NewInt(12345))

// 	z := hash256([]byte("some message"))
// 	signature := privKey.sign(z)

// 	point := g.rmulS256(privKey.secret)
// 	point.repr()
// 	fmt.Println(z.Text(16))
// 	fmt.Println(signature.r.Text(16))
// 	fmt.Println(signature.s.Text(16))
// }

// func main() {
// 	z := fromHex("7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d")
// 	r := fromHex("eff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c")
// 	s := fromHex("c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6")

// 	x := fromHex("887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c")
// 	y := fromHex("61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")

// 	p := newS256Point(x, y)

// 	sig := Signature{r: r, s: s}

// 	fmt.Println(p.verifySignature(sig, z))
// }
