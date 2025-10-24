package composite

import (
	"fmt"
	"strings"
)

// Example1_FileSystem demonstrates the Composite Pattern with a file system.
// This shows how files (leaf) and directories (composite) can be treated uniformly.
func Example1_FileSystem() {
	fmt.Println("=== Example 1: File System ===")

	// Create files (leaf nodes)
	file1 := NewFile("config.json", 5)
	file2 := NewFile("readme.md", 15)
	file3 := NewFile("main.go", 120)
	file4 := NewFile("utils.go", 80)
	file5 := NewFile("test.go", 45)

	// Create directories (composite nodes)
	rootDir := NewDirectory("project")
	srcDir := NewDirectory("src")
	docsDir := NewDirectory("docs")

	// Build the tree structure
	rootDir.Add(file1)        // config.json in root
	rootDir.Add(docsDir)      // docs/ in root
	rootDir.Add(srcDir)       // src/ in root
	docsDir.Add(file2)        // readme.md in docs/
	srcDir.Add(file3)         // main.go in src/
	srcDir.Add(file4)         // utils.go in src/
	srcDir.Add(file5)         // test.go in src/

	// Display the entire tree
	fmt.Println("Directory Structure:")
	fmt.Println(rootDir.Display(""))

	// Calculate total size (demonstrates uniform treatment)
	fmt.Printf("\nTotal size of '%s': %d KB\n", rootDir.GetName(), rootDir.GetSize())
	fmt.Printf("Size of '%s': %d KB\n", srcDir.GetName(), srcDir.GetSize())
	fmt.Printf("Size of '%s': %d KB\n", file3.GetName(), file3.GetSize())

	fmt.Println("\n" + strings.Repeat("─", 50))
}

// Example2_OrganizationChart demonstrates the Composite Pattern with an organization hierarchy.
// This shows how individual employees and managers (with teams) can be treated uniformly.
func Example2_OrganizationChart() {
	fmt.Println("=== Example 2: Organization Chart ===")

	// Create individual employees (leaf nodes)
	emp1 := NewIndividualEmployee("Alice Johnson", "Software Engineer", 95000)
	emp2 := NewIndividualEmployee("Bob Smith", "Software Engineer", 92000)
	emp3 := NewIndividualEmployee("Carol White", "QA Engineer", 85000)
	emp4 := NewIndividualEmployee("David Lee", "DevOps Engineer", 98000)
	emp5 := NewIndividualEmployee("Emma Brown", "Product Designer", 90000)

	// Create managers (composite nodes)
	teamLead := NewManager("Frank Miller", "Engineering Lead", 130000)
	cto := NewManager("Grace Chen", "CTO", 180000)

	// Build the organization hierarchy
	teamLead.AddSubordinate(emp1)
	teamLead.AddSubordinate(emp2)
	teamLead.AddSubordinate(emp3)
	teamLead.AddSubordinate(emp4)

	cto.AddSubordinate(teamLead)
	cto.AddSubordinate(emp5)

	// Display the organization chart
	fmt.Println("Organization Hierarchy:")
	fmt.Println(cto.PrintHierarchy(""))

	// Calculate total salary costs (demonstrates uniform treatment)
	fmt.Printf("\nTotal salary cost for CTO's department: $%d\n", cto.GetSalary())
	fmt.Printf("Total salary cost for Engineering Lead's team: $%d\n", teamLead.GetSalary())
	fmt.Printf("Salary for individual employee: $%d\n", emp1.GetSalary())

	fmt.Println("\n" + strings.Repeat("─", 50))
}

// Example3_UIComponents demonstrates the Composite Pattern with UI components.
// This shows how simple components and container panels can be treated uniformly.
func Example3_UIComponents() {
	fmt.Println("=== Example 3: UI Components ===")

	// Create simple components (leaf nodes)
	usernameField := NewTextBox("Enter username", 200, 30)
	passwordField := NewTextBox("Enter password", 200, 30)
	loginButton := NewButton("Login", 100, 40)
	cancelButton := NewButton("Cancel", 100, 40)
	rememberMe := NewButton("Remember Me", 150, 25)

	// Create panels (composite nodes)
	loginPanel := NewPanel("login-panel", 400, 300)
	buttonPanel := NewPanel("button-panel", 400, 50)

	// Build the UI hierarchy
	loginPanel.Add(usernameField)
	loginPanel.Add(passwordField)
	loginPanel.Add(rememberMe)

	buttonPanel.Add(loginButton)
	buttonPanel.Add(cancelButton)

	loginPanel.Add(buttonPanel) // Nested panels

	// Render the UI (demonstrates uniform treatment)
	fmt.Println("Rendered UI:")
	fmt.Println(loginPanel.Render())

	fmt.Printf("\nTotal UI width: %d pixels\n", loginPanel.GetWidth())
	fmt.Printf("Total UI height: %d pixels\n", loginPanel.GetHeight())

	fmt.Println("\n" + strings.Repeat("─", 50))
}

