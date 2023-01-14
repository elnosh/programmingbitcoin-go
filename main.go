package main

import "fmt"

func main() {
	// a := newFieldElement(12, 97)
	// b := newFieldElement(77, 97)

	// powaresult := a.pow(7)
	// powbresult := b.pow(49)

	// mulresult := powaresult.mul(*powbresult)

	// fmt.Println(mulresult)

	// c := newFieldElement(3, 31)
	// d := newFieldElement(24, 31)

	// divresult := c.div(*d)

	// fmt.Println(divresult)

	// e := newFieldElement(17, 31)
	// fmt.Println(e.pow(-3))

	// f := newFieldElement(4, 31)
	// powf := f.pow(-4)
	// g := newFieldElement(11, 31)
	// fmt.Println(powf.mul(*g))

	p1 := newPoint(2, 5, 5, 7)
	p2 := newPoint(-1, -1, 5, 7)
	//p3 := newPoint(18, 77, 5, 7)
	p4 := newPoint(3, 7, 5, 7)
	p5 := newPoint(3, 7, 5, 7)

	// fmt.Println(p1)
	// fmt.Println(p2)
	//fmt.Println(p3)
	//fmt.Println(p4)

	fmt.Printf("%v + %v = %v\n", p1, p2, p1.add(*p2))
	fmt.Printf("%v + %v = %v\n", p4, p2, p4.add(*p2))

	fmt.Println(p4.eq(*p5))
	fmt.Println(p4.ne(*p2))
}
