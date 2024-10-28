package main

import "fmt"

//Когда использовать: Когда нужно создавать сложные объекты пошагово.

type House struct {
	windowType string
	doorType   string
	floor      int
}

// Интерфейс строителя
type IBuilder interface {
	setWindowType()
	setDoorType()
	setNumFloor()
	getHouse() House
}

// Конкретные строители
type NormalBuilder struct {
	house House
}

func newNormalBuilder() *NormalBuilder {
	return &NormalBuilder{}
}

func (b *NormalBuilder) setWindowType() {
	b.house.windowType = "Деревянные окна"
}

func (b *NormalBuilder) setDoorType() {
	b.house.doorType = "Деревянные двери"
}

func (b *NormalBuilder) setNumFloor() {
	b.house.floor = 2
}

func (b *NormalBuilder) getHouse() House {
	return b.house
}

// Директор
type Director struct {
	builder IBuilder
}

func newDirector(b IBuilder) *Director {
	return &Director{
		builder: b,
	}
}

func (d *Director) setBuilder(b IBuilder) {
	d.builder = b
}

func (d *Director) buildHouse() House {
	d.builder.setWindowType()
	d.builder.setDoorType()
	d.builder.setNumFloor()
	return d.builder.getHouse()
}

func main() {
	normalBuilder := newNormalBuilder()
	director := newDirector(normalBuilder)
	house := director.buildHouse()
	fmt.Printf("Тип окон дома: %s\n", house.windowType)
}
