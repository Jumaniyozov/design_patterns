package observer

import (
	"fmt"
	"time"
)

// Example1_BasicStockObserver demonstrates basic observer pattern with stock prices.
func Example1_BasicStockObserver() {
	fmt.Println("=== Example 1: Basic Stock Observer ===")

	// Create a stock subject
	appleStock := NewStockPrice("AAPL", 150.00)

	// Create trader observers with different trading strategies
	trader1 := NewTraderObserver("Trader-1", 145.00, 160.00)
	trader1.SetNotificationCallback(func(msg string) {
		fmt.Println(msg)
	})

	trader2 := NewTraderObserver("Trader-2", 140.00, 165.00)
	trader2.SetNotificationCallback(func(msg string) {
		fmt.Println(msg)
	})

	// Create analytics observer
	analytics := NewAnalyticsObserver("Analytics-Engine")

	// Attach observers to the stock
	appleStock.Attach(trader1)
	appleStock.Attach(trader2)
	appleStock.Attach(analytics)

	fmt.Println("Initial stock price: $150.00")
	fmt.Println("Trader-1 strategy: Buy at $145, Sell at $160")
	fmt.Println("Trader-2 strategy: Buy at $140, Sell at $165")
	fmt.Println()

	// Simulate price changes
	priceChanges := []float64{155.00, 145.00, 140.00, 155.00, 162.00, 165.00}

	for _, price := range priceChanges {
		fmt.Printf("\n--- Price Update: $%.2f ---\n", price)
		appleStock.SetPrice(price)
		time.Sleep(100 * time.Millisecond) // Simulate time between updates
	}

	// Display analytics
	fmt.Println("\n" + analytics.GetStats())
	fmt.Println()
}

// Example2_AlertSystem demonstrates alert-based observers.
func Example2_AlertSystem() {
	fmt.Println("=== Example 2: Alert System ===")

	// Create stock subject
	teslaStock := NewStockPrice("TSLA", 700.00)

	// Create alert observers for different price targets
	alert1 := NewAlertObserver("Alert-PriceTarget-750", 750.00, func(symbol string, price float64) {
		fmt.Printf("🚨 ALERT: %s reached $%.2f! Price target hit!\n", symbol, price)
	})

	alert2 := NewAlertObserver("Alert-PriceTarget-800", 800.00, func(symbol string, price float64) {
		fmt.Printf("🚨 CRITICAL ALERT: %s reached $%.2f! Major milestone!\n", symbol, price)
	})

	// Attach alerts
	teslaStock.Attach(alert1)
	teslaStock.Attach(alert2)

	// Simulate price movements
	fmt.Println("Initial price: $700.00")
	fmt.Println("Alerts set at $750 and $800")
	fmt.Println()

	prices := []float64{720.00, 740.00, 755.00, 780.00, 805.00}

	for _, price := range prices {
		fmt.Printf("Price update: $%.2f\n", price)
		teslaStock.SetPrice(price)
	}

	fmt.Println()
}

// Example3_DynamicObservers demonstrates attaching and detaching observers at runtime.
func Example3_DynamicObservers() {
	fmt.Println("=== Example 3: Dynamic Observer Management ===")

	// Create stock subject
	googleStock := NewStockPrice("GOOGL", 2800.00)

	// Create observers
	trader := NewTraderObserver("Active-Trader", 2750.00, 2850.00)
	trader.SetNotificationCallback(func(msg string) {
		fmt.Println(msg)
	})

	analytics := NewAnalyticsObserver("Market-Analytics")

	// Initially attach trader only
	googleStock.Attach(trader)
	fmt.Println("Attached: Active-Trader")

	// First price change
	fmt.Println("\n--- Price Update: $2820.00 ---")
	googleStock.SetPrice(2820.00)

	// Now attach analytics
	googleStock.Attach(analytics)
	fmt.Println("\nAttached: Market-Analytics")

	// Second price change (both observers notified)
	fmt.Println("\n--- Price Update: $2760.00 ---")
	googleStock.SetPrice(2760.00)

	// Detach trader
	googleStock.Detach(trader)
	fmt.Println("\nDetached: Active-Trader")

	// Third price change (only analytics notified)
	fmt.Println("\n--- Price Update: $2850.00 ---")
	googleStock.SetPrice(2850.00)

	// Display analytics results
	fmt.Println("\n" + analytics.GetStats())
	fmt.Println()
}

