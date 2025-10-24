package strategy

import "fmt"

// Example1_BasicPaymentProcessing demonstrates the fundamental use of the Strategy pattern.
// This shows how to select different payment methods at runtime without conditional logic.
func Example1_BasicPaymentProcessing() {
	fmt.Println("\n=== Example 1: Basic Payment Processing ===")
	fmt.Println("Using different payment strategies interchangeably")

	// Create a credit card strategy
	creditCardStrategy := NewCreditCardStrategy(
		"4532123456789010",
		"12/25",
		"123",
		"John Doe",
	)

	// Create a payment processor with the credit card strategy
	processor := NewPaymentProcessor(creditCardStrategy)

	// Set the payment details
	processor.SetDetails(PaymentDetails{
		CardNumber: "4532123456789010",
		ExpiryDate: "12/25",
		CVV:        "123",
	})

	// Process a payment
	fmt.Printf("Using: %s\n", processor.GetCurrentStrategyName())
	txID, err := processor.Process(99.99)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Transaction ID: %s\n", txID)

	// Refund the payment
	err = processor.Refund(txID, 99.99)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// Example2_SwitchingStrategiesAtRuntime demonstrates the key advantage of the Strategy pattern:
// changing the algorithm at runtime without modifying client code.
func Example2_SwitchingStrategiesAtRuntime() {
	fmt.Println("\n=== Example 2: Switching Strategies at Runtime ===")
	fmt.Println("Same processor, different payment methods")

	processor := NewPaymentProcessor(nil)

	// Customer 1: Wants to pay with credit card
	fmt.Println("--- Customer 1: Credit Card Payment ---")
	ccStrategy := NewCreditCardStrategy("4532123456789010", "12/25", "123", "Alice")
	processor.SetStrategy(ccStrategy)
	processor.SetDetails(PaymentDetails{CardNumber: "4532123456789010"})

	fmt.Printf("Using: %s\n", processor.GetCurrentStrategyName())
	txID, err := processor.Process(50.00)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Transaction ID: %s\n", txID)

	// Customer 2: Wants to pay with PayPal
	fmt.Println("\n--- Customer 2: PayPal Payment ---")
	ppStrategy := NewPayPalStrategy("bob@example.com")
	processor.SetStrategy(ppStrategy)
	processor.SetDetails(PaymentDetails{Email: "bob@example.com"})

	fmt.Printf("Using: %s\n", processor.GetCurrentStrategyName())
	txID, err = processor.Process(75.50)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Transaction ID: %s\n", txID)

	// Customer 3: Wants to pay with cryptocurrency
	fmt.Println("\n--- Customer 3: Cryptocurrency Payment ---")
	cryptoStrategy := NewCryptoStrategy("1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9", "BTC")
	processor.SetStrategy(cryptoStrategy)
	processor.SetDetails(PaymentDetails{WalletAddr: "1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9"})

	fmt.Printf("Using: %s\n", processor.GetCurrentStrategyName())
	txID, err = processor.Process(0.025)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Transaction ID: %s\n", txID)
}

// Example3_ErrorHandlingPerStrategy demonstrates how each strategy handles validation errors.
// This shows the benefit of encapsulating algorithm-specific logic in separate types.
func Example3_ErrorHandlingPerStrategy() {
	fmt.Println("\n=== Example 3: Error Handling Per Strategy ===")
	fmt.Println("Each strategy validates according to its own rules")

	processor := NewPaymentProcessor(nil)

	// Example 1: Invalid credit card
	fmt.Println("--- Invalid Credit Card ---")
	invalidCC := NewCreditCardStrategy("123", "13/25", "99", "Invalid") // Invalid card
	processor.SetStrategy(invalidCC)
	processor.SetDetails(PaymentDetails{})

	_, err := processor.Process(100.00)
	if err != nil {
		fmt.Printf("Caught error: %v\n", err)
	}

	// Example 2: Invalid PayPal email
	fmt.Println("\n--- Invalid PayPal Email ---")
	invalidPayPal := NewPayPalStrategy("notanemail") // Invalid email
	processor.SetStrategy(invalidPayPal)
	processor.SetDetails(PaymentDetails{})

	_, err = processor.Process(100.00)
	if err != nil {
		fmt.Printf("Caught error: %v\n", err)
	}

	// Example 3: Invalid cryptocurrency wallet
	fmt.Println("\n--- Invalid Cryptocurrency Wallet ---")
	invalidCrypto := NewCryptoStrategy("abc", "BTC") // Too short
	processor.SetStrategy(invalidCrypto)
	processor.SetDetails(PaymentDetails{})

	_, err = processor.Process(100.00)
	if err != nil {
		fmt.Printf("Caught error: %v\n", err)
	}
}

// Example4_PaymentProcessor demonstrates a real-world scenario where different
// payment methods are chosen based on user preferences or system configuration.
func Example4_MultiplePaymentsWithDifferentMethods() {
	fmt.Println("\n=== Example 4: Multiple Payments with Different Methods ===")
	fmt.Println("Processing orders with customer's preferred payment method")

	// Simulate customer data
	customers := []struct {
		name     string
		strategy PaymentStrategy
		details  PaymentDetails
		amount   float64
	}{
		{
			name:     "Alice (Credit Card)",
			strategy: NewCreditCardStrategy("4532123456789010", "12/25", "123", "Alice"),
			details:  PaymentDetails{CardNumber: "4532123456789010"},
			amount:   99.99,
		},
		{
			name:     "Bob (PayPal)",
			strategy: NewPayPalStrategy("bob@example.com"),
			details:  PaymentDetails{Email: "bob@example.com"},
			amount:   149.99,
		},
		{
			name:     "Charlie (Bitcoin)",
			strategy: NewCryptoStrategy("1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9", "BTC"),
			details:  PaymentDetails{WalletAddr: "1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9"},
			amount:   0.05,
		},
	}

	processor := NewPaymentProcessor(nil)

	for i, customer := range customers {
		fmt.Printf("--- Order %d: %s ---\n", i+1, customer.name)

		// Set strategy and details
		processor.SetStrategy(customer.strategy)
		processor.SetDetails(customer.details)

		// Process payment
		fmt.Printf("Payment Method: %s\n", processor.GetCurrentStrategyName())
		fmt.Printf("Amount: %.2f\n", customer.amount)

		txID, err := processor.Process(customer.amount)
		if err != nil {
			fmt.Printf("Payment failed: %v\n", err)
			continue
		}

		fmt.Printf("✓ Payment successful (TX: %s)\n\n", txID)
	}
}

// Example5_DynamicStrategySelection demonstrates how to select strategies
// based on external configuration or business logic.
func Example5_DynamicStrategySelection() {
	fmt.Println("\n=== Example 5: Dynamic Strategy Selection ===")
	fmt.Println("Selecting payment method based on configuration")

	// Simulate a configuration that specifies which payment methods are available
	availablePaymentMethods := map[string]func() PaymentStrategy{
		"credit_card": func() PaymentStrategy {
			return NewCreditCardStrategy("4532123456789010", "12/25", "123", "Default Card")
		},
		"paypal": func() PaymentStrategy {
			return NewPayPalStrategy("default@example.com")
		},
		"crypto": func() PaymentStrategy {
			return NewCryptoStrategy("1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9", "ETH")
		},
	}

	processor := NewPaymentProcessor(nil)

	// Simulate user selecting different payment methods
	selectedMethods := []string{"credit_card", "paypal", "crypto"}

	for _, methodName := range selectedMethods {
		fmt.Printf("--- Processing with %s ---\n", methodName)

		// Dynamically get the strategy
		strategyFactory, exists := availablePaymentMethods[methodName]
		if !exists {
			fmt.Printf("Payment method %s not available\n\n", methodName)
			continue
		}

		strategy := strategyFactory()
		processor.SetStrategy(strategy)
		processor.SetDetails(PaymentDetails{})

		fmt.Printf("Using: %s\n", processor.GetCurrentStrategyName())
		txID, err := processor.Process(100.00)
		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			continue
		}

		fmt.Printf("Transaction ID: %s\n\n", txID)
	}
}

