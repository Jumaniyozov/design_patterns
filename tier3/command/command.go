// Package command implements the Command pattern.
//
// The Command pattern encapsulates a request as an object, thereby allowing you to
// parameterize clients with different requests, queue or log requests, and support
// undoable operations.
//
// Key components:
// - Command: Interface defining Execute() and Undo() methods
// - ConcreteCommand: Specific command implementation that calls receiver
// - Receiver: The object that performs the actual work
// - Invoker: Triggers command execution without knowing details
// - Client: Creates commands and sets their receivers
package command

import (
	"errors"
	"fmt"
	"sync"
)

// Command defines the interface for executing an operation.
type Command interface {
	Execute() error
	Undo() error
}

// TextBuffer is the receiver that performs actual text operations.
type TextBuffer struct {
	content string
	mu      sync.RWMutex
}

// NewTextBuffer creates a new text buffer.
func NewTextBuffer() *TextBuffer {
	return &TextBuffer{
		content: "",
	}
}

// Insert inserts text at the specified position.
func (tb *TextBuffer) Insert(text string, pos int) error {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if pos < 0 || pos > len(tb.content) {
		return fmt.Errorf("invalid position: %d", pos)
	}

	tb.content = tb.content[:pos] + text + tb.content[pos:]
	return nil
}

// Delete removes text from the specified position.
func (tb *TextBuffer) Delete(pos, length int) (string, error) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if pos < 0 || pos >= len(tb.content) {
		return "", fmt.Errorf("invalid position: %d", pos)
	}

	if pos+length > len(tb.content) {
		length = len(tb.content) - pos
	}

	deleted := tb.content[pos : pos+length]
	tb.content = tb.content[:pos] + tb.content[pos+length:]
	return deleted, nil
}

// GetContent returns the current buffer content.
func (tb *TextBuffer) GetContent() string {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	return tb.content
}

// SetContent sets the buffer content (used for undo).
func (tb *TextBuffer) SetContent(content string) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.content = content
}

// InsertCommand implements the Command interface for inserting text.
type InsertCommand struct {
	buffer *TextBuffer
	text   string
	pos    int
}

// NewInsertCommand creates a new insert command.
func NewInsertCommand(buffer *TextBuffer, text string, pos int) *InsertCommand {
	return &InsertCommand{
		buffer: buffer,
		text:   text,
		pos:    pos,
	}
}

// Execute performs the insert operation.
func (c *InsertCommand) Execute() error {
	return c.buffer.Insert(c.text, c.pos)
}

// Undo reverses the insert operation.
func (c *InsertCommand) Undo() error {
	_, err := c.buffer.Delete(c.pos, len(c.text))
	return err
}

// DeleteCommand implements the Command interface for deleting text.
type DeleteCommand struct {
	buffer      *TextBuffer
	pos         int
	length      int
	deletedText string // Stored for undo
}

// NewDeleteCommand creates a new delete command.
func NewDeleteCommand(buffer *TextBuffer, pos, length int) *DeleteCommand {
	return &DeleteCommand{
		buffer: buffer,
		pos:    pos,
		length: length,
	}
}

// Execute performs the delete operation.
func (c *DeleteCommand) Execute() error {
	deleted, err := c.buffer.Delete(c.pos, c.length)
	if err != nil {
		return err
	}
	c.deletedText = deleted // Store for undo
	return nil
}

// Undo reverses the delete operation.
func (c *DeleteCommand) Undo() error {
	return c.buffer.Insert(c.deletedText, c.pos)
}

// ReplaceCommand implements the Command interface for replacing text.
type ReplaceCommand struct {
	buffer      *TextBuffer
	pos         int
	oldText     string
	newText     string
	replacedLen int
}

// NewReplaceCommand creates a new replace command.
func NewReplaceCommand(buffer *TextBuffer, pos int, oldLen int, newText string) *ReplaceCommand {
	return &ReplaceCommand{
		buffer:      buffer,
		pos:         pos,
		newText:     newText,
		replacedLen: oldLen,
	}
}

