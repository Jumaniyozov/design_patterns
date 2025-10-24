package composite

import (
	"strings"
	"testing"
)

// =============================================================================
// File System Tests
// =============================================================================

func TestFile_GetName(t *testing.T) {
	file := NewFile("test.txt", 100)
	if file.GetName() != "test.txt" {
		t.Errorf("Expected name 'test.txt', got '%s'", file.GetName())
	}
}

func TestFile_GetSize(t *testing.T) {
	file := NewFile("test.txt", 100)
	if file.GetSize() != 100 {
		t.Errorf("Expected size 100, got %d", file.GetSize())
	}
}

func TestFile_Display(t *testing.T) {
	file := NewFile("test.txt", 100)
	display := file.Display("")
	if !strings.Contains(display, "test.txt") || !strings.Contains(display, "100") {
		t.Errorf("Display output incorrect: %s", display)
	}
}

func TestDirectory_GetName(t *testing.T) {
	dir := NewDirectory("testdir")
	if dir.GetName() != "testdir" {
		t.Errorf("Expected name 'testdir', got '%s'", dir.GetName())
	}
}

func TestDirectory_EmptySize(t *testing.T) {
	dir := NewDirectory("empty")
	if dir.GetSize() != 0 {
		t.Errorf("Expected empty directory size 0, got %d", dir.GetSize())
	}
}

func TestDirectory_AddAndGetSize(t *testing.T) {
	dir := NewDirectory("parent")
	file1 := NewFile("file1.txt", 100)
	file2 := NewFile("file2.txt", 200)

	dir.Add(file1)
	dir.Add(file2)

	expectedSize := 300
	if dir.GetSize() != expectedSize {
		t.Errorf("Expected directory size %d, got %d", expectedSize, dir.GetSize())
	}
}

func TestDirectory_NestedStructure(t *testing.T) {
	root := NewDirectory("root")
	subdir := NewDirectory("subdir")
	file1 := NewFile("file1.txt", 100)
	file2 := NewFile("file2.txt", 150)

	subdir.Add(file2)
	root.Add(file1)
	root.Add(subdir)

	expectedSize := 250 // 100 + 150
	if root.GetSize() != expectedSize {
		t.Errorf("Expected nested structure size %d, got %d", expectedSize, root.GetSize())
	}
}

func TestDirectory_Remove(t *testing.T) {
	dir := NewDirectory("test")
	file1 := NewFile("file1.txt", 100)
	file2 := NewFile("file2.txt", 200)

	dir.Add(file1)
	dir.Add(file2)

	if dir.GetSize() != 300 {
		t.Errorf("Expected initial size 300, got %d", dir.GetSize())
	}

	dir.Remove(file1)

	if dir.GetSize() != 200 {
		t.Errorf("Expected size after removal 200, got %d", dir.GetSize())
	}

	if len(dir.GetChildren()) != 1 {
		t.Errorf("Expected 1 child after removal, got %d", len(dir.GetChildren()))
	}
}

func TestDirectory_Display(t *testing.T) {
	dir := NewDirectory("root")
	file := NewFile("test.txt", 50)
	dir.Add(file)

	display := dir.Display("")
	if !strings.Contains(display, "root") || !strings.Contains(display, "test.txt") {
		t.Errorf("Display output missing expected content: %s", display)
	}
}

// =============================================================================
// Organization Chart Tests
// =============================================================================

func TestIndividualEmployee_GetName(t *testing.T) {
	emp := NewIndividualEmployee("John Doe", "Engineer", 80000)
	if emp.GetName() != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", emp.GetName())
	}
}

func TestIndividualEmployee_GetRole(t *testing.T) {
	emp := NewIndividualEmployee("John Doe", "Engineer", 80000)
	if emp.GetRole() != "Engineer" {
		t.Errorf("Expected role 'Engineer', got '%s'", emp.GetRole())
	}
}

func TestIndividualEmployee_GetSalary(t *testing.T) {
	emp := NewIndividualEmployee("John Doe", "Engineer", 80000)
	if emp.GetSalary() != 80000 {
		t.Errorf("Expected salary 80000, got %d", emp.GetSalary())
	}
}

func TestManager_GetName(t *testing.T) {
	mgr := NewManager("Jane Smith", "Manager", 100000)
	if mgr.GetName() != "Jane Smith" {
		t.Errorf("Expected name 'Jane Smith', got '%s'", mgr.GetName())
	}
}

func TestManager_GetRole(t *testing.T) {
	mgr := NewManager("Jane Smith", "Manager", 100000)
	if mgr.GetRole() != "Manager" {
		t.Errorf("Expected role 'Manager', got '%s'", mgr.GetRole())
	}
}

