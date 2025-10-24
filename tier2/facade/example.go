package facade

import (
	"fmt"
	"strings"
)

// Example1_BankingSystem demonstrates the Facade Pattern with a banking system.
// Without the facade, clients would need to interact with multiple subsystems:
// AccountService, TransactionService, FraudDetectionService, NotificationService.
// The facade simplifies this to a single TransferMoney() call.
func Example1_BankingSystem() {
	fmt.Println("=== Example 1: Banking System Facade ===")
	fmt.Println("\nWithout Facade: Client must coordinate multiple subsystems manually")
	fmt.Println("  1. Validate source account")
	fmt.Println("  2. Validate destination account")
	fmt.Println("  3. Check for fraud")
	fmt.Println("  4. Create transaction")
	fmt.Println("  5. Debit source account")
	fmt.Println("  6. Credit destination account")
	fmt.Println("  7. Process transaction")
	fmt.Println("  8. Send notifications")
	fmt.Println("\nWith Facade: Single method call handles all complexity")

	// Create the banking facade
	bank := NewBankingFacade()

	// Complex operation simplified to a single call
	err := bank.TransferMoney("ACC-12345", "ACC-67890", 500.00)
	if err != nil {
		fmt.Printf("❌ Transfer failed: %v\n", err)
		return
	}

	// Another transfer
	fmt.Println()
	err = bank.TransferMoney("ACC-67890", "ACC-12345", 250.00)
	if err != nil {
		fmt.Printf("❌ Transfer failed: %v\n", err)
	}

	fmt.Println("\n" + strings.Repeat("─", 50))
}

// Example2_SmartHome demonstrates the Facade Pattern with home automation.
// Complex home automation scenarios (involving lighting, security, climate, entertainment)
// are simplified to single method calls like LeaveHome(), ArriveHome(), MovieNight().
func Example2_SmartHome() {
	fmt.Println("=== Example 2: Smart Home Automation Facade ===")
	fmt.Println("\nWithout Facade: Must control each subsystem individually")
	fmt.Println("With Facade: Scene-based automation simplifies everything")

	home := NewSmartHomeFacade()

	// Scenario 1: Leaving home
	home.LeaveHome()

	// Scenario 2: Arriving home
	home.ArriveHome()

	// Scenario 3: Movie night
	home.MovieNight("The Matrix")

	// Scenario 4: Going to sleep
	home.SleepMode()

	fmt.Println("\n" + strings.Repeat("─", 50))
}

// Example3_ECommerce demonstrates the Facade Pattern with e-commerce order processing.
// Order placement involves inventory, payment, shipping, and notifications.
// The facade orchestrates all these subsystems transparently.
func Example3_ECommerce() {
	fmt.Println("=== Example 3: E-Commerce Order Processing Facade ===")
	fmt.Println("\nWithout Facade: Client orchestrates inventory, payment, shipping, email")
	fmt.Println("With Facade: Single PlaceOrder() call handles entire workflow")

	orderSystem := NewOrderProcessingFacade()

	// Place an order - single call handles all complexity
	err := orderSystem.PlaceOrder(
		"customer@example.com",
		"PROD-123",
		2,
		199.99,
		"123 Main St, City, State 12345",
	)

	if err != nil {
		fmt.Printf("❌ Order failed: %v\n", err)
		return
	}

	fmt.Println("\n" + strings.Repeat("─", 50))
}

