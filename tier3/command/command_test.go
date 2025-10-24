package command

import (
	"sync"
	"testing"
)

func TestTextBuffer_Insert(t *testing.T) {
	buffer := NewTextBuffer()

	err := buffer.Insert("Hello", 0)
	if err != nil {
		t.Errorf("Insert failed: %v", err)
	}

	if buffer.GetContent() != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", buffer.GetContent())
	}

	// Insert in middle
	err = buffer.Insert(" World", 5)
	if err != nil {
		t.Errorf("Insert failed: %v", err)
	}

	if buffer.GetContent() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", buffer.GetContent())
	}
}

func TestTextBuffer_Delete(t *testing.T) {
	buffer := NewTextBuffer()
	buffer.SetContent("Hello World")

	deleted, err := buffer.Delete(5, 6)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	if deleted != " World" {
		t.Errorf("Expected deleted text ' World', got '%s'", deleted)
	}

	if buffer.GetContent() != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", buffer.GetContent())
	}
}

func TestInsertCommand_ExecuteUndo(t *testing.T) {
	buffer := NewTextBuffer()
	cmd := NewInsertCommand(buffer, "Hello", 0)

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	if buffer.GetContent() != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", buffer.GetContent())
	}

	// Undo
	err = cmd.Undo()
	if err != nil {
		t.Errorf("Undo failed: %v", err)
	}

	if buffer.GetContent() != "" {
		t.Errorf("Expected empty buffer, got '%s'", buffer.GetContent())
	}
}

func TestDeleteCommand_ExecuteUndo(t *testing.T) {
	buffer := NewTextBuffer()
	buffer.SetContent("Hello World")

	cmd := NewDeleteCommand(buffer, 5, 6)

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	if buffer.GetContent() != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", buffer.GetContent())
	}

	// Undo
	err = cmd.Undo()
	if err != nil {
		t.Errorf("Undo failed: %v", err)
	}

	if buffer.GetContent() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", buffer.GetContent())
	}
}

func TestReplaceCommand_ExecuteUndo(t *testing.T) {
	buffer := NewTextBuffer()
	buffer.SetContent("Hello World")

	cmd := NewReplaceCommand(buffer, 6, 5, "Go")

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	if buffer.GetContent() != "Hello Go" {
		t.Errorf("Expected 'Hello Go', got '%s'", buffer.GetContent())
	}

	// Undo
	err = cmd.Undo()
	if err != nil {
		t.Errorf("Undo failed: %v", err)
	}

	if buffer.GetContent() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", buffer.GetContent())
	}
}

func TestTextEditor_UndoRedo(t *testing.T) {
	editor := NewTextEditor()

	// Execute commands
	editor.Execute(NewInsertCommand(editor.buffer, "A", 0))
	editor.Execute(NewInsertCommand(editor.buffer, "B", 1))
	editor.Execute(NewInsertCommand(editor.buffer, "C", 2))

	if editor.GetContent() != "ABC" {
		t.Errorf("Expected 'ABC', got '%s'", editor.GetContent())
	}

	// Undo once
	editor.Undo()
	if editor.GetContent() != "AB" {
		t.Errorf("Expected 'AB' after undo, got '%s'", editor.GetContent())
	}

	// Undo twice
	editor.Undo()
	if editor.GetContent() != "A" {
		t.Errorf("Expected 'A' after undo, got '%s'", editor.GetContent())
	}

	// Redo once
	editor.Redo()
	if editor.GetContent() != "AB" {
		t.Errorf("Expected 'AB' after redo, got '%s'", editor.GetContent())
	}

	// Redo twice
	editor.Redo()
	if editor.GetContent() != "ABC" {
		t.Errorf("Expected 'ABC' after redo, got '%s'", editor.GetContent())
	}
}

func TestTextEditor_UndoAfterNewCommand(t *testing.T) {
	editor := NewTextEditor()

	// Execute commands
	editor.Execute(NewInsertCommand(editor.buffer, "A", 0))
	editor.Execute(NewInsertCommand(editor.buffer, "B", 1))
	editor.Execute(NewInsertCommand(editor.buffer, "C", 2))

	// Undo twice
	editor.Undo()
	editor.Undo()

	if editor.GetContent() != "A" {
		t.Errorf("Expected 'A', got '%s'", editor.GetContent())
	}

	// Execute new command - should truncate redo history
	editor.Execute(NewInsertCommand(editor.buffer, "X", 1))

	if editor.GetContent() != "AX" {
		t.Errorf("Expected 'AX', got '%s'", editor.GetContent())
	}

	// Should not be able to redo
	err := editor.Redo()
	if err == nil {
		t.Error("Expected redo to fail after new command")
	}
}

func TestMacroCommand_Execute(t *testing.T) {
	buffer := NewTextBuffer()

	macro := NewMacroCommand("Test",
		NewInsertCommand(buffer, "A", 0),
		NewInsertCommand(buffer, "B", 1),
		NewInsertCommand(buffer, "C", 2),
	)

	err := macro.Execute()
	if err != nil {
		t.Errorf("Macro execute failed: %v", err)
	}

	if buffer.GetContent() != "ABC" {
		t.Errorf("Expected 'ABC', got '%s'", buffer.GetContent())
	}
}

func TestMacroCommand_Undo(t *testing.T) {
	buffer := NewTextBuffer()

	macro := NewMacroCommand("Test",
		NewInsertCommand(buffer, "A", 0),
		NewInsertCommand(buffer, "B", 1),
		NewInsertCommand(buffer, "C", 2),
	)

	macro.Execute()
	err := macro.Undo()
	if err != nil {
		t.Errorf("Macro undo failed: %v", err)
	}

	if buffer.GetContent() != "" {
		t.Errorf("Expected empty buffer, got '%s'", buffer.GetContent())
	}
}

