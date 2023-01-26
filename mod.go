package main

func mod(d, m int) int {
	// modbig := new(big.Int).Mod(big.NewInt(d), big.NewInt(m)).Int64()
	// return modbig
	return (d%m + m) % m

	// or
	// bx, by := big.NewInt(d), big.NewInt(m)
	// return new(big.Int).Mod(bx, by).Int64()
}

// func mod(d, m float64) float64 {
// 	return math.Mod(d, m)
// }
