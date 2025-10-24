package adapter

import (
	"errors"
	"fmt"
)

// Example1_PaymentProcessorWithoutAdapter demonstrates the problem: incompatible interfaces.
func Example1_PaymentProcessorWithoutAdapter() {
	fmt.Println("=== Example 1: The Problem - Incompatible Interfaces ===")
	fmt.Println()

	fmt.Println("❌ Without Adapter Pattern:")
	fmt.Println("─────────────────────────")
	fmt.Println("Each payment gateway has a different interface:")
	fmt.Println()
	fmt.Println("Stripe:  gateway.CreateCharge(token, cents, currency)")
	fmt.Println("PayPal:  service.ExecutePayment(email, amount, code, memo)")
	fmt.Println("Square:  processor.MakePayment(customerId, dollars, ...)")
	fmt.Println()
	fmt.Println("Your code would need if/else logic for each gateway,")
	fmt.Println("making it tightly coupled and difficult to maintain.")
	fmt.Println()
	fmt.Println("💡 The Adapter Pattern solves this by creating a unified interface.")
}

// Example2_PaymentProcessorWithAdapter demonstrates using adapters for payment processing.
func Example2_PaymentProcessorWithAdapter() {
	fmt.Println("\n=== Example 2: Payment Processing with Adapters ===")
	fmt.Println()

	// Create payment object (same format regardless of gateway)
	payment := Payment{
		CustomerID:  "cus_12345",
		Amount:      99.99,
		Currency:    "USD",
		Description: "Premium subscription",
	}

	fmt.Println("Payment Request:")
	fmt.Printf("  Customer: %s\n", payment.CustomerID)
	fmt.Printf("  Amount:   $%.2f %s\n", payment.Amount, payment.Currency)
	fmt.Printf("  Desc:     %s\n\n", payment.Description)

	// Process with Stripe (using adapter)
	fmt.Println("Processing with Stripe...")
	stripeProcessor := NewStripeAdapter("sk_test_stripe_key")
	receipt1, err := stripeProcessor.Process(payment)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Success! Transaction ID: %s\n", receipt1.TransactionID)
		fmt.Printf("  Status: %s, Amount: $%.2f\n", receipt1.Status, receipt1.Amount)
	}

	fmt.Println()

	// Process with PayPal (using adapter - same interface!)
	fmt.Println("Processing with PayPal...")
	paypalProcessor := NewPayPalAdapter("client_id", "client_secret")
	receipt2, err := paypalProcessor.Process(payment)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Success! Transaction ID: %s\n", receipt2.TransactionID)
		fmt.Printf("  Status: %s, Amount: $%.2f\n", receipt2.Status, receipt2.Amount)
	}

	fmt.Println()
	fmt.Println("💡 Notice: Both gateways use the same Process() interface,")
	fmt.Println("   but internally they call completely different APIs.")
}

// processPaymentWithAnyGateway demonstrates polymorphism with adapters.
func processPaymentWithAnyGateway(processor PaymentProcessor, gateway string, payment Payment) {
	fmt.Printf("\nProcessing with %s...\n", gateway)
	receipt, err := processor.Process(payment)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("✓ Transaction %s completed successfully\n", receipt.TransactionID)
}

// Example3_PolymorphicPaymentProcessing demonstrates treating all adapters uniformly.
func Example3_PolymorphicPaymentProcessing() {
	fmt.Println("\n=== Example 3: Polymorphic Payment Processing ===")
	fmt.Println()

	payment := Payment{
		CustomerID:  "cus_67890",
		Amount:      149.99,
		Currency:    "USD",
		Description: "Annual subscription",
	}

	// Store different payment processors in a slice
	processors := map[string]PaymentProcessor{
		"Stripe": NewStripeAdapter("sk_test_key"),
		"PayPal": NewPayPalAdapter("client_id", "secret"),
	}

	fmt.Println("Processing same payment through multiple gateways:")
	for name, processor := range processors {
		processPaymentWithAnyGateway(processor, name, payment)
	}

	fmt.Println()
	fmt.Println("💡 This demonstrates polymorphism: We can treat all payment")
	fmt.Println("   gateways uniformly through the PaymentProcessor interface,")
	fmt.Println("   regardless of their underlying implementation.")
}