// Execute performs the replace operation.
func (c *ReplaceCommand) Execute() error {
	// Delete old text
	deleted, err := c.buffer.Delete(c.pos, c.replacedLen)
	if err != nil {
		return err
	}
	c.oldText = deleted // Store for undo

	// Insert new text
	return c.buffer.Insert(c.newText, c.pos)
}

// Undo reverses the replace operation.
func (c *ReplaceCommand) Undo() error {
	// Delete new text
	_, err := c.buffer.Delete(c.pos, len(c.newText))
	if err != nil {
		return err
	}

	// Restore old text
	return c.buffer.Insert(c.oldText, c.pos)
}

// TextEditor is the invoker that executes commands and manages history.
type TextEditor struct {
	buffer       *TextBuffer
	history      []Command
	currentIndex int // Points to the next position to add command
	mu           sync.Mutex
}

// NewTextEditor creates a new text editor.
func NewTextEditor() *TextEditor {
	return &TextEditor{
		buffer:       NewTextBuffer(),
		history:      make([]Command, 0),
		currentIndex: 0,
	}
}

// Execute executes a command and adds it to history.
func (e *TextEditor) Execute(cmd Command) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := cmd.Execute(); err != nil {
		return err
	}

	// Truncate history if we're not at the end (after undos)
	e.history = e.history[:e.currentIndex]

	// Add command to history
	e.history = append(e.history, cmd)
	e.currentIndex++

	return nil
}

// Undo undoes the last command.
func (e *TextEditor) Undo() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.currentIndex == 0 {
		return errors.New("nothing to undo")
	}

	e.currentIndex--
	return e.history[e.currentIndex].Undo()
}

// Redo redoes the last undone command.
func (e *TextEditor) Redo() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.currentIndex >= len(e.history) {
		return errors.New("nothing to redo")
	}

	err := e.history[e.currentIndex].Execute()
	if err != nil {
		return err
	}

	e.currentIndex++
	return nil
}

// GetContent returns the current buffer content.
func (e *TextEditor) GetContent() string {
	return e.buffer.GetContent()
}

// CanUndo returns true if undo is possible.
func (e *TextEditor) CanUndo() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.currentIndex > 0
}

// CanRedo returns true if redo is possible.
func (e *TextEditor) CanRedo() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.currentIndex < len(e.history)
}

// MacroCommand combines multiple commands into one.
type MacroCommand struct {
	commands []Command
	name     string
}

// NewMacroCommand creates a new macro command.
func NewMacroCommand(name string, commands ...Command) *MacroCommand {
	return &MacroCommand{
		name:     name,
		commands: commands,
	}
}

// Execute executes all commands in sequence.
func (m *MacroCommand) Execute() error {
	for i, cmd := range m.commands {
		if err := cmd.Execute(); err != nil {
			// Rollback previously executed commands
			for j := i - 1; j >= 0; j-- {
				m.commands[j].Undo()
			}
			return fmt.Errorf("macro command '%s' failed at step %d: %w", m.name, i, err)
		}
	}
	return nil
}

// Undo undoes all commands in reverse order.
func (m *MacroCommand) Undo() error {
	// Undo in reverse order
	for i := len(m.commands) - 1; i >= 0; i-- {
		if err := m.commands[i].Undo(); err != nil {
			return fmt.Errorf("macro command '%s' undo failed at step %d: %w", m.name, i, err)
		}
	}
	return nil
}

// CommandQueue implements an async command execution queue.
type CommandQueue struct {
	commands chan Command
	workers  int
	wg       sync.WaitGroup
	mu       sync.Mutex
	results  []error
}

// NewCommandQueue creates a new command queue with specified number of workers.
func NewCommandQueue(workers int, bufferSize int) *CommandQueue {
	cq := &CommandQueue{
		commands: make(chan Command, bufferSize),
		workers:  workers,
		results:  make([]error, 0),
	}

	// Start worker goroutines
	for i := 0; i < workers; i++ {
		cq.wg.Add(1)
		go cq.worker(i)
	}

	return cq
}

// worker processes commands from the queue.
func (cq *CommandQueue) worker(id int) {
	defer cq.wg.Done()

	for cmd := range cq.commands {
		err := cmd.Execute()
		if err != nil {
			cq.mu.Lock()
			cq.results = append(cq.results, err)
			cq.mu.Unlock()
		}
	}
}

