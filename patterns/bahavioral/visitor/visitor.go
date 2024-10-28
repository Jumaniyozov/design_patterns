package main

import "fmt"

// Когда использовать: Когда нужно выполнить операции над объектами сложной структуры, не изменяя их классы.

// Посетитель
type Visitor interface {
	visitForSquare(*Square)
	visitForCircle(*Circle)
}

// Элемент
type Shape interface {
	getType() string
	accept(Visitor)
}

// Конкретные элементы
type Square struct {
	side int
}

func (s *Square) accept(v Visitor) {
	v.visitForSquare(s)
}

func (s *Square) getType() string {
	return "Квадрат"
}

type Circle struct {
	radius int
}

func (c *Circle) accept(v Visitor) {
	v.visitForCircle(c)
}

func (c *Circle) getType() string {
	return "Круг"
}

// Конкретный посетитель
type AreaCalculator struct {
	area int
}

func (a *AreaCalculator) visitForSquare(s *Square) {
	a.area = s.side * s.side
	fmt.Printf("Вычисление площади для %s: %d\n", s.getType(), a.area)
}

func (a *AreaCalculator) visitForCircle(c *Circle) {
	a.area = c.radius * c.radius * 3 // Приблизительное значение π
	fmt.Printf("Вычисление площади для %s: %d\n", c.getType(), a.area)
}

func main() {
	square := &Square{side: 4}
	circle := &Circle{radius: 3}

	areaCalc := &AreaCalculator{}

	square.accept(areaCalc)
	circle.accept(areaCalc)
}
