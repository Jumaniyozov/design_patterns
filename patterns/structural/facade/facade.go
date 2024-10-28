package main

import "fmt"

// Когда использовать: Когда нужно предоставить простой интерфейс к сложной системе.

// Подсистемы
type CPU struct{}

func (c *CPU) Freeze() {
	fmt.Println("Замораживаем CPU.")
}

func (c *CPU) Jump(position int) {
	fmt.Printf("Прыжок на позицию %d.\n", position)
}

func (c *CPU) Execute() {
	fmt.Println("Выполняем инструкции.")
}

type Memory struct{}

func (m *Memory) Load(position int, data []byte) {
	fmt.Printf("Загрузка данных в память на позицию %d.\n", position)
}

type HardDrive struct{}

func (h *HardDrive) Read(lba int, size int) []byte {
	fmt.Printf("Чтение %d байт с LBA %d.\n", size, lba)
	return []byte{}
}

// Фасад
type Computer struct {
	cpu       *CPU
	memory    *Memory
	hardDrive *HardDrive
}

func NewComputer() *Computer {
	return &Computer{
		cpu:       &CPU{},
		memory:    &Memory{},
		hardDrive: &HardDrive{},
	}
}

func (c *Computer) Start() {
	c.cpu.Freeze()
	c.memory.Load(0, c.hardDrive.Read(0, 1024))
	c.cpu.Jump(0)
	c.cpu.Execute()
}

func main() {
	computer := NewComputer()
	computer.Start()
}
