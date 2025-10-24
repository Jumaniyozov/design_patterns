// Package facade demonstrates the Facade Pattern, a structural design pattern
// that provides a simplified, unified interface to a complex subsystem.
//
// The Facade Pattern is essential for:
// - Simplifying complex subsystems with multiple interconnected components
// - Reducing dependencies between clients and implementation details
// - Creating clear architectural boundaries and layers
// - Integrating multiple third-party libraries with a cohesive interface
package facade

import (
	"fmt"
	"time"
)

// =============================================================================
// Example 1: Banking System Facade
// =============================================================================

// Subsystem components for banking operations

// AccountService handles account-related operations.
type AccountService struct{}

// ValidateAccount checks if an account exists and has sufficient balance.
func (a *AccountService) ValidateAccount(accountID string, amount float64) (bool, error) {
	// Simulate account validation
	if accountID == "" {
		return false, fmt.Errorf("invalid account ID")
	}
	if amount < 0 {
		return false, fmt.Errorf("invalid amount")
	}
	// Simulate checking balance (in real system, would query database)
	return true, nil
}

// GetAccountBalance retrieves the current account balance.
func (a *AccountService) GetAccountBalance(accountID string) (float64, error) {
	// Simulate balance lookup
	return 5000.00, nil
}

// DebitAccount removes funds from an account.
func (a *AccountService) DebitAccount(accountID string, amount float64) error {
	fmt.Printf("  ✓ Debited $%.2f from account %s\n", amount, accountID)
	return nil
}

// CreditAccount adds funds to an account.
func (a *AccountService) CreditAccount(accountID string, amount float64) error {
	fmt.Printf("  ✓ Credited $%.2f to account %s\n", amount, accountID)
	return nil
}

// TransactionService handles transaction processing.
type TransactionService struct{}

// CreateTransaction records a new transaction.
func (t *TransactionService) CreateTransaction(from, to string, amount float64) (string, error) {
	transactionID := fmt.Sprintf("TXN-%d", time.Now().Unix())
	fmt.Printf("  ✓ Transaction created: %s\n", transactionID)
	return transactionID, nil
}

// ProcessTransaction processes a pending transaction.
func (t *TransactionService) ProcessTransaction(transactionID string) error {
	fmt.Printf("  ✓ Transaction %s processed\n", transactionID)
	return nil
}

// FraudDetectionService checks for fraudulent activity.
type FraudDetectionService struct{}

// CheckFraud analyzes a transaction for fraud indicators.
func (f *FraudDetectionService) CheckFraud(from, to string, amount float64) (bool, error) {
	// Simulate fraud detection (in real system, would use ML models)
	if amount > 10000 {
		return false, fmt.Errorf("transaction exceeds fraud threshold")
	}
	fmt.Printf("  ✓ Fraud check passed for $%.2f\n", amount)
	return true, nil
}

// NotificationService handles notifications to users.
type NotificationService struct{}

// SendNotification sends a notification to a user.
func (n *NotificationService) SendNotification(accountID, message string) error {
	fmt.Printf("  ✓ Notification sent to %s: %s\n", accountID, message)
	return nil
}

// BankingFacade provides a simplified interface to the banking subsystem.
// This is the Facade that coordinates all subsystem components.
type BankingFacade struct {
	accountService  *AccountService
	transactionSvc  *TransactionService
	fraudDetection  *FraudDetectionService
	notification    *NotificationService
}

// NewBankingFacade creates a new banking facade.
func NewBankingFacade() *BankingFacade {
	return &BankingFacade{
		accountService:  &AccountService{},
		transactionSvc:  &TransactionService{},
		fraudDetection:  &FraudDetectionService{},
		notification:    &NotificationService{},
	}
}

// TransferMoney orchestrates a money transfer between two accounts.
// This single method hides the complexity of coordinating multiple subsystems.
func (b *BankingFacade) TransferMoney(fromAccount, toAccount string, amount float64) error {
	fmt.Printf("\n🏦 Initiating transfer: $%.2f from %s to %s\n", amount, fromAccount, toAccount)

	// Step 1: Validate source account
	valid, err := b.accountService.ValidateAccount(fromAccount, amount)
	if !valid || err != nil {
		return fmt.Errorf("source account validation failed: %w", err)
	}

	// Step 2: Validate destination account
	valid, err = b.accountService.ValidateAccount(toAccount, 0)
	if !valid || err != nil {
		return fmt.Errorf("destination account validation failed: %w", err)
	}

	// Step 3: Check for fraud
	_, err = b.fraudDetection.CheckFraud(fromAccount, toAccount, amount)
	if err != nil {
		return fmt.Errorf("fraud detection failed: %w", err)
	}

	// Step 4: Create transaction record
	txnID, err := b.transactionSvc.CreateTransaction(fromAccount, toAccount, amount)
	if err != nil {
		return fmt.Errorf("transaction creation failed: %w", err)
	}

	// Step 5: Debit source account
	if err := b.accountService.DebitAccount(fromAccount, amount); err != nil {
		return fmt.Errorf("debit failed: %w", err)
	}

	// Step 6: Credit destination account
	if err := b.accountService.CreditAccount(toAccount, amount); err != nil {
		// In real system, would rollback the debit
		return fmt.Errorf("credit failed: %w", err)
	}

	// Step 7: Process transaction
	if err := b.transactionSvc.ProcessTransaction(txnID); err != nil {
		return fmt.Errorf("transaction processing failed: %w", err)
	}

	// Step 8: Send notifications
	b.notification.SendNotification(fromAccount, fmt.Sprintf("$%.2f transferred to %s", amount, toAccount))
	b.notification.SendNotification(toAccount, fmt.Sprintf("$%.2f received from %s", amount, fromAccount))

	fmt.Printf("✅ Transfer completed successfully!\n")
	return nil
}

