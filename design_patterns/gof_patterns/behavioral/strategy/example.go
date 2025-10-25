package strategy

import "fmt"

// Example1_ShippingStrategies demonstrates different shipping calculation strategies.
func Example1_ShippingStrategies() {
	fmt.Println("\n=== Example 1: Shipping Cost Strategies ===")

	// Sample orders
	orders := []Order{
		{
			Items:       []string{"Book", "Pen"},
			TotalPrice:  25.00,
			Weight:      0.5,
			Destination: "USA",
		},
		{
			Items:       []string{"Laptop", "Mouse", "Keyboard"},
			TotalPrice:  1200.00,
			Weight:      3.5,
			Destination: "USA",
		},
		{
			Items:       []string{"Camera"},
			TotalPrice:  800.00,
			Weight:      2.0,
			Destination: "Asia",
		},
	}

	strategies := []ShippingStrategy{
		NewStandardShipping(),
		NewExpressShipping(),
		NewFreeShipping(100.00),
		NewInternationalShipping(),
	}

	// Calculate shipping for each order with each strategy
	for i, order := range orders {
		fmt.Printf("\nOrder %d: $%.2f (%.2fkg) to %s\n",
			i+1, order.TotalPrice, order.Weight, order.Destination)

		calculator := NewShippingCalculator(nil)

		for _, strategy := range strategies {
			calculator.SetStrategy(strategy)
			cost := calculator.Calculate(order)
			fmt.Printf("  %s: $%.2f\n", strategy.GetName(), cost)
		}
	}
}

// Example2_RuntimeStrategyChange demonstrates changing strategies at runtime.
func Example2_RuntimeStrategyChange() {
	fmt.Println("\n=== Example 2: Runtime Strategy Changes ===")

	order := Order{
		Items:       []string{"Product A", "Product B"},
		TotalPrice:  75.00,
		Weight:      1.5,
		Destination: "USA",
	}

	calculator := NewShippingCalculator(NewStandardShipping())

	fmt.Printf("Order: $%.2f, Weight: %.2fkg\n\n", order.TotalPrice, order.Weight)

	// Start with standard shipping
	fmt.Println("Customer selects Standard Shipping:")
	cost := calculator.Calculate(order)
	fmt.Printf("Total: $%.2f\n", cost)

	// Customer changes mind
	fmt.Println("\nCustomer changes to Express Shipping:")
	calculator.SetStrategy(NewExpressShipping())
	cost = calculator.Calculate(order)
	fmt.Printf("Total: $%.2f\n", cost)

	// Customer qualifies for free shipping after adding more items
	order.TotalPrice = 120.00
	fmt.Println("\nCustomer adds more items (total now $120):")
	calculator.SetStrategy(NewFreeShipping(100.00))
	cost = calculator.Calculate(order)
	fmt.Printf("Total: $%.2f\n", cost)

	fmt.Println("\nThis demonstrates the power of runtime strategy switching!")
}

// Example3_CompressionStrategies demonstrates file compression strategies.
func Example3_CompressionStrategies() {
	fmt.Println("\n=== Example 3: Compression Strategies ===")

	data := []byte("This is some sample data that needs to be compressed for storage and transmission")
	fmt.Printf("Original data: %d bytes\n", len(data))

	compressor := NewFileCompressor(&ZipCompression{})

	strategies := []CompressionStrategy{
		&ZipCompression{},
		&GzipCompression{},
		&Bzip2Compression{},
	}

	fmt.Println("\nTrying different compression algorithms:")
	for _, strategy := range strategies {
		compressor.SetStrategy(strategy)
		compressed := compressor.Compress(data)
		fmt.Printf("  %s\n", string(compressed))
	}

	fmt.Println("\nThe same FileCompressor can use different compression algorithms")
	fmt.Println("without changing its code - just swap the strategy!")
}

// Example4_SortingStrategies demonstrates different sorting strategies.
func Example4_SortingStrategies() {
	fmt.Println("\n=== Example 4: Sorting Strategies ===")

	data := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Printf("Original data: %v\n", data)

	strategies := []SortStrategy{
		&BubbleSort{},
		&QuickSort{},
	}

	sorter := NewSorter(&BubbleSort{})

	fmt.Println("\nSorting with different algorithms:")
	for _, strategy := range strategies {
		sorter.SetStrategy(strategy)
		sorted := sorter.Sort(data)
		fmt.Printf("  %s result: %v\n", strategy.GetName(), sorted)
	}

	fmt.Println("\nFor small datasets, BubbleSort is fine.")
	fmt.Println("For large datasets, QuickSort is more efficient.")
	fmt.Println("Strategy pattern lets you choose the right algorithm for the context!")
}