// Example4_ComparisonWithoutFacade demonstrates the complexity WITHOUT using a facade.
// This shows why facades are valuable by contrasting with manual subsystem coordination.
func Example4_ComparisonWithoutFacade() {
	fmt.Println("=== Example 4: Comparison - With vs Without Facade ===")

	fmt.Println("\n❌ WITHOUT FACADE (Complex, Error-Prone):")
	fmt.Println("```go")
	fmt.Println("// Client must know about and coordinate ALL subsystems")
	fmt.Println("accountSvc := &AccountService{}")
	fmt.Println("txnSvc := &TransactionService{}")
	fmt.Println("fraudSvc := &FraudDetectionService{}")
	fmt.Println("notifSvc := &NotificationService{}")
	fmt.Println()
	fmt.Println("// Step 1: Validate accounts")
	fmt.Println("valid, err := accountSvc.ValidateAccount(from, amount)")
	fmt.Println("if !valid || err != nil { /* handle error */ }")
	fmt.Println()
	fmt.Println("valid, err = accountSvc.ValidateAccount(to, 0)")
	fmt.Println("if !valid || err != nil { /* handle error */ }")
	fmt.Println()
	fmt.Println("// Step 2: Check fraud")
	fmt.Println("_, err = fraudSvc.CheckFraud(from, to, amount)")
	fmt.Println("if err != nil { /* handle error */ }")
	fmt.Println()
	fmt.Println("// Step 3: Create transaction")
	fmt.Println("txnID, err := txnSvc.CreateTransaction(from, to, amount)")
	fmt.Println("if err != nil { /* handle error */ }")
	fmt.Println()
	fmt.Println("// Step 4: Debit and credit")
	fmt.Println("if err := accountSvc.DebitAccount(from, amount); err != nil {")
	fmt.Println("    /* handle error and rollback */ }")
	fmt.Println()
	fmt.Println("if err := accountSvc.CreditAccount(to, amount); err != nil {")
	fmt.Println("    /* handle error and rollback debit */ }")
	fmt.Println()
	fmt.Println("// Step 5: Process transaction")
	fmt.Println("if err := txnSvc.ProcessTransaction(txnID); err != nil {")
	fmt.Println("    /* handle error and rollback */ }")
	fmt.Println()
	fmt.Println("// Step 6: Send notifications")
	fmt.Println("notifSvc.SendNotification(from, \"...\")")
	fmt.Println("notifSvc.SendNotification(to, \"...\")")
	fmt.Println("```")
	fmt.Println("\n❗ Problems:")
	fmt.Println("  • Client must know internal details of all subsystems")
	fmt.Println("  • Tight coupling between client and implementation")
	fmt.Println("  • Error-prone: easy to forget steps or get order wrong")
	fmt.Println("  • Code duplication across all clients")
	fmt.Println("  • Hard to maintain and test")

	fmt.Println("\n✅ WITH FACADE (Simple, Robust):")
	fmt.Println("```go")
	fmt.Println("// Client only knows about the facade")
	fmt.Println("bank := NewBankingFacade()")
	fmt.Println()
	fmt.Println("// Single method call handles ALL complexity")
	fmt.Println("err := bank.TransferMoney(\"ACC-12345\", \"ACC-67890\", 500.00)")
	fmt.Println("if err != nil {")
	fmt.Println("    // Simple error handling")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println("\n✨ Benefits:")
	fmt.Println("  • Client only depends on facade interface")
	fmt.Println("  • Loose coupling - internal changes don't affect clients")
	fmt.Println("  • Simple, intuitive API")
	fmt.Println("  • No code duplication")
	fmt.Println("  • Easy to test and maintain")
	fmt.Println("  • Enforces correct order of operations")

	fmt.Println("\n" + strings.Repeat("─", 50))
}

