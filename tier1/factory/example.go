package factory

import "fmt"

// Example1_BasicFactoryUsage demonstrates the basic usage of the factory pattern.
// This shows how the factory function abstracts away the creation logic, allowing
// clients to request a payment processor without knowing the concrete implementation.
func Example1_BasicFactoryUsage() {
	fmt.Println("=== Example 1: Basic Factory Usage ===")

	// Create a credit card processor using the factory
	ccConfig := PaymentConfig{
		Method:     CreditCard,
		CardNumber: "4532015112830366",
		CVV:        "123",
		ExpiryDate: "12/25",
	}

	processor, err := NewPaymentProcessor(ccConfig)
	if err != nil {
		fmt.Printf("Error creating processor: %v\n", err)
		return
	}

	fmt.Printf("Created: %s\n", processor.GetProcessorName())

	// Process a payment - notice we're using the interface, not the concrete type
	transactionID, err := processor.ProcessPayment(99.99)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
		return
	}

	fmt.Printf("Payment successful! Transaction ID: %s\n\n", transactionID)
}

// Example2_MultiplePaymentMethods demonstrates how the factory pattern allows
// working with different payment methods through a common interface.
// This is the key benefit: clients don't need to know about concrete implementations.
func Example2_MultiplePaymentMethods() {
	fmt.Println("=== Example 2: Multiple Payment Methods ===")

	// Define multiple payment configurations
	configs := []PaymentConfig{
		{
			Method:     CreditCard,
			CardNumber: "4532015112830366",
			CVV:        "123",
			ExpiryDate: "12/25",
		},
		{
			Method:   PayPal,
			Email:    "user@example.com",
			APIToken: "paypal-token-abc123",
		},
		{
			Method:        Bitcoin,
			WalletAddress: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
			Network:       "mainnet",
		},
	}

	amount := 150.00

	// Process payment with each method - same code works for all types!
	for i, config := range configs {
		fmt.Printf("Processing payment #%d with %s\n", i+1, config.Method)

		processor, err := NewPaymentProcessor(config)
		if err != nil {
			fmt.Printf("  Error: %v\n\n", err)
			continue
		}

		transactionID, err := processor.ProcessPayment(amount)
		if err != nil {
			fmt.Printf("  Payment failed: %v\n\n", err)
			continue
		}

		fmt.Printf("  ✓ Success using %s\n", processor.GetProcessorName())
		fmt.Printf("  Transaction ID: %s\n\n", transactionID)
	}
}

