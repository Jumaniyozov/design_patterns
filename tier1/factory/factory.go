// Package factory demonstrates the Factory Pattern - a creational design pattern
// that provides an interface for creating objects without specifying their exact classes.
//
// The Factory Pattern encapsulates object creation logic, making code more flexible
// and maintainable. In Go, this is typically implemented through constructor functions
// (functions starting with "New").
package factory

import (
	"fmt"
	"strings"
)

// PaymentMethod represents the type of payment method to use
type PaymentMethod string

const (
	CreditCard PaymentMethod = "creditcard"
	PayPal     PaymentMethod = "paypal"
	Bitcoin    PaymentMethod = "bitcoin"
)

// PaymentProcessor is the common interface that all payment processors implement.
// This abstraction allows clients to work with any payment processor without
// knowing the concrete implementation details.
type PaymentProcessor interface {
	ProcessPayment(amount float64) (string, error)
	ValidateAccount() error
	GetProcessorName() string
}

// CreditCardProcessor handles credit card payments
type CreditCardProcessor struct {
	cardNumber string
	cvv        string
	expiryDate string
}

// ProcessPayment processes a credit card payment
func (c *CreditCardProcessor) ProcessPayment(amount float64) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("invalid amount: %.2f", amount)
	}

	// Simulate payment processing
	transactionID := fmt.Sprintf("CC-%s-%.2f", c.cardNumber[len(c.cardNumber)-4:], amount)
	return transactionID, nil
}

// ValidateAccount validates the credit card details
func (c *CreditCardProcessor) ValidateAccount() error {
	if len(c.cardNumber) < 13 || len(c.cardNumber) > 19 {
		return fmt.Errorf("invalid card number length")
	}
	if len(c.cvv) != 3 && len(c.cvv) != 4 {
		return fmt.Errorf("invalid CVV length")
	}
	return nil
}

// GetProcessorName returns the name of this payment processor
func (c *CreditCardProcessor) GetProcessorName() string {
	return "Credit Card Processor"
}

// PayPalProcessor handles PayPal payments
type PayPalProcessor struct {
	email    string
	apiToken string
}

// ProcessPayment processes a PayPal payment
func (p *PayPalProcessor) ProcessPayment(amount float64) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("invalid amount: %.2f", amount)
	}

	// Simulate payment processing
	transactionID := fmt.Sprintf("PP-%s-%.2f", p.email, amount)
	return transactionID, nil
}

// ValidateAccount validates the PayPal account
func (p *PayPalProcessor) ValidateAccount() error {
	if !strings.Contains(p.email, "@") {
		return fmt.Errorf("invalid email format")
	}
	if len(p.apiToken) == 0 {
		return fmt.Errorf("API token required")
	}
	return nil
}

// GetProcessorName returns the name of this payment processor
func (p *PayPalProcessor) GetProcessorName() string {
	return "PayPal Processor"
}

// BitcoinProcessor handles Bitcoin payments
type BitcoinProcessor struct {
	walletAddress string
	network       string // mainnet, testnet
}

// ProcessPayment processes a Bitcoin payment
func (b *BitcoinProcessor) ProcessPayment(amount float64) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("invalid amount: %.2f", amount)
	}

	// Simulate payment processing
	transactionID := fmt.Sprintf("BTC-%s-%s-%.8f", b.network, b.walletAddress[:8], amount)
	return transactionID, nil
}

// ValidateAccount validates the Bitcoin wallet
func (b *BitcoinProcessor) ValidateAccount() error {
	if len(b.walletAddress) < 26 || len(b.walletAddress) > 35 {
		return fmt.Errorf("invalid wallet address length")
	}
	if b.network != "mainnet" && b.network != "testnet" {
		return fmt.Errorf("invalid network: must be 'mainnet' or 'testnet'")
	}
	return nil
}

// GetProcessorName returns the name of this payment processor
func (b *BitcoinProcessor) GetProcessorName() string {
	return fmt.Sprintf("Bitcoin Processor (%s)", b.network)
}

// PaymentConfig holds the configuration for creating a payment processor
type PaymentConfig struct {
	Method PaymentMethod

	// Credit Card fields
	CardNumber string
	CVV        string
	ExpiryDate string

	// PayPal fields
	Email    string
	APIToken string

	// Bitcoin fields
	WalletAddress string
	Network       string
}

// NewPaymentProcessor is the factory function that creates and returns the appropriate
// PaymentProcessor based on the payment method specified in the configuration.
//
// This is the core of the Factory Pattern - it encapsulates the creation logic
// and returns an interface type, allowing clients to work with any payment processor
// without knowing the concrete implementation.
func NewPaymentProcessor(config PaymentConfig) (PaymentProcessor, error) {
	var processor PaymentProcessor

	switch config.Method {
	case CreditCard:
		processor = &CreditCardProcessor{
			cardNumber: config.CardNumber,
			cvv:        config.CVV,
			expiryDate: config.ExpiryDate,
		}

	case PayPal:
		processor = &PayPalProcessor{
			email:    config.Email,
			apiToken: config.APIToken,
		}

	case Bitcoin:
		network := config.Network
		if network == "" {
			network = "mainnet" // default to mainnet
		}
		processor = &BitcoinProcessor{
			walletAddress: config.WalletAddress,
			network:       network,
		}

	default:
		return nil, fmt.Errorf("unsupported payment method: %s", config.Method)
	}

	// Validate the created processor before returning
	if err := processor.ValidateAccount(); err != nil {
		return nil, fmt.Errorf("validation failed for %s: %w", processor.GetProcessorName(), err)
	}

	return processor, nil
}

// NewCreditCardProcessor is a specialized factory function for creating credit card processors.
// This demonstrates how you can provide both a general factory (NewPaymentProcessor)
// and specialized factories for specific types.
func NewCreditCardProcessor(cardNumber, cvv, expiryDate string) (PaymentProcessor, error) {
	return NewPaymentProcessor(PaymentConfig{
		Method:     CreditCard,
		CardNumber: cardNumber,
		CVV:        cvv,
		ExpiryDate: expiryDate,
	})
}

// NewPayPalProcessor is a specialized factory function for creating PayPal processors
func NewPayPalProcessor(email, apiToken string) (PaymentProcessor, error) {
	return NewPaymentProcessor(PaymentConfig{
		Method:   PayPal,
		Email:    email,
		APIToken: apiToken,
	})
}

// NewBitcoinProcessor is a specialized factory function for creating Bitcoin processors
func NewBitcoinProcessor(walletAddress, network string) (PaymentProcessor, error) {
	return NewPaymentProcessor(PaymentConfig{
		Method:        Bitcoin,
		WalletAddress: walletAddress,
		Network:       network,
	})
}
