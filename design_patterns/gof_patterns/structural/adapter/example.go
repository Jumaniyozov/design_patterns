package adapter

import "fmt"

// Example1_PaymentProcessors demonstrates adapting multiple payment gateways.
func Example1_PaymentProcessors() {
	fmt.Println("\n=== Example 1: Payment Gateway Adapters ===")

	// Create adapters for different payment providers
	processors := []PaymentProcessor{
		NewStripeAdapter("sk_test_stripe_key"),
		NewPayPalAdapter("paypal_client_id"),
		NewSquareAdapter("square_location_id"),
	}

	// Process payments through all gateways using the same interface
	fmt.Println("\nProcessing $50.00 USD through all gateways:")
	for _, processor := range processors {
		fmt.Printf("\n%s:\n", processor.GetProviderName())

		txID, err := processor.ProcessPayment(50.00, "USD")
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}

		fmt.Printf("  Transaction successful: %s\n", txID)

		// Process a refund
		err = processor.RefundPayment(txID, 10.00)
		if err != nil {
			fmt.Printf("  Refund error: %v\n", err)
			continue
		}

		fmt.Printf("  Refund of $10.00 processed successfully\n")
	}
}

// Example2_UniformInterface demonstrates the power of a uniform interface.
func Example2_UniformInterface() {
	fmt.Println("\n=== Example 2: Uniform Interface Benefits ===")

	// Function that works with any PaymentProcessor
	processOrder := func(processor PaymentProcessor, amount float64) {
		fmt.Printf("\nProcessing order with %s\n", processor.GetProviderName())

		txID, err := processor.ProcessPayment(amount, "USD")
		if err != nil {
			fmt.Printf("Payment failed: %v\n", err)
			return
		}

		fmt.Printf("Payment successful! Transaction ID: %s\n", txID)
		fmt.Printf("Order confirmed and shipped.\n")
	}

	// Same function works with all adapters
	fmt.Println("\nProcessing multiple orders:")
	processOrder(NewStripeAdapter("key1"), 99.99)
	processOrder(NewPayPalAdapter("key2"), 149.99)
	processOrder(NewSquareAdapter("key3"), 79.99)

	fmt.Println("\nNotice: The processOrder function doesn't know or care")
	fmt.Println("which payment gateway is being used. It just uses the interface!")
}

// Example3_RuntimeSelection demonstrates selecting adapters at runtime.
func Example3_RuntimeSelection() {
	fmt.Println("\n=== Example 3: Runtime Adapter Selection ===")

	// Simulate user preferences
	users := map[string]string{
		"user_1": "stripe",
		"user_2": "paypal",
		"user_3": "square",
		"user_4": "stripe",
	}

	// Factory function to create appropriate adapter
	createProcessor := func(provider string) PaymentProcessor {
		switch provider {
		case "stripe":
			return NewStripeAdapter("stripe_key")
		case "paypal":
			return NewPayPalAdapter("paypal_key")
		case "square":
			return NewSquareAdapter("square_key")
		default:
			return NewStripeAdapter("default_key") // fallback
		}
	}

	fmt.Println("\nProcessing payments based on user preferences:")
	for user, preferredProvider := range users {
		processor := createProcessor(preferredProvider)

		fmt.Printf("\n%s (prefers %s):\n", user, preferredProvider)
		txID, err := processor.ProcessPayment(25.00, "USD")

		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}

		fmt.Printf("  Processed via %s: %s\n", processor.GetProviderName(), txID)
	}
}

// Example4_LoggerAdapters demonstrates adapting different logging libraries.
func Example4_LoggerAdapters() {
	fmt.Println("\n=== Example 4: Logger Adapters ===")

	loggers := []struct {
		name   string
		logger Logger
	}{
		{"Zap Logger", NewZapAdapter()},
		{"Logrus Logger", NewLogrusAdapter()},
	}

	message := "Application event occurred"

	fmt.Println("\nSame logging interface, different implementations:")
	for _, l := range loggers {
		fmt.Printf("\n%s:\n", l.name)
		l.logger.Info(message)
		l.logger.Error("An error occurred")
		l.logger.Debug("Debug information")
	}

	fmt.Println("\nBenefit: Application code uses one Logger interface,")
	fmt.Println("can switch logging libraries without changing business logic!")
}