func TestMacroCommand_RollbackOnError(t *testing.T) {
	buffer := NewTextBuffer()
	buffer.SetContent("AB")

	// Create macro where middle command will fail
	macro := NewMacroCommand("Test",
		NewInsertCommand(buffer, "X", 2),
		NewDeleteCommand(buffer, 10, 5), // This will fail - invalid position
		NewInsertCommand(buffer, "Y", 3),
	)

	err := macro.Execute()
	if err == nil {
		t.Error("Expected macro to fail")
	}

	// First command should be rolled back
	if buffer.GetContent() != "AB" {
		t.Errorf("Expected 'AB' after rollback, got '%s'", buffer.GetContent())
	}
}

func TestAccount_Operations(t *testing.T) {
	account := NewAccount(100.00)

	// Test deposit
	err := account.Deposit(50.00)
	if err != nil {
		t.Errorf("Deposit failed: %v", err)
	}

	if account.GetBalance() != 150.00 {
		t.Errorf("Expected balance 150.00, got %.2f", account.GetBalance())
	}

	// Test withdraw
	err = account.Withdraw(30.00)
	if err != nil {
		t.Errorf("Withdraw failed: %v", err)
	}

	if account.GetBalance() != 120.00 {
		t.Errorf("Expected balance 120.00, got %.2f", account.GetBalance())
	}

	// Test insufficient funds
	err = account.Withdraw(200.00)
	if err == nil {
		t.Error("Expected withdraw to fail with insufficient funds")
	}
}

func TestDepositCommand(t *testing.T) {
	account := NewAccount(100.00)
	cmd := NewDepositCommand(account, 50.00)

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	if account.GetBalance() != 150.00 {
		t.Errorf("Expected balance 150.00, got %.2f", account.GetBalance())
	}

	// Undo
	err = cmd.Undo()
	if err != nil {
		t.Errorf("Undo failed: %v", err)
	}

	if account.GetBalance() != 100.00 {
		t.Errorf("Expected balance 100.00, got %.2f", account.GetBalance())
	}
}

func TestTransferCommand(t *testing.T) {
	from := NewAccount(200.00)
	to := NewAccount(100.00)

	cmd := NewTransferCommand(from, to, 50.00)

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	if from.GetBalance() != 150.00 {
		t.Errorf("Expected from balance 150.00, got %.2f", from.GetBalance())
	}

	if to.GetBalance() != 150.00 {
		t.Errorf("Expected to balance 150.00, got %.2f", to.GetBalance())
	}

	// Undo
	err = cmd.Undo()
	if err != nil {
		t.Errorf("Undo failed: %v", err)
	}

	if from.GetBalance() != 200.00 {
		t.Errorf("Expected from balance 200.00 after undo, got %.2f", from.GetBalance())
	}

	if to.GetBalance() != 100.00 {
		t.Errorf("Expected to balance 100.00 after undo, got %.2f", to.GetBalance())
	}
}

func TestCommandQueue(t *testing.T) {
	queue := NewCommandQueue(2, 10)

	executed := 0
	var mu sync.Mutex

	// Submit 10 commands
	for i := 0; i < 10; i++ {
		cmd := NewFunctionCommand(
			func() error {
				mu.Lock()
				executed++
				mu.Unlock()
				return nil
			},
			nil,
		)
		queue.Submit(cmd)
	}

	// Wait for completion
	errors := queue.Close()

	if len(errors) > 0 {
		t.Errorf("Expected no errors, got %d", len(errors))
	}

	mu.Lock()
	if executed != 10 {
		t.Errorf("Expected 10 executions, got %d", executed)
	}
	mu.Unlock()
}

func TestFunctionCommand(t *testing.T) {
	executed := false
	undone := false

	cmd := NewFunctionCommand(
		func() error {
			executed = true
			return nil
		},
		func() error {
			undone = true
			return nil
		},
	)

	// Execute
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute failed: %v", err)
	}

	if !executed {
		t.Error("Execute function not called")
	}

	// Undo
	err = cmd.Undo()
	if err != nil {
		t.Errorf("Undo failed: %v", err)
	}

	if !undone {
		t.Error("Undo function not called")
	}
}

func TestConcurrentCommandExecution(t *testing.T) {
	editor := NewTextEditor()
	var wg sync.WaitGroup

	// Execute 100 commands concurrently
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			cmd := NewInsertCommand(editor.buffer, "X", 0)
			editor.Execute(cmd)
		}(i)
	}

	wg.Wait()

	// Should have 100 X's
	content := editor.GetContent()
	if len(content) != 100 {
		t.Errorf("Expected 100 characters, got %d", len(content))
	}
}

// Benchmark tests
func BenchmarkInsertCommand(b *testing.B) {
	buffer := NewTextBuffer()
	cmd := NewInsertCommand(buffer, "test", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd.Execute()
		cmd.Undo()
	}
}

func BenchmarkTextEditor_UndoRedo(b *testing.B) {
	editor := NewTextEditor()

	// Build up some history
	for i := 0; i < 100; i++ {
		editor.Execute(NewInsertCommand(editor.buffer, "X", 0))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			editor.Undo()
		} else {
			editor.Redo()
		}
	}
}

func BenchmarkMacroCommand(b *testing.B) {
	buffer := NewTextBuffer()

	macro := NewMacroCommand("Test",
		NewInsertCommand(buffer, "A", 0),
		NewInsertCommand(buffer, "B", 1),
		NewInsertCommand(buffer, "C", 2),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		macro.Execute()
		macro.Undo()
	}
}
