// Package strategy implements the Strategy design pattern.
//
// The Strategy pattern encapsulates a family of algorithms, making them interchangeable.
// It lets the algorithm vary independently from clients that use it.
//
// Key components:
// - Strategy interface: Defines the contract for all concrete strategies
// - Concrete strategies: Different implementations of the strategy interface
// - Context: Uses the strategy to accomplish its goal
//
// This package demonstrates the pattern with a payment processing system where
// different payment methods are interchangeable strategies.
package strategy

import (
	"fmt"
	"regexp"
	"strings"
)

// PaymentStrategy defines the interface for different payment methods.
// Any type implementing these methods can be used as a payment strategy.
type PaymentStrategy interface {
	// Validate checks if the payment details are valid for this strategy.
	// It should return an error if validation fails.
	Validate(details PaymentDetails) error

	// Process executes the payment using this strategy.
	// Returns a transaction ID on success or an error on failure.
	Process(amount float64) (transactionID string, err error)

	// Refund reverses a previous payment.
	Refund(transactionID string, amount float64) error

	// GetName returns a human-readable name for this strategy.
	GetName() string
}

// PaymentDetails holds information needed for payment processing.
// Different strategies may use different fields.
type PaymentDetails struct {
	CardNumber string // For credit/debit cards
	ExpiryDate string // MM/YY format
	CVV        string // 3-4 digit security code
	Email      string // For PayPal
	WalletAddr string // For cryptocurrency
	AccountNum string // For bank transfers
}

// ============================================================================
// Concrete Strategy 1: Credit Card Payment
// ============================================================================

// CreditCardStrategy implements PaymentStrategy for credit card payments.
type CreditCardStrategy struct {
	cardNumber string
	expiryDate string
	cvv        string
	cardHolder string
}

// NewCreditCardStrategy creates a new credit card payment strategy.
func NewCreditCardStrategy(cardNumber, expiryDate, cvv, cardHolder string) *CreditCardStrategy {
	return &CreditCardStrategy{
		cardNumber: cardNumber,
		expiryDate: expiryDate,
		cvv:        cvv,
		cardHolder: cardHolder,
	}
}

// Validate checks if the credit card details are valid.
func (c *CreditCardStrategy) Validate(details PaymentDetails) error {
	// Validate card number (Luhn algorithm simplified)
	if len(c.cardNumber) < 13 || len(c.cardNumber) > 19 {
		return fmt.Errorf("invalid card number length: %d", len(c.cardNumber))
	}

	// Validate card only contains digits
	if !regexp.MustCompile(`^\d+$`).MatchString(c.cardNumber) {
		return fmt.Errorf("card number contains invalid characters")
	}

	// Validate expiry date format (MM/YY)
	parts := strings.Split(c.expiryDate, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid expiry date format, expected MM/YY")
	}

	// Validate CVV
	if len(c.cvv) < 3 || len(c.cvv) > 4 {
		return fmt.Errorf("invalid CVV length: %d", len(c.cvv))
	}

	if !regexp.MustCompile(`^\d+$`).MatchString(c.cvv) {
		return fmt.Errorf("CVV contains invalid characters")
	}

	return nil
}

// Process executes a credit card payment.
func (c *CreditCardStrategy) Process(amount float64) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("amount must be positive")
	}

	// In a real system, this would call a payment gateway API
	// For demonstration, we simulate success
	transactionID := fmt.Sprintf("CC-%d-%.2f", hashCard(c.cardNumber), amount)
	fmt.Printf("Processing credit card payment for %s\n", c.cardHolder)
	return transactionID, nil
}

// Refund reverses a credit card payment.
func (c *CreditCardStrategy) Refund(transactionID string, amount float64) error {
	fmt.Printf("Refunding %.2f to credit card ending in %s\n", amount, c.cardNumber[len(c.cardNumber)-4:])
	return nil
}

// GetName returns the strategy name.
func (c *CreditCardStrategy) GetName() string {
	return "Credit Card"
}

// ============================================================================
// Concrete Strategy 2: PayPal Payment
// ============================================================================

// PayPalStrategy implements PaymentStrategy for PayPal payments.
type PayPalStrategy struct {
	email string
}

// NewPayPalStrategy creates a new PayPal payment strategy.
func NewPayPalStrategy(email string) *PayPalStrategy {
	return &PayPalStrategy{email: email}
}

// Validate checks if the PayPal account is valid.
func (p *PayPalStrategy) Validate(details PaymentDetails) error {
	// Simple email validation
	if !strings.Contains(p.email, "@") || !strings.Contains(p.email, ".") {
		return fmt.Errorf("invalid PayPal email: %s", p.email)
	}
	return nil
}

// Process executes a PayPal payment.
// In a real system, this would redirect to PayPal for authentication.
func (p *PayPalStrategy) Process(amount float64) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("amount must be positive")
	}

	// Simulate PayPal OAuth flow and payment
	fmt.Printf("Redirecting to PayPal for user %s\n", p.email)
	transactionID := fmt.Sprintf("PP-%d-%.2f", hashEmail(p.email), amount)
	fmt.Printf("PayPal payment confirmed via OAuth\n")
	return transactionID, nil
}

