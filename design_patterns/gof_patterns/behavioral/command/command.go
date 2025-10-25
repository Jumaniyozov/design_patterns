package command

import "fmt"

type Command interface {
	Execute()
	Undo()
}

type Light struct {
	isOn bool
}

func (l *Light) On() {
	l.isOn = true
	fmt.Println("Light is ON")
}

func (l *Light) Off() {
	l.isOn = false
	fmt.Println("Light is OFF")
}

type LightOnCommand struct {
	light *Light
}

func (c *LightOnCommand) Execute() {
	c.light.On()
}

func (c *LightOnCommand) Undo() {
	c.light.Off()
}

type LightOffCommand struct {
	light *Light
}

func (c *LightOffCommand) Execute() {
	c.light.Off()
}

func (c *LightOffCommand) Undo() {
	c.light.On()
}

type RemoteControl struct {
	command Command
	history []Command
}

func (r *RemoteControl) SetCommand(c Command) {
	r.command = c
}

func (r *RemoteControl) PressButton() {
	r.command.Execute()
	r.history = append(r.history, r.command)
}

func (r *RemoteControl) PressUndo() {
	if len(r.history) > 0 {
		cmd := r.history[len(r.history)-1]
		cmd.Undo()
		r.history = r.history[:len(r.history)-1]
	}
}
