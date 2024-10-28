package main

import "fmt"

// Когда использовать: Когда нужно определить общий алгоритм с некоторыми шагами, которые реализуются в подклассах.

// Абстрактный класс
type Game interface {
	initialize()
	startPlay()
	endPlay()
}

func playGame(g Game) {
	g.initialize()
	g.startPlay()
	g.endPlay()
}

// Конкретные классы
type Cricket struct{}

func (c *Cricket) initialize() {
	fmt.Println("Инициализация игры в крикет.")
}

func (c *Cricket) startPlay() {
	fmt.Println("Начало игры в крикет.")
}

func (c *Cricket) endPlay() {
	fmt.Println("Завершение игры в крикет.")
}

func main() {
	game := &Cricket{}
	playGame(game)
}
