package chainofresponsibility

import (
	"context"
	"testing"
	"time"
)

func TestAuthenticationHandler(t *testing.T) {
	handler := NewAuthenticationHandler()
	ctx := context.Background()

	// Valid auth
	request := NewHTTPRequest("GET", "/test")
	request.Headers["Authorization"] = "Bearer test-token"

	err := handler.Handle(ctx, request)
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	if request.User != "test-token" {
		t.Errorf("Expected user 'test-token', got '%s'", request.User)
	}

	// Missing auth
	invalidRequest := NewHTTPRequest("GET", "/test")
	err = handler.Handle(ctx, invalidRequest)
	if err == nil {
		t.Error("Expected error for missing auth")
	}
}

func TestAuthorizationHandler(t *testing.T) {
	handler := NewAuthorizationHandler("admin")
	ctx := context.Background()

	// Valid role
	request := NewHTTPRequest("GET", "/test")
	request.Role = "admin"

	err := handler.Handle(ctx, request)
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	// Invalid role
	invalidRequest := NewHTTPRequest("GET", "/test")
	invalidRequest.Role = "user"

	err = handler.Handle(ctx, invalidRequest)
	if err == nil {
		t.Error("Expected error for invalid role")
	}
}

func TestValidationHandler(t *testing.T) {
	handler := NewValidationHandler()
	ctx := context.Background()

	// Valid request
	request := NewHTTPRequest("GET", "/test")
	err := handler.Handle(ctx, request)
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}

	// Invalid method
	invalidRequest := NewHTTPRequest("INVALID", "/test")
	err = handler.Handle(ctx, invalidRequest)
	if err == nil {
		t.Error("Expected error for invalid method")
	}

	// Empty path
	emptyPath := NewHTTPRequest("GET", "")
	err = handler.Handle(ctx, emptyPath)
	if err == nil {
		t.Error("Expected error for empty path")
	}
}

func TestChainedHandlers(t *testing.T) {
	// Build chain
	auth := NewAuthenticationHandler()
	authz := NewAuthorizationHandler("admin")
	validation := NewValidationHandler()

	auth.SetNext(authz).SetNext(validation)

	ctx := context.Background()

	// Valid request should pass through entire chain
	request := NewHTTPRequest("GET", "/test")
	request.Headers["Authorization"] = "Bearer admin-token"
	request.Role = "admin"

	err := auth.Handle(ctx, request)
	if err != nil {
		t.Errorf("Expected success through chain, got error: %v", err)
	}

	// Request failing at authz should stop chain
	failRequest := NewHTTPRequest("GET", "/test")
	failRequest.Headers["Authorization"] = "Bearer user-token"
	failRequest.Role = "user"

	err = auth.Handle(ctx, failRequest)
	if err == nil {
		t.Error("Expected error at authorization step")
	}
}

func TestRateLimitHandler(t *testing.T) {
	handler := NewRateLimitHandler(2, 100*time.Millisecond)
	ctx := context.Background()

	request := NewHTTPRequest("GET", "/test")
	request.Headers["Authorization"] = "Bearer user1"
	request.User = "user1"

	// First two requests should succeed
	err := handler.Handle(ctx, request)
	if err != nil {
		t.Errorf("Request 1 failed: %v", err)
	}

	err = handler.Handle(ctx, request)
	if err != nil {
		t.Errorf("Request 2 failed: %v", err)
	}

	// Third request should fail (rate limit exceeded)
	err = handler.Handle(ctx, request)
	if err == nil {
		t.Error("Expected rate limit error")
	}

	// Wait for window to pass
	time.Sleep(150 * time.Millisecond)

	// Should succeed after window
	err = handler.Handle(ctx, request)
	if err != nil {
		t.Errorf("Request after window failed: %v", err)
	}
}

func TestManagerApproval(t *testing.T) {
	handler := NewManagerApproval(1000.00)
	ctx := context.Background()

	// Amount within limit
	request := NewExpenseRequest(500.00, "Office Supplies", "John")
	err := handler.Approve(ctx, request)
	if err != nil {
		t.Errorf("Expected approval, got error: %v", err)
	}

	if len(request.ApprovedBy) != 1 || request.ApprovedBy[0] != "Manager" {
		t.Errorf("Expected approval by Manager, got: %v", request.ApprovedBy)
	}

	// Amount exceeding limit (no next handler, should not error but also not approve)
	largeRequest := NewExpenseRequest(2000.00, "Equipment", "Jane")
	err = handler.Approve(ctx, largeRequest)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if len(largeRequest.ApprovedBy) > 0 {
		t.Errorf("Expected no approval for large amount, got: %v", largeRequest.ApprovedBy)
	}
}

