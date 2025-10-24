package adapter

import (
	"errors"
	"testing"
)

// TestStripeAdapter_Process tests Stripe payment processing.
func TestStripeAdapter_Process(t *testing.T) {
	adapter := NewStripeAdapter("test_api_key")

	payment := Payment{
		CustomerID:  "cus_123",
		Amount:      100.50,
		Currency:    "USD",
		Description: "Test payment",
	}

	receipt, err := adapter.Process(payment)

	if err != nil {
		t.Fatalf("Expected successful processing, got error: %v", err)
	}

	if receipt == nil {
		t.Fatal("Expected receipt, got nil")
	}

	if receipt.Amount != payment.Amount {
		t.Errorf("Expected amount %.2f, got %.2f", payment.Amount, receipt.Amount)
	}

	if receipt.Status != "completed" {
		t.Errorf("Expected status 'completed', got '%s'", receipt.Status)
	}

	if receipt.TransactionID == "" {
		t.Error("Expected transaction ID, got empty string")
	}
}

// TestStripeAdapter_ProcessInvalidAmount tests validation.
func TestStripeAdapter_ProcessInvalidAmount(t *testing.T) {
	adapter := NewStripeAdapter("test_api_key")

	payment := Payment{
		CustomerID: "cus_123",
		Amount:     0, // Invalid amount
		Currency:   "USD",
	}

	_, err := adapter.Process(payment)

	if err == nil {
		t.Error("Expected error for zero amount, got nil")
	}
}

// TestStripeAdapter_Refund tests Stripe refund functionality.
func TestStripeAdapter_Refund(t *testing.T) {
	adapter := NewStripeAdapter("test_api_key")

	err := adapter.Refund("ch_test_123", 50.00)

	if err != nil {
		t.Errorf("Expected successful refund, got error: %v", err)
	}
}

// TestPayPalAdapter_Process tests PayPal payment processing.
func TestPayPalAdapter_Process(t *testing.T) {
	adapter := NewPayPalAdapter("client_id", "client_secret")

	payment := Payment{
		CustomerID:  "user@example.com",
		Amount:      75.25,
		Currency:    "USD",
		Description: "Test payment",
	}

	receipt, err := adapter.Process(payment)

	if err != nil {
		t.Fatalf("Expected successful processing, got error: %v", err)
	}

	if receipt == nil {
		t.Fatal("Expected receipt, got nil")
	}

	if receipt.Amount != payment.Amount {
		t.Errorf("Expected amount %.2f, got %.2f", payment.Amount, receipt.Amount)
	}

	if receipt.TransactionID == "" {
		t.Error("Expected transaction ID, got empty string")
	}
}

// TestPayPalAdapter_ProcessMissingCredentials tests credential validation.
func TestPayPalAdapter_ProcessMissingCredentials(t *testing.T) {
	adapter := NewPayPalAdapter("", "") // Empty credentials

	payment := Payment{
		CustomerID: "user@example.com",
		Amount:     100.00,
		Currency:   "USD",
	}

	_, err := adapter.Process(payment)

	if err == nil {
		t.Error("Expected error for missing credentials, got nil")
	}
}

// TestPaymentProcessorPolymorphism tests that all adapters implement the interface.
func TestPaymentProcessorPolymorphism(t *testing.T) {
	payment := Payment{
		CustomerID: "test_customer",
		Amount:     50.00,
		Currency:   "USD",
	}

	processors := []struct {
		name      string
		processor PaymentProcessor
	}{
		{"Stripe", NewStripeAdapter("test_key")},
		{"PayPal", NewPayPalAdapter("id", "secret")},
	}

	for _, tc := range processors {
		t.Run(tc.name, func(t *testing.T) {
			receipt, err := tc.processor.Process(payment)
			if err != nil {
				t.Errorf("Expected successful processing for %s, got error: %v", tc.name, err)
			}
			if receipt == nil {
				t.Errorf("Expected receipt for %s, got nil", tc.name)
			}
		})
	}
}