func TestManager_NoSubordinates(t *testing.T) {
	mgr := NewManager("Jane Smith", "Manager", 100000)
	if mgr.GetSalary() != 100000 {
		t.Errorf("Expected manager salary 100000, got %d", mgr.GetSalary())
	}
}

func TestManager_WithSubordinates(t *testing.T) {
	mgr := NewManager("Jane Smith", "Manager", 100000)
	emp1 := NewIndividualEmployee("Alice", "Engineer", 80000)
	emp2 := NewIndividualEmployee("Bob", "Engineer", 85000)

	mgr.AddSubordinate(emp1)
	mgr.AddSubordinate(emp2)

	expectedTotal := 100000 + 80000 + 85000
	if mgr.GetSalary() != expectedTotal {
		t.Errorf("Expected total salary %d, got %d", expectedTotal, mgr.GetSalary())
	}
}

func TestManager_NestedManagers(t *testing.T) {
	ceo := NewManager("CEO", "Chief Executive", 200000)
	vp := NewManager("VP", "Vice President", 150000)
	emp := NewIndividualEmployee("Worker", "Staff", 70000)

	vp.AddSubordinate(emp)
	ceo.AddSubordinate(vp)

	expectedTotal := 200000 + 150000 + 70000
	if ceo.GetSalary() != expectedTotal {
		t.Errorf("Expected nested total salary %d, got %d", expectedTotal, ceo.GetSalary())
	}
}

func TestManager_RemoveSubordinate(t *testing.T) {
	mgr := NewManager("Manager", "Lead", 100000)
	emp1 := NewIndividualEmployee("Alice", "Engineer", 80000)
	emp2 := NewIndividualEmployee("Bob", "Engineer", 85000)

	mgr.AddSubordinate(emp1)
	mgr.AddSubordinate(emp2)

	initialTotal := mgr.GetSalary()
	mgr.RemoveSubordinate(emp1)

	expectedAfterRemoval := 100000 + 85000
	if mgr.GetSalary() != expectedAfterRemoval {
		t.Errorf("Expected salary after removal %d, got %d", expectedAfterRemoval, mgr.GetSalary())
	}

	if len(mgr.GetSubordinates()) != 1 {
		t.Errorf("Expected 1 subordinate after removal, got %d", len(mgr.GetSubordinates()))
	}

	_ = initialTotal // Used for clarity
}

// =============================================================================
// UI Components Tests
// =============================================================================

func TestButton_Render(t *testing.T) {
	button := NewButton("Click Me", 100, 50)
	rendered := button.Render()
	if !strings.Contains(rendered, "Click Me") {
		t.Errorf("Rendered button missing text: %s", rendered)
	}
}

func TestButton_GetDimensions(t *testing.T) {
	button := NewButton("Click Me", 100, 50)
	if button.GetWidth() != 100 {
		t.Errorf("Expected width 100, got %d", button.GetWidth())
	}
	if button.GetHeight() != 50 {
		t.Errorf("Expected height 50, got %d", button.GetHeight())
	}
}

func TestTextBox_Render(t *testing.T) {
	textbox := NewTextBox("Enter text", 200, 30)
	rendered := textbox.Render()
	if !strings.Contains(rendered, "Enter text") {
		t.Errorf("Rendered textbox missing placeholder: %s", rendered)
	}
}

func TestPanel_EmptyRender(t *testing.T) {
	panel := NewPanel("test-panel", 400, 300)
	rendered := panel.Render()
	if !strings.Contains(rendered, "test-panel") {
		t.Errorf("Rendered panel missing name: %s", rendered)
	}
}

func TestPanel_WithComponents(t *testing.T) {
	panel := NewPanel("login-panel", 400, 300)
	button := NewButton("Login", 100, 40)
	textbox := NewTextBox("Username", 200, 30)

	panel.Add(button)
	panel.Add(textbox)

	rendered := panel.Render()
	if !strings.Contains(rendered, "Login") || !strings.Contains(rendered, "Username") {
		t.Errorf("Rendered panel missing components: %s", rendered)
	}
}

func TestPanel_NestedPanels(t *testing.T) {
	outerPanel := NewPanel("outer", 500, 400)
	innerPanel := NewPanel("inner", 400, 300)
	button := NewButton("Click", 100, 40)

	innerPanel.Add(button)
	outerPanel.Add(innerPanel)

	rendered := outerPanel.Render()
	if !strings.Contains(rendered, "outer") ||
		!strings.Contains(rendered, "inner") ||
		!strings.Contains(rendered, "Click") {
		t.Errorf("Rendered nested panels missing content: %s", rendered)
	}
}

