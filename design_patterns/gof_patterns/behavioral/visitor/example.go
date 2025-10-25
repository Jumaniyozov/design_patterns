package visitor

import "fmt"

func Example1_ShapeVisitor() {
	fmt.Println("\n=== Example 1: Shape Visitor ===")

	shapes := []Shape{
		&Circle{Radius: 5},
		&Rectangle{Width: 10, Height: 20},
	}

	calculator := &AreaCalculator{}
	for _, shape := range shapes {
		shape.Accept(calculator)
	}
}
