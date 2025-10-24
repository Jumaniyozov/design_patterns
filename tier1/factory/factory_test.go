package factory

import (
	"strings"
	"testing"
)

// TestNewPaymentProcessor_CreditCard tests the factory creates credit card processors correctly
func TestNewPaymentProcessor_CreditCard(t *testing.T) {
	config := PaymentConfig{
		Method:     CreditCard,
		CardNumber: "4532015112830366",
		CVV:        "123",
		ExpiryDate: "12/25",
	}

	processor, err := NewPaymentProcessor(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if processor == nil {
		t.Fatal("Expected processor, got nil")
	}

	name := processor.GetProcessorName()
	if name != "Credit Card Processor" {
		t.Errorf("Expected 'Credit Card Processor', got: %s", name)
	}

	// Verify it's the correct concrete type
	if _, ok := processor.(*CreditCardProcessor); !ok {
		t.Errorf("Expected *CreditCardProcessor, got: %T", processor)
	}
}

// TestNewPaymentProcessor_PayPal tests the factory creates PayPal processors correctly
func TestNewPaymentProcessor_PayPal(t *testing.T) {
	config := PaymentConfig{
		Method:   PayPal,
		Email:    "user@example.com",
		APIToken: "token123",
	}

	processor, err := NewPaymentProcessor(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if processor == nil {
		t.Fatal("Expected processor, got nil")
	}

	name := processor.GetProcessorName()
	if name != "PayPal Processor" {
		t.Errorf("Expected 'PayPal Processor', got: %s", name)
	}

	// Verify it's the correct concrete type
	if _, ok := processor.(*PayPalProcessor); !ok {
		t.Errorf("Expected *PayPalProcessor, got: %T", processor)
	}
}

// TestNewPaymentProcessor_Bitcoin tests the factory creates Bitcoin processors correctly
func TestNewPaymentProcessor_Bitcoin(t *testing.T) {
	config := PaymentConfig{
		Method:        Bitcoin,
		WalletAddress: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		Network:       "mainnet",
	}

	processor, err := NewPaymentProcessor(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if processor == nil {
		t.Fatal("Expected processor, got nil")
	}

	name := processor.GetProcessorName()
	expectedName := "Bitcoin Processor (mainnet)"
	if name != expectedName {
		t.Errorf("Expected '%s', got: %s", expectedName, name)
	}

	// Verify it's the correct concrete type
	if _, ok := processor.(*BitcoinProcessor); !ok {
		t.Errorf("Expected *BitcoinProcessor, got: %T", processor)
	}
}

// TestNewPaymentProcessor_Bitcoin_DefaultNetwork tests that Bitcoin defaults to mainnet
func TestNewPaymentProcessor_Bitcoin_DefaultNetwork(t *testing.T) {
	config := PaymentConfig{
		Method:        Bitcoin,
		WalletAddress: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
		// Network not specified - should default to mainnet
	}

	processor, err := NewPaymentProcessor(config)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	btcProcessor, ok := processor.(*BitcoinProcessor)
	if !ok {
		t.Fatalf("Expected *BitcoinProcessor, got: %T", processor)
	}

	if btcProcessor.network != "mainnet" {
		t.Errorf("Expected default network 'mainnet', got: %s", btcProcessor.network)
	}
}

// TestNewPaymentProcessor_UnsupportedMethod tests error handling for unsupported payment methods
func TestNewPaymentProcessor_UnsupportedMethod(t *testing.T) {
	config := PaymentConfig{
		Method: "invalid",
	}

	processor, err := NewPaymentProcessor(config)
	if err == nil {
		t.Fatal("Expected error for unsupported payment method, got nil")
	}

	if processor != nil {
		t.Errorf("Expected nil processor, got: %v", processor)
	}

	expectedErrMsg := "unsupported payment method"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("Expected error containing '%s', got: %v", expectedErrMsg, err)
	}
}

// TestNewPaymentProcessor_ValidationErrors tests that the factory validates input
func TestNewPaymentProcessor_ValidationErrors(t *testing.T) {
	tests := []struct {
		name        string
		config      PaymentConfig
		expectedErr string
	}{
		{
			name: "Invalid credit card number",
			config: PaymentConfig{
				Method:     CreditCard,
				CardNumber: "123", // too short
				CVV:        "123",
				ExpiryDate: "12/25",
			},
			expectedErr: "invalid card number length",
		},
		{
			name: "Invalid CVV",
			config: PaymentConfig{
				Method:     CreditCard,
				CardNumber: "4532015112830366",
				CVV:        "12", // too short
				ExpiryDate: "12/25",
			},
			expectedErr: "invalid CVV length",
		},
		{
			name: "Invalid PayPal email",
			config: PaymentConfig{
				Method:   PayPal,
				Email:    "not-an-email", // missing @
				APIToken: "token123",
			},
			expectedErr: "invalid email format",
		},
		{
			name: "Missing PayPal API token",
			config: PaymentConfig{
				Method:   PayPal,
				Email:    "user@example.com",
				APIToken: "", // empty token
			},
			expectedErr: "API token required",
		},
		{
			name: "Invalid Bitcoin network",
			config: PaymentConfig{
				Method:        Bitcoin,
				WalletAddress: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				Network:       "invalid",
			},
			expectedErr: "invalid network",
		},
		{
			name: "Invalid Bitcoin wallet address",
			config: PaymentConfig{
				Method:        Bitcoin,
				WalletAddress: "short", // too short
				Network:       "mainnet",
			},
			expectedErr: "invalid wallet address length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, err := NewPaymentProcessor(tt.config)

			if err == nil {
				t.Fatalf("Expected error, got nil (processor: %v)", processor)
			}

			if processor != nil {
				t.Errorf("Expected nil processor on validation error, got: %v", processor)
			}

			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("Expected error containing '%s', got: %v", tt.expectedErr, err)
			}
		})
	}
}

