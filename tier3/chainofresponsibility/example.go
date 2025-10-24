package chainofresponsibility

import (
	"context"
	"fmt"
	"time"
)

// Example1_HTTPMiddlewareChain demonstrates HTTP request processing pipeline.
func Example1_HTTPMiddlewareChain() {
	fmt.Println("=== Example 1: HTTP Middleware Chain ===")

	// Build the middleware chain
	auth := NewAuthenticationHandler()
	authz := NewAuthorizationHandler("admin")
	validation := NewValidationHandler()
	logging := NewLoggingHandler()
	business := NewBusinessLogicHandler()

	// Chain handlers together
	auth.SetNext(authz).
		SetNext(validation).
		SetNext(logging).
		SetNext(business)

	// Create a valid request
	fmt.Println("--- Valid Request ---")
	request := NewHTTPRequest("POST", "/api/users")
	request.Headers["Authorization"] = "Bearer admin-token"
	request.Role = "admin"

	ctx := context.Background()
	if err := auth.Handle(ctx, request); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n--- Invalid Request (No Auth) ---")
	invalidRequest := NewHTTPRequest("GET", "/api/users")
	invalidRequest.Role = "admin"

	if err := auth.Handle(ctx, invalidRequest); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println()
}

// Example2_RateLimiting demonstrates rate limiting middleware.
func Example2_RateLimiting() {
	fmt.Println("=== Example 2: Rate Limiting ===")

	// Build chain with rate limiting
	auth := NewAuthenticationHandler()
	rateLimit := NewRateLimitHandler(3, 1*time.Second)
	business := NewBusinessLogicHandler()

	auth.SetNext(rateLimit).SetNext(business)

	// Create user request
	request := NewHTTPRequest("GET", "/api/data")
	request.Headers["Authorization"] = "Bearer user123"
	request.Role = "user"

	ctx := context.Background()

	// Make multiple requests
	fmt.Println("Making 5 requests rapidly...")
	for i := 1; i <= 5; i++ {
		fmt.Printf("\nRequest #%d:\n", i)
		if err := auth.Handle(ctx, request); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println()
}

// Example3_ExpenseApproval demonstrates approval workflow chain.
func Example3_ExpenseApproval() {
	fmt.Println("=== Example 3: Expense Approval Workflow ===")

	// Build approval chain
	manager := NewManagerApproval(1000.00)
	director := NewDirectorApproval(10000.00)
	cfo := NewCFOApproval(100000.00)

	manager.SetNext(director).SetNext(cfo)

	ctx := context.Background()

	// Test different expense amounts
	expenses := []struct {
		amount      float64
		description string
	}{
		{500.00, "Office Supplies"},
		{5000.00, "New Laptops"},
		{50000.00, "Marketing Campaign"},
		{150000.00, "Office Renovation"},
	}

	for _, exp := range expenses {
		fmt.Printf("--- Expense Request: $%.2f for %s ---\n", exp.amount, exp.description)

		request := NewExpenseRequest(exp.amount, exp.description, "John Doe")

		if err := manager.Approve(ctx, request); err != nil {
			fmt.Printf("L Rejected: %v\n", err)
		} else {
			fmt.Printf(" Approved by: %v\n", request.ApprovedBy)
		}

		fmt.Println()
	}
}

// Example4_ValidationChain demonstrates multiple validators.
func Example4_ValidationChain() {
	fmt.Println("=== Example 4: Validation Chain ===")

	// Create validation chain
	auth := NewAuthenticationHandler()
	validation := NewValidationHandler()
	authz := NewAuthorizationHandler("user")

	auth.SetNext(validation).SetNext(authz)

	ctx := context.Background()

	testCases := []struct {
		name    string
		request *HTTPRequest
	}{
		{
			name: "Valid Request",
			request: func() *HTTPRequest {
				r := NewHTTPRequest("GET", "/api/profile")
				r.Headers["Authorization"] = "Bearer valid-token"
				r.Role = "user"
				return r
			}(),
		},
		{
			name: "Missing Auth",
			request: func() *HTTPRequest {
				r := NewHTTPRequest("GET", "/api/profile")
				r.Role = "user"
				return r
			}(),
		},
		{
			name: "Invalid Method",
			request: func() *HTTPRequest {
				r := NewHTTPRequest("INVALID", "/api/profile")
				r.Headers["Authorization"] = "Bearer valid-token"
				r.Role = "user"
				return r
			}(),
		},
		{
			name: "Wrong Role",
			request: func() *HTTPRequest {
				r := NewHTTPRequest("GET", "/api/profile")
				r.Headers["Authorization"] = "Bearer valid-token"
				r.Role = "admin"
				return r
			}(),
		},
	}

	for _, tc := range testCases {
		fmt.Printf("--- %s ---\n", tc.name)
		if err := auth.Handle(ctx, tc.request); err != nil {
			fmt.Printf("L Failed: %v\n", err)
		} else {
			fmt.Printf(" Passed all validations\n")
		}
		fmt.Println()
	}
}

// Example5_ConfigurableChain demonstrates dynamic chain configuration.
func Example5_ConfigurableChain() {
	fmt.Println("=== Example 5: Configurable Chain ===")

	business := NewBusinessLogicHandler()

	// Configuration 1: Public endpoint (no auth)
	fmt.Println("--- Configuration 1: Public Endpoint ---")
	logging1 := NewLoggingHandler()
	validation1 := NewValidationHandler()
	logging1.SetNext(validation1).SetNext(business)

	request1 := NewHTTPRequest("GET", "/api/public/health")
	request1.User = "anonymous"
	request1.Role = "public"

	ctx := context.Background()
	logging1.Handle(ctx, request1)

	// Configuration 2: Protected endpoint (auth + authz)
	fmt.Println("\n--- Configuration 2: Protected Endpoint ---")
	auth2 := NewAuthenticationHandler()
	authz2 := NewAuthorizationHandler("admin")
	logging2 := NewLoggingHandler()
	validation2 := NewValidationHandler()

	auth2.SetNext(authz2).
		SetNext(logging2).
		SetNext(validation2).
		SetNext(business)

	request2 := NewHTTPRequest("DELETE", "/api/admin/users/123")
	request2.Headers["Authorization"] = "Bearer admin-token"
	request2.Role = "admin"

	if err := auth2.Handle(ctx, request2); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println()
}

// Example6_ShortCircuiting demonstrates early termination in chain.
func Example6_ShortCircuiting() {
	fmt.Println("=== Example 6: Short-Circuiting Chain ===")

	// Build chain
	auth := NewAuthenticationHandler()
	authz := NewAuthorizationHandler("admin")
	logging := NewLoggingHandler()
	business := NewBusinessLogicHandler()

	auth.SetNext(authz).
		SetNext(logging).
		SetNext(business)

	ctx := context.Background()

	// Request that fails at authorization
	fmt.Println("--- Request failing at authorization ---")
	request := NewHTTPRequest("POST", "/api/admin/delete-all")
	request.Headers["Authorization"] = "Bearer user-token"
	request.Role = "user" // Not admin

	fmt.Println("Chain: Auth -> Authz -> Logging -> Business")
	fmt.Println("Expected: Stop at Authz (role mismatch)")

	if err := auth.Handle(ctx, request); err != nil {
		fmt.Printf("L Chain stopped: %v\n", err)
		fmt.Println("Notice: Logging and Business handlers never executed")
	}

	fmt.Println()
}

// Example7_MultipleHandlerTypes demonstrates different handler types in one chain.
func Example7_MultipleHandlerTypes() {
	fmt.Println("=== Example 7: Mixed Handler Types ===")

	// Build comprehensive chain
	auth := NewAuthenticationHandler()
	rateLimit := NewRateLimitHandler(5, 10*time.Second)
	authz := NewAuthorizationHandler("")
	validation := NewValidationHandler()
	logging := NewLoggingHandler()
	business := NewBusinessLogicHandler()

	// Chain all handlers
	auth.SetNext(rateLimit).
		SetNext(authz).
		SetNext(validation).
		SetNext(logging).
		SetNext(business)

	ctx := context.Background()

	// Process multiple requests
	requests := []*HTTPRequest{
		{
			Method:  "GET",
			Path:    "/api/products",
			Headers: map[string]string{"Authorization": "Bearer customer-token"},
			Role:    "customer",
		},
		{
			Method:  "POST",
			Path:    "/api/orders",
			Headers: map[string]string{"Authorization": "Bearer customer-token"},
			Role:    "customer",
		},
		{
			Method:  "PUT",
			Path:    "/api/profile",
			Headers: map[string]string{"Authorization": "Bearer customer-token"},
			Role:    "customer",
		},
	}

	for i, req := range requests {
		fmt.Printf("--- Request #%d: %s %s ---\n", i+1, req.Method, req.Path)
		if err := auth.Handle(ctx, req); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println()
		time.Sleep(100 * time.Millisecond)
	}
}