// Refund reverses a PayPal payment.
func (p *PayPalStrategy) Refund(transactionID string, amount float64) error {
	fmt.Printf("Issuing refund of %.2f to PayPal account %s\n", amount, p.email)
	return nil
}

// GetName returns the strategy name.
func (p *PayPalStrategy) GetName() string {
	return "PayPal"
}

// ============================================================================
// Concrete Strategy 3: Cryptocurrency Payment
// ============================================================================

// CryptoStrategy implements PaymentStrategy for cryptocurrency payments.
type CryptoStrategy struct {
	walletAddress string
	cryptoType    string // "BTC", "ETH", etc.
}

// NewCryptoStrategy creates a new cryptocurrency payment strategy.
func NewCryptoStrategy(walletAddress, cryptoType string) *CryptoStrategy {
	return &CryptoStrategy{
		walletAddress: walletAddress,
		cryptoType:    cryptoType,
	}
}

// Validate checks if the wallet address is valid.
func (c *CryptoStrategy) Validate(details PaymentDetails) error {
	if len(c.walletAddress) < 26 {
		return fmt.Errorf("invalid wallet address length")
	}

	// Simple check - real validation would use actual checksums
	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(c.walletAddress) {
		return fmt.Errorf("wallet address contains invalid characters")
	}

	return nil
}

// Process executes a cryptocurrency payment.
// This would involve blockchain monitoring in a real system.
func (c *CryptoStrategy) Process(amount float64) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("amount must be positive")
	}

	fmt.Printf("Generating blockchain transaction for %s wallet\n", c.cryptoType)
	fmt.Printf("Monitoring blockchain for confirmation...\n")
	transactionID := fmt.Sprintf("CRYPTO-%s-%d", c.cryptoType, hashAddress(c.walletAddress))
	fmt.Printf("Blockchain transaction confirmed\n")
	return transactionID, nil
}

// Refund reverses a cryptocurrency payment.
// Note: In reality, crypto transactions are typically irreversible.
func (c *CryptoStrategy) Refund(transactionID string, amount float64) error {
	fmt.Printf("Initiating reverse transaction to %s wallet\n", c.walletAddress)
	return nil
}

// GetName returns the strategy name.
func (c *CryptoStrategy) GetName() string {
	return fmt.Sprintf("Cryptocurrency (%s)", c.cryptoType)
}

// ============================================================================
// Context: PaymentProcessor
// ============================================================================

// PaymentProcessor is the context that uses a payment strategy.
// It encapsulates the workflow of payment processing while delegating
// the specific payment logic to the chosen strategy.
type PaymentProcessor struct {
	strategy PaymentStrategy
	details  PaymentDetails
}

// NewPaymentProcessor creates a new payment processor with the given strategy.
func NewPaymentProcessor(strategy PaymentStrategy) *PaymentProcessor {
	return &PaymentProcessor{
		strategy: strategy,
	}
}

// SetStrategy allows changing the payment strategy at runtime.
// This demonstrates the flexibility of the Strategy pattern.
func (p *PaymentProcessor) SetStrategy(strategy PaymentStrategy) {
	p.strategy = strategy
}

// SetDetails sets the payment details for the processor.
func (p *PaymentProcessor) SetDetails(details PaymentDetails) {
	p.details = details
}

// Process executes the payment using the current strategy.
// It handles the common workflow: validation, processing, and error handling.
func (p *PaymentProcessor) Process(amount float64) (string, error) {
	if p.strategy == nil {
		return "", fmt.Errorf("no payment strategy set")
	}

	// Step 1: Validate the payment method
	if err := p.strategy.Validate(p.details); err != nil {
		return "", fmt.Errorf("validation failed: %w", err)
	}

	// Step 2: Process the payment using the selected strategy
	transactionID, err := p.strategy.Process(amount)
	if err != nil {
		return "", fmt.Errorf("payment processing failed: %w", err)
	}

	// Step 3: Return transaction ID for future reference
	return transactionID, nil
}

// Refund refunds a previous payment using the current strategy.
func (p *PaymentProcessor) Refund(transactionID string, amount float64) error {
	if p.strategy == nil {
		return fmt.Errorf("no payment strategy set")
	}

	return p.strategy.Refund(transactionID, amount)
}

// GetCurrentStrategyName returns the name of the current payment strategy.
func (p *PaymentProcessor) GetCurrentStrategyName() string {
	if p.strategy == nil {
		return "None"
	}
	return p.strategy.GetName()
}

// ============================================================================
// Helper Functions
// ============================================================================

// hashCard generates a simple hash of the card number for transaction IDs.
func hashCard(cardNumber string) int {
	hash := 0
	for _, ch := range cardNumber {
		hash = ((hash << 5) - hash) + int(ch)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash % 1000000
}

// hashEmail generates a simple hash of an email for transaction IDs.
func hashEmail(email string) int {
	hash := 0
	for _, ch := range email {
		hash = ((hash << 5) - hash) + int(ch)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash % 1000000
}

// hashAddress generates a simple hash of a wallet address for transaction IDs.
func hashAddress(address string) int {
	hash := 0
	for _, ch := range address {
		hash = ((hash << 5) - hash) + int(ch)
	}
	if hash < 0 {
		hash = -hash
	}
	return hash % 1000000
}
