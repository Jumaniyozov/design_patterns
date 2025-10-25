// Package abstractfactory demonstrates the Abstract Factory pattern.
// It provides an interface for creating families of related objects without
// specifying their concrete classes, ensuring product families are used together.
package abstractfactory

import "fmt"

// Product interfaces - define what products can do

// Button represents a UI button component
type Button interface {
	Render() string
	OnClick() string
}

// Checkbox represents a UI checkbox component
type Checkbox interface {
	Render() string
	Toggle() string
}

// Abstract Factory interface

// UIFactory creates families of related UI components
type UIFactory interface {
	CreateButton() Button
	CreateCheckbox() Checkbox
}

// Concrete Products - Windows family

// WindowsButton is a Windows-style button
type WindowsButton struct{}

func (b *WindowsButton) Render() string {
	return "[Windows Button]"
}

func (b *WindowsButton) OnClick() string {
	return "Windows button clicked"
}

// WindowsCheckbox is a Windows-style checkbox
type WindowsCheckbox struct {
	checked bool
}

func (c *WindowsCheckbox) Render() string {
	if c.checked {
		return "[X] Windows Checkbox"
	}
	return "[ ] Windows Checkbox"
}

func (c *WindowsCheckbox) Toggle() string {
	c.checked = !c.checked
	return fmt.Sprintf("Windows checkbox toggled: %v", c.checked)
}

// Concrete Products - macOS family

// MacButton is a macOS-style button
type MacButton struct{}

func (b *MacButton) Render() string {
	return "( macOS Button )"
}

func (b *MacButton) OnClick() string {
	return "macOS button clicked with animation"
}

// MacCheckbox is a macOS-style checkbox
type MacCheckbox struct {
	checked bool
}

func (c *MacCheckbox) Render() string {
	if c.checked {
		return "☑ macOS Checkbox"
	}
	return "☐ macOS Checkbox"
}

func (c *MacCheckbox) Toggle() string {
	c.checked = !c.checked
	return fmt.Sprintf("macOS checkbox toggled with transition: %v", c.checked)
}

// Concrete Products - Linux family

// LinuxButton is a Linux-style button
type LinuxButton struct{}

func (b *LinuxButton) Render() string {
	return "{ Linux Button }"
}

func (b *LinuxButton) OnClick() string {
	return "Linux button clicked"
}

// LinuxCheckbox is a Linux-style checkbox
type LinuxCheckbox struct {
	checked bool
}

func (c *LinuxCheckbox) Render() string {
	if c.checked {
		return "[*] Linux Checkbox"
	}
	return "[-] Linux Checkbox"
}

func (c *LinuxCheckbox) Toggle() string {
	c.checked = !c.checked
	return fmt.Sprintf("Linux checkbox toggled: %v", c.checked)
}

// Concrete Factories

// WindowsFactory creates Windows UI components
type WindowsFactory struct{}

func (f *WindowsFactory) CreateButton() Button {
	return &WindowsButton{}
}

func (f *WindowsFactory) CreateCheckbox() Checkbox {
	return &WindowsCheckbox{}
}

// MacFactory creates macOS UI components
type MacFactory struct{}

func (f *MacFactory) CreateButton() Button {
	return &MacButton{}
}

func (f *MacFactory) CreateCheckbox() Checkbox {
	return &MacCheckbox{}
}

// LinuxFactory creates Linux UI components
type LinuxFactory struct{}

func (f *LinuxFactory) CreateButton() Button {
	return &LinuxButton{}
}

func (f *LinuxFactory) CreateCheckbox() Checkbox {
	return &LinuxCheckbox{}
}

// Factory Provider - returns appropriate factory based on platform

// GetUIFactory returns the appropriate UI factory for the given platform
func GetUIFactory(platform string) UIFactory {
	switch platform {
	case "windows":
		return &WindowsFactory{}
	case "macos":
		return &MacFactory{}
	case "linux":
		return &LinuxFactory{}
	default:
		// Default to Linux
		return &LinuxFactory{}
	}
}

// Application demonstrates using the factory

// Application uses UI factory to create consistent UI components
type Application struct {
	factory  UIFactory
	button   Button
	checkbox Checkbox
}

// NewApplication creates an application with platform-specific UI factory
func NewApplication(factory UIFactory) *Application {
	return &Application{
		factory:  factory,
		button:   factory.CreateButton(),
		checkbox: factory.CreateCheckbox(),
	}
}

// Render renders all UI components
func (a *Application) Render() string {
	return fmt.Sprintf("Application UI:\n  %s\n  %s",
		a.button.Render(),
		a.checkbox.Render())
}

// InteractWithUI simulates user interaction
func (a *Application) InteractWithUI() string {
	result := a.button.OnClick() + "\n"
	result += a.checkbox.Toggle()
	return result
}
