package facade

import "fmt"

type CPU struct{}

func (c *CPU) Freeze()           { fmt.Println("CPU: Freezing") }
func (c *CPU) Jump(position int) { fmt.Println("CPU: Jumping to", position) }
func (c *CPU) Execute()          { fmt.Println("CPU: Executing") }

type Memory struct{}

func (m *Memory) Load(position int, data string) { fmt.Println("Memory: Loading data at", position) }

type HardDrive struct{}

func (h *HardDrive) Read(lba int, size int) string {
	fmt.Println("HardDrive: Reading", size, "bytes from", lba)
	return "data"
}

type ComputerFacade struct {
	cpu    *CPU
	memory *Memory
	hd     *HardDrive
}

func NewComputerFacade() *ComputerFacade {
	return &ComputerFacade{&CPU{}, &Memory{}, &HardDrive{}}
}

func (c *ComputerFacade) Start() {
	fmt.Println("\nStarting computer...")
	c.cpu.Freeze()
	data := c.hd.Read(0, 1024)
	c.memory.Load(0, data)
	c.cpu.Jump(0)
	c.cpu.Execute()
	fmt.Println("Computer started!")
}
