// Package adapter demonstrates the Adapter pattern in Go.
// The Adapter pattern converts the interface of a class into another interface
// clients expect, allowing incompatible interfaces to work together.
package adapter

import (
	"fmt"
	"strings"
	"time"
)

// PaymentProcessor is our target interface that all payment gateways should implement.
type PaymentProcessor interface {
	ProcessPayment(amount float64, currency string) (string, error)
	RefundPayment(transactionID string, amount float64) error
	GetProviderName() string
}

// StripeAPI simulates the Stripe third-party library with its own interface.
type StripeAPI struct {
	apiKey string
}

// Charge is Stripe's method (different name and parameters).
func (s *StripeAPI) Charge(amountCents int, curr string, token string) (string, error) {
	txID := fmt.Sprintf("stripe_tx_%d", time.Now().Unix())
	fmt.Printf("[Stripe] Charging %d cents in %s\n", amountCents, curr)
	return txID, nil
}

// CreateRefund is Stripe's refund method.
func (s *StripeAPI) CreateRefund(chargeID string, amountCents int) error {
	fmt.Printf("[Stripe] Refunding %d cents for charge %s\n", amountCents, chargeID)
	return nil
}

// StripeAdapter adapts StripeAPI to PaymentProcessor interface.
type StripeAdapter struct {
	stripe *StripeAPI
	token  string
}

// NewStripeAdapter creates a new Stripe adapter.
func NewStripeAdapter(apiKey string) PaymentProcessor {
	return &StripeAdapter{
		stripe: &StripeAPI{apiKey: apiKey},
		token:  "tok_visa_test",
	}
}

func (s *StripeAdapter) ProcessPayment(amount float64, currency string) (string, error) {
	// Adapt: convert dollars to cents, call Stripe's Charge method
	amountCents := int(amount * 100)
	return s.stripe.Charge(amountCents, strings.ToUpper(currency), s.token)
}

func (s *StripeAdapter) RefundPayment(transactionID string, amount float64) error {
	amountCents := int(amount * 100)
	return s.stripe.CreateRefund(transactionID, amountCents)
}

func (s *StripeAdapter) GetProviderName() string {
	return "Stripe"
}

// PayPalAPI simulates the PayPal third-party library.
type PayPalAPI struct {
	clientID string
}

// MakePayment is PayPal's payment method (different interface).
func (p *PayPalAPI) MakePayment(price float64, currencyCode string, method string) (map[string]interface{}, error) {
	txID := fmt.Sprintf("paypal_tx_%d", time.Now().Unix())
	fmt.Printf("[PayPal] Processing payment of %.2f %s\n", price, currencyCode)
	return map[string]interface{}{
		"transaction_id": txID,
		"status":         "completed",
	}, nil
}

// IssueRefund is PayPal's refund method.
func (p *PayPalAPI) IssueRefund(txID string, refundAmount float64) error {
	fmt.Printf("[PayPal] Issuing refund of %.2f for transaction %s\n", refundAmount, txID)
	return nil
}

// PayPalAdapter adapts PayPalAPI to PaymentProcessor interface.
type PayPalAdapter struct {
	paypal *PayPalAPI
}

// NewPayPalAdapter creates a new PayPal adapter.
func NewPayPalAdapter(clientID string) PaymentProcessor {
	return &PayPalAdapter{
		paypal: &PayPalAPI{clientID: clientID},
	}
}

func (p *PayPalAdapter) ProcessPayment(amount float64, currency string) (string, error) {
	// Adapt: call PayPal's MakePayment and extract transaction ID
	result, err := p.paypal.MakePayment(amount, currency, "paypal_account")
	if err != nil {
		return "", err
	}
	return result["transaction_id"].(string), nil
}

func (p *PayPalAdapter) RefundPayment(transactionID string, amount float64) error {
	return p.paypal.IssueRefund(transactionID, amount)
}

func (p *PayPalAdapter) GetProviderName() string {
	return "PayPal"
}

// SquareAPI simulates the Square third-party library.
type SquareAPI struct {
	locationID string
}

// CreatePayment is Square's payment method (yet another different interface).
func (s *SquareAPI) CreatePayment(amountMoney map[string]interface{}, sourceID string) (string, error) {
	amount := amountMoney["amount"].(int)
	currency := amountMoney["currency"].(string)
	txID := fmt.Sprintf("square_tx_%d", time.Now().Unix())
	fmt.Printf("[Square] Creating payment of %d %s\n", amount, currency)
	return txID, nil
}

