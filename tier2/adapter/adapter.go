// Package adapter demonstrates the Adapter Pattern, a structural design pattern
// that converts the interface of a class into another interface that clients expect.
//
// The Adapter Pattern is essential for:
// - Integrating third-party libraries with incompatible interfaces
// - Creating stable internal APIs that hide external complexity
// - Supporting multiple implementations with different interfaces
// - Isolating code from external API changes
package adapter

import (
	"errors"
	"fmt"
	"time"
)

// =============================================================================
// Example 1: Payment Gateway Adapters
// =============================================================================

// Payment represents a standardized payment request across all gateways.
type Payment struct {
	CustomerID string
	Amount     float64
	Currency   string
	Description string
}

// Receipt represents a standardized payment receipt.
type Receipt struct {
	TransactionID string
	Status        string
	ProcessedAt   time.Time
	Amount        float64
}

// PaymentProcessor is the target interface that our application expects.
// This is the stable internal interface that insulates us from external changes.
type PaymentProcessor interface {
	Process(payment Payment) (*Receipt, error)
	Refund(transactionID string, amount float64) error
}

// --- Stripe External Library (Adaptee) ---

// StripePaymentGateway represents a third-party Stripe library with its own interface.
// In reality, this would be an external package we don't control.
type StripePaymentGateway struct {
	APIKey string
}

// CreateCharge is Stripe's method for processing payments (different from our interface).
func (s *StripePaymentGateway) CreateCharge(customerToken string, amountInCents int, currency string) (string, error) {
	if s.APIKey == "" {
		return "", errors.New("stripe: API key required")
	}
	if amountInCents <= 0 {
		return "", errors.New("stripe: amount must be positive")
	}
	// Simulate Stripe API call
	chargeID := fmt.Sprintf("ch_stripe_%d", time.Now().Unix())
	return chargeID, nil
}

// RefundCharge is Stripe's method for refunding (different signature).
func (s *StripePaymentGateway) RefundCharge(chargeID string, amountInCents int) error {
	if chargeID == "" {
		return errors.New("stripe: charge ID required")
	}
	// Simulate refund
	return nil
}

// StripeAdapter adapts StripePaymentGateway to our PaymentProcessor interface.
type StripeAdapter struct {
	gateway *StripePaymentGateway
}

// NewStripeAdapter creates a new adapter for Stripe.
func NewStripeAdapter(apiKey string) *StripeAdapter {
	return &StripeAdapter{
		gateway: &StripePaymentGateway{APIKey: apiKey},
	}
}

// Process adapts our Payment format to Stripe's CreateCharge format.
func (a *StripeAdapter) Process(payment Payment) (*Receipt, error) {
	// Convert dollars to cents (Stripe uses cents)
	amountInCents := int(payment.Amount * 100)

	// Call Stripe's API with adapted parameters
	chargeID, err := a.gateway.CreateCharge(payment.CustomerID, amountInCents, payment.Currency)
	if err != nil {
		return nil, fmt.Errorf("stripe adapter: %w", err)
	}

	// Convert Stripe's response to our Receipt format
	return &Receipt{
		TransactionID: chargeID,
		Status:        "completed",
		ProcessedAt:   time.Now(),
		Amount:        payment.Amount,
	}, nil
}

// Refund adapts our refund interface to Stripe's format.
func (a *StripeAdapter) Refund(transactionID string, amount float64) error {
	amountInCents := int(amount * 100)
	return a.gateway.RefundCharge(transactionID, amountInCents)
}

// --- PayPal External Library (Adaptee) ---

// PayPalPaymentService represents a third-party PayPal library.
type PayPalPaymentService struct {
	ClientID     string
	ClientSecret string
}

// ExecutePayment is PayPal's method (completely different from Stripe and our interface).
func (p *PayPalPaymentService) ExecutePayment(accountEmail string, paymentAmount float64, currencyCode string, memo string) (*PayPalTransaction, error) {
	if p.ClientID == "" || p.ClientSecret == "" {
		return nil, errors.New("paypal: credentials required")
	}
	if paymentAmount <= 0 {
		return nil, errors.New("paypal: amount must be positive")
	}
	// Simulate PayPal API call
	return &PayPalTransaction{
		ID:        fmt.Sprintf("PAYPAL-%d", time.Now().Unix()),
		Status:    "COMPLETED",
		Timestamp: time.Now(),
	}, nil
}

// ReversePayment is PayPal's refund method.
func (p *PayPalPaymentService) ReversePayment(txnID string, refundAmount float64) error {
	if txnID == "" {
		return errors.New("paypal: transaction ID required")
	}
	return nil
}

// PayPalTransaction represents PayPal's transaction response.
type PayPalTransaction struct {
	ID        string
	Status    string
	Timestamp time.Time
}

