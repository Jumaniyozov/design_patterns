package abstract_factory

import "fmt"

// Example1_CrossPlatformUI demonstrates creating UI for different platforms.
func Example1_CrossPlatformUI() {
	fmt.Println("\n=== Example 1: Cross-Platform UI ===")

	platforms := []GUIFactory{
		&WindowsFactory{},
		&MacFactory{},
		&LinuxFactory{},
	}

	for _, factory := range platforms {
		app := NewApplication(factory)
		app.RenderUI()
		app.HandleUserInput()
	}
}

// Example2_DatabaseAbstraction demonstrates database abstraction.
func Example2_DatabaseAbstraction() {
	fmt.Println("\n=== Example 2: Database Abstraction ===")

	databases := []DatabaseFactory{
		&MySQLFactory{},
		&PostgreSQLFactory{},
	}

	for _, dbFactory := range databases {
		fmt.Printf("\n--- Using %s ---\n", dbFactory.GetDBType())

		conn := dbFactory.CreateConnection()
		conn.Connect()

		query := dbFactory.CreateQueryBuilder()
		sql := query.Select("*").From("users").Build()
		fmt.Printf("Query: %s\n", sql)

		tx := dbFactory.CreateTransaction()
		tx.Begin()
		fmt.Println("Executing query...")
		tx.Commit()

		conn.Close()
	}
}

// Example3_FactorySelection demonstrates runtime factory selection.
func Example3_FactorySelection() {
	fmt.Println("\n=== Example 3: Runtime Factory Selection ===")

	// Simulate runtime platform detection
	platforms := map[string]GUIFactory{
		"windows": &WindowsFactory{},
		"mac":     &MacFactory{},
		"linux":   &LinuxFactory{},
	}

	// Select factory based on runtime condition
	currentPlatform := "mac" // This would be detected at runtime
	factory := platforms[currentPlatform]

	fmt.Printf("Detected platform: %s\n", currentPlatform)

	app := NewApplication(factory)
	app.RenderUI()
}

// Example4_ConsistentFamilies demonstrates that products work together.
func Example4_ConsistentFamilies() {
	fmt.Println("\n=== Example 4: Consistent Product Families ===")

	fmt.Println("\nKey benefit: All products from a factory are compatible\n")

	factory := &WindowsFactory{}

	// All these components are designed to work together
	button := factory.CreateButton()
	checkbox := factory.CreateCheckbox()
	textField := factory.CreateTextField()

	fmt.Println("Creating Windows form with consistent styling:")
	button.Render()
	checkbox.Render()
	textField.SetText("Windows styled input")
	textField.Render()

	fmt.Println("\nAll components share:")
	fmt.Println("  • Same visual style")
	fmt.Println("  • Consistent behavior")
	fmt.Println("  • Compatible event handling")
	fmt.Println("  • Unified theme")
}

// Example5_AbstractFactoryBenefits summarizes the benefits.
func Example5_AbstractFactoryBenefits() {
	fmt.Println("\n=== Example 5: Abstract Factory Benefits ===")

	fmt.Println("\nKEY BENEFITS:")

	fmt.Println("\n1. GUARANTEED COMPATIBILITY")
	fmt.Println("   Products from same factory are designed to work together")
	fmt.Println("   No risk of mixing incompatible components")

	fmt.Println("\n2. EASY PLATFORM SWITCHING")
	fmt.Println("   Change one line of code to switch entire product family:")
	fmt.Println("   factory := &WindowsFactory{}")
	fmt.Println("   ↓")
	fmt.Println("   factory := &MacFactory{}")

	fmt.Println("\n3. ISOLATION FROM CONCRETE CLASSES")
	fmt.Println("   Client code only knows about interfaces")
	fmt.Println("   Doesn't depend on Windows, Mac, Linux specifics")

	fmt.Println("\n4. SINGLE RESPONSIBILITY")
	fmt.Println("   Product creation logic centralized in factories")
	fmt.Println("   Easy to maintain and extend")

	fmt.Println("\nREAL-WORLD APPLICATIONS:")
	fmt.Println("  • Cross-platform GUI toolkits")
	fmt.Println("  • Database abstraction layers")
	fmt.Println("  • Theme systems (dark/light mode)")
	fmt.Println("  • Document generators (PDF/HTML/Markdown)")
	fmt.Println("  • Cloud provider abstraction (AWS/Azure/GCP)")
	fmt.Println("  • Game engines (DirectX/OpenGL/Vulkan)")

	// Demonstrate
	fmt.Println("\n\nPRACTICAL EXAMPLE:")
	fmt.Println("Application works with any factory:\n")

	factories := []GUIFactory{
		&WindowsFactory{},
		&MacFactory{},
	}

	for _, factory := range factories {
		app := NewApplication(factory)
		fmt.Printf("Running on %s:\n", factory.GetName())
		button := app.factory.CreateButton()
		button.Render()
		fmt.Println()
	}

	fmt.Println("Same application code, different product families!")
}
