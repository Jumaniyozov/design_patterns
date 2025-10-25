// Package saga demonstrates the Saga pattern.
// It manages distributed transactions across services using local transactions
// and compensating actions, maintaining consistency without distributed locks.
package saga

import (
	"context"
	"fmt"
	"time"
)

// Step represents a saga step with execute and compensate actions
type Step interface {
	Execute(ctx context.Context) error
	Compensate(ctx context.Context) error
	Name() string
}

// Saga coordinates a series of steps
type Saga struct {
	steps     []Step
	completed []Step
}

// NewSaga creates a new saga
func NewSaga() *Saga {
	return &Saga{
		steps:     make([]Step, 0),
		completed: make([]Step, 0),
	}
}

// AddStep adds a step to the saga
func (s *Saga) AddStep(step Step) {
	s.steps = append(s.steps, step)
}

// Execute executes the saga
func (s *Saga) Execute(ctx context.Context) error {
	fmt.Println("=== Starting Saga ===")

	for _, step := range s.steps {
		fmt.Printf("Executing step: %s\n", step.Name())

		if err := step.Execute(ctx); err != nil {
			fmt.Printf("Step %s failed: %v\n", step.Name(), err)
			fmt.Println("=== Compensating ===")

			// Compensate in reverse order
			if compErr := s.compensate(ctx); compErr != nil {
				return fmt.Errorf("compensation failed: %v (original error: %v)", compErr, err)
			}

			return fmt.Errorf("saga failed at step %s: %v", step.Name(), err)
		}

		s.completed = append(s.completed, step)
		fmt.Printf("Step %s completed\n", step.Name())
	}

	fmt.Println("=== Saga Completed Successfully ===")
	return nil
}

func (s *Saga) compensate(ctx context.Context) error {
	// Compensate in reverse order
	for i := len(s.completed) - 1; i >= 0; i-- {
		step := s.completed[i]
		fmt.Printf("Compensating step: %s\n", step.Name())

		if err := step.Compensate(ctx); err != nil {
			return fmt.Errorf("compensation failed for %s: %v", step.Name(), err)
		}

		fmt.Printf("Step %s compensated\n", step.Name())
	}
	return nil
}

// Example: Order Processing Saga

// ReserveInventoryStep reserves inventory
type ReserveInventoryStep struct {
	OrderID string
	Items   []string
	reserved bool
}

func (s *ReserveInventoryStep) Name() string {
	return "ReserveInventory"
}

func (s *ReserveInventoryStep) Execute(ctx context.Context) error {
	fmt.Printf("  Reserving inventory for order %s: %v\n", s.OrderID, s.Items)
	time.Sleep(100 * time.Millisecond) // Simulate work
	s.reserved = true
	return nil
}

func (s *ReserveInventoryStep) Compensate(ctx context.Context) error {
	if s.reserved {
		fmt.Printf("  Releasing reserved inventory for order %s\n", s.OrderID)
		time.Sleep(50 * time.Millisecond)
		s.reserved = false
	}
	return nil
}

// ChargePaymentStep charges payment
type ChargePaymentStep struct {
	OrderID string
	Amount  float64
	charged bool
}

func (s *ChargePaymentStep) Name() string {
	return "ChargePayment"
}

func (s *ChargePaymentStep) Execute(ctx context.Context) error {
	fmt.Printf("  Charging $%.2f for order %s\n", s.Amount, s.OrderID)
	time.Sleep(100 * time.Millisecond)
	s.charged = true
	return nil
}

func (s *ChargePaymentStep) Compensate(ctx context.Context) error {
	if s.charged {
		fmt.Printf("  Refunding $%.2f for order %s\n", s.Amount, s.OrderID)
		time.Sleep(50 * time.Millisecond)
		s.charged = false
	}
	return nil
}

// CreateOrderStep creates the order
type CreateOrderStep struct {
	OrderID string
	created bool
}

func (s *CreateOrderStep) Name() string {
	return "CreateOrder"
}

func (s *CreateOrderStep) Execute(ctx context.Context) error {
	fmt.Printf("  Creating order %s in database\n", s.OrderID)
	time.Sleep(100 * time.Millisecond)
	s.created = true
	return nil
}

func (s *CreateOrderStep) Compensate(ctx context.Context) error {
	if s.created {
		fmt.Printf("  Deleting order %s from database\n", s.OrderID)
		time.Sleep(50 * time.Millisecond)
		s.created = false
	}
	return nil
}

// SendNotificationStep sends notification
type SendNotificationStep struct {
	OrderID string
	Email   string
	sent    bool
}

func (s *SendNotificationStep) Name() string {
	return "SendNotification"
}

func (s *SendNotificationStep) Execute(ctx context.Context) error {
	fmt.Printf("  Sending confirmation email to %s for order %s\n", s.Email, s.OrderID)
	time.Sleep(100 * time.Millisecond)
	s.sent = true
	return nil
}

func (s *SendNotificationStep) Compensate(ctx context.Context) error {
	if s.sent {
		fmt.Printf("  Sending cancellation email to %s for order %s\n", s.Email, s.OrderID)
		time.Sleep(50 * time.Millisecond)
		s.sent = false
	}
	return nil
}

// FailingStep simulates a failing step
type FailingStep struct {
	name string
}

func (s *FailingStep) Name() string {
	return s.name
}

func (s *FailingStep) Execute(ctx context.Context) error {
	return fmt.Errorf("simulated failure in %s", s.name)
}

func (s *FailingStep) Compensate(ctx context.Context) error {
	return nil
}

// OrderService uses saga for order processing
type OrderService struct{}

// ProcessOrder processes an order using saga
func (os *OrderService) ProcessOrder(orderID string, items []string, amount float64, email string) error {
	saga := NewSaga()

	// Add steps in order
	saga.AddStep(&ReserveInventoryStep{
		OrderID: orderID,
		Items:   items,
	})

	saga.AddStep(&ChargePaymentStep{
		OrderID: orderID,
		Amount:  amount,
	})

	saga.AddStep(&CreateOrderStep{
		OrderID: orderID,
	})

	saga.AddStep(&SendNotificationStep{
		OrderID: orderID,
		Email:   email,
	})

	return saga.Execute(context.Background())
}

// ProcessOrderWithFailure demonstrates saga with compensation
func (os *OrderService) ProcessOrderWithFailure(orderID string, items []string, amount float64) error {
	saga := NewSaga()

	saga.AddStep(&ReserveInventoryStep{
		OrderID: orderID,
		Items:   items,
	})

	saga.AddStep(&ChargePaymentStep{
		OrderID: orderID,
		Amount:  amount,
	})

	// This step will fail
	saga.AddStep(&FailingStep{name: "ShipOrder"})

	return saga.Execute(context.Background())
}

// SagaOrchestrator orchestrates multiple sagas
type SagaOrchestrator struct {
	sagas map[string]*Saga
}

// NewSagaOrchestrator creates a saga orchestrator
func NewSagaOrchestrator() *SagaOrchestrator {
	return &SagaOrchestrator{
		sagas: make(map[string]*Saga),
	}
}

// RegisterSaga registers a saga with an ID
func (so *SagaOrchestrator) RegisterSaga(id string, saga *Saga) {
	so.sagas[id] = saga
}

// ExecuteSaga executes a registered saga
func (so *SagaOrchestrator) ExecuteSaga(id string, ctx context.Context) error {
	saga, exists := so.sagas[id]
	if !exists {
		return fmt.Errorf("saga %s not found", id)
	}
	return saga.Execute(ctx)
}
