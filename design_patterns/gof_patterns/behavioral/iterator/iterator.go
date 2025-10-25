package iterator

type Iterator interface {
	HasNext() bool
	Next() interface{}
}

type Collection interface {
	CreateIterator() Iterator
}

type Book struct {
	Title string
}

type BookCollection struct {
	books []*Book
}

func (c *BookCollection) CreateIterator() Iterator {
	return &BookIterator{collection: c, index: 0}
}

func (c *BookCollection) AddBook(book *Book) {
	c.books = append(c.books, book)
}

type BookIterator struct {
	collection *BookCollection
	index      int
}

func (i *BookIterator) HasNext() bool {
	return i.index < len(i.collection.books)
}

func (i *BookIterator) Next() interface{} {
	if i.HasNext() {
		book := i.collection.books[i.index]
		i.index++
		return book
	}
	return nil
}
