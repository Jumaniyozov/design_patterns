package main

import "fmt"

type Animal interface {
	Speak() string
}

type Dog struct{}

func (d Dog) Speak() string {
	return "Гав!"
}

type Cat struct{}

func (c Cat) Speak() string {
	return "Мяу!"
}

func AnimalFactory(animalType string) Animal {
	switch animalType {
	case "dog":
		return &Dog{}
	case "cat":
		return &Cat{}
	}

	return nil
}

func main() {
	animal := AnimalFactory("dog")
	fmt.Println(animal.Speak()) // Output: Гав!
}
