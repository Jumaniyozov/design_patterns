package main

import "fmt"

// Когда использовать: Когда нужно превратить запросы в объекты, позволяя передавать их как аргументы.

// Интерфейс команды
type Command interface {
	execute()
}

// Получатель
type Light struct{}

func (l *Light) on() {
	fmt.Println("Свет включен")
}

func (l *Light) off() {
	fmt.Println("Свет выключен")
}

// Конкретные команды
type OnCommand struct {
	light *Light
}

func (c *OnCommand) execute() {
	c.light.on()
}

type OffCommand struct {
	light *Light
}

func (c *OffCommand) execute() {
	c.light.off()
}

// Инвокер
type Switch struct {
	commands []Command
}

func (s *Switch) storeAndExecute(cmd Command) {
	s.commands = append(s.commands, cmd)
	cmd.execute()
}

func main() {
	light := &Light{}
	onCommand := &OnCommand{light}
	offCommand := &OffCommand{light}

	sw := &Switch{}
	sw.storeAndExecute(onCommand)
	sw.storeAndExecute(offCommand)
}