func TestPanel_Remove(t *testing.T) {
	panel := NewPanel("test", 400, 300)
	button1 := NewButton("Button1", 100, 40)
	button2 := NewButton("Button2", 100, 40)

	panel.Add(button1)
	panel.Add(button2)

	if len(panel.GetComponents()) != 2 {
		t.Errorf("Expected 2 components, got %d", len(panel.GetComponents()))
	}

	panel.Remove(button1)

	if len(panel.GetComponents()) != 1 {
		t.Errorf("Expected 1 component after removal, got %d", len(panel.GetComponents()))
	}
}

// =============================================================================
// Menu System Tests
// =============================================================================

func TestMenuItem_GetName(t *testing.T) {
	item := NewMenuItem("File", "Open file")
	if item.GetName() != "File" {
		t.Errorf("Expected name 'File', got '%s'", item.GetName())
	}
}

func TestMenuItem_Print(t *testing.T) {
	item := NewMenuItem("Save", "Save file")
	printed := item.Print("")
	if !strings.Contains(printed, "Save") {
		t.Errorf("Printed menu item missing name: %s", printed)
	}
}

func TestSubMenu_GetName(t *testing.T) {
	menu := NewSubMenu("File Menu")
	if menu.GetName() != "File Menu" {
		t.Errorf("Expected name 'File Menu', got '%s'", menu.GetName())
	}
}

func TestSubMenu_AddItems(t *testing.T) {
	menu := NewSubMenu("File")
	item1 := NewMenuItem("New", "Create new")
	item2 := NewMenuItem("Open", "Open file")

	menu.Add(item1)
	menu.Add(item2)

	if len(menu.GetItems()) != 2 {
		t.Errorf("Expected 2 items, got %d", len(menu.GetItems()))
	}
}

func TestSubMenu_NestedMenus(t *testing.T) {
	mainMenu := NewSubMenu("Main")
	fileMenu := NewSubMenu("File")
	newItem := NewMenuItem("New", "Create new")

	fileMenu.Add(newItem)
	mainMenu.Add(fileMenu)

	printed := mainMenu.Print("")
	if !strings.Contains(printed, "Main") ||
		!strings.Contains(printed, "File") ||
		!strings.Contains(printed, "New") {
		t.Errorf("Printed nested menu missing content: %s", printed)
	}
}

func TestSubMenu_Remove(t *testing.T) {
	menu := NewSubMenu("Edit")
	item1 := NewMenuItem("Copy", "Copy text")
	item2 := NewMenuItem("Paste", "Paste text")

	menu.Add(item1)
	menu.Add(item2)

	menu.Remove(item1)

	if len(menu.GetItems()) != 1 {
		t.Errorf("Expected 1 item after removal, got %d", len(menu.GetItems()))
	}
}

// =============================================================================
// Integration Tests
// =============================================================================

func TestComposite_DeepNesting(t *testing.T) {
	// Create a deeply nested structure
	root := NewDirectory("root")
	level1 := NewDirectory("level1")
	level2 := NewDirectory("level2")
	level3 := NewDirectory("level3")
	file := NewFile("deep.txt", 100)

	level3.Add(file)
	level2.Add(level3)
	level1.Add(level2)
	root.Add(level1)

	// Should correctly calculate size through all levels
	if root.GetSize() != 100 {
		t.Errorf("Expected size 100 through deep nesting, got %d", root.GetSize())
	}
}

func TestComposite_MultipleChildren(t *testing.T) {
	// Test with many children at each level
	root := NewDirectory("root")

	for i := 0; i < 10; i++ {
		root.Add(NewFile("file", 10))
	}

	if root.GetSize() != 100 {
		t.Errorf("Expected size 100 with 10 children, got %d", root.GetSize())
	}
}

func TestComposite_MixedStructure(t *testing.T) {
	// Create a mixed structure with files and directories at same level
	root := NewDirectory("root")
	root.Add(NewFile("file1.txt", 50))
	root.Add(NewDirectory("dir1"))
	root.Add(NewFile("file2.txt", 75))

	subdir := NewDirectory("dir2")
	subdir.Add(NewFile("file3.txt", 25))
	root.Add(subdir)

	expectedSize := 50 + 75 + 25 // Files only
	if root.GetSize() != expectedSize {
		t.Errorf("Expected mixed structure size %d, got %d", expectedSize, root.GetSize())
	}
}
