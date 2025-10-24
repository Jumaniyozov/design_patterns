package command

import "fmt"

// Example1_TextEditorBasic demonstrates basic command execution with undo/redo.
func Example1_TextEditorBasic() {
	fmt.Println("=== Example 1: Basic Text Editor with Undo/Redo ===")

	editor := NewTextEditor()

	// Execute insert commands
	fmt.Println("Typing 'Hello World'...")
	editor.Execute(NewInsertCommand(editor.buffer, "Hello", 0))
	fmt.Printf("Content: '%s'\n", editor.GetContent())

	editor.Execute(NewInsertCommand(editor.buffer, " ", 5))
	fmt.Printf("Content: '%s'\n", editor.GetContent())

	editor.Execute(NewInsertCommand(editor.buffer, "World", 6))
	fmt.Printf("Content: '%s'\n\n", editor.GetContent())

	// Undo operations
	fmt.Println("Undo x3:")
	for i := 0; i < 3; i++ {
		editor.Undo()
		fmt.Printf("  After undo: '%s'\n", editor.GetContent())
	}

	// Redo operations
	fmt.Println("\nRedo x2:")
	for i := 0; i < 2; i++ {
		editor.Redo()
		fmt.Printf("  After redo: '%s'\n", editor.GetContent())
	}

	fmt.Println()
}

// Example2_TextEditorAdvanced demonstrates delete and replace operations.
func Example2_TextEditorAdvanced() {
	fmt.Println("=== Example 2: Advanced Text Editing ===")

	editor := NewTextEditor()

	// Build initial text
	fmt.Println("Building text: 'The quick brown fox'")
	editor.Execute(NewInsertCommand(editor.buffer, "The quick brown fox", 0))
	fmt.Printf("Content: '%s'\n\n", editor.GetContent())

	// Delete operation
	fmt.Println("Deleting 'quick ' (6 chars at position 4)...")
	editor.Execute(NewDeleteCommand(editor.buffer, 4, 6))
	fmt.Printf("Content: '%s'\n\n", editor.GetContent())

	// Replace operation
	fmt.Println("Replacing 'brown' with 'red'...")
	editor.Execute(NewReplaceCommand(editor.buffer, 4, 5, "red"))
	fmt.Printf("Content: '%s'\n\n", editor.GetContent())

	// Undo all
	fmt.Println("Undo all operations:")
	for editor.CanUndo() {
		editor.Undo()
		fmt.Printf("  '%s'\n", editor.GetContent())
	}

	fmt.Println()
}

// Example3_MacroCommands demonstrates combining multiple commands.
func Example3_MacroCommands() {
	fmt.Println("=== Example 3: Macro Commands ===")

	editor := NewTextEditor()

	// Create a macro that formats text as a title
	// Title format: "=== Text ==="
	fmt.Println("Creating macro: Format as Title")

	buffer := editor.buffer
	titleMacro := NewMacroCommand("FormatTitle",
		NewInsertCommand(buffer, "=== ", 0),
		NewInsertCommand(buffer, "My Title", 4),
		NewInsertCommand(buffer, " ===", 12),
	)

	fmt.Println("Executing macro...")
	editor.Execute(titleMacro)
	fmt.Printf("Content: '%s'\n\n", editor.GetContent())

	fmt.Println("Undoing macro (single undo)...")
	editor.Undo()
	fmt.Printf("Content: '%s'\n\n", editor.GetContent())

	fmt.Println("Redoing macro...")
	editor.Redo()
	fmt.Printf("Content: '%s'\n\n", editor.GetContent())

	// Create another macro for code block
	editor2 := NewTextEditor()
	codeBlockMacro := NewMacroCommand("CodeBlock",
		NewInsertCommand(editor2.buffer, "```go\n", 0),
		NewInsertCommand(editor2.buffer, "func main() {}\n", 6),
		NewInsertCommand(editor2.buffer, "```", 21),
	)

	fmt.Println("Creating code block macro...")
	editor2.Execute(codeBlockMacro)
	fmt.Printf("Content:\n%s\n", editor2.GetContent())

	fmt.Println()
}

// Example4_BankingTransactions demonstrates commands with business logic.
func Example4_BankingTransactions() {
	fmt.Println("=== Example 4: Banking Transactions ===")

	// Create accounts
	checking := NewAccount(1000.00)
	savings := NewAccount(500.00)

	fmt.Printf("Initial balances:\n")
	fmt.Printf("  Checking: $%.2f\n", checking.GetBalance())
	fmt.Printf("  Savings:  $%.2f\n\n", savings.GetBalance())

	// Create transaction history
	transactions := make([]Command, 0)

	// Execute some transactions
	fmt.Println("Executing transactions...")

	deposit := NewDepositCommand(checking, 200.00)
	deposit.Execute()
	transactions = append(transactions, deposit)
	fmt.Printf("  Deposit $200 to checking: $%.2f\n", checking.GetBalance())

	withdraw := NewWithdrawCommand(checking, 150.00)
	withdraw.Execute()
	transactions = append(transactions, withdraw)
	fmt.Printf("  Withdraw $150 from checking: $%.2f\n", checking.GetBalance())

	transfer := NewTransferCommand(checking, savings, 300.00)
	transfer.Execute()
	transactions = append(transactions, transfer)
	fmt.Printf("  Transfer $300 checking -> savings\n")
	fmt.Printf("    Checking: $%.2f\n", checking.GetBalance())
	fmt.Printf("    Savings:  $%.2f\n\n", savings.GetBalance())

	// Rollback last transaction
	fmt.Println("Rolling back last transaction (transfer)...")
	transactions[len(transactions)-1].Undo()
	fmt.Printf("  Checking: $%.2f\n", checking.GetBalance())
	fmt.Printf("  Savings:  $%.2f\n\n", savings.GetBalance())

	// Rollback all transactions
	fmt.Println("Rolling back all transactions...")
	for i := len(transactions) - 2; i >= 0; i-- {
		transactions[i].Undo()
	}
	fmt.Printf("  Checking: $%.2f\n", checking.GetBalance())
	fmt.Printf("  Savings:  $%.2f\n", savings.GetBalance())

	fmt.Println()
}

