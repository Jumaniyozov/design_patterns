package main

import "fmt"

// Когда использовать: Когда нужно динамически добавлять новую функциональность объектам без изменения их кода.

// Компонент
type Beverage interface {
	getDescription() string
	cost() float64
}

// Конкретный компонент
type Espresso struct{}

func (e *Espresso) getDescription() string {
	return "Эспрессо"
}

func (e *Espresso) cost() float64 {
	return 1.99
}

// Декоратор
type Mocha struct {
	beverage Beverage
}

func (m *Mocha) getDescription() string {
	return m.beverage.getDescription() + ", Мокка"
}

func (m *Mocha) cost() float64 {
	return m.beverage.cost() + 0.20
}

func main() {
	beverage := &Espresso{}
	fmt.Println(beverage.getDescription(), " $", beverage.cost())

	beverageWithMocha := &Mocha{beverage}
	fmt.Println(beverageWithMocha.getDescription(), " $", beverageWithMocha.cost())
}
