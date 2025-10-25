package iterator

import "fmt"

func Example1_BookCollection() {
	fmt.Println("\n=== Example 1: Book Collection Iterator ===")

	collection := &BookCollection{}
	collection.AddBook(&Book{Title: "Design Patterns"})
	collection.AddBook(&Book{Title: "Clean Code"})
	collection.AddBook(&Book{Title: "The Pragmatic Programmer"})

	iter := collection.CreateIterator()
	for iter.HasNext() {
		book := iter.Next().(*Book)
		fmt.Println("Book:", book.Title)
	}
}
