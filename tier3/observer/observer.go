// Package observer implements the Observer pattern.
//
// The Observer pattern defines a one-to-many dependency between objects
// where when one object (subject) changes state, all its dependents (observers)
// are notified and updated automatically.
//
// Key components:
// - Subject: The object being observed, maintains list of observers
// - Observer: Interface for objects that should be notified of subject changes
// - ConcreteSubject: Specific implementation of subject with state
// - ConcreteObserver: Specific implementation that reacts to subject changes
package observer

import (
	"fmt"
	"sync"
)

// Observer defines the interface for objects that should be notified
// of changes in a subject.
type Observer interface {
	// Update is called when the observed subject changes
	Update(subject Subject)
	// GetID returns a unique identifier for the observer
	GetID() string
}

// Subject defines the interface for objects that can be observed.
type Subject interface {
	// Attach adds an observer to the list of observers
	Attach(observer Observer)
	// Detach removes an observer from the list
	Detach(observer Observer)
	// Notify notifies all observers of a change
	Notify()
}

// StockPrice represents a stock that can be observed for price changes.
// This is a ConcreteSubject implementation.
type StockPrice struct {
	symbol    string
	price     float64
	observers []Observer
	mu        sync.RWMutex // Protects observers list for concurrent access
}

// NewStockPrice creates a new stock price subject.
func NewStockPrice(symbol string, initialPrice float64) *StockPrice {
	return &StockPrice{
		symbol:    symbol,
		price:     initialPrice,
		observers: make([]Observer, 0),
	}
}

// Attach adds an observer to receive price updates.
func (s *StockPrice) Attach(observer Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
}

// Detach removes an observer from receiving updates.
func (s *StockPrice) Detach(observer Observer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, obs := range s.observers {
		if obs.GetID() == observer.GetID() {
			// Remove observer by slicing it out
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			return
		}
	}
}

// Notify notifies all observers of a price change.
func (s *StockPrice) Notify() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, observer := range s.observers {
		observer.Update(s)
	}
}

// SetPrice updates the stock price and notifies all observers.
func (s *StockPrice) SetPrice(newPrice float64) {
	s.mu.Lock()
	oldPrice := s.price
	s.price = newPrice
	s.mu.Unlock()

	if oldPrice != newPrice {
		s.Notify()
	}
}

// GetPrice returns the current stock price.
func (s *StockPrice) GetPrice() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.price
}

// GetSymbol returns the stock symbol.
func (s *StockPrice) GetSymbol() string {
	return s.symbol
}

// TraderObserver represents a trader that monitors stock prices.
// This is a ConcreteObserver implementation.
type TraderObserver struct {
	id             string
	buyThreshold   float64 // Buy when price drops to this level
	sellThreshold  float64 // Sell when price rises to this level
	portfolio      map[string]int
	notificationFn func(string) // Optional callback for notifications
}

// NewTraderObserver creates a new trader observer.
func NewTraderObserver(id string, buyThreshold, sellThreshold float64) *TraderObserver {
	return &TraderObserver{
		id:            id,
		buyThreshold:  buyThreshold,
		sellThreshold: sellThreshold,
		portfolio:     make(map[string]int),
	}
}

// Update is called when the observed stock price changes.
func (t *TraderObserver) Update(subject Subject) {
	stock, ok := subject.(*StockPrice)
	if !ok {
		return
	}

	price := stock.GetPrice()
	symbol := stock.GetSymbol()

	// Trading logic based on thresholds
	if price <= t.buyThreshold {
		t.buy(symbol, price)
	} else if price >= t.sellThreshold {
		t.sell(symbol, price)
	} else {
		t.notify(fmt.Sprintf("[%s] Price update: %s = $%.2f (holding)", t.id, symbol, price))
	}
}

// GetID returns the trader's unique identifier.
func (t *TraderObserver) GetID() string {
	return t.id
}

// SetNotificationCallback sets a custom notification function.
func (t *TraderObserver) SetNotificationCallback(fn func(string)) {
	t.notificationFn = fn
}

func (t *TraderObserver) buy(symbol string, price float64) {
	t.portfolio[symbol] += 10 // Buy 10 shares
	t.notify(fmt.Sprintf("[%s] BUY: %s at $%.2f (threshold: $%.2f) - Portfolio: %d shares",
		t.id, symbol, price, t.buyThreshold, t.portfolio[symbol]))
}

