package main

import "fmt"

// Когда использовать: Когда нужно работать с большим количеством объектов, которые имеют общие внутренние состояния.

type Dress interface {
	getColor() string
}

// Конкретные объекты
type TerroristDress struct {
	color string
}

func (t *TerroristDress) getColor() string {
	return t.color
}

type CounterTerroristDress struct {
	color string
}

func (c *CounterTerroristDress) getColor() string {
	return c.color
}

// Фабрика
type DressFactory struct {
	dressMap map[string]Dress
}

func (d *DressFactory) getDressByType(dressType string) (Dress, error) {
	if d.dressMap == nil {
		d.dressMap = make(map[string]Dress)
	}

	if dress, ok := d.dressMap[dressType]; ok {
		return dress, nil
	}

	var dress Dress
	if dressType == "TERRORIST" {
		dress = &TerroristDress{color: "красный"}
	} else if dressType == "COUNTER_TERRORIST" {
		dress = &CounterTerroristDress{color: "зеленый"}
	} else {
		return nil, fmt.Errorf("Неправильный тип одежды")
	}

	d.dressMap[dressType] = dress
	return dress, nil
}

func main() {
	factory := &DressFactory{}
	dress1, _ := factory.getDressByType("TERRORIST")
	dress2, _ := factory.getDressByType("TERRORIST")
	fmt.Println("dress1 и dress2 одинаковы?", dress1 == dress2) // true
}
