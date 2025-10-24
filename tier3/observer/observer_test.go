package observer

import (
	"sync"
	"testing"
	"time"
)

func TestStockPrice_AttachDetach(t *testing.T) {
	stock := NewStockPrice("TEST", 100.0)
	observer := NewTraderObserver("test-trader", 90.0, 110.0)

	// Test attach
	stock.Attach(observer)
	if len(stock.observers) != 1 {
		t.Errorf("Expected 1 observer, got %d", len(stock.observers))
	}

	// Test detach
	stock.Detach(observer)
	if len(stock.observers) != 0 {
		t.Errorf("Expected 0 observers after detach, got %d", len(stock.observers))
	}
}

func TestStockPrice_Notify(t *testing.T) {
	stock := NewStockPrice("TEST", 100.0)
	updateCount := 0

	// Create a custom observer that counts updates
	observer := &mockObserver{
		id: "test-observer",
		updateFn: func(s Subject) {
			updateCount++
		},
	}

	stock.Attach(observer)
	stock.SetPrice(105.0)

	if updateCount != 1 {
		t.Errorf("Expected 1 update notification, got %d", updateCount)
	}

	// Setting same price should not trigger notification
	stock.SetPrice(105.0)
	if updateCount != 1 {
		t.Errorf("Expected no additional update for same price, got %d", updateCount)
	}
}

