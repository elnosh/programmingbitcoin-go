package main

func mod(d, m int) int {
	return (d%m + m) % m

	// or
	// bx, by := big.NewInt(d), big.NewInt(m)
	// return new(big.Int).Mod(bx, by).Int64()
}
