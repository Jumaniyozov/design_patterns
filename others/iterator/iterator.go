// Package iterator demonstrates the Iterator pattern.
// It provides a way to access elements sequentially without exposing
// the underlying representation, useful for complex data structures.
package iterator

// TreeNode represents a node in a binary tree
type TreeNode struct {
	Value int
	Left  *TreeNode
	Right *TreeNode
}

// InOrderIterator performs in-order traversal (Left-Root-Right)
func (t *TreeNode) InOrderIterator() <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		t.inOrderTraverse(ch)
	}()
	return ch
}

func (t *TreeNode) inOrderTraverse(ch chan<- int) {
	if t == nil {
		return
	}
	t.Left.inOrderTraverse(ch)
	ch <- t.Value
	t.Right.inOrderTraverse(ch)
}

// PreOrderIterator performs pre-order traversal (Root-Left-Right)
func (t *TreeNode) PreOrderIterator() <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		t.preOrderTraverse(ch)
	}()
	return ch
}

func (t *TreeNode) preOrderTraverse(ch chan<- int) {
	if t == nil {
		return
	}
	ch <- t.Value
	t.Left.preOrderTraverse(ch)
	t.Right.preOrderTraverse(ch)
}

// PostOrderIterator performs post-order traversal (Left-Right-Root)
func (t *TreeNode) PostOrderIterator() <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		t.postOrderTraverse(ch)
	}()
	return ch
}

func (t *TreeNode) postOrderTraverse(ch chan<- int) {
	if t == nil {
		return
	}
	t.Left.postOrderTraverse(ch)
	t.Right.postOrderTraverse(ch)
	ch <- t.Value
}

// RangeIterator demonstrates a custom range iterator
type RangeIterator struct {
	start   int
	end     int
	step    int
	current int
}

// NewRangeIterator creates an iterator over a numeric range
func NewRangeIterator(start, end, step int) *RangeIterator {
	return &RangeIterator{
		start:   start,
		end:     end,
		step:    step,
		current: start,
	}
}

// HasNext checks if there are more elements
func (r *RangeIterator) HasNext() bool {
	if r.step > 0 {
		return r.current < r.end
	}
	return r.current > r.end
}

// Next returns the next element
func (r *RangeIterator) Next() int {
	value := r.current
	r.current += r.step
	return value
}

// Channel creates a channel-based iterator
func (r *RangeIterator) Channel() <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for r.HasNext() {
			ch <- r.Next()
		}
	}()
	return ch
}