// Example4_ChannelBasedObserver demonstrates Go's channel-based observer pattern.
func Example4_ChannelBasedObserver() {
	fmt.Println("=== Example 4: Channel-Based Observer (Event Bus) ===")

	// Create event bus
	bus := NewEventBus()

	// Create buffered channels for subscribers
	subscriber1 := make(chan interface{}, 10)
	subscriber2 := make(chan interface{}, 10)
	subscriber3 := make(chan interface{}, 10)

	// Subscribe to events
	bus.Subscribe("stock.price.update", subscriber1)
	bus.Subscribe("stock.price.update", subscriber2)
	bus.Subscribe("stock.alert", subscriber3)

	// Start goroutines to handle events
	done := make(chan bool, 3)

	go func() {
		fmt.Println("Subscriber-1 started (Trading Bot)")
		for event := range subscriber1 {
			priceUpdate, ok := event.(map[string]interface{})
			if !ok {
				continue
			}
			fmt.Printf("  [Trading Bot] Processing price update: %s = $%.2f\n",
				priceUpdate["symbol"], priceUpdate["price"])
		}
		done <- true
	}()

	go func() {
		fmt.Println("Subscriber-2 started (Analytics)")
		for event := range subscriber2 {
			priceUpdate, ok := event.(map[string]interface{})
			if !ok {
				continue
			}
			fmt.Printf("  [Analytics] Recording data: %s = $%.2f\n",
				priceUpdate["symbol"], priceUpdate["price"])
		}
		done <- true
	}()

	go func() {
		fmt.Println("Subscriber-3 started (Alert System)")
		for event := range subscriber3 {
			alert, ok := event.(string)
			if !ok {
				continue
			}
			fmt.Printf("  [Alert System] 🚨 %s\n", alert)
		}
		done <- true
	}()

	fmt.Println()

	// Publish price updates
	time.Sleep(100 * time.Millisecond) // Let subscribers start

	fmt.Println("Publishing price updates...")
	bus.Publish("stock.price.update", map[string]interface{}{
		"symbol": "MSFT",
		"price":  380.00,
	})

	time.Sleep(50 * time.Millisecond)

	bus.Publish("stock.price.update", map[string]interface{}{
		"symbol": "MSFT",
		"price":  385.00,
	})

	time.Sleep(50 * time.Millisecond)

	// Publish an alert
	fmt.Println("\nPublishing alert...")
	bus.Publish("stock.alert", "MSFT price exceeded $385!")

	time.Sleep(50 * time.Millisecond)

	// Cleanup
	close(subscriber1)
	close(subscriber2)
	close(subscriber3)

	// Wait for all goroutines to finish
	for i := 0; i < 3; i++ {
		<-done
	}

	fmt.Println("\nAll subscribers processed events")
	fmt.Println()
}

// Example5_MultipleSubjects demonstrates observing multiple subjects.
func Example5_MultipleSubjects() {
	fmt.Println("=== Example 5: Multiple Subjects (Portfolio Tracker) ===")

	// Create portfolio analytics that observes multiple stocks
	portfolioAnalytics := &PortfolioAnalytics{
		id:     "Portfolio-Tracker",
		stocks: make(map[string]float64),
	}

	// Create multiple stock subjects
	stocks := []*StockPrice{
		NewStockPrice("AAPL", 150.00),
		NewStockPrice("GOOGL", 2800.00),
		NewStockPrice("MSFT", 380.00),
		NewStockPrice("AMZN", 3400.00),
	}

	// Attach the portfolio tracker to all stocks
	for _, stock := range stocks {
		stock.Attach(portfolioAnalytics)
		fmt.Printf("Tracking: %s at $%.2f\n", stock.GetSymbol(), stock.GetPrice())
	}

	fmt.Println()

	// Simulate price changes across different stocks
	fmt.Println("--- Price Updates ---")
	stocks[0].SetPrice(155.00) // AAPL
	stocks[2].SetPrice(385.00) // MSFT
	stocks[1].SetPrice(2850.00) // GOOGL
	stocks[3].SetPrice(3350.00) // AMZN

	// Display portfolio summary
	fmt.Println()
	portfolioAnalytics.DisplayPortfolio()
	fmt.Println()
}

// PortfolioAnalytics observes multiple stock subjects.
type PortfolioAnalytics struct {
	id     string
	stocks map[string]float64
}

// Update records price updates from any stock subject.
func (p *PortfolioAnalytics) Update(subject Subject) {
	stock, ok := subject.(*StockPrice)
	if !ok {
		return
	}

	symbol := stock.GetSymbol()
	price := stock.GetPrice()
	oldPrice := p.stocks[symbol]

	p.stocks[symbol] = price

	if oldPrice > 0 {
		change := ((price - oldPrice) / oldPrice) * 100
		fmt.Printf("[%s] %s: $%.2f (%.2f%%)\n", p.id, symbol, price, change)
	} else {
		fmt.Printf("[%s] %s: $%.2f (initial)\n", p.id, symbol, price)
	}
}

// GetID returns the portfolio analytics identifier.
func (p *PortfolioAnalytics) GetID() string {
	return p.id
}

// DisplayPortfolio shows the current portfolio state.
func (p *PortfolioAnalytics) DisplayPortfolio() {
	fmt.Println("=== Portfolio Summary ===")
	totalValue := 0.0
	for symbol, price := range p.stocks {
		fmt.Printf("%s: $%.2f\n", symbol, price)
		totalValue += price
	}
	fmt.Printf("Total tracked value: $%.2f\n", totalValue)
}
