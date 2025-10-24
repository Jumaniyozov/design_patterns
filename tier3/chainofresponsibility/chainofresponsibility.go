// Package chainofresponsibility implements the Chain of Responsibility pattern.
//
// The Chain of Responsibility pattern passes requests along a chain of handlers.
// Each handler decides whether to process the request and/or pass it to the next
// handler in the chain.
//
// Key components:
// - Handler: Interface defining the handle method
// - ConcreteHandler: Specific handler implementation
// - Chain: Sequence of linked handlers
// - Request: Data passed through the chain
package chainofresponsibility

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Handler defines the interface for handling requests.
type Handler interface {
	Handle(ctx context.Context, request interface{}) error
	SetNext(handler Handler) Handler
}

// BaseHandler provides common chain management functionality.
type BaseHandler struct {
	next Handler
}

// SetNext sets the next handler in the chain.
func (h *BaseHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

// CallNext calls the next handler if it exists.
func (h *BaseHandler) CallNext(ctx context.Context, request interface{}) error {
	if h.next != nil {
		return h.next.Handle(ctx, request)
	}
	return nil
}

// HTTPRequest represents an HTTP request to be processed.
type HTTPRequest struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    string
	User    string
	Role    string
}

// NewHTTPRequest creates a new HTTP request.
func NewHTTPRequest(method, path string) *HTTPRequest {
	return &HTTPRequest{
		Method:  method,
		Path:    path,
		Headers: make(map[string]string),
	}
}

// AuthenticationHandler validates user authentication.
type AuthenticationHandler struct {
	BaseHandler
}

// NewAuthenticationHandler creates a new authentication handler.
func NewAuthenticationHandler() *AuthenticationHandler {
	return &AuthenticationHandler{}
}

// Handle processes authentication.
func (h *AuthenticationHandler) Handle(ctx context.Context, request interface{}) error {
	req, ok := request.(*HTTPRequest)
	if !ok {
		return errors.New("invalid request type")
	}

	// Check for auth token
	token := req.Headers["Authorization"]
	if token == "" {
		return errors.New("authentication failed: no token provided")
	}

	// Simulate token validation
	if !strings.HasPrefix(token, "Bearer ") {
		return errors.New("authentication failed: invalid token format")
	}

	// Extract user from token (simplified)
	req.User = strings.TrimPrefix(token, "Bearer ")

	fmt.Printf("[Auth]  Authenticated user: %s\n", req.User)

	// Pass to next handler
	return h.CallNext(ctx, request)
}

// AuthorizationHandler validates user permissions.
type AuthorizationHandler struct {
	BaseHandler
	requiredRole string
}

// NewAuthorizationHandler creates a new authorization handler.
func NewAuthorizationHandler(requiredRole string) *AuthorizationHandler {
	return &AuthorizationHandler{
		requiredRole: requiredRole,
	}
}

// Handle processes authorization.
func (h *AuthorizationHandler) Handle(ctx context.Context, request interface{}) error {
	req, ok := request.(*HTTPRequest)
	if !ok {
		return errors.New("invalid request type")
	}

	// Check user role
	if req.Role == "" {
		return errors.New("authorization failed: no role assigned")
	}

	if h.requiredRole != "" && req.Role != h.requiredRole {
		return fmt.Errorf("authorization failed: requires %s role, got %s", h.requiredRole, req.Role)
	}

	fmt.Printf("[Authz]  Authorized role: %s\n", req.Role)

	// Pass to next handler
	return h.CallNext(ctx, request)
}

// ValidationHandler validates request data.
type ValidationHandler struct {
	BaseHandler
}

// NewValidationHandler creates a new validation handler.
func NewValidationHandler() *ValidationHandler {
	return &ValidationHandler{}
}

