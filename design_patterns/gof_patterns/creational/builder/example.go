package builder

import (
	"fmt"
	"time"
)

// Example1_HTTPRequestBuilder demonstrates building HTTP requests.
func Example1_HTTPRequestBuilder() {
	fmt.Println("\n=== Example 1: HTTP Request Builder ===")

	// Simple GET request
	fmt.Println("\n1. Simple GET request:")
	req1, _ := NewHTTPRequestBuilder().
		URL("https://api.example.com/users").
		Build()
	fmt.Println(req1)

	// Complex POST request with all options
	fmt.Println("2. Complex POST request:")
	req2, _ := NewHTTPRequestBuilder().
		URL("https://api.example.com/users").
		Method("POST").
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer secret-token").
		Body(`{"name": "John Doe", "email": "john@example.com"}`).
		Timeout(60 * time.Second).
		RetryCount(5).
		FollowRedirects(false).
		Build()
	fmt.Println(req2)

	// Request with validation error
	fmt.Println("3. Request without URL (validation error):")
	_, err := NewHTTPRequestBuilder().
		Method("GET").
		Build()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// Example2_SQLQueryBuilder demonstrates building SQL queries.
func Example2_SQLQueryBuilder() {
	fmt.Println("\n=== Example 2: SQL Query Builder ===")

	// Simple query
	fmt.Println("\n1. Simple SELECT query:")
	query1, _ := NewQueryBuilder().
		Select("*").
		From("users").
		Build()
	fmt.Println(query1)

	// Complex query with all clauses
	fmt.Println("\n2. Complex query with WHERE, JOIN, ORDER BY, LIMIT:")
	query2, _ := NewQueryBuilder().
		Select("u.id", "u.name", "u.email", "p.title").
		From("users u").
		Join("INNER JOIN posts p ON u.id = p.user_id").
		Where("u.age > 18").
		Where("u.active = true").
		OrderBy("u.name ASC").
		Limit(10).
		Offset(20).
		Build()
	fmt.Println(query2)

	// Query with specific columns
	fmt.Println("\n3. Query selecting specific columns:")
	query3, _ := NewQueryBuilder().
		Select("id", "name", "email").
		From("users").
		Where("created_at > '2024-01-01'").
		OrderBy("created_at DESC").
		Limit(5).
		Build()
	fmt.Println(query3)
}

// Example3_EmailBuilder demonstrates building email messages.
func Example3_EmailBuilder() {
	fmt.Println("\n=== Example 3: Email Builder ===")

	// Simple email
	fmt.Println("\n1. Simple email:")
	email1, _ := NewEmailBuilder().
		From("sender@example.com").
		To("recipient@example.com").
		Subject("Hello").
		Body("This is a simple email.").
		Build()
	fmt.Println(email1)

	// Complex email with all options
	fmt.Println("\n2. Complex email with CC, BCC, attachments:")
	email2, _ := NewEmailBuilder().
		From("sender@example.com").
		To("recipient1@example.com", "recipient2@example.com").
		CC("cc@example.com").
		BCC("bcc@example.com").
		Subject("Monthly Report").
		Body("Please find the monthly report attached.").
		HTMLBody("<h1>Monthly Report</h1><p>Please find attached.</p>").
		Attachment("/path/to/report.pdf").
		Attachment("/path/to/data.xlsx").
		Priority("high").
		Build()
	fmt.Println(email2)

	// Email with validation error
	fmt.Println("\n3. Email without required fields (validation error):")
	_, err := NewEmailBuilder().
		From("sender@example.com").
		Subject("Test").
		Build()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// Example4_ComputerBuilder demonstrates building computers.
func Example4_ComputerBuilder() {
	fmt.Println("\n=== Example 4: Computer Builder ===")

	// Budget computer
	fmt.Println("\n1. Budget computer (minimal configuration):")
	budget, _ := NewComputerBuilder().
		CPU("Intel i3").
		RAM(8).
		Storage(256).
		Build()
	fmt.Println(budget)

	// Gaming computer
	fmt.Println("2. Gaming computer (high-end):")
	gaming, _ := NewComputerBuilder().
		CPU("Intel i9-13900K").
		RAM(32).
		Storage(2000).
		GPU("NVIDIA RTX 4090").
		OS("Windows 11").
		Monitor("ASUS ROG 27\" 4K 144Hz").
		Keyboard("Mechanical RGB").
		Mouse("Logitech G Pro").
		Build()
	fmt.Println(gaming)

	// Workstation
	fmt.Println("3. Professional workstation:")
	workstation, _ := NewComputerBuilder().
		CPU("AMD Ryzen 9 7950X").
		RAM(64).
		Storage(4000).
		GPU("NVIDIA RTX A6000").
		OS("Ubuntu 22.04").
		Monitor("Dell UltraSharp 32\" 4K").
		Keyboard("Logitech MX Keys").
		Mouse("Logitech MX Master 3").
		Build()
	fmt.Println(workstation)
}

// Example5_FluentInterface demonstrates the fluent interface pattern.
func Example5_FluentInterface() {
	fmt.Println("\n=== Example 5: Fluent Interface Demonstration ===")

	fmt.Println("\nFluent interface enables readable, self-documenting code:\n")

	// Show the fluent interface in action
	fmt.Println("Building an HTTP request with fluent interface:")
	fmt.Println("  NewHTTPRequestBuilder()")
	fmt.Println("    .URL(\"https://api.example.com/data\")")
	fmt.Println("    .Method(\"POST\")")
	fmt.Println("    .Header(\"Content-Type\", \"application/json\")")
	fmt.Println("    .Body(`{\"key\": \"value\"}`)")
	fmt.Println("    .Timeout(30 * time.Second)")
	fmt.Println("    .Build()")

	req, _ := NewHTTPRequestBuilder().
		URL("https://api.example.com/data").
		Method("POST").
		Header("Content-Type", "application/json").
		Body(`{"key": "value"}`).
		Timeout(30 * time.Second).
		Build()

	fmt.Println("\nResult:")
	fmt.Println(req)

	fmt.Println("Benefits:")
	fmt.Println("  ✓ Self-documenting code")
	fmt.Println("  ✓ Method chaining for readability")
	fmt.Println("  ✓ Optional parameters easy to add/remove")
	fmt.Println("  ✓ Order doesn't matter (except Build())")
}

// Example6_ValidationInBuilder demonstrates validation during building.
func Example6_ValidationInBuilder() {
	fmt.Println("\n=== Example 6: Validation in Builder ===")

	fmt.Println("\nBuilders validate before creating objects:\n")

	// Test various validation scenarios
	testCases := []struct {
		name    string
		builder func() error
	}{
		{
			"HTTP request without URL",
			func() error {
				_, err := NewHTTPRequestBuilder().Method("GET").Build()
				return err
			},
		},
		{
			"SQL query without table",
			func() error {
				_, err := NewQueryBuilder().Select("*").Build()
				return err
			},
		},
		{
			"Email without sender",
			func() error {
				_, err := NewEmailBuilder().
					To("recipient@example.com").
					Subject("Test").
					Body("Test body").
					Build()
				return err
			},
		},
		{
			"Email without recipients",
			func() error {
				_, err := NewEmailBuilder().
					From("sender@example.com").
					Subject("Test").
					Body("Test body").
					Build()
				return err
			},
		},
		{
			"Computer without CPU",
			func() error {
				_, err := NewComputerBuilder().RAM(16).Storage(512).Build()
				return err
			},
		},
	}

	for _, tc := range testCases {
		fmt.Printf("%s:\n", tc.name)
		err := tc.builder()
		if err != nil {
			fmt.Printf("  ✓ Validation caught error: %v\n", err)
		} else {
			fmt.Printf("  ✗ Should have failed validation\n")
		}
	}

	fmt.Println("\nValidation ensures objects are always in valid state!")
}

// Example7_BuilderPatternBenefits summarizes the benefits.
func Example7_BuilderPatternBenefits() {
	fmt.Println("\n=== Example 7: Builder Pattern Benefits ===")

	fmt.Println("\nKEY BENEFITS:")

	fmt.Println("\n1. HANDLES MANY OPTIONAL PARAMETERS")
	fmt.Println("   Before: NewHTTPRequest(url, method, body, contentType, ...20 params)")
	fmt.Println("   After:  NewHTTPRequestBuilder().URL(url).Method(method).Body(body).Build()")

	fmt.Println("\n2. READABLE AND SELF-DOCUMENTING")
	fmt.Println("   Each method call clearly shows what's being set")
	fmt.Println("   No need to remember parameter order")

	fmt.Println("\n3. IMMUTABILITY")
	fmt.Println("   Builder is mutable during construction")
	fmt.Println("   Product (result) can be immutable")

	fmt.Println("\n4. VALIDATION")
	fmt.Println("   Centralized validation in Build() method")
	fmt.Println("   Ensures objects are always in valid state")

	fmt.Println("\n5. DEFAULT VALUES")
	fmt.Println("   Easy to provide sensible defaults")
	fmt.Println("   Only override what you need to change")

	fmt.Println("\n6. FLEXIBILITY")
	fmt.Println("   Add/remove parameters without breaking existing code")
	fmt.Println("   No need for telescoping constructors")

	fmt.Println("\nREAL-WORLD USE CASES:")
	fmt.Println("  • HTTP/API clients and requests")
	fmt.Println("  • SQL query builders (ORMs)")
	fmt.Println("  • Email/message builders")
	fmt.Println("  • Configuration objects")
	fmt.Println("  • Test data builders")
	fmt.Println("  • Complex UI components")
	fmt.Println("  • Document builders (HTML, XML, JSON)")

	fmt.Println("\nWHEN TO USE:")
	fmt.Println("  ✓ Object has 4+ optional parameters")
	fmt.Println("  ✓ Construction is complex or multi-step")
	fmt.Println("  ✓ Need immutable objects with complex initialization")
	fmt.Println("  ✓ Want readable, self-documenting code")

	fmt.Println("\nWHEN NOT TO USE:")
	fmt.Println("  ✗ Simple objects with few parameters")
	fmt.Println("  ✗ No optional parameters")
	fmt.Println("  ✗ Performance is critical (slight overhead)")

	// Demonstrate
	fmt.Println("\n\nPRACTICAL EXAMPLE:")
	fmt.Println("Building a complex email in one fluent chain:\n")

	email, _ := NewEmailBuilder().
		From("system@company.com").
		To("user1@example.com", "user2@example.com").
		CC("manager@example.com").
		Subject("Weekly Report - Q4 2024").
		HTMLBody("<h2>Weekly Report</h2><p>Attached is the report.</p>").
		Attachment("/reports/weekly_q4.pdf").
		Priority("high").
		Build()

	fmt.Println(email)

	fmt.Println("This would be painful with a traditional constructor!")
}