// Example5_PaymentStrategies demonstrates different payment methods.
func Example5_PaymentStrategies() {
	fmt.Println("\n=== Example 5: Payment Strategies ===")

	amount := 299.99

	paymentMethods := []PaymentStrategy{
		NewCreditCardPayment("4532-1234-5678-9010", "123"),
		NewPayPalPayment("user@example.com"),
		NewCryptoPayment("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "BTC"),
	}

	processor := NewPaymentProcessor(paymentMethods[0])

	fmt.Printf("Processing $%.2f payment with different methods:\n\n", amount)

	for _, method := range paymentMethods {
		processor.SetStrategy(method)
		err := processor.ProcessPayment(amount)

		if err != nil {
			fmt.Printf("  Error: %v\n", err)
		}
		fmt.Println()
	}

	fmt.Println("Same payment processing logic works with any payment method!")
}

// Example6_DiscountStrategies demonstrates pricing with different discount strategies.
func Example6_DiscountStrategies() {
	fmt.Println("\n=== Example 6: Discount Strategies ===")

	basePrice := 100.00
	fmt.Printf("Base price: $%.2f\n\n", basePrice)

	discounts := []DiscountStrategy{
		&NoDiscount{},
		NewPercentageDiscount(10),
		NewPercentageDiscount(25),
		NewPercentageDiscount(50),
		NewFixedAmountDiscount(15.00),
		NewFixedAmountDiscount(30.00),
	}

	calculator := NewPricingCalculator(&NoDiscount{})

	fmt.Println("Applying different discounts:")
	for _, discount := range discounts {
		calculator.SetStrategy(discount)
		finalPrice := calculator.CalculateFinalPrice(basePrice)
		savings := basePrice - finalPrice

		fmt.Printf("  %-20s → $%.2f (saved $%.2f)\n",
			discount.GetDescription(), finalPrice, savings)
	}

	fmt.Println("\nDifferent discount strategies can be applied based on:")
	fmt.Println("  - Customer loyalty level")
	fmt.Println("  - Seasonal promotions")
	fmt.Println("  - Bulk purchase quantities")
	fmt.Println("  - Special coupon codes")
}

// Example7_StrategySelection demonstrates selecting strategies based on context.
func Example7_StrategySelection() {
	fmt.Println("\n=== Example 7: Context-Based Strategy Selection ===")

	// Function to select shipping strategy based on order characteristics
	selectShippingStrategy := func(order Order) ShippingStrategy {
		// Free shipping for high-value orders
		if order.TotalPrice >= 100 {
			return NewFreeShipping(100.00)
		}

		// International shipping for non-USA destinations
		if order.Destination != "USA" {
			return NewInternationalShipping()
		}

		// Express if customer is willing to pay
		if order.TotalPrice >= 50 {
			fmt.Println("  [Suggestion] Consider Express Shipping for $10 more!")
		}

		// Default to standard
		return NewStandardShipping()
	}

	orders := []Order{
		{TotalPrice: 30, Weight: 1.0, Destination: "USA"},
		{TotalPrice: 150, Weight: 2.0, Destination: "USA"},
		{TotalPrice: 80, Weight: 1.5, Destination: "Europe"},
	}

	calculator := NewShippingCalculator(nil)

	fmt.Println("Intelligent strategy selection:\n")
	for i, order := range orders {
		fmt.Printf("Order %d: $%.2f to %s\n", i+1, order.TotalPrice, order.Destination)

		strategy := selectShippingStrategy(order)
		calculator.SetStrategy(strategy)
		cost := calculator.Calculate(order)

		fmt.Printf("  Selected: %s, Cost: $%.2f\n\n", strategy.GetName(), cost)
	}
}

// Example8_StrategyPatternBenefits summarizes the benefits.
func Example8_StrategyPatternBenefits() {
	fmt.Println("\n=== Example 8: Strategy Pattern Benefits ===")

	fmt.Println("\nKEY BENEFITS:")

	fmt.Println("\n1. ELIMINATE COMPLEX CONDITIONALS")
	fmt.Println("   Before: if-else chains for every algorithm variation")
	fmt.Println("   After:  Clean polymorphism through interfaces")

	fmt.Println("\n2. RUNTIME FLEXIBILITY")
	fmt.Println("   - Change algorithm behavior on the fly")
	fmt.Println("   - Respond to user preferences")
	fmt.Println("   - Adapt to runtime conditions")

	fmt.Println("\n3. OPEN/CLOSED PRINCIPLE")
	fmt.Println("   - Add new strategies without modifying existing code")
	fmt.Println("   - Extend functionality by implementing interface")

	fmt.Println("\n4. TESTABILITY")
	fmt.Println("   - Test each strategy independently")
	fmt.Println("   - Mock strategies for unit testing")
	fmt.Println("   - Easier to verify correctness")

	fmt.Println("\n5. REUSABILITY")
	fmt.Println("   - Strategies can be shared across contexts")
	fmt.Println("   - Same algorithm in different parts of application")

	fmt.Println("\nREAL-WORLD APPLICATIONS:")
	fmt.Println("  ✓ Shipping calculations (as demonstrated)")
	fmt.Println("  ✓ Compression algorithms (ZIP, GZIP, BZIP2)")
	fmt.Println("  ✓ Sorting algorithms (QuickSort, MergeSort, BubbleSort)")
	fmt.Println("  ✓ Payment processing (CreditCard, PayPal, Crypto)")
	fmt.Println("  ✓ Discount calculations (Percentage, Fixed, Tiered)")
	fmt.Println("  ✓ Validation rules (Email, Phone, CreditCard)")
	fmt.Println("  ✓ Image rendering (PNG, JPEG, WebP)")
	fmt.Println("  ✓ Notification delivery (Email, SMS, Push)")

	fmt.Println("\nWHEN TO USE:")
	fmt.Println("  • Multiple ways to perform the same task")
	fmt.Println("  • Algorithm selection depends on runtime conditions")
	fmt.Println("  • Complex conditional logic based on type")
	fmt.Println("  • Need to swap implementations easily")

	fmt.Println("\nCONCLUSION:")
	fmt.Println("Strategy pattern is one of the most useful behavioral patterns.")
	fmt.Println("It promotes flexibility, testability, and clean code architecture.")
	fmt.Println("Go's interfaces make this pattern particularly elegant!")
}
