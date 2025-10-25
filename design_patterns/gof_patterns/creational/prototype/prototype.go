// Package prototype demonstrates the Prototype pattern in Go.
// The Prototype pattern creates new objects by cloning existing instances.
package prototype

import (
	"fmt"
	"time"
)

// Document represents a complex document that can be cloned.
type Document struct {
	Title    string
	Content  string
	Author   string
	Created  time.Time
	Tags     []string
	Metadata map[string]string
}

// Clone creates a deep copy of the document.
func (d *Document) Clone() *Document {
	// Deep copy slices and maps
	tags := make([]string, len(d.Tags))
	copy(tags, d.Tags)

	metadata := make(map[string]string)
	for k, v := range d.Metadata {
		metadata[k] = v
	}

	return &Document{
		Title:    d.Title,
		Content:  d.Content,
		Author:   d.Author,
		Created:  d.Created,
		Tags:     tags,
		Metadata: metadata,
	}
}

func (d *Document) String() string {
	return fmt.Sprintf("Document{Title: %s, Author: %s, Tags: %v}", d.Title, d.Author, d.Tags)
}

// Shape interface for geometric shapes.
type Shape interface {
	Clone() Shape
	Draw()
	GetInfo() string
}

// Circle represents a circle shape.
type Circle struct {
	X      int
	Y      int
	Radius int
	Color  string
}

func (c *Circle) Clone() Shape {
	return &Circle{
		X:      c.X,
		Y:      c.Y,
		Radius: c.Radius,
		Color:  c.Color,
	}
}

func (c *Circle) Draw() {
	fmt.Printf("Drawing Circle at (%d,%d) with radius %d, color: %s\n", c.X, c.Y, c.Radius, c.Color)
}

func (c *Circle) GetInfo() string {
	return fmt.Sprintf("Circle(x=%d, y=%d, r=%d, color=%s)", c.X, c.Y, c.Radius, c.Color)
}

// Rectangle represents a rectangle shape.
type Rectangle struct {
	X      int
	Y      int
	Width  int
	Height int
	Color  string
}

func (r *Rectangle) Clone() Shape {
	return &Rectangle{
		X:      r.X,
		Y:      r.Y,
		Width:  r.Width,
		Height: r.Height,
		Color:  r.Color,
	}
}

func (r *Rectangle) Draw() {
	fmt.Printf("Drawing Rectangle at (%d,%d) with size %dx%d, color: %s\n", r.X, r.Y, r.Width, r.Height, r.Color)
}

func (r *Rectangle) GetInfo() string {
	return fmt.Sprintf("Rectangle(x=%d, y=%d, w=%d, h=%d, color=%s)", r.X, r.Y, r.Width, r.Height, r.Color)
}

// GameCharacter represents a game character that can be cloned.
type GameCharacter struct {
	Name      string
	Health    int
	Mana      int
	Level     int
	Inventory []string
	Skills    map[string]int
	Equipment *Equipment
}

type Equipment struct {
	Weapon string
	Armor  string
	Shield string
}

// Clone creates a deep copy of the character.
func (g *GameCharacter) Clone() *GameCharacter {
	inventory := make([]string, len(g.Inventory))
	copy(inventory, g.Inventory)

	skills := make(map[string]int)
	for k, v := range g.Skills {
		skills[k] = v
	}

	equipment := &Equipment{
		Weapon: g.Equipment.Weapon,
		Armor:  g.Equipment.Armor,
		Shield: g.Equipment.Shield,
	}

	return &GameCharacter{
		Name:      g.Name,
		Health:    g.Health,
		Mana:      g.Mana,
		Level:     g.Level,
		Inventory: inventory,
		Skills:    skills,
		Equipment: equipment,
	}
}

func (g *GameCharacter) String() string {
	return fmt.Sprintf("%s (Lvl %d, HP: %d, MP: %d)", g.Name, g.Level, g.Health, g.Mana)
}

// PrototypeRegistry manages prototype instances.
type PrototypeRegistry struct {
	prototypes map[string]Shape
}

func NewPrototypeRegistry() *PrototypeRegistry {
	return &PrototypeRegistry{
		prototypes: make(map[string]Shape),
	}
}

func (r *PrototypeRegistry) Register(name string, prototype Shape) {
	r.prototypes[name] = prototype
}

func (r *PrototypeRegistry) Create(name string) Shape {
	if proto, ok := r.prototypes[name]; ok {
		return proto.Clone()
	}
	return nil
}

func (r *PrototypeRegistry) List() []string {
	var names []string
	for name := range r.prototypes {
		names = append(names, name)
	}
	return names
}
