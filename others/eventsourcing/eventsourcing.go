// Package eventsourcing demonstrates the Event Sourcing pattern.
// It stores application state as a sequence of events, enabling complete
// audit trails, time travel, and event replay capabilities.
package eventsourcing

import (
	"fmt"
	"sync"
	"time"
)

// Event interface for all events
type Event interface {
	EventType() string
	EventTimestamp() time.Time
	AggregateID() string
}

// BaseEvent provides common event fields
type BaseEvent struct {
	Type        string
	Timestamp   time.Time
	AggregateId string
}

func (e *BaseEvent) EventType() string           { return e.Type }
func (e *BaseEvent) EventTimestamp() time.Time   { return e.Timestamp }
func (e *BaseEvent) AggregateID() string         { return e.AggregateId }

// Domain Events

type AccountCreatedEvent struct {
	BaseEvent
	Owner string
}

type MoneyDepositedEvent struct {
	BaseEvent
	Amount float64
}

type MoneyWithdrawnEvent struct {
	BaseEvent
	Amount float64
}

type AccountClosedEvent struct {
	BaseEvent
	Reason string
}

// EventStore stores and retrieves events
type EventStore interface {
	AppendEvent(event Event) error
	GetEvents(aggregateID string) ([]Event, error)
	GetAllEvents() ([]Event, error)
}

// InMemoryEventStore implements EventStore
type InMemoryEventStore struct {
	events map[string][]Event
	mu     sync.RWMutex
}

// NewInMemoryEventStore creates an event store
func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events: make(map[string][]Event),
	}
}

func (es *InMemoryEventStore) AppendEvent(event Event) error {
	es.mu.Lock()
	defer es.mu.Unlock()

	aggregateID := event.AggregateID()
	es.events[aggregateID] = append(es.events[aggregateID], event)
	return nil
}

func (es *InMemoryEventStore) GetEvents(aggregateID string) ([]Event, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	events, exists := es.events[aggregateID]
	if !exists {
		return []Event{}, nil
	}
	return append([]Event{}, events...), nil
}

func (es *InMemoryEventStore) GetAllEvents() ([]Event, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	allEvents := make([]Event, 0)
	for _, events := range es.events {
		allEvents = append(allEvents, events...)
	}
	return allEvents, nil
}

// BankAccount aggregate reconstructed from events
type BankAccount struct {
	ID      string
	Owner   string
	Balance float64
	Closed  bool
	version int
}

// NewBankAccount creates a bank account
func NewBankAccount(id, owner string) *BankAccount {
	return &BankAccount{
		ID:      id,
		Owner:   owner,
		Balance: 0,
		Closed:  false,
		version: 0,
	}
}

// Apply applies an event to the account
func (ba *BankAccount) Apply(event Event) {
	switch e := event.(type) {
	case *AccountCreatedEvent:
		ba.Owner = e.Owner
	case *MoneyDepositedEvent:
		ba.Balance += e.Amount
	case *MoneyWithdrawnEvent:
		ba.Balance -= e.Amount
	case *AccountClosedEvent:
		ba.Closed = true
	}
	ba.version++
}

// Rebuild rebuilds account state from events
func (ba *BankAccount) Rebuild(events []Event) {
	for _, event := range events {
		ba.Apply(event)
	}
}

// GetState returns current state
func (ba *BankAccount) GetState() string {
	status := "Active"
	if ba.Closed {
		status = "Closed"
	}
	return fmt.Sprintf("Account %s | Owner: %s | Balance: $%.2f | Status: %s | Version: %d",
		ba.ID, ba.Owner, ba.Balance, status, ba.version)
}

// BankAccountService manages accounts using event sourcing
type BankAccountService struct {
	eventStore EventStore
}

// NewBankAccountService creates a service
func NewBankAccountService(eventStore EventStore) *BankAccountService {
	return &BankAccountService{eventStore: eventStore}
}

// CreateAccount creates a new account
func (s *BankAccountService) CreateAccount(id, owner string) error {
	event := &AccountCreatedEvent{
		BaseEvent: BaseEvent{
			Type:        "AccountCreated",
			Timestamp:   time.Now(),
			AggregateId: id,
		},
		Owner: owner,
	}
	return s.eventStore.AppendEvent(event)
}

// Deposit deposits money
func (s *BankAccountService) Deposit(id string, amount float64) error {
	event := &MoneyDepositedEvent{
		BaseEvent: BaseEvent{
			Type:        "MoneyDeposited",
			Timestamp:   time.Now(),
			AggregateId: id,
		},
		Amount: amount,
	}
	return s.eventStore.AppendEvent(event)
}

// Withdraw withdraws money
func (s *BankAccountService) Withdraw(id string, amount float64) error {
	// Check balance first by replaying events
	account, err := s.GetAccount(id)
	if err != nil {
		return err
	}

	if account.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	event := &MoneyWithdrawnEvent{
		BaseEvent: BaseEvent{
			Type:        "MoneyWithdrawn",
			Timestamp:   time.Now(),
			AggregateId: id,
		},
		Amount: amount,
	}
	return s.eventStore.AppendEvent(event)
}

// CloseAccount closes an account
func (s *BankAccountService) CloseAccount(id, reason string) error {
	event := &AccountClosedEvent{
		BaseEvent: BaseEvent{
			Type:        "AccountClosed",
			Timestamp:   time.Now(),
			AggregateId: id,
		},
		Reason: reason,
	}
	return s.eventStore.AppendEvent(event)
}

// GetAccount reconstructs account from events
func (s *BankAccountService) GetAccount(id string) (*BankAccount, error) {
	events, err := s.eventStore.GetEvents(id)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("account not found")
	}

	// Reconstruct from first event
	firstEvent := events[0].(*AccountCreatedEvent)
	account := NewBankAccount(id, firstEvent.Owner)

	// Replay all events
	account.Rebuild(events)

	return account, nil
}

// GetAccountHistory returns event history
func (s *BankAccountService) GetAccountHistory(id string) ([]string, error) {
	events, err := s.eventStore.GetEvents(id)
	if err != nil {
		return nil, err
	}

	history := make([]string, 0, len(events))
	for _, event := range events {
		var description string
		switch e := event.(type) {
		case *AccountCreatedEvent:
			description = fmt.Sprintf("[%s] Account created for %s",
				e.Timestamp.Format("15:04:05"), e.Owner)
		case *MoneyDepositedEvent:
			description = fmt.Sprintf("[%s] Deposited $%.2f",
				e.Timestamp.Format("15:04:05"), e.Amount)
		case *MoneyWithdrawnEvent:
			description = fmt.Sprintf("[%s] Withdrew $%.2f",
				e.Timestamp.Format("15:04:05"), e.Amount)
		case *AccountClosedEvent:
			description = fmt.Sprintf("[%s] Account closed: %s",
				e.Timestamp.Format("15:04:05"), e.Reason)
		}
		history = append(history, description)
	}

	return history, nil
}