// Example4_RefundWithAdapters demonstrates refund functionality across adapters.
func Example4_RefundWithAdapters() {
	fmt.Println("\n=== Example 4: Refund Processing with Adapters ===")
	fmt.Println()

	// Process a payment first
	payment := Payment{
		CustomerID: "cus_99999",
		Amount:     79.99,
		Currency:   "USD",
	}

	processor := NewStripeAdapter("sk_test_key")
	receipt, _ := processor.Process(payment)

	fmt.Printf("Original charge: %s for $%.2f\n", receipt.TransactionID, receipt.Amount)

	// Refund the payment
	refundAmount := 79.99
	fmt.Printf("\nIssuing refund of $%.2f...\n", refundAmount)
	err := processor.Refund(receipt.TransactionID, refundAmount)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Refund successful\n")
	}

	fmt.Println()
	fmt.Println("💡 The adapter translates refund calls to each gateway's")
	fmt.Println("   specific refund method and parameter format.")
}

// Example5_DatabaseAdapters demonstrates database driver adapters.
func Example5_DatabaseAdapters() {
	fmt.Println("\n=== Example 5: Database Driver Adapters ===")
	fmt.Println()

	query := "SELECT id, name, email FROM users WHERE active = true"

	fmt.Println("Querying SQLite database...")
	sqliteDB := NewSQLiteAdapter("users.db")
	result1, err := sqliteDB.Query(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Retrieved %d rows from SQLite\n", result1.RowCount)
		for i, row := range result1.Rows {
			fmt.Printf("  Row %d: %v\n", i+1, row)
		}
	}

	fmt.Println()

	fmt.Println("Querying PostgreSQL database...")
	postgresDB := NewPostgresAdapter("postgres://localhost/mydb")
	result2, err := postgresDB.Query(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Retrieved %d rows from PostgreSQL\n", result2.RowCount)
		for i, row := range result2.Rows {
			fmt.Printf("  Row %d: %v\n", i+1, row)
		}
	}

	fmt.Println()
	fmt.Println("💡 Both databases use the same Query() interface, but")
	fmt.Println("   internally they call different library methods and")
	fmt.Println("   transform data formats differently.")
}

// queryAnyDatabase demonstrates polymorphic database operations.
func queryAnyDatabase(db Database, dbName string) {
	fmt.Printf("\nQuerying %s...\n", dbName)
	result, err := db.Query("SELECT * FROM products LIMIT 10")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("✓ Retrieved %d rows\n", result.RowCount)
}

// Example6_PolymorphicDatabaseOperations demonstrates treating all databases uniformly.
func Example6_PolymorphicDatabaseOperations() {
	fmt.Println("\n=== Example 6: Polymorphic Database Operations ===")
	fmt.Println()

	// Create different database adapters
	databases := map[string]Database{
		"SQLite":     NewSQLiteAdapter("app.db"),
		"PostgreSQL": NewPostgresAdapter("postgres://localhost/app"),
	}

	fmt.Println("Executing same query across multiple databases:")
	for name, db := range databases {
		queryAnyDatabase(db, name)
	}

	fmt.Println()
	fmt.Println("💡 This allows your application to be database-agnostic.")
	fmt.Println("   Switching databases only requires changing the adapter,")
	fmt.Println("   not the application code.")
}

