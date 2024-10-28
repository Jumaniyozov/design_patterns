package main

import "fmt"

// Когда использовать: Когда нужно последовательно получить доступ к элементам коллекции без раскрытия ее внутреннего представления.

// Коллекция
type Collection struct {
	items []int
}

func (c *Collection) createIterator() *Iterator {
	return &Iterator{collection: c}
}

// Итератор
type Iterator struct {
	collection *Collection
	index      int
}

func (i *Iterator) hasNext() bool {
	return i.index < len(i.collection.items)
}

func (i *Iterator) next() int {
	item := i.collection.items[i.index]
	i.index++
	return item
}

func main() {
	collection := &Collection{items: []int{1, 2, 3, 4}}
	iterator := collection.createIterator()
	for iterator.hasNext() {
		fmt.Println(iterator.next())
	}
}