// TestSpecializedFactories tests the specialized factory functions
func TestSpecializedFactories(t *testing.T) {
	t.Run("NewCreditCardProcessor", func(t *testing.T) {
		processor, err := NewCreditCardProcessor("4532015112830366", "123", "12/25")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if _, ok := processor.(*CreditCardProcessor); !ok {
			t.Errorf("Expected *CreditCardProcessor, got: %T", processor)
		}
	})

	t.Run("NewPayPalProcessor", func(t *testing.T) {
		processor, err := NewPayPalProcessor("user@example.com", "token123")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if _, ok := processor.(*PayPalProcessor); !ok {
			t.Errorf("Expected *PayPalProcessor, got: %T", processor)
		}
	})

	t.Run("NewBitcoinProcessor", func(t *testing.T) {
		processor, err := NewBitcoinProcessor("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "testnet")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if _, ok := processor.(*BitcoinProcessor); !ok {
			t.Errorf("Expected *BitcoinProcessor, got: %T", processor)
		}
	})
}

// TestProcessPayment tests payment processing for all processor types
func TestProcessPayment(t *testing.T) {
	tests := []struct {
		name   string
		config PaymentConfig
		amount float64
		valid  bool
	}{
		{
			name: "Valid credit card payment",
			config: PaymentConfig{
				Method:     CreditCard,
				CardNumber: "4532015112830366",
				CVV:        "123",
				ExpiryDate: "12/25",
			},
			amount: 99.99,
			valid:  true,
		},
		{
			name: "Valid PayPal payment",
			config: PaymentConfig{
				Method:   PayPal,
				Email:    "user@example.com",
				APIToken: "token123",
			},
			amount: 150.00,
			valid:  true,
		},
		{
			name: "Valid Bitcoin payment",
			config: PaymentConfig{
				Method:        Bitcoin,
				WalletAddress: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				Network:       "mainnet",
			},
			amount: 0.01,
			valid:  true,
		},
		{
			name: "Invalid amount (zero)",
			config: PaymentConfig{
				Method:     CreditCard,
				CardNumber: "4532015112830366",
				CVV:        "123",
				ExpiryDate: "12/25",
			},
			amount: 0,
			valid:  false,
		},
		{
			name: "Invalid amount (negative)",
			config: PaymentConfig{
				Method:     CreditCard,
				CardNumber: "4532015112830366",
				CVV:        "123",
				ExpiryDate: "12/25",
			},
			amount: -50.00,
			valid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, err := NewPaymentProcessor(tt.config)
			if err != nil {
				t.Fatalf("Failed to create processor: %v", err)
			}

			transactionID, err := processor.ProcessPayment(tt.amount)

			if tt.valid {
				if err != nil {
					t.Errorf("Expected successful payment, got error: %v", err)
				}
				if transactionID == "" {
					t.Error("Expected transaction ID, got empty string")
				}
			} else {
				if err == nil {
					t.Errorf("Expected payment error, got success with ID: %s", transactionID)
				}
				if transactionID != "" {
					t.Errorf("Expected empty transaction ID on error, got: %s", transactionID)
				}
			}
		})
	}
}

