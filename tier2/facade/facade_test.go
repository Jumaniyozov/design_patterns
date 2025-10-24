package facade

import (
	"testing"
)

// =============================================================================
// Banking System Facade Tests
// =============================================================================

func TestBankingFacade_TransferMoney_Success(t *testing.T) {
	bank := NewBankingFacade()

	err := bank.TransferMoney("ACC-12345", "ACC-67890", 500.00)
	if err != nil {
		t.Errorf("Expected successful transfer, got error: %v", err)
	}
}

func TestBankingFacade_TransferMoney_InvalidSourceAccount(t *testing.T) {
	bank := NewBankingFacade()

	err := bank.TransferMoney("", "ACC-67890", 500.00)
	if err == nil {
		t.Error("Expected error for invalid source account, got nil")
	}
}

func TestBankingFacade_TransferMoney_InvalidAmount(t *testing.T) {
	bank := NewBankingFacade()

	err := bank.TransferMoney("ACC-12345", "ACC-67890", -100.00)
	if err == nil {
		t.Error("Expected error for negative amount, got nil")
	}
}

func TestBankingFacade_TransferMoney_FraudDetection(t *testing.T) {
	bank := NewBankingFacade()

	// Amount exceeds fraud threshold (>10000)
	err := bank.TransferMoney("ACC-12345", "ACC-67890", 15000.00)
	if err == nil {
		t.Error("Expected fraud detection error for large amount, got nil")
	}
}

func TestBankingFacade_GetAccountBalance(t *testing.T) {
	bank := NewBankingFacade()

	balance, err := bank.GetAccountBalance("ACC-12345")
	if err != nil {
		t.Errorf("Expected successful balance retrieval, got error: %v", err)
	}

	if balance <= 0 {
		t.Errorf("Expected positive balance, got %.2f", balance)
	}
}

func TestAccountService_ValidateAccount_InvalidID(t *testing.T) {
	svc := &AccountService{}

	valid, err := svc.ValidateAccount("", 100.00)
	if valid {
		t.Error("Expected validation to fail for empty account ID")
	}
	if err == nil {
		t.Error("Expected error for empty account ID, got nil")
	}
}

func TestAccountService_ValidateAccount_InvalidAmount(t *testing.T) {
	svc := &AccountService{}

	valid, err := svc.ValidateAccount("ACC-123", -50.00)
	if valid {
		t.Error("Expected validation to fail for negative amount")
	}
	if err == nil {
		t.Error("Expected error for negative amount, got nil")
	}
}

