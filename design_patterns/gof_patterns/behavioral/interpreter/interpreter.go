package interpreter

type Expression interface {
	Interpret() int
}

type Number struct {
	value int
}

func (n *Number) Interpret() int {
	return n.value
}

type Addition struct {
	left, right Expression
}

func (a *Addition) Interpret() int {
	return a.left.Interpret() + a.right.Interpret()
}

type Subtraction struct {
	left, right Expression
}

func (s *Subtraction) Interpret() int {
	return s.left.Interpret() - s.right.Interpret()
}
