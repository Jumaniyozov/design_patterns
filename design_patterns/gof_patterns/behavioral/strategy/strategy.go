// Package strategy demonstrates the Strategy pattern in Go.
// The Strategy pattern defines a family of algorithms, encapsulates each one,
// and makes them interchangeable.
package strategy

import (
	"fmt"
	"math"
	"strings"
)

// ShippingStrategy defines the interface for calculating shipping costs.
type ShippingStrategy interface {
	Calculate(order Order) float64
	GetName() string
}

// Order represents a customer order.
type Order struct {
	Items       []string
	TotalPrice  float64
	Weight      float64 // in kg
	Destination string
}

// StandardShipping provides flat-rate standard shipping.
type StandardShipping struct {
	flatRate float64
}

func NewStandardShipping() ShippingStrategy {
	return &StandardShipping{flatRate: 5.00}
}

func (s *StandardShipping) Calculate(order Order) float64 {
	fmt.Printf("  [Standard Shipping] Flat rate: $%.2f\n", s.flatRate)
	return s.flatRate
}

func (s *StandardShipping) GetName() string {
	return "Standard Shipping"
}

// ExpressShipping provides faster delivery at higher cost.
type ExpressShipping struct {
	flatRate float64
}

func NewExpressShipping() ShippingStrategy {
	return &ExpressShipping{flatRate: 15.00}
}

func (e *ExpressShipping) Calculate(order Order) float64 {
	fmt.Printf("  [Express Shipping] Flat rate: $%.2f\n", e.flatRate)
	return e.flatRate
}

func (e *ExpressShipping) GetName() string {
	return "Express Shipping"
}

// InternationalShipping calculates based on weight and destination.
type InternationalShipping struct {
	ratePerKg      float64
	baseFee        float64
	zoneMultiplier map[string]float64
}

func NewInternationalShipping() ShippingStrategy {
	return &InternationalShipping{
		ratePerKg: 10.00,
		baseFee:   20.00,
		zoneMultiplier: map[string]float64{
			"Europe":  1.0,
			"Asia":    1.5,
			"Africa":  2.0,
			"Oceania": 1.8,
		},
	}
}

func (i *InternationalShipping) Calculate(order Order) float64 {
	multiplier := i.zoneMultiplier[order.Destination]
	if multiplier == 0 {
		multiplier = 1.0 // default
	}

	cost := i.baseFee + (order.Weight * i.ratePerKg * multiplier)
	fmt.Printf("  [International] Base: $%.2f + Weight(%.2fkg) * Rate($%.2f) * Zone(%.1fx) = $%.2f\n",
		i.baseFee, order.Weight, i.ratePerKg, multiplier, cost)
	return cost
}

func (i *InternationalShipping) GetName() string {
	return "International Shipping"
}

// FreeShipping applies when order meets minimum threshold.
type FreeShipping struct {
	minimumOrder float64
}

func NewFreeShipping(minimum float64) ShippingStrategy {
	return &FreeShipping{minimumOrder: minimum}
}

func (f *FreeShipping) Calculate(order Order) float64 {
	if order.TotalPrice >= f.minimumOrder {
		fmt.Printf("  [Free Shipping] Order $%.2f >= $%.2f threshold\n",
			order.TotalPrice, f.minimumOrder)
		return 0.00
	}

	// Fallback to standard shipping
	fmt.Printf("  [Free Shipping] Order $%.2f < $%.2f threshold, using standard rate\n",
		order.TotalPrice, f.minimumOrder)
	return 5.00
}

func (f *FreeShipping) GetName() string {
	return "Free Shipping"
}

// ShippingCalculator is the context that uses shipping strategies.
type ShippingCalculator struct {
	strategy ShippingStrategy
}

// NewShippingCalculator creates a calculator with a strategy.
func NewShippingCalculator(strategy ShippingStrategy) *ShippingCalculator {
	return &ShippingCalculator{strategy: strategy}
}

// SetStrategy allows changing the strategy at runtime.
func (sc *ShippingCalculator) SetStrategy(strategy ShippingStrategy) {
	sc.strategy = strategy
}

// Calculate executes the current strategy.
func (sc *ShippingCalculator) Calculate(order Order) float64 {
	if sc.strategy == nil {
		fmt.Println("  [Warning] No strategy set, using standard")
		sc.strategy = NewStandardShipping()
	}
	return sc.strategy.Calculate(order)
}