// TestPaymentProcessor_InterfaceCompliance tests that all processors implement the interface correctly
func TestPaymentProcessor_InterfaceCompliance(t *testing.T) {
	processors := []struct {
		name      string
		processor PaymentProcessor
	}{
		{
			name: "CreditCardProcessor",
			processor: &CreditCardProcessor{
				cardNumber: "4532015112830366",
				cvv:        "123",
				expiryDate: "12/25",
			},
		},
		{
			name: "PayPalProcessor",
			processor: &PayPalProcessor{
				email:    "user@example.com",
				apiToken: "token123",
			},
		},
		{
			name: "BitcoinProcessor",
			processor: &BitcoinProcessor{
				walletAddress: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				network:       "mainnet",
			},
		},
	}

	for _, tt := range processors {
		t.Run(tt.name, func(t *testing.T) {
			// Test ValidateAccount
			if err := tt.processor.ValidateAccount(); err != nil {
				t.Errorf("ValidateAccount() error = %v", err)
			}

			// Test GetProcessorName
			name := tt.processor.GetProcessorName()
			if name == "" {
				t.Error("GetProcessorName() returned empty string")
			}

			// Test ProcessPayment
			transactionID, err := tt.processor.ProcessPayment(100.00)
			if err != nil {
				t.Errorf("ProcessPayment() error = %v", err)
			}
			if transactionID == "" {
				t.Error("ProcessPayment() returned empty transaction ID")
			}
		})
	}
}

// TestValidateAccount tests account validation for all processor types
func TestValidateAccount(t *testing.T) {
	tests := []struct {
		name      string
		processor PaymentProcessor
		shouldErr bool
	}{
		{
			name: "Valid credit card",
			processor: &CreditCardProcessor{
				cardNumber: "4532015112830366",
				cvv:        "123",
				expiryDate: "12/25",
			},
			shouldErr: false,
		},
		{
			name: "Invalid credit card - short number",
			processor: &CreditCardProcessor{
				cardNumber: "123",
				cvv:        "123",
				expiryDate: "12/25",
			},
			shouldErr: true,
		},
		{
			name: "Invalid credit card - bad CVV",
			processor: &CreditCardProcessor{
				cardNumber: "4532015112830366",
				cvv:        "12",
				expiryDate: "12/25",
			},
			shouldErr: true,
		},
		{
			name: "Valid PayPal",
			processor: &PayPalProcessor{
				email:    "user@example.com",
				apiToken: "token123",
			},
			shouldErr: false,
		},
		{
			name: "Invalid PayPal - bad email",
			processor: &PayPalProcessor{
				email:    "not-an-email",
				apiToken: "token123",
			},
			shouldErr: true,
		},
		{
			name: "Valid Bitcoin",
			processor: &BitcoinProcessor{
				walletAddress: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				network:       "mainnet",
			},
			shouldErr: false,
		},
		{
			name: "Invalid Bitcoin - bad network",
			processor: &BitcoinProcessor{
				walletAddress: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
				network:       "invalid",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.processor.ValidateAccount()

			if tt.shouldErr && err == nil {
				t.Error("Expected validation error, got nil")
			}

			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no validation error, got: %v", err)
			}
		})
	}
}

// BenchmarkNewPaymentProcessor benchmarks the factory creation performance
func BenchmarkNewPaymentProcessor(b *testing.B) {
	config := PaymentConfig{
		Method:     CreditCard,
		CardNumber: "4532015112830366",
		CVV:        "123",
		ExpiryDate: "12/25",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewPaymentProcessor(config)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// BenchmarkProcessPayment benchmarks payment processing performance
func BenchmarkProcessPayment(b *testing.B) {
	processor, err := NewCreditCardProcessor("4532015112830366", "123", "12/25")
	if err != nil {
		b.Fatalf("Failed to create processor: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := processor.ProcessPayment(99.99)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// BenchmarkFactoryVsDirectCreation compares factory creation vs direct instantiation
func BenchmarkFactoryVsDirectCreation(b *testing.B) {
	b.Run("Factory", func(b *testing.B) {
		config := PaymentConfig{
			Method:     CreditCard,
			CardNumber: "4532015112830366",
			CVV:        "123",
			ExpiryDate: "12/25",
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = NewPaymentProcessor(config)
		}
	})

	b.Run("Direct", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = &CreditCardProcessor{
				cardNumber: "4532015112830366",
				cvv:        "123",
				expiryDate: "12/25",
			}
		}
	})
}