// RefundTransaction is Square's refund method.
func (s *SquareAPI) RefundTransaction(paymentID string, amountMoney map[string]interface{}) error {
	amount := amountMoney["amount"].(int)
	currency := amountMoney["currency"].(string)
	fmt.Printf("[Square] Refunding %d %s for payment %s\n", amount, currency, paymentID)
	return nil
}

// SquareAdapter adapts SquareAPI to PaymentProcessor interface.
type SquareAdapter struct {
	square   *SquareAPI
	sourceID string
}

// NewSquareAdapter creates a new Square adapter.
func NewSquareAdapter(locationID string) PaymentProcessor {
	return &SquareAdapter{
		square:   &SquareAPI{locationID: locationID},
		sourceID: "cnon:card-nonce-ok",
	}
}

func (s *SquareAdapter) ProcessPayment(amount float64, currency string) (string, error) {
	// Adapt: create the money map Square expects
	amountMoney := map[string]interface{}{
		"amount":   int(amount * 100), // Square uses cents
		"currency": strings.ToUpper(currency),
	}
	return s.square.CreatePayment(amountMoney, s.sourceID)
}

func (s *SquareAdapter) RefundPayment(transactionID string, amount float64) error {
	amountMoney := map[string]interface{}{
		"amount":   int(amount * 100),
		"currency": "USD",
	}
	return s.square.RefundTransaction(transactionID, amountMoney)
}

func (s *SquareAdapter) GetProviderName() string {
	return "Square"
}

// Logger is our target interface for logging.
type Logger interface {
	Info(message string)
	Error(message string)
	Debug(message string)
}

// ZapLogger simulates the Zap logging library (different interface).
type ZapLogger struct {
	level string
}

func (z *ZapLogger) InfoLevel(msg string, fields ...interface{}) {
	fmt.Printf("[Zap INFO] %s %v\n", msg, fields)
}

func (z *ZapLogger) ErrorLevel(msg string, fields ...interface{}) {
	fmt.Printf("[Zap ERROR] %s %v\n", msg, fields)
}

func (z *ZapLogger) DebugLevel(msg string, fields ...interface{}) {
	fmt.Printf("[Zap DEBUG] %s %v\n", msg, fields)
}

// ZapAdapter adapts ZapLogger to our Logger interface.
type ZapAdapter struct {
	zap *ZapLogger
}

// NewZapAdapter creates a new Zap logger adapter.
func NewZapAdapter() Logger {
	return &ZapAdapter{
		zap: &ZapLogger{level: "info"},
	}
}

func (z *ZapAdapter) Info(message string) {
	z.zap.InfoLevel(message)
}

func (z *ZapAdapter) Error(message string) {
	z.zap.ErrorLevel(message)
}

func (z *ZapAdapter) Debug(message string) {
	z.zap.DebugLevel(message)
}

// LogrusLogger simulates the Logrus logging library (another different interface).
type LogrusLogger struct {
	formatter string
}

func (l *LogrusLogger) WithFields(fields map[string]interface{}) *LogrusLogger {
	return l
}

func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	fmt.Printf("[Logrus INFO] "+format+"\n", args...)
}

func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf("[Logrus ERROR] "+format+"\n", args...)
}

func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	fmt.Printf("[Logrus DEBUG] "+format+"\n", args...)
}

// LogrusAdapter adapts LogrusLogger to our Logger interface.
type LogrusAdapter struct {
	logrus *LogrusLogger
}

// NewLogrusAdapter creates a new Logrus logger adapter.
func NewLogrusAdapter() Logger {
	return &LogrusAdapter{
		logrus: &LogrusLogger{formatter: "json"},
	}
}

func (l *LogrusAdapter) Info(message string) {
	l.logrus.Infof("%s", message)
}

func (l *LogrusAdapter) Error(message string) {
	l.logrus.Errorf("%s", message)
}

func (l *LogrusAdapter) Debug(message string) {
	l.logrus.Debugf("%s", message)
}

// MockPaymentProcessor is a mock implementation for testing.
type MockPaymentProcessor struct{}

// NewMockPaymentProcessor creates a new mock payment processor.
func NewMockPaymentProcessor() PaymentProcessor {
	return &MockPaymentProcessor{}
}

func (m *MockPaymentProcessor) ProcessPayment(amount float64, currency string) (string, error) {
	return "mock_transaction_12345", nil
}

func (m *MockPaymentProcessor) RefundPayment(transactionID string, amount float64) error {
	return nil
}

func (m *MockPaymentProcessor) GetProviderName() string {
	return "Mock (for testing)"
}