// GetAccountBalance provides a simple way to check account balance.
func (b *BankingFacade) GetAccountBalance(accountID string) (float64, error) {
	return b.accountService.GetAccountBalance(accountID)
}

// =============================================================================
// Example 2: Smart Home Automation Facade
// =============================================================================

// Subsystem components for smart home

// LightingSystem controls home lighting.
type LightingSystem struct{}

// TurnOn turns lights on.
func (l *LightingSystem) TurnOn(room string) {
	fmt.Printf("  💡 Lights turned ON in %s\n", room)
}

// TurnOff turns lights off.
func (l *LightingSystem) TurnOff(room string) {
	fmt.Printf("  💡 Lights turned OFF in %s\n", room)
}

// Dim dims the lights to a percentage.
func (l *LightingSystem) Dim(room string, level int) {
	fmt.Printf("  💡 Lights dimmed to %d%% in %s\n", level, room)
}

// SecuritySystem handles home security.
type SecuritySystem struct {
	armed bool
}

// Arm activates the security system.
func (s *SecuritySystem) Arm() {
	s.armed = true
	fmt.Println("  🔒 Security system ARMED")
}

// Disarm deactivates the security system.
func (s *SecuritySystem) Disarm() {
	s.armed = false
	fmt.Println("  🔓 Security system DISARMED")
}

// IsArmed checks if the system is armed.
func (s *SecuritySystem) IsArmed() bool {
	return s.armed
}

// ClimateControl manages temperature and HVAC.
type ClimateControl struct{}

// SetTemperature sets the target temperature.
func (c *ClimateControl) SetTemperature(temp int) {
	fmt.Printf("  🌡️  Temperature set to %d°F\n", temp)
}

// TurnOnHeating activates heating.
func (c *ClimateControl) TurnOnHeating() {
	fmt.Println("  🔥 Heating turned ON")
}

// TurnOffHeating deactivates heating.
func (c *ClimateControl) TurnOffHeating() {
	fmt.Println("  ❄️  Heating turned OFF")
}

// EntertainmentSystem controls audio/video equipment.
type EntertainmentSystem struct{}

// TurnOn powers on the entertainment system.
func (e *EntertainmentSystem) TurnOn() {
	fmt.Println("  📺 Entertainment system turned ON")
}

// TurnOff powers off the entertainment system.
func (e *EntertainmentSystem) TurnOff() {
	fmt.Println("  📺 Entertainment system turned OFF")
}

// PlayMovie starts movie playback.
func (e *EntertainmentSystem) PlayMovie(title string) {
	fmt.Printf("  🎬 Playing movie: %s\n", title)
}

// SmartHomeFacade provides simple commands for complex home automation scenarios.
type SmartHomeFacade struct {
	lighting      *LightingSystem
	security      *SecuritySystem
	climate       *ClimateControl
	entertainment *EntertainmentSystem
}

// NewSmartHomeFacade creates a new smart home facade.
func NewSmartHomeFacade() *SmartHomeFacade {
	return &SmartHomeFacade{
		lighting:      &LightingSystem{},
		security:      &SecuritySystem{},
		climate:       &ClimateControl{},
		entertainment: &EntertainmentSystem{},
	}
}

// LeaveHome activates all systems for leaving home.
func (s *SmartHomeFacade) LeaveHome() {
	fmt.Println("\n🏠 Activating 'Leave Home' mode...")
	s.lighting.TurnOff("Living Room")
	s.lighting.TurnOff("Bedroom")
	s.lighting.TurnOff("Kitchen")
	s.security.Arm()
	s.climate.SetTemperature(65)
	s.entertainment.TurnOff()
	fmt.Println("✅ Home secured for departure!")
}

// ArriveHome activates all systems for arriving home.
func (s *SmartHomeFacade) ArriveHome() {
	fmt.Println("\n🏠 Activating 'Arrive Home' mode...")
	s.security.Disarm()
	s.lighting.TurnOn("Living Room")
	s.lighting.TurnOn("Kitchen")
	s.climate.SetTemperature(72)
	fmt.Println("✅ Welcome home!")
}