// Example7_LoggerAdapters demonstrates adapting different logging libraries.
func Example7_LoggerAdapters() {
	fmt.Println("\n=== Example 7: Logger Adapters ===")
	fmt.Println()

	message := "User logged in successfully"
	err := errors.New("authentication failed")

	fmt.Println("Using Standard Logger Adapter:")
	stdLogger := NewStdLoggerAdapter()
	stdLogger.Info(message)
	stdLogger.Error("Login failed", err)
	stdLogger.Debug("Checking credentials")

	fmt.Println()

	fmt.Println("Using Structured Logger Adapter:")
	structLogger := NewStructuredLoggerAdapter()
	structLogger.Info(message)
	structLogger.Error("Login failed", err)
	structLogger.Debug("Checking credentials")

	fmt.Println()
	fmt.Println("💡 Both loggers implement the same Logger interface,")
	fmt.Println("   allowing you to swap logging implementations without")
	fmt.Println("   changing application code.")
}

// logWithAnyLogger demonstrates polymorphic logging.
func logWithAnyLogger(logger Logger, loggerName string) {
	fmt.Printf("\nUsing %s:\n", loggerName)
	logger.Info("Application started")
	logger.Debug("Loading configuration")
	logger.Error("Failed to connect", errors.New("connection timeout"))
}

// Example8_PolymorphicLogging demonstrates treating all loggers uniformly.
func Example8_PolymorphicLogging() {
	fmt.Println("\n=== Example 8: Polymorphic Logging ===")
	fmt.Println()

	loggers := map[string]Logger{
		"Standard Logger":    NewStdLoggerAdapter(),
		"Structured Logger":  NewStructuredLoggerAdapter(),
	}

	fmt.Println("Logging same messages with different implementations:")
	for name, logger := range loggers {
		logWithAnyLogger(logger, name)
	}

	fmt.Println()
	fmt.Println("💡 This demonstrates how adapters enable you to switch")
	fmt.Println("   third-party libraries without touching business logic.")
}

// Example9_AdapterBenefits demonstrates key benefits of the Adapter pattern.
func Example9_AdapterBenefits() {
	fmt.Println("\n=== Example 9: Key Benefits of Adapter Pattern ===")
	fmt.Println()

	fmt.Println("1. INTEGRATION")
	fmt.Println("   ✓ Integrate libraries with incompatible interfaces")
	fmt.Println("   ✓ No need to modify external library code")
	fmt.Println()

	fmt.Println("2. DECOUPLING")
	fmt.Println("   ✓ Application code depends on stable interface")
	fmt.Println("   ✓ External library changes affect only adapter")
	fmt.Println()

	fmt.Println("3. FLEXIBILITY")
	fmt.Println("   ✓ Switch implementations without changing client code")
	fmt.Println("   ✓ Support multiple backends simultaneously")
	fmt.Println()

	fmt.Println("4. TESTABILITY")
	fmt.Println("   ✓ Easy to create mock adapters for testing")
	fmt.Println("   ✓ Test without depending on external services")
	fmt.Println()

	fmt.Println("5. SINGLE RESPONSIBILITY")
	fmt.Println("   ✓ Adapter focuses solely on translation")
	fmt.Println("   ✓ Clear separation of concerns")
}

// Example10_WhenNotToUseAdapter demonstrates anti-patterns.
func Example10_WhenNotToUseAdapter() {
	fmt.Println("\n=== Example 10: When NOT to Use Adapter Pattern ===")
	fmt.Println()

	fmt.Println("❌ AVOID adapters when:")
	fmt.Println()
	fmt.Println("1. Interfaces are already compatible")
	fmt.Println("   → Unnecessary indirection adds complexity")
	fmt.Println()
	fmt.Println("2. Only one implementation ever needed")
	fmt.Println("   → Direct usage is simpler and clearer")
	fmt.Println()
	fmt.Println("3. Adapter becomes more complex than direct integration")
	fmt.Println("   → Defeats the purpose of simplification")
	fmt.Println()
	fmt.Println("4. Performance is critical")
	fmt.Println("   → Adapter adds extra function call overhead")
	fmt.Println()
	fmt.Println("5. Adapter leaks too many implementation details")
	fmt.Println("   → Breaks abstraction and couples to specifics")
	fmt.Println()
	fmt.Println("💡 Use adapters strategically when they provide clear value,")
	fmt.Println("   not reflexively for every third-party integration.")
}
