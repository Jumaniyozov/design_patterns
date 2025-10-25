package composite

import "fmt"

type Component interface {
	Operation() string
	Add(Component)
	Remove(Component)
	GetChildren() []Component
}

type Leaf struct {
	name string
}

func NewLeaf(name string) *Leaf {
	return &Leaf{name: name}
}

func (l *Leaf) Operation() string {
	return fmt.Sprintf("Leaf: %s", l.name)
}

func (l *Leaf) Add(Component)            {}
func (l *Leaf) Remove(Component)         {}
func (l *Leaf) GetChildren() []Component { return nil }

type Composite struct {
	name     string
	children []Component
}

func NewComposite(name string) *Composite {
	return &Composite{name: name, children: make([]Component, 0)}
}

func (c *Composite) Operation() string {
	result := fmt.Sprintf("Composite: %s\n", c.name)
	for _, child := range c.children {
		result += "  " + child.Operation() + "\n"
	}
	return result
}

func (c *Composite) Add(component Component) {
	c.children = append(c.children, component)
}

func (c *Composite) Remove(component Component) {
	for i, child := range c.children {
		if child == component {
			c.children = append(c.children[:i], c.children[i+1:]...)
			break
		}
	}
}

func (c *Composite) GetChildren() []Component {
	return c.children
}
