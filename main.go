package main

import (
	"fmt"
)

func main() {
	// hexval := "c7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6"
	// hexbytes, _ := hex.DecodeString(hexval)

	// base58str := base58encode(hexbytes)
	// fmt.Println(base58str)

	privKey := newPrivateKey(fromHex("54321deadbeef"))
	fmt.Println(privKey.wif(true, false))
}

// func main() {
// 	privKey := newPrivateKey(big.NewInt(12345))

// 	z := hash256([]byte("some message"))
// 	zint := new(big.Int).SetBytes(z[:])
// 	signature := privKey.sign(zint)

// 	point := g.rmulS256(privKey.secret)
// 	point.repr()
// 	fmt.Println(zint.Text(16))
// 	fmt.Println(signature.r.Text(16))
// 	fmt.Println(signature.s.Text(16))
// }