// MovieNight sets up the perfect movie-watching environment.
func (s *SmartHomeFacade) MovieNight(movieTitle string) {
	fmt.Println("\n🎬 Activating 'Movie Night' mode...")
	s.lighting.Dim("Living Room", 20)
	s.lighting.TurnOff("Kitchen")
	s.climate.SetTemperature(70)
	s.entertainment.TurnOn()
	s.entertainment.PlayMovie(movieTitle)
	fmt.Println("✅ Enjoy your movie!")
}

// SleepMode prepares the home for sleep.
func (s *SmartHomeFacade) SleepMode() {
	fmt.Println("\n😴 Activating 'Sleep' mode...")
	s.lighting.TurnOff("Living Room")
	s.lighting.TurnOff("Kitchen")
	s.lighting.Dim("Bedroom", 10)
	s.security.Arm()
	s.climate.SetTemperature(68)
	s.entertainment.TurnOff()
	fmt.Println("✅ Goodnight!")
}

// =============================================================================
// Example 3: E-Commerce Order Processing Facade
// =============================================================================

// InventoryService manages product inventory.
type InventoryService struct{}

// CheckStock verifies product availability.
func (i *InventoryService) CheckStock(productID string, quantity int) (bool, error) {
	fmt.Printf("  ✓ Checked stock for product %s: %d available\n", productID, quantity)
	return true, nil
}

// ReserveStock reserves items for an order.
func (i *InventoryService) ReserveStock(productID string, quantity int) error {
	fmt.Printf("  ✓ Reserved %d units of product %s\n", quantity, productID)
	return nil
}

// PaymentService processes payments.
type PaymentService struct{}

// ProcessPayment charges the customer.
func (p *PaymentService) ProcessPayment(customerID string, amount float64) (string, error) {
	paymentID := fmt.Sprintf("PAY-%d", time.Now().Unix())
	fmt.Printf("  ✓ Payment processed: %s ($%.2f)\n", paymentID, amount)
	return paymentID, nil
}

// ShippingService handles order shipping.
type ShippingService struct{}

// CreateShipment creates a shipping label and arranges delivery.
func (s *ShippingService) CreateShipment(orderID, address string) (string, error) {
	trackingNum := fmt.Sprintf("TRACK-%d", time.Now().Unix())
	fmt.Printf("  ✓ Shipment created with tracking: %s\n", trackingNum)
	return trackingNum, nil
}

// EmailService sends email notifications.
type EmailService struct{}

// SendEmail sends an email to a customer.
func (e *EmailService) SendEmail(to, subject, body string) error {
	fmt.Printf("  ✓ Email sent to %s: %s\n", to, subject)
	return nil
}

// OrderProcessingFacade simplifies the complex order processing workflow.
type OrderProcessingFacade struct {
	inventory *InventoryService
	payment   *PaymentService
	shipping  *ShippingService
	email     *EmailService
}

// NewOrderProcessingFacade creates a new order processing facade.
func NewOrderProcessingFacade() *OrderProcessingFacade {
	return &OrderProcessingFacade{
		inventory: &InventoryService{},
		payment:   &PaymentService{},
		shipping:  &ShippingService{},
		email:     &EmailService{},
	}
}

// PlaceOrder orchestrates the entire order placement process.
func (o *OrderProcessingFacade) PlaceOrder(customerID, productID string, quantity int, amount float64, address string) error {
	fmt.Printf("\n📦 Processing order for customer %s...\n", customerID)

	// Step 1: Check inventory
	available, err := o.inventory.CheckStock(productID, quantity)
	if !available || err != nil {
		return fmt.Errorf("product not available: %w", err)
	}

	// Step 2: Reserve stock
	if err := o.inventory.ReserveStock(productID, quantity); err != nil {
		return fmt.Errorf("stock reservation failed: %w", err)
	}

	// Step 3: Process payment
	paymentID, err := o.payment.ProcessPayment(customerID, amount)
	if err != nil {
		return fmt.Errorf("payment failed: %w", err)
	}

	// Step 4: Create shipment
	trackingNum, err := o.shipping.CreateShipment(fmt.Sprintf("ORD-%d", time.Now().Unix()), address)
	if err != nil {
		return fmt.Errorf("shipment creation failed: %w", err)
	}

	// Step 5: Send confirmation email
	emailBody := fmt.Sprintf("Your order has been confirmed. Payment ID: %s, Tracking: %s", paymentID, trackingNum)
	if err := o.email.SendEmail(customerID, "Order Confirmation", emailBody); err != nil {
		return fmt.Errorf("email notification failed: %w", err)
	}

	fmt.Println("✅ Order placed successfully!")
	return nil
}