// TestSQLiteAdapter_Query tests SQLite query functionality.
func TestSQLiteAdapter_Query(t *testing.T) {
	adapter := NewSQLiteAdapter("test.db")

	result, err := adapter.Query("SELECT * FROM users")

	if err != nil {
		t.Fatalf("Expected successful query, got error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.RowCount != 2 {
		t.Errorf("Expected 2 rows, got %d", result.RowCount)
	}

	if len(result.Rows) != 2 {
		t.Errorf("Expected 2 rows in slice, got %d", len(result.Rows))
	}
}

// TestSQLiteAdapter_Execute tests SQLite execute functionality.
func TestSQLiteAdapter_Execute(t *testing.T) {
	adapter := NewSQLiteAdapter("test.db")

	affected, err := adapter.Execute("INSERT INTO users VALUES (1, 'Test', 'test@example.com')")

	if err != nil {
		t.Fatalf("Expected successful execution, got error: %v", err)
	}

	if affected != 1 {
		t.Errorf("Expected 1 affected row, got %d", affected)
	}
}

// TestSQLiteAdapter_Close tests closing the connection.
func TestSQLiteAdapter_Close(t *testing.T) {
	adapter := NewSQLiteAdapter("test.db")

	err := adapter.Close()

	if err != nil {
		t.Errorf("Expected successful close, got error: %v", err)
	}

	// Verify connection is closed
	if adapter.conn.IsOpen {
		t.Error("Expected connection to be closed")
	}
}

// TestPostgresAdapter_Query tests Postgres query functionality.
func TestPostgresAdapter_Query(t *testing.T) {
	adapter := NewPostgresAdapter("postgres://localhost/test")

	result, err := adapter.Query("SELECT * FROM users")

	if err != nil {
		t.Fatalf("Expected successful query, got error: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.RowCount != 2 {
		t.Errorf("Expected 2 rows, got %d", result.RowCount)
	}
}

// TestPostgresAdapter_Execute tests Postgres execute functionality.
func TestPostgresAdapter_Execute(t *testing.T) {
	adapter := NewPostgresAdapter("postgres://localhost/test")

	affected, err := adapter.Execute("UPDATE users SET active = true")

	if err != nil {
		t.Fatalf("Expected successful execution, got error: %v", err)
	}

	if affected != 2 {
		t.Errorf("Expected 2 affected rows, got %d", affected)
	}
}

// TestDatabasePolymorphism tests that all database adapters implement the interface.
func TestDatabasePolymorphism(t *testing.T) {
	databases := []struct {
		name string
		db   Database
	}{
		{"SQLite", NewSQLiteAdapter("test.db")},
		{"PostgreSQL", NewPostgresAdapter("postgres://localhost/test")},
	}

	for _, tc := range databases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Query
			result, err := tc.db.Query("SELECT * FROM users")
			if err != nil {
				t.Errorf("Expected successful query for %s, got error: %v", tc.name, err)
			}
			if result == nil {
				t.Errorf("Expected result for %s, got nil", tc.name)
			}

			// Test Execute
			affected, err := tc.db.Execute("UPDATE users SET active = true")
			if err != nil {
				t.Errorf("Expected successful execute for %s, got error: %v", tc.name, err)
			}
			if affected == 0 {
				t.Errorf("Expected affected rows for %s, got 0", tc.name)
			}

			// Test Close
			err = tc.db.Close()
			if err != nil {
				t.Errorf("Expected successful close for %s, got error: %v", tc.name, err)
			}
		})
	}
}

// TestStdLoggerAdapter tests standard logger adapter.
func TestStdLoggerAdapter(t *testing.T) {
	adapter := NewStdLoggerAdapter()

	// These should not panic
	adapter.Info("Test info message")
	adapter.Error("Test error", errors.New("test error"))
	adapter.Debug("Test debug message")
}

// TestStructuredLoggerAdapter tests structured logger adapter.
func TestStructuredLoggerAdapter(t *testing.T) {
	adapter := NewStructuredLoggerAdapter()

	// These should not panic
	adapter.Info("Test info message")
	adapter.Error("Test error", errors.New("test error"))
	adapter.Debug("Test debug message")
}

// TestLoggerPolymorphism tests that all logger adapters implement the interface.
func TestLoggerPolymorphism(t *testing.T) {
	loggers := []struct {
		name   string
		logger Logger
	}{
		{"StdLogger", NewStdLoggerAdapter()},
		{"StructuredLogger", NewStructuredLoggerAdapter()},
	}

	for _, tc := range loggers {
		t.Run(tc.name, func(t *testing.T) {
			// These should not panic
			tc.logger.Info("Test message")
			tc.logger.Error("Error occurred", errors.New("test error"))
			tc.logger.Debug("Debug info")
		})
	}
}

// TestAdapterIsolatesExternalChanges demonstrates how adapters isolate code from external changes.
func TestAdapterIsolatesExternalChanges(t *testing.T) {
	// Simulate external library API change
	// If Stripe changes their API, only the adapter needs updating

	adapter := NewStripeAdapter("test_key")
	payment := Payment{
		CustomerID: "cus_123",
		Amount:     100.00,
		Currency:   "USD",
	}

	// Client code remains unchanged even if Stripe's internal API changes
	receipt, err := adapter.Process(payment)

	if err != nil {
		t.Fatalf("Expected successful processing, got error: %v", err)
	}

	if receipt == nil {
		t.Fatal("Expected receipt, got nil")
	}
}

// BenchmarkStripeAdapter benchmarks Stripe adapter performance.
func BenchmarkStripeAdapter(b *testing.B) {
	adapter := NewStripeAdapter("test_key")
	payment := Payment{
		CustomerID: "cus_bench",
		Amount:     50.00,
		Currency:   "USD",
	}

	for i := 0; i < b.N; i++ {
		_, _ = adapter.Process(payment)
	}
}

// BenchmarkPayPalAdapter benchmarks PayPal adapter performance.
func BenchmarkPayPalAdapter(b *testing.B) {
	adapter := NewPayPalAdapter("client_id", "secret")
	payment := Payment{
		CustomerID: "user@example.com",
		Amount:     50.00,
		Currency:   "USD",
	}

	for i := 0; i < b.N; i++ {
		_, _ = adapter.Process(payment)
	}
}

// BenchmarkSQLiteAdapter benchmarks SQLite adapter performance.
func BenchmarkSQLiteAdapter(b *testing.B) {
	adapter := NewSQLiteAdapter("bench.db")

	for i := 0; i < b.N; i++ {
		_, _ = adapter.Query("SELECT * FROM users")
	}
}
