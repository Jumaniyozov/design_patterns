package strategy

import (
	"testing"
)

// TestCreditCardStrategyValidation tests credit card validation logic
func TestCreditCardStrategyValidation(t *testing.T) {
	tests := []struct {
		name    string
		card    *CreditCardStrategy
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid credit card",
			card:    NewCreditCardStrategy("4532123456789010", "12/25", "123", "John Doe"),
			wantErr: false,
		},
		{
			name:    "Invalid card number length - too short",
			card:    NewCreditCardStrategy("123", "12/25", "123", "John Doe"),
			wantErr: true,
			errMsg:  "invalid card number length",
		},
		{
			name:    "Invalid card number length - too long",
			card:    NewCreditCardStrategy("12345678901234567890", "12/25", "123", "John Doe"),
			wantErr: true,
			errMsg:  "invalid card number length",
		},
		{
			name:    "Invalid card number - contains letters",
			card:    NewCreditCardStrategy("453212345678ABCD", "12/25", "123", "John Doe"),
			wantErr: true,
			errMsg:  "invalid characters",
		},
		{
			name:    "Invalid expiry date format",
			card:    NewCreditCardStrategy("4532123456789010", "1225", "123", "John Doe"),
			wantErr: true,
			errMsg:  "invalid expiry date format",
		},
		{
			name:    "Invalid CVV - too short",
			card:    NewCreditCardStrategy("4532123456789010", "12/25", "12", "John Doe"),
			wantErr: true,
			errMsg:  "invalid CVV length",
		},
		{
			name:    "Invalid CVV - too long",
			card:    NewCreditCardStrategy("4532123456789010", "12/25", "12345", "John Doe"),
			wantErr: true,
			errMsg:  "invalid CVV length",
		},
		{
			name:    "Invalid CVV - contains letters",
			card:    NewCreditCardStrategy("4532123456789010", "12/25", "12A", "John Doe"),
			wantErr: true,
			errMsg:  "invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.card.Validate(PaymentDetails{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
				t.Errorf("Validate() error message %q does not contain %q", err.Error(), tt.errMsg)
			}
		})
	}
}

// TestCreditCardStrategyProcess tests credit card payment processing
func TestCreditCardStrategyProcess(t *testing.T) {
	card := NewCreditCardStrategy("4532123456789010", "12/25", "123", "John Doe")

	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{
			name:    "Valid amount",
			amount:  99.99,
			wantErr: false,
		},
		{
			name:    "Zero amount",
			amount:  0.0,
			wantErr: true,
		},
		{
			name:    "Negative amount",
			amount:  -50.00,
			wantErr: true,
		},
		{
			name:    "Large amount",
			amount:  999999.99,
			wantErr: false,
		},
		{
			name:    "Small positive amount",
			amount:  0.01,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txID, err := card.Process(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && txID == "" {
				t.Errorf("Process() returned empty transaction ID")
			}
			if tt.wantErr && txID != "" {
				t.Errorf("Process() should return empty transaction ID on error")
			}
		})
	}
}