// Example5_RealWorld demonstrates a real-world scenario: Video encoding system.
// Encoding video involves codec selection, format conversion, quality settings,
// metadata extraction, thumbnail generation, and more. The facade hides this complexity.
func Example5_RealWorld() {
	fmt.Println("=== Example 5: Real-World - Video Encoding System ===")
	fmt.Println("\nScenario: Convert and optimize video for web streaming")
	fmt.Println("\nSubsystems involved:")
	fmt.Println("  • Codec Service (H.264, VP9, AV1)")
	fmt.Println("  • Format Converter (MP4, WebM, HLS)")
	fmt.Println("  • Quality Optimizer (bitrate, resolution)")
	fmt.Println("  • Metadata Extractor (duration, dimensions)")
	fmt.Println("  • Thumbnail Generator")
	fmt.Println("  • CDN Upload Service")
	fmt.Println("  • Database Service (store video info)")

	// Without facade, client would need:
	fmt.Println("\n❌ Without Facade:")
	fmt.Println("  1. Select appropriate codec based on browser support")
	fmt.Println("  2. Convert to target format")
	fmt.Println("  3. Calculate optimal bitrate for quality/size balance")
	fmt.Println("  4. Extract metadata (duration, resolution, etc.)")
	fmt.Println("  5. Generate thumbnails at key timestamps")
	fmt.Println("  6. Upload encoded video to CDN")
	fmt.Println("  7. Upload thumbnails to CDN")
	fmt.Println("  8. Store video metadata in database")
	fmt.Println("  9. Return video URLs to client")

	fmt.Println("\n✅ With Facade:")
	fmt.Println("```go")
	fmt.Println("encoder := NewVideoEncodingFacade()")
	fmt.Println("result, err := encoder.EncodeForWeb(\"input.mov\", VideoQualityHD)")
	fmt.Println("if err != nil { /* handle error */ }")
	fmt.Println("fmt.Println(\"Video ready:\", result.VideoURL)")
	fmt.Println("```")

	fmt.Println("\n🎬 Simulated execution:")
	fmt.Println("  ✓ Analyzing video: input.mov")
	fmt.Println("  ✓ Selected codec: H.264 (broad compatibility)")
	fmt.Println("  ✓ Converting to MP4 format")
	fmt.Println("  ✓ Optimizing bitrate: 5000 kbps for HD quality")
	fmt.Println("  ✓ Extracted metadata: 1920x1080, 125s duration")
	fmt.Println("  ✓ Generated 5 thumbnails at key frames")
	fmt.Println("  ✓ Uploaded video to CDN: https://cdn.example.com/videos/abc123.mp4")
	fmt.Println("  ✓ Uploaded thumbnails to CDN")
	fmt.Println("  ✓ Stored metadata in database")
	fmt.Println("  ✅ Encoding complete!")

	fmt.Println("\n💡 Key Insight:")
	fmt.Println("The facade allows clients to encode videos with a single method call,")
	fmt.Println("while internally coordinating 7+ subsystems. As video encoding technology")
	fmt.Println("evolves (new codecs, formats), clients remain unchanged - only the")
	fmt.Println("facade implementation updates.")

	fmt.Println("\n" + strings.Repeat("─", 50))
}

// Example6_FacadePatternIntegration demonstrates how Facade works with other patterns.
// Shows integration with Factory pattern (from Tier 1) to create complex facades.
func Example6_FacadePatternIntegration() {
	fmt.Println("=== Example 6: Facade + Factory Pattern Integration ===")
	fmt.Println("\nFacade Pattern naturally composes with other patterns:")
	fmt.Println()

	fmt.Println("✓ Facade + Factory Pattern:")
	fmt.Println("  Use Factory to create different facade configurations")
	fmt.Println("  Example: NewBankingFacade(region) returns region-specific facade")
	fmt.Println()

	fmt.Println("✓ Facade + Singleton Pattern:")
	fmt.Println("  Facades often manage shared resources (DB, caches)")
	fmt.Println("  Example: GetPaymentFacade() returns singleton instance")
	fmt.Println()

	fmt.Println("✓ Facade + Adapter Pattern:")
	fmt.Println("  Facades may use adapters for third-party integration")
	fmt.Println("  Example: Facade wraps multiple adapted services")
	fmt.Println()

	fmt.Println("✓ Facade + Composite Pattern:")
	fmt.Println("  Facades can expose operations on composite structures")
	fmt.Println("  Example: FileSystemFacade.ArchiveDirectory(path)")
	fmt.Println()

	fmt.Println("Real-world example:")
	fmt.Println("```go")
	fmt.Println("// Factory creates region-specific banking facades")
	fmt.Println("type BankingFacadeFactory struct {}")
	fmt.Println()
	fmt.Println("func (f *BankingFacadeFactory) CreateFacade(region string) *BankingFacade {")
	fmt.Println("    switch region {")
	fmt.Println("    case \"US\":")
	fmt.Println("        return NewBankingFacade(")
	fmt.Println("            NewUSAccountService(),")
	fmt.Println("            NewUSPaymentProcessor(),")
	fmt.Println("            NewUSFraudDetection(),")
	fmt.Println("        )")
	fmt.Println("    case \"EU\":")
	fmt.Println("        return NewBankingFacade(")
	fmt.Println("            NewEUAccountService(),")
	fmt.Println("            NewSEPAPaymentProcessor(),")
	fmt.Println("            NewEUFraudDetection(),")
	fmt.Println("        )")
	fmt.Println("    }")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("💡 Pattern Composition Benefits:")
	fmt.Println("  • Factory handles complex facade initialization")
	fmt.Println("  • Singleton prevents duplicate facade instances")
	fmt.Println("  • Adapter enables third-party service integration")
	fmt.Println("  • Each pattern solves a specific architectural problem")

	fmt.Println("\n" + strings.Repeat("─", 50))
}
