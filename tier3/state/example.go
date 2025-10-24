package state

import "fmt"

// Example1_TCPConnection demonstrates TCP connection state machine.
func Example1_TCPConnection() {
	fmt.Println("=== Example 1: TCP Connection State Machine ===")

	conn := NewTCPConnection()
	fmt.Printf("Initial state: %s\n\n", conn.GetState())

	// Try to send while closed
	fmt.Println("--- Attempt to send while closed ---")
	err := conn.Send([]byte("data"))
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
	}

	// Open connection
	fmt.Println("--- Open connection ---")
	conn.Open()
	fmt.Println()

	// Send data
	fmt.Println("--- Send data ---")
	conn.Send([]byte("Hello, World!"))
	fmt.Println()

	// Close connection
	fmt.Println("--- Close connection ---")
	conn.Close()
	fmt.Println()

	// Try to send after closed
	fmt.Println("--- Attempt to send after close ---")
	err = conn.Send([]byte("more data"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println()
}

// Example2_DocumentWorkflow demonstrates document approval workflow.
func Example2_DocumentWorkflow() {
	fmt.Println("=== Example 2: Document Approval Workflow ===")

	doc := NewDocument("Important Proposal")
	fmt.Printf("Initial state: %s\n\n", doc.GetState())

	// Try to approve while in draft
	fmt.Println("--- Attempt to approve draft ---")
	err := doc.Approve()
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
	}

	// Submit for review
	fmt.Println("--- Submit for review ---")
	doc.Publish()
	fmt.Println()

	// Reject and send back to draft
	fmt.Println("--- Reject document ---")
	doc.Reject()
	fmt.Println()

	// Submit again
	fmt.Println("--- Submit for review again ---")
	doc.Publish()
	fmt.Println()

	// Approve
	fmt.Println("--- Approve document ---")
	doc.Approve()
	fmt.Println()

	// Try to reject published document
	fmt.Println("--- Attempt to reject published document ---")
	err = doc.Reject()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println()
}

// Example3_PlayerCharacter demonstrates game character state.
func Example3_PlayerCharacter() {
	fmt.Println("=== Example 3: Game Character States ===")

	player := NewPlayer()
	fmt.Printf("Initial state: %s (Health: %d, Position: (%d, %d))\n\n",
		player.GetState(), player.health, player.x, player.y)

	// Start walking
	fmt.Println("--- Move right ---")
	player.MoveRight()
	fmt.Printf("Position: (%d, %d)\n\n", player.x, player.y)

	// Continue walking
	fmt.Println("--- Move right again ---")
	player.MoveRight()
	fmt.Printf("Position: (%d, %d)\n\n", player.x, player.y)

	// Jump
	fmt.Println("--- Jump ---")
	player.Jump()
	fmt.Printf("Position: (%d, %d)\n\n", player.x, player.y)

	// Try to jump again
	fmt.Println("--- Attempt double jump ---")
	err := player.Jump()
	if err != nil {
		fmt.Printf("Error: %v\n\n", err)
	}

	// Take damage while jumping
	fmt.Println("--- Take damage (50) while jumping ---")
	player.TakeDamage(50)
	fmt.Printf("Health: %d, State: %s, Position: (%d, %d)\n\n",
		player.health, player.GetState(), player.x, player.y)

	// Take fatal damage
	fmt.Println("--- Take fatal damage (60) ---")
	player.TakeDamage(60)
	fmt.Printf("Health: %d, State: %s\n\n", player.health, player.GetState())

	// Try to move when dead
	fmt.Println("--- Attempt to move when dead ---")
	err = player.MoveRight()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println()
}

// Example4_StateTransitions demonstrates all state transitions.
func Example4_StateTransitions() {
	fmt.Println("=== Example 4: Complete State Transitions ===")

	fmt.Println("--- TCP Connection Lifecycle ---")
	conn := NewTCPConnection()
	conn.Open()
	conn.Send([]byte("test"))
	conn.Close()
	fmt.Println()

	fmt.Println("--- Document Approval Lifecycle ---")
	doc := NewDocument("Report")
	doc.Publish()  // Draft -> Review
	doc.Approve()  // Review -> Published
	fmt.Println()

	fmt.Println("--- Player Combat Sequence ---")
	player := NewPlayer()
	player.MoveRight()     // Idle -> Walking
	player.Jump()          // Walking -> Jumping
	player.TakeDamage(25)  // Jumping -> Idle (landing)
	player.TakeDamage(80)  // Idle -> Dead
	fmt.Println()
}