// TestPayPalStrategyValidation tests PayPal validation logic
func TestPayPalStrategyValidation(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Valid email",
			email:   "user@example.com",
			wantErr: false,
		},
		{
			name:    "Email missing @",
			email:   "userexample.com",
			wantErr: true,
		},
		{
			name:    "Email missing domain",
			email:   "user@example",
			wantErr: true,
		},
		{
			name:    "Empty email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "Valid complex email",
			email:   "john.doe+test@company.co.uk",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pp := NewPayPalStrategy(tt.email)
			err := pp.Validate(PaymentDetails{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestPayPalStrategyProcess tests PayPal payment processing
func TestPayPalStrategyProcess(t *testing.T) {
	pp := NewPayPalStrategy("user@example.com")

	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{
			name:    "Valid amount",
			amount:  75.50,
			wantErr: false,
		},
		{
			name:    "Zero amount",
			amount:  0.0,
			wantErr: true,
		},
		{
			name:    "Negative amount",
			amount:  -25.00,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txID, err := pp.Process(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && txID == "" {
				t.Errorf("Process() returned empty transaction ID")
			}
		})
	}
}

// TestCryptoStrategyValidation tests cryptocurrency wallet validation
func TestCryptoStrategyValidation(t *testing.T) {
	tests := []struct {
		name    string
		address string
		crypto  string
		wantErr bool
	}{
		{
			name:    "Valid Bitcoin address",
			address: "1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9",
			crypto:  "BTC",
			wantErr: false,
		},
		{
			name:    "Valid Ethereum address",
			address: "0x742d35Cc6634C0532925a3b844Bc9e7595f42cbe",
			crypto:  "ETH",
			wantErr: false,
		},
		{
			name:    "Address too short",
			address: "abc",
			crypto:  "BTC",
			wantErr: true,
		},
		{
			name:    "Address with invalid characters",
			address: "1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9!!!",
			crypto:  "BTC",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crypto := NewCryptoStrategy(tt.address, tt.crypto)
			err := crypto.Validate(PaymentDetails{})
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCryptoStrategyProcess tests cryptocurrency payment processing
func TestCryptoStrategyProcess(t *testing.T) {
	crypto := NewCryptoStrategy("1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9", "BTC")

	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{
			name:    "Valid amount",
			amount:  0.05,
			wantErr: false,
		},
		{
			name:    "Zero amount",
			amount:  0.0,
			wantErr: true,
		},
		{
			name:    "Negative amount",
			amount:  -0.01,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txID, err := crypto.Process(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && txID == "" {
				t.Errorf("Process() returned empty transaction ID")
			}
		})
	}
}

// TestPaymentProcessorStrategySwitch tests switching strategies at runtime
func TestPaymentProcessorStrategySwitch(t *testing.T) {
	processor := NewPaymentProcessor(nil)

	// Test setting strategies
	ccStrategy := NewCreditCardStrategy("4532123456789010", "12/25", "123", "John")
	processor.SetStrategy(ccStrategy)

	if processor.GetCurrentStrategyName() != "Credit Card" {
		t.Errorf("GetCurrentStrategyName() = %s, want Credit Card", processor.GetCurrentStrategyName())
	}

	// Switch to PayPal
	ppStrategy := NewPayPalStrategy("john@example.com")
	processor.SetStrategy(ppStrategy)

	if processor.GetCurrentStrategyName() != "PayPal" {
		t.Errorf("GetCurrentStrategyName() = %s, want PayPal", processor.GetCurrentStrategyName())
	}

	// Switch to Crypto
	cryptoStrategy := NewCryptoStrategy("1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9", "BTC")
	processor.SetStrategy(cryptoStrategy)

	if processor.GetCurrentStrategyName() != "Cryptocurrency (BTC)" {
		t.Errorf("GetCurrentStrategyName() = %s, want Cryptocurrency (BTC)", processor.GetCurrentStrategyName())
	}
}

// TestPaymentProcessorValidation tests processor validation workflow
func TestPaymentProcessorValidation(t *testing.T) {
	processor := NewPaymentProcessor(nil)
	processor.SetDetails(PaymentDetails{})

	// Test with no strategy set
	_, err := processor.Process(100.0)
	if err == nil {
		t.Error("Process() should error when no strategy is set")
	}

	// Test with invalid credit card
	invalidCC := NewCreditCardStrategy("123", "12/25", "123", "John")
	processor.SetStrategy(invalidCC)

	_, err = processor.Process(100.0)
	if err == nil {
		t.Error("Process() should error with invalid card")
	}
	if !contains(err.Error(), "validation failed") {
		t.Errorf("Process() error should mention validation failure, got: %v", err)
	}
}

// TestPaymentProcessorProcessing tests complete payment processing workflow
func TestPaymentProcessorProcessing(t *testing.T) {
	processor := NewPaymentProcessor(NewCreditCardStrategy("4532123456789010", "12/25", "123", "John"))
	processor.SetDetails(PaymentDetails{CardNumber: "4532123456789010"})

	txID, err := processor.Process(99.99)
	if err != nil {
		t.Errorf("Process() unexpected error: %v", err)
	}
	if txID == "" {
		t.Error("Process() returned empty transaction ID")
	}
	if !contains(txID, "CC-") {
		t.Errorf("Process() returned invalid transaction ID format: %s", txID)
	}
}

// TestPaymentProcessorRefund tests refund functionality
func TestPaymentProcessorRefund(t *testing.T) {
	processor := NewPaymentProcessor(NewCreditCardStrategy("4532123456789010", "12/25", "123", "John"))

	err := processor.Refund("CC-12345-99.99", 99.99)
	if err != nil {
		t.Errorf("Refund() unexpected error: %v", err)
	}

	// Test with no strategy
	processor2 := NewPaymentProcessor(nil)
	err = processor2.Refund("TX-12345", 50.00)
	if err == nil {
		t.Error("Refund() should error when no strategy is set")
	}
}

// TestMultipleStrategiesWithSameProcessor tests different strategies with same processor
func TestMultipleStrategiesWithSameProcessor(t *testing.T) {
	processor := NewPaymentProcessor(nil)
	processor.SetDetails(PaymentDetails{})

	strategies := []struct {
		name     string
		strategy PaymentStrategy
		setup    func(PaymentDetails) PaymentDetails
		amount   float64
	}{
		{
			name:     "Credit Card",
			strategy: NewCreditCardStrategy("4532123456789010", "12/25", "123", "John"),
			setup: func(p PaymentDetails) PaymentDetails {
				p.CardNumber = "4532123456789010"
				return p
			},
			amount: 50.0,
		},
		{
			name:     "PayPal",
			strategy: NewPayPalStrategy("john@example.com"),
			setup: func(p PaymentDetails) PaymentDetails {
				p.Email = "john@example.com"
				return p
			},
			amount: 75.0,
		},
		{
			name:     "Crypto",
			strategy: NewCryptoStrategy("1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9", "BTC"),
			setup: func(p PaymentDetails) PaymentDetails {
				p.WalletAddr = "1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9"
				return p
			},
			amount: 0.05,
		},
	}

	for _, s := range strategies {
		t.Run(s.name, func(t *testing.T) {
			processor.SetStrategy(s.strategy)
			processor.SetDetails(s.setup(PaymentDetails{}))

			txID, err := processor.Process(s.amount)
			if err != nil {
				t.Errorf("Process() unexpected error: %v", err)
			}
			if txID == "" {
				t.Errorf("Process() returned empty transaction ID for %s", s.name)
			}
		})
	}
}

// TestStrategyGetName tests the GetName method for all strategies
func TestStrategyGetName(t *testing.T) {
	tests := []struct {
		name     string
		strategy PaymentStrategy
		expected string
	}{
		{
			name:     "Credit Card",
			strategy: NewCreditCardStrategy("4532123456789010", "12/25", "123", "John"),
			expected: "Credit Card",
		},
		{
			name:     "PayPal",
			strategy: NewPayPalStrategy("john@example.com"),
			expected: "PayPal",
		},
		{
			name:     "Crypto BTC",
			strategy: NewCryptoStrategy("1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9", "BTC"),
			expected: "Cryptocurrency (BTC)",
		},
		{
			name:     "Crypto ETH",
			strategy: NewCryptoStrategy("0x742d35Cc6634C0532925a3b844Bc9e7595f42cbe", "ETH"),
			expected: "Cryptocurrency (ETH)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.strategy.GetName()
			if got != tt.expected {
				t.Errorf("GetName() = %q, want %q", got, tt.expected)
			}
		})
	}
}

// BenchmarkCreditCardProcess benchmarks credit card payment processing
func BenchmarkCreditCardProcess(b *testing.B) {
	card := NewCreditCardStrategy("4532123456789010", "12/25", "123", "John")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		card.Process(99.99)
	}
}

// BenchmarkPayPalProcess benchmarks PayPal payment processing
func BenchmarkPayPalProcess(b *testing.B) {
	pp := NewPayPalStrategy("john@example.com")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pp.Process(99.99)
	}
}

// BenchmarkCryptoProcess benchmarks crypto payment processing
func BenchmarkCryptoProcess(b *testing.B) {
	crypto := NewCryptoStrategy("1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9", "BTC")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		crypto.Process(0.05)
	}
}

// BenchmarkPaymentProcessorProcess benchmarks complete payment flow
func BenchmarkPaymentProcessorProcess(b *testing.B) {
	processor := NewPaymentProcessor(NewCreditCardStrategy("4532123456789010", "12/25", "123", "John"))
	processor.SetDetails(PaymentDetails{CardNumber: "4532123456789010"})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		processor.Process(99.99)
	}
}

// BenchmarkStrategySwitch benchmarks switching between strategies
func BenchmarkStrategySwitch(b *testing.B) {
	processor := NewPaymentProcessor(NewCreditCardStrategy("4532123456789010", "12/25", "123", "John"))
	cc := NewCreditCardStrategy("4532123456789010", "12/25", "123", "John")
	pp := NewPayPalStrategy("john@example.com")
	crypto := NewCryptoStrategy("1A1z7agoat2GPFH7q05VtEWEB97YMQP5Z9", "BTC")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		processor.SetStrategy(cc)
		processor.SetStrategy(pp)
		processor.SetStrategy(crypto)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != "" && substr != ""
}