// CompressionStrategy defines the interface for data compression algorithms.
type CompressionStrategy interface {
	Compress(data []byte) []byte
	GetAlgorithm() string
}

// ZipCompression implements ZIP compression.
type ZipCompression struct{}

func (z *ZipCompression) Compress(data []byte) []byte {
	compressed := []byte(fmt.Sprintf("[ZIP compressed: %d bytes -> %d bytes]",
		len(data), len(data)/2))
	return compressed
}

func (z *ZipCompression) GetAlgorithm() string {
	return "ZIP"
}

// GzipCompression implements GZIP compression.
type GzipCompression struct{}

func (g *GzipCompression) Compress(data []byte) []byte {
	compressed := []byte(fmt.Sprintf("[GZIP compressed: %d bytes -> %d bytes]",
		len(data), int(float64(len(data))*0.4)))
	return compressed
}

func (g *GzipCompression) GetAlgorithm() string {
	return "GZIP"
}

// Bzip2Compression implements BZIP2 compression.
type Bzip2Compression struct{}

func (b *Bzip2Compression) Compress(data []byte) []byte {
	compressed := []byte(fmt.Sprintf("[BZIP2 compressed: %d bytes -> %d bytes]",
		len(data), int(float64(len(data))*0.3)))
	return compressed
}

func (b *Bzip2Compression) GetAlgorithm() string {
	return "BZIP2"
}

// FileCompressor uses compression strategies.
type FileCompressor struct {
	strategy CompressionStrategy
}

func NewFileCompressor(strategy CompressionStrategy) *FileCompressor {
	return &FileCompressor{strategy: strategy}
}

func (fc *FileCompressor) SetStrategy(strategy CompressionStrategy) {
	fc.strategy = strategy
}

func (fc *FileCompressor) Compress(data []byte) []byte {
	fmt.Printf("Compressing with %s algorithm...\n", fc.strategy.GetAlgorithm())
	return fc.strategy.Compress(data)
}

// SortStrategy defines the interface for sorting algorithms.
type SortStrategy interface {
	Sort(data []int) []int
	GetName() string
}

// BubbleSort implements bubble sort algorithm.
type BubbleSort struct{}

func (b *BubbleSort) Sort(data []int) []int {
	fmt.Println("  [Bubble Sort] Sorting... (simple, O(n²))")
	result := make([]int, len(data))
	copy(result, data)

	n := len(result)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}

func (b *BubbleSort) GetName() string {
	return "Bubble Sort"
}

// QuickSort implements quick sort algorithm.
type QuickSort struct{}

func (q *QuickSort) Sort(data []int) []int {
	fmt.Println("  [Quick Sort] Sorting... (efficient, O(n log n))")
	result := make([]int, len(data))
	copy(result, data)
	q.quickSort(result, 0, len(result)-1)
	return result
}

func (q *QuickSort) quickSort(arr []int, low, high int) {
	if low < high {
		pi := q.partition(arr, low, high)
		q.quickSort(arr, low, pi-1)
		q.quickSort(arr, pi+1, high)
	}
}

func (q *QuickSort) partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1

	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

func (q *QuickSort) GetName() string {
	return "Quick Sort"
}

// Sorter uses sorting strategies.
type Sorter struct {
	strategy SortStrategy
}

func NewSorter(strategy SortStrategy) *Sorter {
	return &Sorter{strategy: strategy}
}

func (s *Sorter) SetStrategy(strategy SortStrategy) {
	s.strategy = strategy
}

func (s *Sorter) Sort(data []int) []int {
	return s.strategy.Sort(data)
}

// PaymentStrategy defines the interface for payment processing.
type PaymentStrategy interface {
	Pay(amount float64) error
	GetMethod() string
}

// CreditCardPayment processes credit card payments.
type CreditCardPayment struct {
	cardNumber string
	cvv        string
}

func NewCreditCardPayment(cardNumber, cvv string) PaymentStrategy {
	return &CreditCardPayment{
		cardNumber: maskCardNumber(cardNumber),
		cvv:        cvv,
	}
}