// Example6_StrategyComparison shows how the Strategy pattern eliminates
// the need for conditional logic that would be required without the pattern.
func Example6_StrategyComparison() {
	fmt.Println("\n=== Example 6: Strategy Pattern vs Conditional Approach ===")
	fmt.Println("Why Strategy pattern is better than conditionals")

	fmt.Println("--- WITHOUT Strategy Pattern (Bad) ---")
	fmt.Println(`
func processPaymentOld(method string, amount float64) {
    if method == "credit_card" {
        // 50 lines of credit card logic
    } else if method == "paypal" {
        // 40 lines of PayPal logic
    } else if method == "crypto" {
        // 60 lines of crypto logic
    }
    // Problems: violates SRP, OCP, hard to test, maintenance nightmare
}`)

	fmt.Println("\n--- WITH Strategy Pattern (Good) ---")
	fmt.Println(`
type PaymentStrategy interface {
    Process(amount float64) (txID string, err error)
}

func processPaymentNew(strategy PaymentStrategy, amount float64) {
    txID, err := strategy.Process(amount)
    // Clean, testable, extensible - add new methods without changing this code
}`)

	fmt.Println("\n--- Demonstrating the benefit ---")

	processor := NewPaymentProcessor(NewCreditCardStrategy("4532123456789010", "12/25", "123", "Demo"))
	processor.SetDetails(PaymentDetails{})

	// Same code works with ANY strategy - no if/else needed!
	strategies := []PaymentStrategy{
		NewCreditCardStrategy("4532123456789010", "12/25", "123", "Demo"),
		NewPayPalStrategy("demo@example.com"),
		NewCryptoStrategy("1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9", "BTC"),
	}

	for _, strategy := range strategies {
		processor.SetStrategy(strategy)
		fmt.Printf("\nUsing %s\n", processor.GetCurrentStrategyName())
		txID, _ := processor.Process(50.00)
		fmt.Printf("Result: %s\n", txID)
	}

	fmt.Println("\n✓ Same code, different behaviors - that's the power of Strategy!")
}