// Submit adds a command to the queue for async execution.
func (cq *CommandQueue) Submit(cmd Command) {
	cq.commands <- cmd
}

// Close closes the queue and waits for all commands to complete.
func (cq *CommandQueue) Close() []error {
	close(cq.commands)
	cq.wg.Wait()

	cq.mu.Lock()
	defer cq.mu.Unlock()
	return cq.results
}

// FunctionCommand wraps a function as a Command.
// This demonstrates Go's idiomatic approach using closures.
type FunctionCommand struct {
	execute func() error
	undo    func() error
}

// NewFunctionCommand creates a command from execute and undo functions.
func NewFunctionCommand(execute, undo func() error) *FunctionCommand {
	return &FunctionCommand{
		execute: execute,
		undo:    undo,
	}
}

// Execute executes the function.
func (f *FunctionCommand) Execute() error {
	if f.execute != nil {
		return f.execute()
	}
	return nil
}

// Undo undoes the function.
func (f *FunctionCommand) Undo() error {
	if f.undo != nil {
		return f.undo()
	}
	return errors.New("undo not supported")
}

// Account is a receiver for banking operations.
type Account struct {
	balance float64
	mu      sync.Mutex
}

// NewAccount creates a new account with initial balance.
func NewAccount(initialBalance float64) *Account {
	return &Account{
		balance: initialBalance,
	}
}

// Deposit adds money to the account.
func (a *Account) Deposit(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if amount <= 0 {
		return errors.New("deposit amount must be positive")
	}

	a.balance += amount
	return nil
}

// Withdraw removes money from the account.
func (a *Account) Withdraw(amount float64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if amount <= 0 {
		return errors.New("withdraw amount must be positive")
	}

	if a.balance < amount {
		return errors.New("insufficient funds")
	}

	a.balance -= amount
	return nil
}

// GetBalance returns the current balance.
func (a *Account) GetBalance() float64 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

// DepositCommand implements a deposit operation.
type DepositCommand struct {
	account *Account
	amount  float64
}

// NewDepositCommand creates a new deposit command.
func NewDepositCommand(account *Account, amount float64) *DepositCommand {
	return &DepositCommand{
		account: account,
		amount:  amount,
	}
}

// Execute performs the deposit.
func (c *DepositCommand) Execute() error {
	return c.account.Deposit(c.amount)
}

// Undo reverses the deposit.
func (c *DepositCommand) Undo() error {
	return c.account.Withdraw(c.amount)
}

// WithdrawCommand implements a withdraw operation.
type WithdrawCommand struct {
	account *Account
	amount  float64
}

// NewWithdrawCommand creates a new withdraw command.
func NewWithdrawCommand(account *Account, amount float64) *WithdrawCommand {
	return &WithdrawCommand{
		account: account,
		amount:  amount,
	}
}

// Execute performs the withdrawal.
func (c *WithdrawCommand) Execute() error {
	return c.account.Withdraw(c.amount)
}

// Undo reverses the withdrawal.
func (c *WithdrawCommand) Undo() error {
	return c.account.Deposit(c.amount)
}

// TransferCommand implements a transfer between accounts.
type TransferCommand struct {
	from   *Account
	to     *Account
	amount float64
}

// NewTransferCommand creates a new transfer command.
func NewTransferCommand(from, to *Account, amount float64) *TransferCommand {
	return &TransferCommand{
		from:   from,
		to:     to,
		amount: amount,
	}
}

// Execute performs the transfer.
func (c *TransferCommand) Execute() error {
	if err := c.from.Withdraw(c.amount); err != nil {
		return err
	}

	if err := c.to.Deposit(c.amount); err != nil {
		// Rollback withdrawal if deposit fails
		c.from.Deposit(c.amount)
		return err
	}

	return nil
}

// Undo reverses the transfer.
func (c *TransferCommand) Undo() error {
	if err := c.to.Withdraw(c.amount); err != nil {
		return err
	}

	if err := c.from.Deposit(c.amount); err != nil {
		// Rollback if restoration fails
		c.to.Deposit(c.amount)
		return err
	}

	return nil
}
