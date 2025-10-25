// Package anticorruptionlayer demonstrates the Anti-Corruption Layer pattern.
// It isolates domain models from external systems by translating between
// different representations, preventing external concepts from corrupting the domain.
package anticorruptionlayer

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Clean Domain Model (Our system)

type Customer struct {
	ID    string
	Name  string
	Email string
}

type LineItem struct {
	ProductID string
	Name      string
	Quantity  int
	Price     float64
}

type Order struct {
	ID          string
	Customer    Customer
	Items       []LineItem
	Total       float64
	Status      OrderStatus
	PlacedAt    time.Time
}

type OrderStatus int

const (
	OrderPending OrderStatus = iota
	OrderConfirmed
	OrderShipped
	OrderDelivered
	OrderCancelled
)

func (os OrderStatus) String() string {
	return []string{"Pending", "Confirmed", "Shipped", "Delivered", "Cancelled"}[os]
}

// Legacy External System (Poor design, different model)

type LegacyCustomerRecord struct {
	CustNum  int
	FullName string
	EmailAddr string
}

type LegacyOrderRecord struct {
	OrderNum    string
	CustNum     int
	ItemList    string  // Comma-separated: "id:qty:price,id:qty:price"
	TotalAmt    float64
	StatusCode  int     // 0=new, 1=paid, 2=sent, 3=done, 99=cancel
	OrderDate   string  // Format: "2006-01-02T15:04:05"
}

type LegacyProductRecord struct {
	ProdID   string
	ProdName string
}

// Legacy External Service (represents external system)

type LegacyOrderService struct {
	orders   map[string]*LegacyOrderRecord
	customers map[int]*LegacyCustomerRecord
	products map[string]*LegacyProductRecord
}

func NewLegacyOrderService() *LegacyOrderService {
	return &LegacyOrderService{
		orders:   make(map[string]*LegacyOrderRecord),
		customers: make(map[int]*LegacyCustomerRecord),
		products: make(map[string]*LegacyProductRecord),
	}
}