// PayPalAdapter adapts PayPalPaymentService to our PaymentProcessor interface.
type PayPalAdapter struct {
	service *PayPalPaymentService
}

// NewPayPalAdapter creates a new adapter for PayPal.
func NewPayPalAdapter(clientID, clientSecret string) *PayPalAdapter {
	return &PayPalAdapter{
		service: &PayPalPaymentService{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		},
	}
}

// Process adapts our Payment format to PayPal's ExecutePayment format.
func (a *PayPalAdapter) Process(payment Payment) (*Receipt, error) {
	// Call PayPal's API with adapted parameters
	txn, err := a.service.ExecutePayment(
		payment.CustomerID,
		payment.Amount,
		payment.Currency,
		payment.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("paypal adapter: %w", err)
	}

	// Convert PayPal's response to our Receipt format
	return &Receipt{
		TransactionID: txn.ID,
		Status:        txn.Status,
		ProcessedAt:   txn.Timestamp,
		Amount:        payment.Amount,
	}, nil
}

// Refund adapts our refund interface to PayPal's format.
func (a *PayPalAdapter) Refund(transactionID string, amount float64) error {
	return a.service.ReversePayment(transactionID, amount)
}

// =============================================================================
// Example 2: Database Driver Adapters
// =============================================================================

// QueryResult represents standardized query results.
type QueryResult struct {
	Rows     []map[string]interface{}
	RowCount int
}

// Database is the target interface that our application expects.
type Database interface {
	Query(sql string, args ...interface{}) (*QueryResult, error)
	Execute(sql string, args ...interface{}) (int64, error)
	Close() error
}

// --- SQLite External Library (Adaptee) ---

// SQLiteConnection represents an external SQLite library.
type SQLiteConnection struct {
	FilePath string
	IsOpen   bool
}

// RunQuery is SQLite's query method (different signature).
func (s *SQLiteConnection) RunQuery(queryString string) ([][]string, error) {
	if !s.IsOpen {
		return nil, errors.New("sqlite: connection not open")
	}
	// Simulate query execution
	mockData := [][]string{
		{"1", "John Doe", "john@example.com"},
		{"2", "Jane Smith", "jane@example.com"},
	}
	return mockData, nil
}

// ExecuteStatement is SQLite's execute method.
func (s *SQLiteConnection) ExecuteStatement(stmt string) (int, error) {
	if !s.IsOpen {
		return 0, errors.New("sqlite: connection not open")
	}
	// Simulate execution
	return 1, nil
}

// Disconnect closes the connection.
func (s *SQLiteConnection) Disconnect() error {
	s.IsOpen = false
	return nil
}

// SQLiteAdapter adapts SQLiteConnection to our Database interface.
type SQLiteAdapter struct {
	conn *SQLiteConnection
}

// NewSQLiteAdapter creates a new SQLite adapter.
func NewSQLiteAdapter(filePath string) *SQLiteAdapter {
	return &SQLiteAdapter{
		conn: &SQLiteConnection{
			FilePath: filePath,
			IsOpen:   true,
		},
	}
}

// Query adapts SQLite's RunQuery to our interface.
func (a *SQLiteAdapter) Query(sql string, args ...interface{}) (*QueryResult, error) {
	// Call SQLite's method
	rows, err := a.conn.RunQuery(sql)
	if err != nil {
		return nil, fmt.Errorf("sqlite adapter: %w", err)
	}

	// Convert SQLite's string slices to our map format
	result := &QueryResult{
		Rows:     make([]map[string]interface{}, len(rows)),
		RowCount: len(rows),
	}

	for i, row := range rows {
		result.Rows[i] = map[string]interface{}{
			"id":    row[0],
			"name":  row[1],
			"email": row[2],
		}
	}

	return result, nil
}

// Execute adapts SQLite's ExecuteStatement to our interface.
func (a *SQLiteAdapter) Execute(sql string, args ...interface{}) (int64, error) {
	affected, err := a.conn.ExecuteStatement(sql)
	if err != nil {
		return 0, fmt.Errorf("sqlite adapter: %w", err)
	}
	return int64(affected), nil
}

// Close adapts SQLite's Disconnect to our interface.
func (a *SQLiteAdapter) Close() error {
	return a.conn.Disconnect()
}

// --- PostgreSQL External Library (Adaptee) ---

// PostgresDB represents an external Postgres library.
type PostgresDB struct {
	ConnectionString string
	Connected        bool
}

// ExecSQL is Postgres's method for executing queries.
func (p *PostgresDB) ExecSQL(command string, params []interface{}) (map[string][]interface{}, error) {
	if !p.Connected {
		return nil, errors.New("postgres: not connected")
	}
	// Simulate query execution
	return map[string][]interface{}{
		"id":    {1, 2},
		"name":  {"Alice", "Bob"},
		"email": {"alice@example.com", "bob@example.com"},
	}, nil
}