// Example5_AsyncCommandQueue demonstrates async command processing.
func Example5_AsyncCommandQueue() {
	fmt.Println("=== Example 5: Async Command Queue ===")

	// Create a command queue with 3 workers
	queue := NewCommandQueue(3, 10)

	fmt.Println("Submitting 10 commands to queue with 3 workers...")

	// Create 10 commands that simulate work
	for i := 0; i < 10; i++ {
		cmdNum := i + 1
		cmd := NewFunctionCommand(
			func() error {
				// Simulate work
				fmt.Printf("  [Worker] Executing command #%d\n", cmdNum)
				return nil
			},
			func() error {
				return nil
			},
		)
		queue.Submit(cmd)
	}

	// Close queue and wait for completion
	fmt.Println("\nWaiting for all commands to complete...")
	errors := queue.Close()

	if len(errors) > 0 {
		fmt.Printf("Completed with %d errors\n", len(errors))
	} else {
		fmt.Println("All commands completed successfully")
	}

	fmt.Println()
}

// Example6_FunctionCommands demonstrates using closures as commands.
func Example6_FunctionCommands() {
	fmt.Println("=== Example 6: Function Commands (Closures) ===")

	editor := NewTextEditor()

	// Create commands using closures
	fmt.Println("Creating commands with closures...")

	var insertedText string
	var insertPos int

	appendHello := NewFunctionCommand(
		func() error {
			insertedText = "Hello "
			insertPos = len(editor.buffer.GetContent())
			return editor.buffer.Insert(insertedText, insertPos)
		},
		func() error {
			_, err := editor.buffer.Delete(insertPos, len(insertedText))
			return err
		},
	)

	var insertedText2 string
	var insertPos2 int

	appendWorld := NewFunctionCommand(
		func() error {
			insertedText2 = "World!"
			insertPos2 = len(editor.buffer.GetContent())
			return editor.buffer.Insert(insertedText2, insertPos2)
		},
		func() error {
			_, err := editor.buffer.Delete(insertPos2, len(insertedText2))
			return err
		},
	)

	// Execute commands
	fmt.Println("Executing: Append 'Hello '")
	editor.Execute(appendHello)
	fmt.Printf("Content: '%s'\n\n", editor.GetContent())

	fmt.Println("Executing: Append 'World!'")
	editor.Execute(appendWorld)
	fmt.Printf("Content: '%s'\n\n", editor.GetContent())

	// Undo
	fmt.Println("Undo last command:")
	editor.Undo()
	fmt.Printf("Content: '%s'\n\n", editor.GetContent())

	// Redo
	fmt.Println("Redo:")
	editor.Redo()
	fmt.Printf("Content: '%s'\n\n", editor.GetContent())

	fmt.Println()
}

// Example7_ComplexMacro demonstrates a complex multi-step macro.
func Example7_ComplexMacro() {
	fmt.Println("=== Example 7: Complex Banking Macro ===")

	checking := NewAccount(1000.00)
	savings := NewAccount(500.00)
	investment := NewAccount(2000.00)

	fmt.Printf("Initial balances:\n")
	fmt.Printf("  Checking:   $%.2f\n", checking.GetBalance())
	fmt.Printf("  Savings:    $%.2f\n", savings.GetBalance())
	fmt.Printf("  Investment: $%.2f\n\n", investment.GetBalance())

	// Create a "month-end" macro
	// 1. Deposit paycheck to checking
	// 2. Transfer savings amount
	// 3. Transfer investment amount
	monthEndMacro := NewMacroCommand("MonthEnd",
		NewDepositCommand(checking, 3000.00),        // Paycheck
		NewTransferCommand(checking, savings, 500.00),    // Save
		NewTransferCommand(checking, investment, 1000.00), // Invest
	)

	fmt.Println("Executing 'Month-End' macro:")
	fmt.Println("  1. Deposit $3000 paycheck")
	fmt.Println("  2. Transfer $500 to savings")
	fmt.Println("  3. Transfer $1000 to investment")
	fmt.Println()

	if err := monthEndMacro.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("After month-end processing:\n")
	fmt.Printf("  Checking:   $%.2f\n", checking.GetBalance())
	fmt.Printf("  Savings:    $%.2f\n", savings.GetBalance())
	fmt.Printf("  Investment: $%.2f\n\n", investment.GetBalance())

	fmt.Println("Undoing month-end macro...")
	monthEndMacro.Undo()

	fmt.Printf("After undo:\n")
	fmt.Printf("  Checking:   $%.2f\n", checking.GetBalance())
	fmt.Printf("  Savings:    $%.2f\n", savings.GetBalance())
	fmt.Printf("  Investment: $%.2f\n\n", investment.GetBalance())

	fmt.Println()
}