// Handle processes validation.
func (h *ValidationHandler) Handle(ctx context.Context, request interface{}) error {
	req, ok := request.(*HTTPRequest)
	if !ok {
		return errors.New("invalid request type")
	}

	// Validate method
	validMethods := []string{"GET", "POST", "PUT", "DELETE"}
	valid := false
	for _, method := range validMethods {
		if req.Method == method {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("validation failed: invalid method %s", req.Method)
	}

	// Validate path
	if req.Path == "" {
		return errors.New("validation failed: empty path")
	}

	fmt.Printf("[Validation]  Request valid: %s %s\n", req.Method, req.Path)

	// Pass to next handler
	return h.CallNext(ctx, request)
}

// RateLimitHandler implements rate limiting.
type RateLimitHandler struct {
	BaseHandler
	maxRequests int
	window      time.Duration
	requests    map[string][]time.Time
}

// NewRateLimitHandler creates a new rate limit handler.
func NewRateLimitHandler(maxRequests int, window time.Duration) *RateLimitHandler {
	return &RateLimitHandler{
		maxRequests: maxRequests,
		window:      window,
		requests:    make(map[string][]time.Time),
	}
}

// Handle processes rate limiting.
func (h *RateLimitHandler) Handle(ctx context.Context, request interface{}) error {
	req, ok := request.(*HTTPRequest)
	if !ok {
		return errors.New("invalid request type")
	}

	// Get user's request history
	now := time.Now()
	userRequests := h.requests[req.User]

	// Filter requests within window
	var validRequests []time.Time
	for _, t := range userRequests {
		if now.Sub(t) < h.window {
			validRequests = append(validRequests, t)
		}
	}

	// Check rate limit
	if len(validRequests) >= h.maxRequests {
		return fmt.Errorf("rate limit exceeded: %d requests in %v", len(validRequests), h.window)
	}

	// Add current request
	validRequests = append(validRequests, now)
	h.requests[req.User] = validRequests

	fmt.Printf("[RateLimit]  Request allowed: %d/%d in %v\n", len(validRequests), h.maxRequests, h.window)

	// Pass to next handler
	return h.CallNext(ctx, request)
}

// LoggingHandler logs requests.
type LoggingHandler struct {
	BaseHandler
}

// NewLoggingHandler creates a new logging handler.
func NewLoggingHandler() *LoggingHandler {
	return &LoggingHandler{}
}

// Handle processes logging.
func (h *LoggingHandler) Handle(ctx context.Context, request interface{}) error {
	req, ok := request.(*HTTPRequest)
	if !ok {
		return errors.New("invalid request type")
	}

	fmt.Printf("[Logging] Request: %s %s (User: %s, Role: %s)\n",
		req.Method, req.Path, req.User, req.Role)

	// Pass to next handler
	return h.CallNext(ctx, request)
}

// BusinessLogicHandler processes the actual business logic.
type BusinessLogicHandler struct {
	BaseHandler
}

// NewBusinessLogicHandler creates a new business logic handler.
func NewBusinessLogicHandler() *BusinessLogicHandler {
	return &BusinessLogicHandler{}
}

// Handle processes business logic.
func (h *BusinessLogicHandler) Handle(ctx context.Context, request interface{}) error {
	req, ok := request.(*HTTPRequest)
	if !ok {
		return errors.New("invalid request type")
	}

	fmt.Printf("[BusinessLogic]  Processing request: %s %s\n", req.Method, req.Path)
	fmt.Printf("[BusinessLogic] Response: 200 OK\n")

	// This is the final handler, no need to call next
	return nil
}

// ExpenseRequest represents an expense approval request.
type ExpenseRequest struct {
	Amount      float64
	Description string
	Requester   string
	ApprovedBy  []string
}

// NewExpenseRequest creates a new expense request.
func NewExpenseRequest(amount float64, description, requester string) *ExpenseRequest {
	return &ExpenseRequest{
		Amount:      amount,
		Description: description,
		Requester:   requester,
		ApprovedBy:  make([]string, 0),
	}
}

// ApprovalHandler defines the interface for approval handlers.
type ApprovalHandler interface {
	Approve(ctx context.Context, request *ExpenseRequest) error
	SetNext(handler ApprovalHandler) ApprovalHandler
}

// BaseApprovalHandler provides common approval functionality.
type BaseApprovalHandler struct {
	next ApprovalHandler
}

// SetNext sets the next approval handler.
func (h *BaseApprovalHandler) SetNext(handler ApprovalHandler) ApprovalHandler {
	h.next = handler
	return handler
}

// CallNext calls the next approval handler.
func (h *BaseApprovalHandler) CallNext(ctx context.Context, request *ExpenseRequest) error {
	if h.next != nil {
		return h.next.Approve(ctx, request)
	}
	return nil
}

// ManagerApproval handles manager-level approvals.
type ManagerApproval struct {
	BaseApprovalHandler
	approvalLimit float64
}

// NewManagerApproval creates a new manager approval handler.
func NewManagerApproval(limit float64) *ManagerApproval {
	return &ManagerApproval{
		approvalLimit: limit,
	}
}

// Approve processes manager approval.
func (h *ManagerApproval) Approve(ctx context.Context, request *ExpenseRequest) error {
	if request.Amount <= h.approvalLimit {
		request.ApprovedBy = append(request.ApprovedBy, "Manager")
		fmt.Printf("[Manager]  Approved $%.2f for %s\n", request.Amount, request.Description)
		return nil // Approved, no need to continue chain
	}

	fmt.Printf("[Manager] Amount $%.2f exceeds limit $%.2f, escalating...\n",
		request.Amount, h.approvalLimit)

	// Escalate to next level
	return h.CallNext(ctx, request)
}

// DirectorApproval handles director-level approvals.
type DirectorApproval struct {
	BaseApprovalHandler
	approvalLimit float64
}

// NewDirectorApproval creates a new director approval handler.
func NewDirectorApproval(limit float64) *DirectorApproval {
	return &DirectorApproval{
		approvalLimit: limit,
	}
}

// Approve processes director approval.
func (h *DirectorApproval) Approve(ctx context.Context, request *ExpenseRequest) error {
	if request.Amount <= h.approvalLimit {
		request.ApprovedBy = append(request.ApprovedBy, "Director")
		fmt.Printf("[Director]  Approved $%.2f for %s\n", request.Amount, request.Description)
		return nil
	}

	fmt.Printf("[Director] Amount $%.2f exceeds limit $%.2f, escalating...\n",
		request.Amount, h.approvalLimit)

	// Escalate to next level
	return h.CallNext(ctx, request)
}

// CFOApproval handles CFO-level approvals.
type CFOApproval struct {
	BaseApprovalHandler
	approvalLimit float64
}

// NewCFOApproval creates a new CFO approval handler.
func NewCFOApproval(limit float64) *CFOApproval {
	return &CFOApproval{
		approvalLimit: limit,
	}
}

// Approve processes CFO approval.
func (h *CFOApproval) Approve(ctx context.Context, request *ExpenseRequest) error {
	if request.Amount <= h.approvalLimit {
		request.ApprovedBy = append(request.ApprovedBy, "CFO")
		fmt.Printf("[CFO]  Approved $%.2f for %s\n", request.Amount, request.Description)
		return nil
	}

	// CFO is the final authority
	return fmt.Errorf("amount $%.2f exceeds maximum approval limit $%.2f",
		request.Amount, h.approvalLimit)
}

// HTTPMiddleware is a function type for HTTP middleware.
type HTTPMiddleware func(http.HandlerFunc) http.HandlerFunc

// ChainMiddleware chains multiple middleware functions.
func ChainMiddleware(middlewares ...HTTPMiddleware) HTTPMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		// Build chain in reverse order
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// LoggingMiddleware logs HTTP requests.
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[HTTP Logging] %s %s\n", r.Method, r.URL.Path)
		next(w, r)
	}
}

// AuthMiddleware validates authentication.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// CORSMiddleware handles CORS headers.
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		next(w, r)
	}
}