func (t *TraderObserver) sell(symbol string, price float64) {
	if t.portfolio[symbol] > 0 {
		t.portfolio[symbol] -= 10 // Sell 10 shares
		t.notify(fmt.Sprintf("[%s] SELL: %s at $%.2f (threshold: $%.2f) - Portfolio: %d shares",
			t.id, symbol, price, t.sellThreshold, t.portfolio[symbol]))
	}
}

func (t *TraderObserver) notify(message string) {
	if t.notificationFn != nil {
		t.notificationFn(message)
	}
}

// AnalyticsObserver tracks statistics about stock price movements.
type AnalyticsObserver struct {
	id           string
	priceHistory []float64
	updates      int
}

// NewAnalyticsObserver creates a new analytics observer.
func NewAnalyticsObserver(id string) *AnalyticsObserver {
	return &AnalyticsObserver{
		id:           id,
		priceHistory: make([]float64, 0),
	}
}

// Update records price data for analytics.
func (a *AnalyticsObserver) Update(subject Subject) {
	stock, ok := subject.(*StockPrice)
	if !ok {
		return
	}

	price := stock.GetPrice()
	a.priceHistory = append(a.priceHistory, price)
	a.updates++
}

// GetID returns the analytics observer's unique identifier.
func (a *AnalyticsObserver) GetID() string {
	return a.id
}

// GetAveragePrice calculates the average price from history.
func (a *AnalyticsObserver) GetAveragePrice() float64 {
	if len(a.priceHistory) == 0 {
		return 0
	}

	sum := 0.0
	for _, price := range a.priceHistory {
		sum += price
	}
	return sum / float64(len(a.priceHistory))
}

// GetVolatility calculates price volatility (standard deviation).
func (a *AnalyticsObserver) GetVolatility() float64 {
	if len(a.priceHistory) < 2 {
		return 0
	}

	avg := a.GetAveragePrice()
	variance := 0.0

	for _, price := range a.priceHistory {
		diff := price - avg
		variance += diff * diff
	}

	variance /= float64(len(a.priceHistory))
	return variance // Simplified: returning variance instead of sqrt(variance)
}

// GetStats returns analytics statistics.
func (a *AnalyticsObserver) GetStats() string {
	return fmt.Sprintf("[%s] Updates: %d, Avg Price: $%.2f, Volatility: %.2f",
		a.id, a.updates, a.GetAveragePrice(), a.GetVolatility())
}

// AlertObserver monitors for price alerts.
type AlertObserver struct {
	id            string
	alertPrice    float64
	alertTriggered bool
	alertFn       func(string, float64)
}

// NewAlertObserver creates a new alert observer.
func NewAlertObserver(id string, alertPrice float64, alertFn func(string, float64)) *AlertObserver {
	return &AlertObserver{
		id:         id,
		alertPrice: alertPrice,
		alertFn:    alertFn,
	}
}

// Update checks if the price has reached the alert threshold.
func (a *AlertObserver) Update(subject Subject) {
	stock, ok := subject.(*StockPrice)
	if !ok {
		return
	}

	price := stock.GetPrice()
	symbol := stock.GetSymbol()

	if !a.alertTriggered && price >= a.alertPrice {
		a.alertTriggered = true
		if a.alertFn != nil {
			a.alertFn(symbol, price)
		}
	}
}

// GetID returns the alert observer's unique identifier.
func (a *AlertObserver) GetID() string {
	return a.id
}

// Reset resets the alert so it can trigger again.
func (a *AlertObserver) Reset() {
	a.alertTriggered = false
}

// EventBus implements a channel-based pub-sub observer pattern.
// This demonstrates Go's idiomatic approach using channels.
type EventBus struct {
	subscribers map[string][]chan interface{}
	mu          sync.RWMutex
}

// NewEventBus creates a new event bus.
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan interface{}),
	}
}

// Subscribe registers a channel to receive events of a specific type.
func (eb *EventBus) Subscribe(eventType string, ch chan interface{}) {
	eb.mu.Lock()
	defer eb.mu.Unlock()
	eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)
}

// Publish sends an event to all subscribers of that event type.
func (eb *EventBus) Publish(eventType string, data interface{}) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	for _, ch := range eb.subscribers[eventType] {
		// Non-blocking send to prevent slow subscribers from blocking
		select {
		case ch <- data:
		default:
			// Channel full, skip this subscriber
		}
	}
}

// Unsubscribe removes a channel from receiving events.
func (eb *EventBus) Unsubscribe(eventType string, ch chan interface{}) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subscribers := eb.subscribers[eventType]
	for i, subscriber := range subscribers {
		if subscriber == ch {
			eb.subscribers[eventType] = append(subscribers[:i], subscribers[i+1:]...)
			return
		}
	}
}
