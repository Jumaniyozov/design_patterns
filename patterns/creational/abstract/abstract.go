package main

import "fmt"

// Когда использовать: Когда нужно создавать семейства связанных или зависимых объектов без указания их конкретных классов.

// Абстрактные продукты
type Shoe interface {
	SetLogo(logo string)
	GetLogo() string
}

type Shirt interface {
	SetLogo(logo string)
	GetLogo() string
}

// Конкретные продукты
type NikeShoe struct {
	logo string
}

func (s *NikeShoe) SetLogo(logo string) {
	s.logo = logo
}

func (s *NikeShoe) GetLogo() string {
	return s.logo
}

type AdidasShoe struct {
	logo string
}

func (s *AdidasShoe) SetLogo(logo string) {
	s.logo = logo
}

func (s *AdidasShoe) GetLogo() string {
	return s.logo
}

// Абстрактная фабрика
type SportsFactory interface {
	MakeShoe() Shoe
	MakeShirt() Shirt
}

// Конкретные фабрики
type NikeFactory struct{}

func (n *NikeFactory) MakeShoe() Shoe {
	shoe := &NikeShoe{}
	shoe.SetLogo("Nike")
	return shoe
}

func (n *NikeFactory) MakeShirt() Shirt {
	// Аналогично
	return nil
}

type AdidasFactory struct{}

func (a *AdidasFactory) MakeShoe() Shoe {
	shoe := &AdidasShoe{}
	shoe.SetLogo("Adidas")
	return shoe
}

func (a *AdidasFactory) MakeShirt() Shirt {
	// Аналогично
	return nil
}

// Клиентский код
func GetSportsFactory(brand string) SportsFactory {
	if brand == "nike" {
		return &NikeFactory{}
	}
	if brand == "adidas" {
		return &AdidasFactory{}
	}
	return nil
}

func main() {
	nikeFactory := GetSportsFactory("nike")
	nikeShoe := nikeFactory.MakeShoe()
	adidasFactory := GetSportsFactory("adidas")
	adidasShoe := adidasFactory.MakeShoe()
	fmt.Println(nikeShoe.GetLogo())   // Output: Nike
	fmt.Println(adidasShoe.GetLogo()) // Output: Adidas
}
