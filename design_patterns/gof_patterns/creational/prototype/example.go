package prototype

import (
	"fmt"
	"time"
)

// Example1_DocumentCloning demonstrates cloning documents.
func Example1_DocumentCloning() {
	fmt.Println("\n=== Example 1: Document Cloning ===")

	// Create a complex original document
	original := &Document{
		Title:   "Design Patterns",
		Content: "A comprehensive guide to design patterns...",
		Author:  "Original Author",
		Created: time.Now(),
		Tags:    []string{"programming", "patterns", "software"},
		Metadata: map[string]string{
			"version": "1.0",
			"format":  "PDF",
		},
	}

	fmt.Println("Original:", original)

	// Clone and modify
	clone1 := original.Clone()
	clone1.Title = "Design Patterns - Copy 1"
	clone1.Author = "Editor 1"

	clone2 := original.Clone()
	clone2.Title = "Design Patterns - Copy 2"
	clone2.Author = "Editor 2"

	fmt.Println("Clone 1:", clone1)
	fmt.Println("Clone 2:", clone2)
	fmt.Println("\nOriginal unchanged:", original)
}

// Example2_ShapeCloning demonstrates cloning shapes.
func Example2_ShapeCloning() {
	fmt.Println("\n=== Example 2: Shape Cloning ===")

	// Create prototype shapes
	circleProto := &Circle{X: 0, Y: 0, Radius: 10, Color: "red"}
	rectProto := &Rectangle{X: 0, Y: 0, Width: 20, Height: 10, Color: "blue"}

	// Clone and position
	fmt.Println("\nCloning circles at different positions:")
	for i := 0; i < 3; i++ {
		circle := circleProto.Clone().(*Circle)
		circle.X = i * 50
		circle.Draw()
	}

	fmt.Println("\nCloning rectangles:")
	for i := 0; i < 3; i++ {
		rect := rectProto.Clone().(*Rectangle)
		rect.Y = i * 30
		rect.Draw()
	}
}

// Example3_GameCharacters demonstrates cloning game characters.
func Example3_GameCharacters() {
	fmt.Println("\n=== Example 3: Game Character Cloning ===")

	// Create a warrior template
	warriorTemplate := &GameCharacter{
		Name:      "Warrior",
		Health:    100,
		Mana:      50,
		Level:     1,
		Inventory: []string{"Sword", "Shield"},
		Skills: map[string]int{
			"Attack":  10,
			"Defense": 15,
		},
		Equipment: &Equipment{
			Weapon: "Iron Sword",
			Armor:  "Chain Mail",
			Shield: "Wood Shield",
		},
	}

	// Spawn multiple warriors
	fmt.Println("\nSpawning warriors from template:")
	for i := 1; i <= 3; i++ {
		warrior := warriorTemplate.Clone()
		warrior.Name = fmt.Sprintf("Warrior_%d", i)
		fmt.Println(warrior)
	}
}

// Example4_PrototypeRegistry demonstrates using a registry.
func Example4_PrototypeRegistry() {
	fmt.Println("\n=== Example 4: Prototype Registry ===")

	registry := NewPrototypeRegistry()

	// Register prototypes
	registry.Register("red-circle", &Circle{Radius: 10, Color: "red"})
	registry.Register("blue-circle", &Circle{Radius: 15, Color: "blue"})
	registry.Register("small-rect", &Rectangle{Width: 20, Height: 10, Color: "green"})

	fmt.Println("\nRegistered prototypes:", registry.List())

	// Create instances from registry
	fmt.Println("\nCreating shapes from registry:")
	shape1 := registry.Create("red-circle")
	shape1.Draw()

	shape2 := registry.Create("blue-circle")
	shape2.Draw()

	shape3 := registry.Create("small-rect")
	shape3.Draw()
}

// Example5_PrototypeBenefits demonstrates the benefits.
func Example5_PrototypeBenefits() {
	fmt.Println("\n=== Example 5: Prototype Pattern Benefits ===")

	fmt.Println("\nKEY BENEFITS:")
	fmt.Println("  ✓ Fast object creation through cloning")
	fmt.Println("  ✓ Avoid expensive initialization")
	fmt.Println("  ✓ Create objects at runtime")
	fmt.Println("  ✓ Reduce subclassing")

	fmt.Println("\nREAL-WORLD USE CASES:")
	fmt.Println("  • Document templates and cloning")
	fmt.Println("  • Game object spawning")
	fmt.Println("  • Configuration presets")
	fmt.Println("  • Graphic shape libraries")
	fmt.Println("  • Test data generation")

	// Demonstrate performance benefit
	fmt.Println("\n\nPerformance Example:")
	fmt.Println("Creating complex objects from scratch vs cloning:\n")

	// Simulate expensive creation
	template := &GameCharacter{
		Name:      "Template",
		Health:    100,
		Mana:      100,
		Level:     10,
		Inventory: []string{"Item1", "Item2", "Item3"},
		Skills: map[string]int{
			"Skill1": 5,
			"Skill2": 10,
		},
		Equipment: &Equipment{Weapon: "Sword", Armor: "Plate", Shield: "Tower"},
	}

	fmt.Println("Cloning is faster than recreating complex objects")
	for i := 1; i <= 5; i++ {
		clone := template.Clone()
		clone.Name = fmt.Sprintf("Clone_%d", i)
		fmt.Printf("  Created: %s\n", clone)
	}
}
