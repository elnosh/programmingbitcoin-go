package main

import "fmt"

func main() {
	var prime float64 = 223
	a := newFieldElement(0, prime)
	b := newFieldElement(7, prime)

	x1 := newFieldElement(15, prime)
	y1 := newFieldElement(86, prime)

	// inf := math.Inf(int(0))
	// infelement := newFieldElement(inf, prime)

	p1 := newPoint(*x1, *y1, *a, *b)
	//p2 := newPoint(*infelement, *infelement, *a, *b)

	// 	product := p1
	// 	count := 1

	// 	for product.ne(*p2) {
	// 		product = product.add(*p1)
	// 		count++
	// 	}

	// 	fmt.Println(product)
	// 	fmt.Println(count)

	product := p1.rmul(1000000000000)
	fmt.Println(product)
}