// Example5_AdapterComposition demonstrates composing adapters.
func Example5_AdapterComposition() {
	fmt.Println("\n=== Example 5: Adapter Composition ===")

	// Simulate an application service that uses both payment and logging
	type OrderService struct {
		processor PaymentProcessor
		logger    Logger
	}

	processOrder := func(service OrderService, amount float64) {
		service.logger.Info(fmt.Sprintf("Processing order for $%.2f", amount))

		txID, err := service.processor.ProcessPayment(amount, "USD")
		if err != nil {
			service.logger.Error(fmt.Sprintf("Payment failed: %v", err))
			return
		}

		service.logger.Info(fmt.Sprintf("Payment successful: %s", txID))
		service.logger.Debug(fmt.Sprintf("Using %s", service.processor.GetProviderName()))
	}

	// Create different service configurations
	services := []OrderService{
		{
			processor: NewStripeAdapter("stripe_key"),
			logger:    NewZapAdapter(),
		},
		{
			processor: NewPayPalAdapter("paypal_key"),
			logger:    NewLogrusAdapter(),
		},
	}

	fmt.Println("\nDifferent service configurations:")
	for i, service := range services {
		fmt.Printf("\nConfiguration %d:\n", i+1)
		processOrder(service, 100.00)
	}

	fmt.Println("\nThis demonstrates how adapters enable flexible composition")
	fmt.Println("of different third-party services with uniform interfaces.")
}

// Example6_AdapterBenefits summarizes the key benefits.
func Example6_AdapterBenefits() {
	fmt.Println("\n=== Example 6: Adapter Pattern Benefits ===")

	fmt.Println("\nKEY BENEFITS:")

	fmt.Println("\n1. INTERFACE COMPATIBILITY")
	fmt.Println("   - Makes incompatible interfaces work together")
	fmt.Println("   - No modification of third-party code required")
	fmt.Println("   - Clean separation between your code and external libraries")

	fmt.Println("\n2. FLEXIBILITY")
	fmt.Println("   - Easy to switch implementations")
	fmt.Println("   - Support multiple providers simultaneously")
	fmt.Println("   - Runtime selection of adapters")

	fmt.Println("\n3. TESTABILITY")
	fmt.Println("   - Mock adapters for unit testing")
	fmt.Println("   - Test business logic without real external services")
	fmt.Println("   - Consistent interface makes testing easier")

	fmt.Println("\n4. MAINTAINABILITY")
	fmt.Println("   - Translation logic isolated in adapters")
	fmt.Println("   - Changes to third-party APIs affect only adapters")
	fmt.Println("   - Business logic remains clean and focused")

	fmt.Println("\n5. OPEN/CLOSED PRINCIPLE")
	fmt.Println("   - Add new adapters without changing existing code")
	fmt.Println("   - Extend functionality by adding new implementations")

	// Demonstrate
	fmt.Println("\n\nPRACTICAL DEMONSTRATION:")

	// Same code works with real and mock implementations
	processors := []PaymentProcessor{
		NewStripeAdapter("key"),
		NewMockPaymentProcessor(),
	}

	for _, p := range processors {
		fmt.Printf("\nUsing %s:\n", p.GetProviderName())
		txID, _ := p.ProcessPayment(10.00, "USD")
		fmt.Printf("  Transaction: %s\n", txID)
	}

	fmt.Println("\nCONCLUSION:")
	fmt.Println("Adapter pattern is essential for integrating third-party services")
	fmt.Println("and maintaining clean, testable, flexible code architecture.")
}
