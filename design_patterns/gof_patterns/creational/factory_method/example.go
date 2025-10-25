package factory_method

import "fmt"

// Example1_NotificationFactory demonstrates creating different notifiers.
func Example1_NotificationFactory() {
	fmt.Println("\n=== Example 1: Notification Factory ===")

	config := NotificationConfig{
		SMTPServer:   "smtp.example.com",
		SMSGateway:   "twilio.com",
		PushService:  "firebase",
		SlackWebhook: "https://hooks.slack.com/...",
		From:         "noreply@example.com",
		AppID:        "my-app-123",
		SlackChannel: "alerts",
	}

	// Create different notifiers using the factory
	channels := []string{"email", "sms", "push", "slack"}

	fmt.Println("\nCreating notifiers for different channels:")
	for _, channel := range channels {
		notifier, err := NewNotifier(channel, config)
		if err != nil {
			fmt.Printf("Error creating %s notifier: %v\n", channel, err)
			continue
		}

		fmt.Printf("\nCreated %s notifier\n", notifier.GetChannel())
		err = notifier.Send("System alert: CPU usage high!")
		if err != nil {
			fmt.Printf("Error sending notification: %v\n", err)
		}
	}

	// Try an unsupported channel
	fmt.Println("\n\nTrying unsupported channel:")
	_, err := NewNotifier("telegram", config)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}
}

// Example2_ParserFactory demonstrates creating parsers based on file type.
func Example2_ParserFactory() {
	fmt.Println("\n=== Example 2: Parser Factory ===")

	files := []string{
		"config.json",
		"data.xml",
		"settings.yaml",
		"values.yml",
		"unknown.txt",
	}

	fmt.Println("\nParsing different file formats:")
	for _, filename := range files {
		fmt.Printf("\n%s:\n", filename)

		parser, err := NewParser(filename)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}

		fmt.Printf("  Created %s parser\n", parser.GetFormat())

		// Simulate parsing
		data := []byte("sample data for " + filename)
		result, err := parser.Parse(data)
		if err != nil {
			fmt.Printf("  Parse error: %v\n", err)
			continue
		}

		fmt.Printf("  Parsed result: %v\n", result)
	}
}

// Example3_HTTPClientFactory demonstrates environment-based client creation.
func Example3_HTTPClientFactory() {
	fmt.Println("\n=== Example 3: HTTP Client Factory ===")

	environments := []string{"production", "development", "test"}

	fmt.Println("\nCreating HTTP clients for different environments:")
	for _, env := range environments {
		fmt.Printf("\n%s environment:\n", env)

		client := NewHTTPClient(env)
		fmt.Printf("  Client type: %s\n", client.GetEnvironment())

		// Make requests
		client.Get("https://api.example.com/users")
		client.Post("https://api.example.com/data", `{"key": "value"}`)
	}
}

// Example4_RuntimeDecision demonstrates runtime type selection.
func Example4_RuntimeDecision() {
	fmt.Println("\n=== Example 4: Runtime Decision Making ===")

	// Simulate user preferences
	userPreferences := map[string]string{
		"user1": "email",
		"user2": "sms",
		"user3": "slack",
		"user4": "push",
	}

	config := NotificationConfig{
		SMTPServer:   "smtp.example.com",
		SMSGateway:   "twilio.com",
		PushService:  "firebase",
		SlackWebhook: "https://hooks.slack.com/...",
		From:         "alerts@example.com",
		AppID:        "app-456",
		SlackChannel: "notifications",
	}

	fmt.Println("\nSending notifications based on user preferences:")
	for user, preferredChannel := range userPreferences {
		fmt.Printf("\n%s prefers %s:\n", user, preferredChannel)

		notifier, err := NewNotifier(preferredChannel, config)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}

		message := fmt.Sprintf("Hello %s! You have a new message.", user)
		err = notifier.Send(message)
		if err != nil {
			fmt.Printf("  Send error: %v\n", err)
		}
	}
}

// Example5_FactoryWithValidation demonstrates proper error handling.
func Example5_FactoryWithValidation() {
	fmt.Println("\n=== Example 5: Factory with Validation ===")

	testCases := []struct {
		name     string
		channel  string
		filename string
	}{
		{"Valid email", "email", ""},
		{"Valid SMS", "sms", ""},
		{"Invalid channel", "fax", ""},
		{"Case insensitive", "EMAIL", ""},
		{"Empty channel", "", ""},
	}

	config := NotificationConfig{
		SMTPServer: "smtp.example.com",
		SMSGateway: "twilio.com",
		From:       "test@example.com",
	}

	fmt.Println("\nTesting notifier factory with various inputs:")
	for _, tc := range testCases {
		fmt.Printf("\n%s (channel: '%s'):\n", tc.name, tc.channel)

		notifier, err := NewNotifier(tc.channel, config)
		if err != nil {
			fmt.Printf("  ✗ Error (expected for invalid inputs): %v\n", err)
			continue
		}

		fmt.Printf("  ✓ Successfully created %s notifier\n", notifier.GetChannel())
	}

	// Test parser factory
	fmt.Println("\n\nTesting parser factory with various inputs:")
	parserTests := []string{
		"valid.json",
		"valid.xml",
		"valid.yaml",
		"noextension",
		"multiple.dots.in.name.yml",
	}

	for _, filename := range parserTests {
		fmt.Printf("\n%s:\n", filename)

		parser, err := NewParser(filename)
		if err != nil {
			fmt.Printf("  ✗ Error: %v\n", err)
			continue
		}

		fmt.Printf("  ✓ Successfully created %s parser\n", parser.GetFormat())
	}
}

// Example6_ExtensibilityDemo shows how easy it is to add new types.
func Example6_ExtensibilityDemo() {
	fmt.Println("\n=== Example 6: Extensibility Demonstration ===")

	fmt.Println("\nKey Benefits of Factory Method Pattern:")
	fmt.Println("\n1. CENTRALIZED CREATION LOGIC")
	fmt.Println("   - All object creation in one place (factory function)")
	fmt.Println("   - Easy to find and modify")
	fmt.Println("   - Consistent initialization across application")

	fmt.Println("\n2. EASY TO EXTEND")
	fmt.Println("   - Adding a new type requires:")
	fmt.Println("     a) Create new concrete type implementing interface")
	fmt.Println("     b) Add case to factory switch statement")
	fmt.Println("   - Client code unchanged!")

	fmt.Println("\n3. TYPE SAFETY")
	fmt.Println("   - Factory returns interface type")
	fmt.Println("   - Compiler ensures all concrete types implement interface")
	fmt.Println("   - Runtime errors for invalid types (better than silent failures)")

	fmt.Println("\n4. DEPENDENCY INVERSION")
	fmt.Println("   - Clients depend on Notifier interface")
	fmt.Println("   - Not on EmailNotifier, SMSNotifier, etc.")
	fmt.Println("   - Loose coupling = easier testing and maintenance")

	// Demonstrate
	fmt.Println("\n\nPractical Demonstration:")
	config := NotificationConfig{
		SMTPServer:   "smtp.example.com",
		SlackWebhook: "https://hooks.slack.com/...",
		SlackChannel: "engineering",
		From:         "system@example.com",
	}

	// Client code stays the same regardless of notification type
	channels := []string{"email", "slack"}
	for _, ch := range channels {
		notifier, _ := NewNotifier(ch, config)
		notifier.Send("Factory pattern makes this code clean and extensible!")
	}

	fmt.Println("\nConclusion:")
	fmt.Println("Factory Method pattern provides a flexible, maintainable way")
	fmt.Println("to create objects while keeping client code clean and simple.")
}
