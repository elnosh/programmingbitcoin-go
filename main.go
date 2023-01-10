package main

import "fmt"

func main() {
	// a := newFieldElement(12, 97)
	// b := newFieldElement(77, 97)

	// powaresult := a.pow(7)
	// powbresult := b.pow(49)

	// mulresult := powaresult.mul(*powbresult)

	// fmt.Println(mulresult)

	c := newFieldElement(3, 31)
	d := newFieldElement(24, 31)

	divresult := c.div(*d)

	fmt.Println(divresult)

	e := newFieldElement(17, 31)
	fmt.Println(e.pow(-3))

	f := newFieldElement(4, 31)
	powf := f.pow(-4)
	g := newFieldElement(11, 31)
	fmt.Println(powf.mul(*g))
}
