// Package unitofwork demonstrates the Unit of Work pattern.
// It maintains a list of objects affected by a business transaction and
// coordinates writing changes as a single atomic operation.
package unitofwork

import (
	"errors"
	"fmt"
)

// Entity represents a domain entity
type Entity interface {
	GetID() string
}

// User entity
type User struct {
	ID    string
	Name  string
	Email string
}

func (u *User) GetID() string { return u.ID }

// Order entity
type Order struct {
	ID       string
	UserID   string
	Total    float64
	Items    []string
}

func (o *Order) GetID() string { return o.ID }

// UnitOfWork tracks changes and commits them atomically
type UnitOfWork struct {
	newEntities     []Entity
	dirtyEntities   []Entity
	deletedEntities []Entity
	committed       bool
}

// NewUnitOfWork creates a new unit of work
func NewUnitOfWork() *UnitOfWork {
	return &UnitOfWork{
		newEntities:     make([]Entity, 0),
		dirtyEntities:   make([]Entity, 0),
		deletedEntities: make([]Entity, 0),
		committed:       false,
	}
}

// RegisterNew registers a new entity
func (uow *UnitOfWork) RegisterNew(entity Entity) {
	uow.newEntities = append(uow.newEntities, entity)
}

// RegisterDirty registers a modified entity
func (uow *UnitOfWork) RegisterDirty(entity Entity) {
	uow.dirtyEntities = append(uow.dirtyEntities, entity)
}

// RegisterDeleted registers a deleted entity
func (uow *UnitOfWork) RegisterDeleted(entity Entity) {
	uow.deletedEntities = append(uow.deletedEntities, entity)
}

// Commit commits all changes atomically
func (uow *UnitOfWork) Commit() error {
	if uow.committed {
		return errors.New("unit of work already committed")
	}

	// Begin transaction
	fmt.Println("=== Begin Transaction ===")

	// Insert new entities
	for _, entity := range uow.newEntities {
		fmt.Printf("INSERT: %T with ID %s\n", entity, entity.GetID())
	}

	// Update dirty entities
	for _, entity := range uow.dirtyEntities {
		fmt.Printf("UPDATE: %T with ID %s\n", entity, entity.GetID())
	}

	// Delete entities
	for _, entity := range uow.deletedEntities {
		fmt.Printf("DELETE: %T with ID %s\n", entity, entity.GetID())
	}

	// Commit transaction
	fmt.Println("=== Commit Transaction ===")
	uow.committed = true

	return nil
}

// Rollback discards all changes
func (uow *UnitOfWork) Rollback() {
	fmt.Println("=== Rollback Transaction ===")
	uow.newEntities = make([]Entity, 0)
	uow.dirtyEntities = make([]Entity, 0)
	uow.deletedEntities = make([]Entity, 0)
}

// GetStats returns statistics about pending changes
func (uow *UnitOfWork) GetStats() string {
	return fmt.Sprintf("New: %d, Modified: %d, Deleted: %d",
		len(uow.newEntities),
		len(uow.dirtyEntities),
		len(uow.deletedEntities))
}

// OrderService demonstrates using unit of work
type OrderService struct {
	uow *UnitOfWork
}

// NewOrderService creates an order service
func NewOrderService(uow *UnitOfWork) *OrderService {
	return &OrderService{uow: uow}
}

// CreateOrder creates an order with related entities
func (os *OrderService) CreateOrder(userID string, items []string, total float64) (*Order, error) {
	// Create user if needed
	user := &User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example.com",
	}
	os.uow.RegisterNew(user)

	// Create order
	order := &Order{
		ID:     fmt.Sprintf("order-%s", userID),
		UserID: userID,
		Total:  total,
		Items:  items,
	}
	os.uow.RegisterNew(order)

	return order, nil
}

// UpdateOrder updates an order
func (os *OrderService) UpdateOrder(order *Order) {
	os.uow.RegisterDirty(order)
}

// CancelOrder cancels an order
func (os *OrderService) CancelOrder(order *Order) {
	os.uow.RegisterDeleted(order)
}

// Complete completes the transaction
func (os *OrderService) Complete() error {
	return os.uow.Commit()
}

// Transaction represents a database transaction
type Transaction struct {
	operations []func() error
	committed  bool
}

// NewTransaction creates a transaction
func NewTransaction() *Transaction {
	return &Transaction{
		operations: make([]func() error, 0),
		committed:  false,
	}
}

// Add adds an operation to the transaction
func (t *Transaction) Add(operation func() error) {
	t.operations = append(t.operations, operation)
}

// Execute executes all operations atomically
func (t *Transaction) Execute() error {
	if t.committed {
		return errors.New("transaction already committed")
	}

	fmt.Println("=== Starting Transaction ===")

	for i, op := range t.operations {
		if err := op(); err != nil {
			fmt.Printf("Operation %d failed: %v\n", i, err)
			fmt.Println("=== Rolling Back ===")
			return err
		}
	}

	fmt.Println("=== Transaction Committed ===")
	t.committed = true
	return nil
}

// WorkManager manages multiple units of work
type WorkManager struct {
	currentUoW *UnitOfWork
}

// NewWorkManager creates a work manager
func NewWorkManager() *WorkManager {
	return &WorkManager{
		currentUoW: NewUnitOfWork(),
	}
}

// GetCurrent returns current unit of work
func (wm *WorkManager) GetCurrent() *UnitOfWork {
	return wm.currentUoW
}

// Commit commits current unit of work and starts new one
func (wm *WorkManager) Commit() error {
	if err := wm.currentUoW.Commit(); err != nil {
		return err
	}
	wm.currentUoW = NewUnitOfWork()
	return nil
}

// Rollback rolls back and starts new unit of work
func (wm *WorkManager) Rollback() {
	wm.currentUoW.Rollback()
	wm.currentUoW = NewUnitOfWork()
}
