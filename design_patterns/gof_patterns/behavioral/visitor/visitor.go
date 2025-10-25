package visitor

import "fmt"

type Visitor interface {
	VisitCircle(*Circle)
	VisitRectangle(*Rectangle)
}

type Shape interface {
	Accept(Visitor)
}

type Circle struct {
	Radius int
}

func (c *Circle) Accept(v Visitor) {
	v.VisitCircle(c)
}

type Rectangle struct {
	Width, Height int
}

func (r *Rectangle) Accept(v Visitor) {
	v.VisitRectangle(r)
}

type AreaCalculator struct {
	area float64
}

func (a *AreaCalculator) VisitCircle(c *Circle) {
	a.area = 3.14 * float64(c.Radius) * float64(c.Radius)
	fmt.Printf("Circle area: %.2f\n", a.area)
}

func (a *AreaCalculator) VisitRectangle(r *Rectangle) {
	a.area = float64(r.Width * r.Height)
	fmt.Printf("Rectangle area: %.2f\n", a.area)
}
