// Package composite demonstrates the Composite Pattern, a structural design pattern
// that lets you compose objects into tree structures to represent part-whole hierarchies.
//
// The Composite Pattern is essential for:
// - Tree or hierarchical structures (file systems, org charts, UI components)
// - Operations that work uniformly on individual items and collections
// - Recursive data structures where parts can contain other parts
// - Building complex structures from simpler components
package composite

import (
	"fmt"
	"strings"
)

// =============================================================================
// Example 1: File System (Files and Directories)
// =============================================================================

// FileSystemComponent is the component interface for both files and directories.
type FileSystemComponent interface {
	GetName() string
	GetSize() int
	Display(indent string) string
}

// File represents a leaf node in the file system (cannot contain other components).
type File struct {
	name string
	size int
}

// NewFile creates a new file.
func NewFile(name string, size int) *File {
	return &File{name: name, size: size}
}

// GetName returns the file name.
func (f *File) GetName() string {
	return f.name
}

// GetSize returns the file size.
func (f *File) GetSize() int {
	return f.size
}

// Display returns a string representation of the file.
func (f *File) Display(indent string) string {
	return fmt.Sprintf("%s📄 %s (%d KB)", indent, f.name, f.size)
}

// Directory represents a composite node that can contain files and other directories.
type Directory struct {
	name     string
	children []FileSystemComponent
}

// NewDirectory creates a new directory.
func NewDirectory(name string) *Directory {
	return &Directory{
		name:     name,
		children: make([]FileSystemComponent, 0),
	}
}

// GetName returns the directory name.
func (d *Directory) GetName() string {
	return d.name
}

// GetSize returns the total size of the directory (sum of all children).
func (d *Directory) GetSize() int {
	totalSize := 0
	for _, child := range d.children {
		totalSize += child.GetSize()
	}
	return totalSize
}

// Display returns a tree representation of the directory and its contents.
func (d *Directory) Display(indent string) string {
	result := fmt.Sprintf("%s📁 %s/", indent, d.name)
	for _, child := range d.children {
		result += "\n" + child.Display(indent+"  ")
	}
	return result
}

// Add adds a component (file or directory) to this directory.
func (d *Directory) Add(component FileSystemComponent) {
	d.children = append(d.children, component)
}

// Remove removes a component from this directory.
func (d *Directory) Remove(component FileSystemComponent) {
	for i, child := range d.children {
		if child == component {
			d.children = append(d.children[:i], d.children[i+1:]...)
			return
		}
	}
}

// GetChildren returns all children of this directory.
func (d *Directory) GetChildren() []FileSystemComponent {
	return d.children
}

// =============================================================================
// Example 2: Organization Chart (Employees and Departments)
// =============================================================================

// Employee represents a component in the organization chart.
type Employee interface {
	GetName() string
	GetRole() string
	GetSalary() int
	PrintHierarchy(indent string) string
}

// IndividualEmployee represents a leaf node (employee without subordinates).
type IndividualEmployee struct {
	name   string
	role   string
	salary int
}

// NewIndividualEmployee creates a new individual employee.
func NewIndividualEmployee(name, role string, salary int) *IndividualEmployee {
	return &IndividualEmployee{
		name:   name,
		role:   role,
		salary: salary,
	}
}

// GetName returns the employee's name.
func (e *IndividualEmployee) GetName() string {
	return e.name
}

// GetRole returns the employee's role.
func (e *IndividualEmployee) GetRole() string {
	return e.role
}

// GetSalary returns the employee's salary.
func (e *IndividualEmployee) GetSalary() int {
	return e.salary
}

// PrintHierarchy returns a string representation of the employee.
func (e *IndividualEmployee) PrintHierarchy(indent string) string {
	return fmt.Sprintf("%s👤 %s - %s ($%d)", indent, e.name, e.role, e.salary)
}

// Manager represents a composite node (employee with subordinates).
type Manager struct {
	employee     *IndividualEmployee
	subordinates []Employee
}

// NewManager creates a new manager.
func NewManager(name, role string, salary int) *Manager {
	return &Manager{
		employee:     NewIndividualEmployee(name, role, salary),
		subordinates: make([]Employee, 0),
	}
}

// GetName returns the manager's name.
func (m *Manager) GetName() string {
	return m.employee.GetName()
}

// GetRole returns the manager's role.
func (m *Manager) GetRole() string {
	return m.employee.GetRole()
}

// GetSalary returns the total salary cost (manager + all subordinates).
func (m *Manager) GetSalary() int {
	total := m.employee.GetSalary()
	for _, sub := range m.subordinates {
		total += sub.GetSalary()
	}
	return total
}

// PrintHierarchy returns a tree representation of the manager and subordinates.
func (m *Manager) PrintHierarchy(indent string) string {
	result := fmt.Sprintf("%s👔 %s - %s ($%d)",
		indent, m.employee.GetName(), m.employee.GetRole(), m.employee.GetSalary())

	if len(m.subordinates) > 0 {
		result += " [manages " + fmt.Sprintf("%d", len(m.subordinates)) + "]"
	}

	for _, sub := range m.subordinates {
		result += "\n" + sub.PrintHierarchy(indent+"  ")
	}
	return result
}

// AddSubordinate adds an employee to this manager's team.
func (m *Manager) AddSubordinate(employee Employee) {
	m.subordinates = append(m.subordinates, employee)
}