// Example3_SpecializedFactories demonstrates using specialized factory functions
// for specific types. This provides a more convenient API for common use cases
// while still maintaining the flexibility of the general factory.
func Example3_SpecializedFactories() {
	fmt.Println("=== Example 3: Specialized Factory Functions ===")

	// Using specialized factory functions for cleaner API
	ccProcessor, err := NewCreditCardProcessor("4532015112830366", "123", "12/25")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	ppProcessor, err := NewPayPalProcessor("merchant@shop.com", "api-key-xyz789")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	btcProcessor, err := NewBitcoinProcessor("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "testnet")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// All processors implement the same interface
	processors := []PaymentProcessor{ccProcessor, ppProcessor, btcProcessor}

	for _, processor := range processors {
		fmt.Printf("Processor: %s\n", processor.GetProcessorName())

		if err := processor.ValidateAccount(); err != nil {
			fmt.Printf("  Validation failed: %v\n\n", err)
			continue
		}

		fmt.Printf("  ✓ Account validated\n\n")
	}
}

// Example4_ErrorHandling demonstrates how the factory pattern handles errors
// during object creation and validation. This shows why factories are valuable:
// they centralize validation and error handling logic.
func Example4_ErrorHandling() {
	fmt.Println("=== Example 4: Error Handling ===")

	// Test various error conditions
	errorCases := []struct {
		name   string
		config PaymentConfig
	}{
		{
			name: "Invalid payment method",
			config: PaymentConfig{
				Method: "invalid",
			},
		},
		{
			name: "Invalid credit card",
			config: PaymentConfig{
				Method:     CreditCard,
				CardNumber: "123", // too short
				CVV:        "12",  // too short
				ExpiryDate: "12/25",
			},
		},
		{
			name: "Invalid PayPal email",
			config: PaymentConfig{
				Method:   PayPal,
				Email:    "not-an-email", // missing @
				APIToken: "token123",
			},
		},
		{
			name: "Invalid Bitcoin network",
			config: PaymentConfig{
				Method:        Bitcoin,
				WalletAddress: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				Network:       "invalid-network",
			},
		},
	}

	for _, tc := range errorCases {
		fmt.Printf("Test: %s\n", tc.name)

		processor, err := NewPaymentProcessor(tc.config)
		if err != nil {
			fmt.Printf("  ✓ Expected error caught: %v\n\n", err)
		} else {
			fmt.Printf("  ✗ Should have failed but got: %s\n\n", processor.GetProcessorName())
		}
	}
}

// Example5_RealWorldScenario demonstrates a realistic e-commerce checkout scenario
// where the payment method is determined at runtime based on user selection.
// This showcases why the factory pattern is essential for real-world applications.
func Example5_RealWorldScenario() {
	fmt.Println("=== Example 5: Real-World E-Commerce Checkout ===")

	// Simulate an e-commerce checkout process
	type Order struct {
		OrderID       string
		Amount        float64
		PaymentMethod string
		UserData      map[string]string
	}

	orders := []Order{
		{
			OrderID:       "ORD-001",
			Amount:        299.99,
			PaymentMethod: "creditcard",
			UserData: map[string]string{
				"cardNumber": "4532015112830366",
				"cvv":        "123",
				"expiryDate": "12/25",
			},
		},
		{
			OrderID:       "ORD-002",
			Amount:        149.50,
			PaymentMethod: "paypal",
			UserData: map[string]string{
				"email":    "customer@email.com",
				"apiToken": "paypal-customer-token",
			},
		},
		{
			OrderID:       "ORD-003",
			Amount:        599.00,
			PaymentMethod: "bitcoin",
			UserData: map[string]string{
				"walletAddress": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				"network":       "mainnet",
			},
		},
	}

	// Process each order
	for _, order := range orders {
		fmt.Printf("Processing Order: %s\n", order.OrderID)
		fmt.Printf("Amount: $%.2f\n", order.Amount)
		fmt.Printf("Payment Method: %s\n", order.PaymentMethod)

		// Build configuration from order data
		config := buildPaymentConfig(order.PaymentMethod, order.UserData)

		// Create processor using factory - the same code handles all payment methods!
		processor, err := NewPaymentProcessor(config)
		if err != nil {
			fmt.Printf("  ✗ Failed to create payment processor: %v\n\n", err)
			continue
		}

		// Process the payment
		transactionID, err := processor.ProcessPayment(order.Amount)
		if err != nil {
			fmt.Printf("  ✗ Payment processing failed: %v\n\n", err)
			continue
		}

		fmt.Printf("  ✓ Payment successful!\n")
		fmt.Printf("  Processor: %s\n", processor.GetProcessorName())
		fmt.Printf("  Transaction ID: %s\n\n", transactionID)
	}
}

// buildPaymentConfig is a helper function that demonstrates how you might
// convert runtime data (like user input or database records) into a PaymentConfig.
func buildPaymentConfig(method string, userData map[string]string) PaymentConfig {
	config := PaymentConfig{
		Method: PaymentMethod(method),
	}

	switch PaymentMethod(method) {
	case CreditCard:
		config.CardNumber = userData["cardNumber"]
		config.CVV = userData["cvv"]
		config.ExpiryDate = userData["expiryDate"]

	case PayPal:
		config.Email = userData["email"]
		config.APIToken = userData["apiToken"]

	case Bitcoin:
		config.WalletAddress = userData["walletAddress"]
		config.Network = userData["network"]
	}

	return config
}

// Example6_FactoryWithDependencyInjection demonstrates a more advanced pattern
// where the factory itself can be injected with dependencies. This is common
// in production systems where payment processors need database connections,
// HTTP clients, or configuration services.
func Example6_FactoryWithDependencyInjection() {
	fmt.Println("=== Example 6: Factory with Dependency Injection ===")

	// In a real application, these might be database connections, HTTP clients, etc.
	type PaymentGateway struct {
		apiEndpoint string
		apiKey      string
	}

	// Create a factory closure that captures dependencies
	makeProcessorWithGateway := func(gateway *PaymentGateway, config PaymentConfig) (PaymentProcessor, error) {
		fmt.Printf("Using gateway: %s\n", gateway.apiEndpoint)
		// The factory can inject the gateway into processors that need it
		return NewPaymentProcessor(config)
	}

	// Setup
	gateway := &PaymentGateway{
		apiEndpoint: "https://api.payment-gateway.com",
		apiKey:      "secret-key-abc123",
	}

	config := PaymentConfig{
		Method:     CreditCard,
		CardNumber: "4532015112830366",
		CVV:        "123",
		ExpiryDate: "12/25",
	}

	// Use the factory with injected dependencies
	processor, err := makeProcessorWithGateway(gateway, config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Created processor: %s\n", processor.GetProcessorName())
	transactionID, err := processor.ProcessPayment(499.99)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
		return
	}

	fmt.Printf("Payment successful! Transaction ID: %s\n\n", transactionID)
}