func (c *CreditCardPayment) Pay(amount float64) error {
	fmt.Printf("  [Credit Card] Processing $%.2f payment with card %s\n",
		amount, c.cardNumber)
	return nil
}

func (c *CreditCardPayment) GetMethod() string {
	return "Credit Card"
}

func maskCardNumber(number string) string {
	if len(number) <= 4 {
		return number
	}
	return strings.Repeat("*", len(number)-4) + number[len(number)-4:]
}

// PayPalPayment processes PayPal payments.
type PayPalPayment struct {
	email string
}

func NewPayPalPayment(email string) PaymentStrategy {
	return &PayPalPayment{email: email}
}

func (p *PayPalPayment) Pay(amount float64) error {
	fmt.Printf("  [PayPal] Processing $%.2f payment for %s\n", amount, p.email)
	return nil
}

func (p *PayPalPayment) GetMethod() string {
	return "PayPal"
}

// CryptoPayment processes cryptocurrency payments.
type CryptoPayment struct {
	walletAddress string
	currency      string
}

func NewCryptoPayment(walletAddress, currency string) PaymentStrategy {
	return &CryptoPayment{
		walletAddress: walletAddress,
		currency:      currency,
	}
}

func (c *CryptoPayment) Pay(amount float64) error {
	// Simplified conversion rate
	rate := 30000.0 // $30k per BTC
	cryptoAmount := amount / rate

	fmt.Printf("  [Crypto] Processing $%.2f (%.8f %s) to wallet %s\n",
		amount, cryptoAmount, c.currency, c.walletAddress[:10]+"...")
	return nil
}

func (c *CryptoPayment) GetMethod() string {
	return "Cryptocurrency"
}

// PaymentProcessor uses payment strategies.
type PaymentProcessor struct {
	strategy PaymentStrategy
}

func NewPaymentProcessor(strategy PaymentStrategy) *PaymentProcessor {
	return &PaymentProcessor{strategy: strategy}
}

func (pp *PaymentProcessor) SetStrategy(strategy PaymentStrategy) {
	pp.strategy = strategy
}

func (pp *PaymentProcessor) ProcessPayment(amount float64) error {
	fmt.Printf("Processing payment of $%.2f using %s\n",
		amount, pp.strategy.GetMethod())
	return pp.strategy.Pay(amount)
}

// DiscountStrategy defines the interface for discount calculations.
type DiscountStrategy interface {
	ApplyDiscount(price float64) float64
	GetDescription() string
}

// NoDiscount applies no discount.
type NoDiscount struct{}

func (n *NoDiscount) ApplyDiscount(price float64) float64 {
	return price
}

func (n *NoDiscount) GetDescription() string {
	return "No discount"
}

// PercentageDiscount applies a percentage discount.
type PercentageDiscount struct {
	percentage float64
}

func NewPercentageDiscount(percentage float64) DiscountStrategy {
	return &PercentageDiscount{percentage: percentage}
}

func (p *PercentageDiscount) ApplyDiscount(price float64) float64 {
	discount := price * p.percentage / 100
	return math.Round((price-discount)*100) / 100
}

func (p *PercentageDiscount) GetDescription() string {
	return fmt.Sprintf("%.0f%% off", p.percentage)
}

// FixedAmountDiscount applies a fixed dollar amount discount.
type FixedAmountDiscount struct {
	amount float64
}

func NewFixedAmountDiscount(amount float64) DiscountStrategy {
	return &FixedAmountDiscount{amount: amount}
}

func (f *FixedAmountDiscount) ApplyDiscount(price float64) float64 {
	discounted := price - f.amount
	if discounted < 0 {
		return 0
	}
	return math.Round(discounted*100) / 100
}

func (f *FixedAmountDiscount) GetDescription() string {
	return fmt.Sprintf("$%.2f off", f.amount)
}

// PricingCalculator uses discount strategies.
type PricingCalculator struct {
	strategy DiscountStrategy
}

func NewPricingCalculator(strategy DiscountStrategy) *PricingCalculator {
	return &PricingCalculator{strategy: strategy}
}

func (pc *PricingCalculator) SetStrategy(strategy DiscountStrategy) {
	pc.strategy = strategy
}

func (pc *PricingCalculator) CalculateFinalPrice(basePrice float64) float64 {
	return pc.strategy.ApplyDiscount(basePrice)
}