// Example4_MenuSystem demonstrates the Composite Pattern with a menu system.
// This shows how menu items and submenus can be treated uniformly.
func Example4_MenuSystem() {
	fmt.Println("=== Example 4: Menu System ===")

	// Create menu items (leaf nodes)
	newFile := NewMenuItem("New", "Create new file")
	openFile := NewMenuItem("Open", "Open existing file")
	save := NewMenuItem("Save", "Save current file")
	saveAs := NewMenuItem("Save As", "Save with new name")
	exit := NewMenuItem("Exit", "Close application")

	copy := NewMenuItem("Copy", "Copy selection")
	paste := NewMenuItem("Paste", "Paste from clipboard")
	cut := NewMenuItem("Cut", "Cut selection")

	preferences := NewMenuItem("Preferences", "Open settings")
	about := NewMenuItem("About", "Show about dialog")

	// Create submenus (composite nodes)
	fileMenu := NewSubMenu("File")
	editMenu := NewSubMenu("Edit")
	helpMenu := NewSubMenu("Help")

	// Build the menu hierarchy
	fileMenu.Add(newFile)
	fileMenu.Add(openFile)
	fileMenu.Add(save)
	fileMenu.Add(saveAs)
	fileMenu.Add(exit)

	editMenu.Add(copy)
	editMenu.Add(paste)
	editMenu.Add(cut)

	helpMenu.Add(preferences)
	helpMenu.Add(about)

	// Create main menu
	mainMenu := NewSubMenu("Main Menu")
	mainMenu.Add(fileMenu)
	mainMenu.Add(editMenu)
	mainMenu.Add(helpMenu)

	// Display the menu structure
	fmt.Println("Menu Structure:")
	fmt.Println(mainMenu.Print(""))

	fmt.Println("\nExecuting 'File' submenu:")
	fileMenu.Execute()

	fmt.Println("\n" + strings.Repeat("─", 50))
}

// Example5_RealWorld demonstrates a real-world scenario: Document structure with sections.
// This shows how the Composite Pattern can represent complex document hierarchies.
func Example5_RealWorld() {
	fmt.Println("=== Example 5: Real-World Document Structure ===")

	// Create a complex document structure for a technical manual
	manual := NewDirectory("TechnicalManual")

	// Chapter 1: Introduction
	intro := NewDirectory("Chapter1_Introduction")
	intro.Add(NewFile("overview.md", 25))
	intro.Add(NewFile("prerequisites.md", 18))
	intro.Add(NewFile("conventions.md", 12))

	// Chapter 2: Getting Started
	gettingStarted := NewDirectory("Chapter2_GettingStarted")
	gettingStarted.Add(NewFile("installation.md", 45))
	gettingStarted.Add(NewFile("configuration.md", 38))
	gettingStarted.Add(NewFile("first_steps.md", 30))

	// Chapter 3: Advanced Topics (with subsections)
	advanced := NewDirectory("Chapter3_Advanced")

	advancedSection1 := NewDirectory("Subsection_Performance")
	advancedSection1.Add(NewFile("optimization.md", 55))
	advancedSection1.Add(NewFile("benchmarking.md", 42))

	advancedSection2 := NewDirectory("Subsection_Security")
	advancedSection2.Add(NewFile("authentication.md", 48))
	advancedSection2.Add(NewFile("authorization.md", 52))

	advanced.Add(advancedSection1)
	advanced.Add(advancedSection2)
	advanced.Add(NewFile("best_practices.md", 65))

	// Appendices
	appendices := NewDirectory("Appendices")
	appendices.Add(NewFile("glossary.md", 20))
	appendices.Add(NewFile("references.md", 15))
	appendices.Add(NewFile("troubleshooting.md", 35))

	// Build the complete manual
	manual.Add(intro)
	manual.Add(gettingStarted)
	manual.Add(advanced)
	manual.Add(appendices)
	manual.Add(NewFile("README.md", 10))

	// Display the document structure
	fmt.Println("Technical Manual Structure:")
	fmt.Println(manual.Display(""))

	// Demonstrate uniform operations
	fmt.Printf("\nTotal manual size: %d KB\n", manual.GetSize())
	fmt.Printf("Chapter 3 (Advanced) size: %d KB\n", advanced.GetSize())
	fmt.Printf("Performance subsection size: %d KB\n", advancedSection1.GetSize())

	// Show how we can work with any level uniformly
	fmt.Println("\nAll chapters:")
	for _, chapter := range manual.GetChildren() {
		if dir, ok := chapter.(*Directory); ok {
			fmt.Printf("  - %s: %d KB (%d items)\n",
				chapter.GetName(),
				chapter.GetSize(),
				len(dir.GetChildren()))
		} else {
			fmt.Printf("  - %s: %d KB (file)\n", chapter.GetName(), chapter.GetSize())
		}
	}

	fmt.Println("\n" + strings.Repeat("─", 50))
}