// RunCommand is Postgres's execute method.
func (p *PostgresDB) RunCommand(cmd string, params []interface{}) (int, error) {
	if !p.Connected {
		return 0, errors.New("postgres: not connected")
	}
	return 2, nil
}

// Terminate closes the connection.
func (p *PostgresDB) Terminate() error {
	p.Connected = false
	return nil
}

// PostgresAdapter adapts PostgresDB to our Database interface.
type PostgresAdapter struct {
	db *PostgresDB
}

// NewPostgresAdapter creates a new Postgres adapter.
func NewPostgresAdapter(connectionString string) *PostgresAdapter {
	return &PostgresAdapter{
		db: &PostgresDB{
			ConnectionString: connectionString,
			Connected:        true,
		},
	}
}

// Query adapts Postgres's ExecSQL to our interface.
func (a *PostgresAdapter) Query(sql string, args ...interface{}) (*QueryResult, error) {
	// Call Postgres's method
	data, err := a.db.ExecSQL(sql, args)
	if err != nil {
		return nil, fmt.Errorf("postgres adapter: %w", err)
	}

	// Convert Postgres's map format to our slice of maps
	result := &QueryResult{
		Rows: make([]map[string]interface{}, 0),
	}

	// Transform column-oriented data to row-oriented
	if len(data) > 0 {
		numRows := len(data["id"])
		result.RowCount = numRows

		for i := 0; i < numRows; i++ {
			row := make(map[string]interface{})
			for key, values := range data {
				row[key] = values[i]
			}
			result.Rows = append(result.Rows, row)
		}
	}

	return result, nil
}

// Execute adapts Postgres's RunCommand to our interface.
func (a *PostgresAdapter) Execute(sql string, args ...interface{}) (int64, error) {
	affected, err := a.db.RunCommand(sql, args)
	if err != nil {
		return 0, fmt.Errorf("postgres adapter: %w", err)
	}
	return int64(affected), nil
}

// Close adapts Postgres's Terminate to our interface.
func (a *PostgresAdapter) Close() error {
	return a.db.Terminate()
}

// =============================================================================
// Example 3: Logger Adapters
// =============================================================================

// Logger is the target interface that our application expects.
type Logger interface {
	Info(message string)
	Error(message string, err error)
	Debug(message string)
}

// --- Standard Library Logger (Adaptee) ---

// StdLogger represents Go's standard library logger (different interface).
type StdLogger struct{}

// Print is the standard logger's method.
func (s *StdLogger) Print(level string, msg string) {
	fmt.Printf("[%s] %s\n", level, msg)
}

// StdLoggerAdapter adapts StdLogger to our Logger interface.
type StdLoggerAdapter struct {
	logger *StdLogger
}

// NewStdLoggerAdapter creates a new standard logger adapter.
func NewStdLoggerAdapter() *StdLoggerAdapter {
	return &StdLoggerAdapter{
		logger: &StdLogger{},
	}
}

// Info adapts to standard logger's Print method.
func (a *StdLoggerAdapter) Info(message string) {
	a.logger.Print("INFO", message)
}

// Error adapts to standard logger's Print method with error formatting.
func (a *StdLoggerAdapter) Error(message string, err error) {
	a.logger.Print("ERROR", fmt.Sprintf("%s: %v", message, err))
}

// Debug adapts to standard logger's Print method.
func (a *StdLoggerAdapter) Debug(message string) {
	a.logger.Print("DEBUG", message)
}

// --- Structured Logger (Adaptee) ---

// StructuredLogger represents a third-party structured logger.
type StructuredLogger struct{}

// Log is the structured logger's method (takes key-value pairs).
func (s *StructuredLogger) Log(level string, fields map[string]interface{}) {
	fmt.Printf("[%s]", level)
	for k, v := range fields {
		fmt.Printf(" %s=%v", k, v)
	}
	fmt.Println()
}

// StructuredLoggerAdapter adapts StructuredLogger to our Logger interface.
type StructuredLoggerAdapter struct {
	logger *StructuredLogger
}

// NewStructuredLoggerAdapter creates a new structured logger adapter.
func NewStructuredLoggerAdapter() *StructuredLoggerAdapter {
	return &StructuredLoggerAdapter{
		logger: &StructuredLogger{},
	}
}

// Info adapts to structured logger's Log method.
func (a *StructuredLoggerAdapter) Info(message string) {
	a.logger.Log("INFO", map[string]interface{}{"message": message})
}

// Error adapts to structured logger's Log method with error field.
func (a *StructuredLoggerAdapter) Error(message string, err error) {
	a.logger.Log("ERROR", map[string]interface{}{
		"message": message,
		"error":   err.Error(),
	})
}

// Debug adapts to structured logger's Log method.
func (a *StructuredLoggerAdapter) Debug(message string) {
	a.logger.Log("DEBUG", map[string]interface{}{"message": message})
}
