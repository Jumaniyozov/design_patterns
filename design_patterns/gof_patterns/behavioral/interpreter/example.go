package interpreter

import "fmt"

func Example1_Calculator() {
	fmt.Println("\n=== Example 1: Simple Calculator ===")

	// (5 + 10) - 3
	expr := &Subtraction{
		left: &Addition{
			left:  &Number{value: 5},
			right: &Number{value: 10},
		},
		right: &Number{value: 3},
	}

	result := expr.Interpret()
	fmt.Printf("(5 + 10) - 3 = %d\n", result)
}