func TestAccountService_ValidateAccount_Success(t *testing.T) {
	svc := &AccountService{}

	valid, err := svc.ValidateAccount("ACC-123", 100.00)
	if !valid {
		t.Error("Expected validation to succeed for valid inputs")
	}
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestFraudDetectionService_CheckFraud_ExceedsThreshold(t *testing.T) {
	svc := &FraudDetectionService{}

	valid, err := svc.CheckFraud("ACC-1", "ACC-2", 15000.00)
	if valid {
		t.Error("Expected fraud check to fail for amount exceeding threshold")
	}
	if err == nil {
		t.Error("Expected error for fraudulent transaction, got nil")
	}
}

func TestFraudDetectionService_CheckFraud_BelowThreshold(t *testing.T) {
	svc := &FraudDetectionService{}

	valid, err := svc.CheckFraud("ACC-1", "ACC-2", 5000.00)
	if !valid {
		t.Error("Expected fraud check to pass for normal amount")
	}
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// =============================================================================
// Smart Home Facade Tests
// =============================================================================

func TestSmartHomeFacade_LeaveHome(t *testing.T) {
	home := NewSmartHomeFacade()

	// Should not panic or error
	home.LeaveHome()

	// Verify security is armed
	if !home.security.IsArmed() {
		t.Error("Expected security system to be armed after leaving home")
	}
}

func TestSmartHomeFacade_ArriveHome(t *testing.T) {
	home := NewSmartHomeFacade()

	// First arm the system
	home.security.Arm()

	// Then arrive home
	home.ArriveHome()

	// Verify security is disarmed
	if home.security.IsArmed() {
		t.Error("Expected security system to be disarmed after arriving home")
	}
}

func TestSmartHomeFacade_MovieNight(t *testing.T) {
	home := NewSmartHomeFacade()

	// Should not panic or error
	home.MovieNight("Inception")
}

func TestSmartHomeFacade_SleepMode(t *testing.T) {
	home := NewSmartHomeFacade()

	// Should not panic or error
	home.SleepMode()

	// Verify security is armed
	if !home.security.IsArmed() {
		t.Error("Expected security system to be armed in sleep mode")
	}
}

func TestSecuritySystem_ArmDisarm(t *testing.T) {
	security := &SecuritySystem{}

	if security.IsArmed() {
		t.Error("Expected security to be initially disarmed")
	}

	security.Arm()
	if !security.IsArmed() {
		t.Error("Expected security to be armed after Arm()")
	}

	security.Disarm()
	if security.IsArmed() {
		t.Error("Expected security to be disarmed after Disarm()")
	}
}

// =============================================================================
// Order Processing Facade Tests
// =============================================================================

func TestOrderProcessingFacade_PlaceOrder_Success(t *testing.T) {
	orderSystem := NewOrderProcessingFacade()

	err := orderSystem.PlaceOrder(
		"customer@example.com",
		"PROD-123",
		2,
		199.99,
		"123 Main St",
	)

	if err != nil {
		t.Errorf("Expected successful order placement, got error: %v", err)
	}
}

func TestInventoryService_CheckStock(t *testing.T) {
	svc := &InventoryService{}

	available, err := svc.CheckStock("PROD-123", 5)
	if !available {
		t.Error("Expected stock to be available")
	}
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestInventoryService_ReserveStock(t *testing.T) {
	svc := &InventoryService{}

	err := svc.ReserveStock("PROD-123", 5)
	if err != nil {
		t.Errorf("Expected successful stock reservation, got error: %v", err)
	}
}

func TestPaymentService_ProcessPayment(t *testing.T) {
	svc := &PaymentService{}

	paymentID, err := svc.ProcessPayment("customer@example.com", 99.99)
	if err != nil {
		t.Errorf("Expected successful payment processing, got error: %v", err)
	}

	if paymentID == "" {
		t.Error("Expected non-empty payment ID")
	}
}

func TestShippingService_CreateShipment(t *testing.T) {
	svc := &ShippingService{}

	trackingNum, err := svc.CreateShipment("ORD-123", "123 Main St")
	if err != nil {
		t.Errorf("Expected successful shipment creation, got error: %v", err)
	}

	if trackingNum == "" {
		t.Error("Expected non-empty tracking number")
	}
}

func TestEmailService_SendEmail(t *testing.T) {
	svc := &EmailService{}

	err := svc.SendEmail("customer@example.com", "Test Subject", "Test Body")
	if err != nil {
		t.Errorf("Expected successful email sending, got error: %v", err)
	}
}

// =============================================================================
// Integration Tests
// =============================================================================

func TestBankingFacade_MultipleTransfers(t *testing.T) {
	bank := NewBankingFacade()

	// Perform multiple transfers
	transfers := []struct {
		from   string
		to     string
		amount float64
	}{
		{"ACC-1", "ACC-2", 100.00},
		{"ACC-2", "ACC-3", 50.00},
		{"ACC-3", "ACC-1", 25.00},
	}

	for _, transfer := range transfers {
		err := bank.TransferMoney(transfer.from, transfer.to, transfer.amount)
		if err != nil {
			t.Errorf("Transfer from %s to %s failed: %v",
				transfer.from, transfer.to, err)
		}
	}
}

func TestSmartHomeFacade_DailyRoutine(t *testing.T) {
	home := NewSmartHomeFacade()

	// Simulate a daily routine
	home.LeaveHome()
	if !home.security.IsArmed() {
		t.Error("Security should be armed when leaving")
	}

	home.ArriveHome()
	if home.security.IsArmed() {
		t.Error("Security should be disarmed when arriving")
	}

	home.MovieNight("The Matrix")

	home.SleepMode()
	if !home.security.IsArmed() {
		t.Error("Security should be armed in sleep mode")
	}
}

func TestOrderProcessingFacade_MultipleOrders(t *testing.T) {
	orderSystem := NewOrderProcessingFacade()

	orders := []struct {
		customer  string
		product   string
		quantity  int
		amount    float64
		address   string
	}{
		{"alice@example.com", "PROD-1", 1, 49.99, "123 Main St"},
		{"bob@example.com", "PROD-2", 2, 99.99, "456 Oak Ave"},
		{"carol@example.com", "PROD-3", 3, 149.99, "789 Pine Rd"},
	}

	for _, order := range orders {
		err := orderSystem.PlaceOrder(
			order.customer,
			order.product,
			order.quantity,
			order.amount,
			order.address,
		)
		if err != nil {
			t.Errorf("Order for %s failed: %v", order.customer, err)
		}
	}
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkBankingFacade_TransferMoney(b *testing.B) {
	bank := NewBankingFacade()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bank.TransferMoney("ACC-12345", "ACC-67890", 500.00)
	}
}

func BenchmarkSmartHomeFacade_LeaveHome(b *testing.B) {
	home := NewSmartHomeFacade()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		home.LeaveHome()
	}
}

func BenchmarkOrderProcessingFacade_PlaceOrder(b *testing.B) {
	orderSystem := NewOrderProcessingFacade()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		orderSystem.PlaceOrder(
			"customer@example.com",
			"PROD-123",
			2,
			199.99,
			"123 Main St",
		)
	}
}