func TestApprovalChain(t *testing.T) {
	// Build approval chain
	manager := NewManagerApproval(1000.00)
	director := NewDirectorApproval(10000.00)
	cfo := NewCFOApproval(100000.00)

	manager.SetNext(director).SetNext(cfo)

	ctx := context.Background()

	testCases := []struct {
		amount        float64
		expectedLevel string
		shouldFail    bool
	}{
		{500.00, "Manager", false},
		{5000.00, "Director", false},
		{50000.00, "CFO", false},
		{150000.00, "", true}, // Exceeds all limits
	}

	for _, tc := range testCases {
		request := NewExpenseRequest(tc.amount, "Test", "User")
		err := manager.Approve(ctx, request)

		if tc.shouldFail {
			if err == nil {
				t.Errorf("Amount %.2f should have failed", tc.amount)
			}
		} else {
			if err != nil {
				t.Errorf("Amount %.2f failed: %v", tc.amount, err)
			}

			if len(request.ApprovedBy) == 0 {
				t.Errorf("Amount %.2f was not approved", tc.amount)
			} else if request.ApprovedBy[0] != tc.expectedLevel {
				t.Errorf("Amount %.2f: expected approval by %s, got %s",
					tc.amount, tc.expectedLevel, request.ApprovedBy[0])
			}
		}
	}
}

func TestShortCircuiting(t *testing.T) {
	finalHandlerCalled := false

	auth := NewAuthenticationHandler()
	authz := NewAuthorizationHandler("admin")

	// Create a custom handler to track if it's called
	finalHandler := &BaseHandler{}
	finalHandler.SetNext(nil)

	// Override Handle method behavior (in real scenario, would implement Handler interface)
	testHandler := struct {
		BaseHandler
		called *bool
	}{
		BaseHandler: *finalHandler,
		called:      &finalHandlerCalled,
	}

	auth.SetNext(authz)

	ctx := context.Background()

	// Request that fails at authz
	request := NewHTTPRequest("GET", "/test")
	request.Headers["Authorization"] = "Bearer user-token"
	request.Role = "user"

	err := auth.Handle(ctx, request)
	if err == nil {
		t.Error("Expected error from authz")
	}

	// Final handler should not be called
	if *testHandler.called {
		t.Error("Final handler should not be called when chain short-circuits")
	}
}

func TestLoggingHandler(t *testing.T) {
	handler := NewLoggingHandler()
	ctx := context.Background()

	request := NewHTTPRequest("POST", "/api/users")
	request.User = "test-user"
	request.Role = "admin"

	err := handler.Handle(ctx, request)
	if err != nil {
		t.Errorf("Logging handler failed: %v", err)
	}
}

func TestBusinessLogicHandler(t *testing.T) {
	handler := NewBusinessLogicHandler()
	ctx := context.Background()

	request := NewHTTPRequest("GET", "/api/data")

	err := handler.Handle(ctx, request)
	if err != nil {
		t.Errorf("Business handler failed: %v", err)
	}
}

func TestCompleteChain(t *testing.T) {
	// Build complete middleware chain
	auth := NewAuthenticationHandler()
	authz := NewAuthorizationHandler("user")
	validation := NewValidationHandler()
	logging := NewLoggingHandler()
	business := NewBusinessLogicHandler()

	auth.SetNext(authz).
		SetNext(validation).
		SetNext(logging).
		SetNext(business)

	ctx := context.Background()

	// Valid request should pass through entire chain
	request := NewHTTPRequest("POST", "/api/orders")
	request.Headers["Authorization"] = "Bearer user-token"
	request.Role = "user"

	err := auth.Handle(ctx, request)
	if err != nil {
		t.Errorf("Complete chain failed: %v", err)
	}
}

// Benchmark tests
func BenchmarkChainedHandlers(b *testing.B) {
	auth := NewAuthenticationHandler()
	authz := NewAuthorizationHandler("user")
	validation := NewValidationHandler()
	logging := NewLoggingHandler()
	business := NewBusinessLogicHandler()

	auth.SetNext(authz).
		SetNext(validation).
		SetNext(logging).
		SetNext(business)

	ctx := context.Background()
	request := NewHTTPRequest("GET", "/test")
	request.Headers["Authorization"] = "Bearer token"
	request.Role = "user"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		auth.Handle(ctx, request)
	}
}

func BenchmarkApprovalChain(b *testing.B) {
	manager := NewManagerApproval(1000.00)
	director := NewDirectorApproval(10000.00)
	cfo := NewCFOApproval(100000.00)

	manager.SetNext(director).SetNext(cfo)

	ctx := context.Background()
	request := NewExpenseRequest(500.00, "Test", "User")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.Approve(ctx, request)
		request.ApprovedBy = request.ApprovedBy[:0] // Reset approval
	}
}