func TestTraderObserver_BuyThreshold(t *testing.T) {
	stock := NewStockPrice("TEST", 150.0)
	trader := NewTraderObserver("test-trader", 145.0, 160.0)

	buyTriggered := false
	trader.SetNotificationCallback(func(msg string) {
		if len(msg) > 0 {
			// Check if message contains buy signal
			if containsSubstring(msg, "BUY") && containsSubstring(msg, "TEST") {
				buyTriggered = true
			}
		}
	})

	stock.Attach(trader)
	stock.SetPrice(145.0) // Should trigger buy

	if !buyTriggered {
		t.Error("Expected buy to be triggered at threshold price")
	}

	if trader.portfolio["TEST"] != 10 {
		t.Errorf("Expected 10 shares in portfolio, got %d", trader.portfolio["TEST"])
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestTraderObserver_SellThreshold(t *testing.T) {
	stock := NewStockPrice("TEST", 150.0)
	trader := NewTraderObserver("test-trader", 140.0, 160.0)

	// First buy some shares
	trader.portfolio["TEST"] = 20

	sellTriggered := false
	trader.SetNotificationCallback(func(msg string) {
		if len(msg) > 0 {
			// Check if message contains sell signal
			if containsSubstring(msg, "SELL") && containsSubstring(msg, "TEST") {
				sellTriggered = true
			}
		}
	})

	stock.Attach(trader)
	stock.SetPrice(160.0) // Should trigger sell

	if !sellTriggered {
		t.Error("Expected sell to be triggered at threshold price")
	}

	if trader.portfolio["TEST"] != 10 {
		t.Errorf("Expected 10 shares in portfolio after sell, got %d", trader.portfolio["TEST"])
	}
}

func TestAnalyticsObserver_Statistics(t *testing.T) {
	stock := NewStockPrice("TEST", 100.0)
	analytics := NewAnalyticsObserver("test-analytics")

	stock.Attach(analytics)

	// Simulate price changes (note: first price is same as initial, so no update)
	prices := []float64{100.0, 110.0, 105.0, 115.0, 120.0}
	for _, price := range prices {
		stock.SetPrice(price)
	}

	// Check update count (only 4 updates because first price equals initial price)
	expectedUpdates := 4
	if analytics.updates != expectedUpdates {
		t.Errorf("Expected %d updates, got %d", expectedUpdates, analytics.updates)
	}

	// Check average (only of prices that triggered updates: 110, 105, 115, 120)
	expectedAvg := 112.50 // (110 + 105 + 115 + 120) / 4
	avg := analytics.GetAveragePrice()
	if avg != expectedAvg {
		t.Errorf("Expected average price %.2f, got %.2f", expectedAvg, avg)
	}

	// Check that volatility is calculated (non-zero)
	volatility := analytics.GetVolatility()
	if volatility == 0 {
		t.Error("Expected non-zero volatility")
	}
}

func TestAlertObserver_Trigger(t *testing.T) {
	stock := NewStockPrice("TEST", 100.0)
	alertTriggered := false
	alertPrice := 0.0

	alert := NewAlertObserver("test-alert", 150.0, func(symbol string, price float64) {
		alertTriggered = true
		alertPrice = price
	})

	stock.Attach(alert)

	// Price below threshold - should not trigger
	stock.SetPrice(140.0)
	if alertTriggered {
		t.Error("Alert should not trigger below threshold")
	}

	// Price at threshold - should trigger
	stock.SetPrice(150.0)
	if !alertTriggered {
		t.Error("Alert should trigger at threshold")
	}
	if alertPrice != 150.0 {
		t.Errorf("Expected alert price 150.0, got %.2f", alertPrice)
	}

	// Reset and test again
	alert.Reset()
	alertTriggered = false
	stock.SetPrice(160.0)
	if !alertTriggered {
		t.Error("Alert should trigger again after reset")
	}
}

func TestEventBus_Subscribe(t *testing.T) {
	bus := NewEventBus()
	ch := make(chan interface{}, 1)

	bus.Subscribe("test.event", ch)

	// Publish event
	testData := "test data"
	bus.Publish("test.event", testData)

	// Receive event
	select {
	case data := <-ch:
		if data != testData {
			t.Errorf("Expected %v, got %v", testData, data)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timeout waiting for event")
	}
}

func TestEventBus_MultipleSubscribers(t *testing.T) {
	bus := NewEventBus()
	ch1 := make(chan interface{}, 1)
	ch2 := make(chan interface{}, 1)
	ch3 := make(chan interface{}, 1)

	bus.Subscribe("test.event", ch1)
	bus.Subscribe("test.event", ch2)
	bus.Subscribe("test.event", ch3)

	// Publish event
	testData := 42
	bus.Publish("test.event", testData)

	// All subscribers should receive the event
	channels := []chan interface{}{ch1, ch2, ch3}
	for i, ch := range channels {
		select {
		case data := <-ch:
			if data != testData {
				t.Errorf("Subscriber %d: Expected %v, got %v", i+1, testData, data)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Subscriber %d: Timeout waiting for event", i+1)
		}
	}
}

func TestEventBus_Unsubscribe(t *testing.T) {
	bus := NewEventBus()
	ch := make(chan interface{}, 1)

	bus.Subscribe("test.event", ch)
	bus.Unsubscribe("test.event", ch)

	// Publish event
	bus.Publish("test.event", "test data")

	// Should not receive event after unsubscribe
	select {
	case <-ch:
		t.Error("Should not receive event after unsubscribe")
	case <-time.After(50 * time.Millisecond):
		// Expected - no event received
	}
}

func TestConcurrentObservers(t *testing.T) {
	stock := NewStockPrice("TEST", 100.0)
	observerCount := 10
	updateCount := 100
	var wg sync.WaitGroup

	// Create multiple observers concurrently
	for i := 0; i < observerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			observer := &mockObserver{
				id: string(rune('A' + id)),
				updateFn: func(s Subject) {
					// Simulate some work
					time.Sleep(time.Microsecond)
				},
			}
			stock.Attach(observer)
		}(i)
	}

	wg.Wait()

	// Verify all observers attached
	if len(stock.observers) != observerCount {
		t.Errorf("Expected %d observers, got %d", observerCount, len(stock.observers))
	}

	// Trigger updates concurrently
	for i := 0; i < updateCount; i++ {
		wg.Add(1)
		go func(price float64) {
			defer wg.Done()
			stock.SetPrice(price)
		}(100.0 + float64(i))
	}

	wg.Wait()

	// If we got here without deadlock or panic, concurrent access works
}

func TestMultipleSubjects(t *testing.T) {
	portfolio := &PortfolioAnalytics{
		id:     "test-portfolio",
		stocks: make(map[string]float64),
	}

	stocks := []*StockPrice{
		NewStockPrice("AAPL", 150.0),
		NewStockPrice("GOOGL", 2800.0),
		NewStockPrice("MSFT", 380.0),
	}

	// Attach portfolio to all stocks
	for _, stock := range stocks {
		stock.Attach(portfolio)
	}

	// Update prices
	stocks[0].SetPrice(155.0)
	stocks[1].SetPrice(2850.0)
	stocks[2].SetPrice(385.0)

	// Verify portfolio tracked all updates
	if portfolio.stocks["AAPL"] != 155.0 {
		t.Errorf("Expected AAPL price 155.0, got %.2f", portfolio.stocks["AAPL"])
	}
	if portfolio.stocks["GOOGL"] != 2850.0 {
		t.Errorf("Expected GOOGL price 2850.0, got %.2f", portfolio.stocks["GOOGL"])
	}
	if portfolio.stocks["MSFT"] != 385.0 {
		t.Errorf("Expected MSFT price 385.0, got %.2f", portfolio.stocks["MSFT"])
	}
}

// Mock observer for testing
type mockObserver struct {
	id       string
	updateFn func(Subject)
}

func (m *mockObserver) Update(subject Subject) {
	if m.updateFn != nil {
		m.updateFn(subject)
	}
}

func (m *mockObserver) GetID() string {
	return m.id
}

// Benchmark tests
func BenchmarkStockPrice_Notify(b *testing.B) {
	stock := NewStockPrice("TEST", 100.0)

	// Attach 100 observers
	for i := 0; i < 100; i++ {
		observer := NewAnalyticsObserver(string(rune('A' + i)))
		stock.Attach(observer)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stock.SetPrice(100.0 + float64(i))
	}
}

func BenchmarkEventBus_Publish(b *testing.B) {
	bus := NewEventBus()

	// Create 100 subscribers
	for i := 0; i < 100; i++ {
		ch := make(chan interface{}, 100)
		bus.Subscribe("test.event", ch)

		// Start goroutine to drain channel
		go func(ch chan interface{}) {
			for range ch {
			}
		}(ch)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Publish("test.event", i)
	}
}