func (s *LegacyOrderService) GetOrder(orderNum string) (*LegacyOrderRecord, error) {
	order, exists := s.orders[orderNum]
	if !exists {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

func (s *LegacyOrderService) GetCustomer(custNum int) (*LegacyCustomerRecord, error) {
	customer, exists := s.customers[custNum]
	if !exists {
		return nil, fmt.Errorf("customer not found")
	}
	return customer, nil
}

func (s *LegacyOrderService) GetProduct(prodID string) (*LegacyProductRecord, error) {
	product, exists := s.products[prodID]
	if !exists {
		return nil, fmt.Errorf("product not found")
	}
	return product, nil
}

// Seed with sample data
func (s *LegacyOrderService) SeedData() {
	s.customers[12345] = &LegacyCustomerRecord{
		CustNum:  12345,
		FullName: "John Doe",
		EmailAddr: "john.doe@example.com",
	}

	s.products["PROD001"] = &LegacyProductRecord{
		ProdID:   "PROD001",
		ProdName: "Widget",
	}

	s.products["PROD002"] = &LegacyProductRecord{
		ProdID:   "PROD002",
		ProdName: "Gadget",
	}

	s.orders["ORD-2024-001"] = &LegacyOrderRecord{
		OrderNum:   "ORD-2024-001",
		CustNum:    12345,
		ItemList:   "PROD001:2:19.99,PROD002:1:29.99",
		TotalAmt:   69.97,
		StatusCode: 1,
		OrderDate:  "2024-01-15T10:30:00",
	}
}

// Anti-Corruption Layer Components

// OrderTranslator translates between legacy and domain models
type OrderTranslator struct {
	legacyService *LegacyOrderService
}

func NewOrderTranslator(legacyService *LegacyOrderService) *OrderTranslator {
	return &OrderTranslator{legacyService: legacyService}
}

// TranslateOrder translates legacy order to domain order
func (t *OrderTranslator) TranslateOrder(legacyOrder *LegacyOrderRecord) (*Order, error) {
	// Translate customer
	legacyCustomer, err := t.legacyService.GetCustomer(legacyOrder.CustNum)
	if err != nil {
		return nil, err
	}

	customer := Customer{
		ID:    fmt.Sprintf("CUST-%d", legacyCustomer.CustNum),
		Name:  legacyCustomer.FullName,
		Email: legacyCustomer.EmailAddr,
	}

	// Translate items
	items, err := t.translateItems(legacyOrder.ItemList)
	if err != nil {
		return nil, err
	}

	// Translate status
	status := t.translateStatus(legacyOrder.StatusCode)

	// Parse date
	placedAt, _ := time.Parse("2006-01-02T15:04:05", legacyOrder.OrderDate)

	return &Order{
		ID:       legacyOrder.OrderNum,
		Customer: customer,
		Items:    items,
		Total:    legacyOrder.TotalAmt,
		Status:   status,
		PlacedAt: placedAt,
	}, nil
}

func (t *OrderTranslator) translateItems(itemList string) ([]LineItem, error) {
	items := make([]LineItem, 0)

	if itemList == "" {
		return items, nil
	}

	// Parse: "PROD001:2:19.99,PROD002:1:29.99"
	parts := strings.Split(itemList, ",")
	for _, part := range parts {
		fields := strings.Split(part, ":")
		if len(fields) != 3 {
			continue
		}

		productID := fields[0]
		quantity, _ := strconv.Atoi(fields[1])
		price, _ := strconv.ParseFloat(fields[2], 64)

		// Get product name
		product, err := t.legacyService.GetProduct(productID)
		if err != nil {
			return nil, err
		}

		items = append(items, LineItem{
			ProductID: productID,
			Name:      product.ProdName,
			Quantity:  quantity,
			Price:     price,
		})
	}

	return items, nil
}

func (t *OrderTranslator) translateStatus(statusCode int) OrderStatus {
	switch statusCode {
	case 0:
		return OrderPending
	case 1:
		return OrderConfirmed
	case 2:
		return OrderShipped
	case 3:
		return OrderDelivered
	case 99:
		return OrderCancelled
	default:
		return OrderPending
	}
}

// OrderACL is the Anti-Corruption Layer facade
type OrderACL struct {
	legacyService *LegacyOrderService
	translator    *OrderTranslator
}

// NewOrderACL creates an order ACL
func NewOrderACL(legacyService *LegacyOrderService) *OrderACL {
	return &OrderACL{
		legacyService: legacyService,
		translator:    NewOrderTranslator(legacyService),
	}
}

// GetOrder retrieves an order in domain model
func (acl *OrderACL) GetOrder(orderID string) (*Order, error) {
	// Call legacy system
	legacyOrder, err := acl.legacyService.GetOrder(orderID)
	if err != nil {
		return nil, err
	}

	// Translate to domain model
	return acl.translator.TranslateOrder(legacyOrder)
}

// OrderService uses ACL to work with clean domain model
type OrderService struct {
	acl *OrderACL
}

// NewOrderService creates an order service
func NewOrderService(acl *OrderACL) *OrderService {
	return &OrderService{acl: acl}
}

// GetOrderDetails gets order details in domain terms
func (os *OrderService) GetOrderDetails(orderID string) (string, error) {
	order, err := os.acl.GetOrder(orderID)
	if err != nil {
		return "", err
	}

	details := fmt.Sprintf("Order: %s\n", order.ID)
	details += fmt.Sprintf("Customer: %s (%s)\n", order.Customer.Name, order.Customer.Email)
	details += fmt.Sprintf("Status: %s\n", order.Status)
	details += fmt.Sprintf("Placed: %s\n", order.PlacedAt.Format("Jan 02, 2006"))
	details += "Items:\n"

	for _, item := range order.Items {
		details += fmt.Sprintf("  - %s (x%d) @ $%.2f\n",
			item.Name, item.Quantity, item.Price)
	}

	details += fmt.Sprintf("Total: $%.2f\n", order.Total)

	return details, nil
}

// Example: Third-Party Payment Gateway ACL

type PaymentGatewayRequest struct {
	MerchantID string
	TxnAmount  int // Amount in cents
	CardNum    string
	CVV        string
}

type PaymentGatewayResponse struct {
	TxnID      string
	ResultCode int // 0=success, 1=declined, 2=error
	Message    string
}

// Our clean payment model
type Payment struct {
	ID     string
	Amount float64
	Status PaymentStatus
}

type PaymentStatus int

const (
	PaymentSuccess PaymentStatus = iota
	PaymentDeclined
	PaymentError
)

// PaymentGatewayACL translates between models
type PaymentGatewayACL struct {
	merchantID string
}

func NewPaymentGatewayACL(merchantID string) *PaymentGatewayACL {
	return &PaymentGatewayACL{merchantID: merchantID}
}

// ProcessPayment processes payment through gateway
func (acl *PaymentGatewayACL) ProcessPayment(amount float64, cardNumber, cvv string) (*Payment, error) {
	// Translate to gateway format (cents)
	gatewayReq := &PaymentGatewayRequest{
		MerchantID: acl.merchantID,
		TxnAmount:  int(amount * 100),
		CardNum:    cardNumber,
		CVV:        cvv,
	}

	// Call gateway (simulated)
	gatewayResp := acl.callGateway(gatewayReq)

	// Translate response to domain model
	return &Payment{
		ID:     gatewayResp.TxnID,
		Amount: amount,
		Status: acl.translateStatus(gatewayResp.ResultCode),
	}, nil
}

func (acl *PaymentGatewayACL) callGateway(req *PaymentGatewayRequest) *PaymentGatewayResponse {
	// Simulate gateway call
	return &PaymentGatewayResponse{
		TxnID:      fmt.Sprintf("TXN-%d", time.Now().Unix()),
		ResultCode: 0,
		Message:    "Approved",
	}
}

func (acl *PaymentGatewayACL) translateStatus(code int) PaymentStatus {
	switch code {
	case 0:
		return PaymentSuccess
	case 1:
		return PaymentDeclined
	default:
		return PaymentError
	}
}