// RemoveSubordinate removes an employee from this manager's team.
func (m *Manager) RemoveSubordinate(employee Employee) {
	for i, sub := range m.subordinates {
		if sub == employee {
			m.subordinates = append(m.subordinates[:i], m.subordinates[i+1:]...)
			return
		}
	}
}

// GetSubordinates returns all subordinates.
func (m *Manager) GetSubordinates() []Employee {
	return m.subordinates
}

// =============================================================================
// Example 3: UI Components (Simple and Container Components)
// =============================================================================

// UIComponent represents a component in the UI hierarchy.
type UIComponent interface {
	Render() string
	GetWidth() int
	GetHeight() int
}

// Button represents a simple UI component (leaf).
type Button struct {
	text   string
	width  int
	height int
}

// NewButton creates a new button.
func NewButton(text string, width, height int) *Button {
	return &Button{text: text, width: width, height: height}
}

// Render returns the button's HTML representation.
func (b *Button) Render() string {
	return fmt.Sprintf("<button width=%d height=%d>%s</button>", b.width, b.height, b.text)
}

// GetWidth returns the button's width.
func (b *Button) GetWidth() int {
	return b.width
}

// GetHeight returns the button's height.
func (b *Button) GetHeight() int {
	return b.height
}

// TextBox represents a simple UI component (leaf).
type TextBox struct {
	placeholder string
	width       int
	height      int
}

// NewTextBox creates a new text box.
func NewTextBox(placeholder string, width, height int) *TextBox {
	return &TextBox{placeholder: placeholder, width: width, height: height}
}

// Render returns the text box's HTML representation.
func (t *TextBox) Render() string {
	return fmt.Sprintf("<input type='text' placeholder='%s' width=%d height=%d />",
		t.placeholder, t.width, t.height)
}

// GetWidth returns the text box's width.
func (t *TextBox) GetWidth() int {
	return t.width
}

// GetHeight returns the text box's height.
func (t *TextBox) GetHeight() int {
	return t.height
}

// Panel represents a composite UI component that can contain other components.
type Panel struct {
	name       string
	components []UIComponent
	width      int
	height     int
}

// NewPanel creates a new panel.
func NewPanel(name string, width, height int) *Panel {
	return &Panel{
		name:       name,
		components: make([]UIComponent, 0),
		width:      width,
		height:     height,
	}
}

// Render returns the panel's HTML representation with all children.
func (p *Panel) Render() string {
	var children []string
	for _, component := range p.components {
		children = append(children, component.Render())
	}

	return fmt.Sprintf("<div class='%s' width=%d height=%d>\n  %s\n</div>",
		p.name, p.width, p.height, strings.Join(children, "\n  "))
}

// GetWidth returns the panel's width.
func (p *Panel) GetWidth() int {
	return p.width
}

// GetHeight returns the panel's height.
func (p *Panel) GetHeight() int {
	return p.height
}

// Add adds a UI component to this panel.
func (p *Panel) Add(component UIComponent) {
	p.components = append(p.components, component)
}

// Remove removes a UI component from this panel.
func (p *Panel) Remove(component UIComponent) {
	for i, comp := range p.components {
		if comp == component {
			p.components = append(p.components[:i], p.components[i+1:]...)
			return
		}
	}
}

// GetComponents returns all components in this panel.
func (p *Panel) GetComponents() []UIComponent {
	return p.components
}

// =============================================================================
// Example 4: Menu System (Menu Items and Submenus)
// =============================================================================

// MenuComponent represents a component in the menu system.
type MenuComponent interface {
	GetName() string
	Execute()
	Print(indent string) string
}

// MenuItem represents a leaf menu item that performs an action.
type MenuItem struct {
	name   string
	action string
}

// NewMenuItem creates a new menu item.
func NewMenuItem(name, action string) *MenuItem {
	return &MenuItem{name: name, action: action}
}

// GetName returns the menu item's name.
func (m *MenuItem) GetName() string {
	return m.name
}

// Execute performs the menu item's action.
func (m *MenuItem) Execute() {
	fmt.Printf("Executing: %s\n", m.action)
}

// Print returns a string representation of the menu item.
func (m *MenuItem) Print(indent string) string {
	return fmt.Sprintf("%s▸ %s", indent, m.name)
}

// SubMenu represents a composite that can contain menu items and other submenus.
type SubMenu struct {
	name  string
	items []MenuComponent
}

// NewSubMenu creates a new submenu.
func NewSubMenu(name string) *SubMenu {
	return &SubMenu{
		name:  name,
		items: make([]MenuComponent, 0),
	}
}

// GetName returns the submenu's name.
func (s *SubMenu) GetName() string {
	return s.name
}

// Execute executes all items in the submenu.
func (s *SubMenu) Execute() {
	fmt.Printf("Opening submenu: %s\n", s.name)
	for _, item := range s.items {
		item.Execute()
	}
}

// Print returns a tree representation of the submenu.
func (s *SubMenu) Print(indent string) string {
	result := fmt.Sprintf("%s📂 %s", indent, s.name)
	for _, item := range s.items {
		result += "\n" + item.Print(indent+"  ")
	}
	return result
}

// Add adds a menu component to this submenu.
func (s *SubMenu) Add(component MenuComponent) {
	s.items = append(s.items, component)
}

// Remove removes a menu component from this submenu.
func (s *SubMenu) Remove(component MenuComponent) {
	for i, item := range s.items {
		if item == component {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return
		}
	}
}

// GetItems returns all items in this submenu.
func (s *SubMenu) GetItems() []MenuComponent {
	return s.items
}
