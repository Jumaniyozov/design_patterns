package state

import "testing"

func TestTCPConnection_StateTransitions(t *testing.T) {
	conn := NewTCPConnection()

	// Initial state should be CLOSED
	if conn.GetState() != "CLOSED" {
		t.Errorf("Expected initial state CLOSED, got %s", conn.GetState())
	}

	// Open connection
	err := conn.Open()
	if err != nil {
		t.Errorf("Open failed: %v", err)
	}

	if conn.GetState() != "ESTABLISHED" {
		t.Errorf("Expected state ESTABLISHED after open, got %s", conn.GetState())
	}

	// Send data
	err = conn.Send([]byte("test"))
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}

	// Close connection
	err = conn.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	if conn.GetState() != "CLOSED" {
		t.Errorf("Expected state CLOSED after close, got %s", conn.GetState())
	}
}

func TestTCPConnection_InvalidOperations(t *testing.T) {
	conn := NewTCPConnection()

	// Send while closed should fail
	err := conn.Send([]byte("test"))
	if err == nil {
		t.Error("Expected error when sending on closed connection")
	}

	// Receive while closed should fail
	_, err = conn.Receive()
	if err == nil {
		t.Error("Expected error when receiving on closed connection")
	}
}

func TestDocument_Workflow(t *testing.T) {
	doc := NewDocument("test")

	// Initial state should be DRAFT
	if doc.GetState() != "DRAFT" {
		t.Errorf("Expected initial state DRAFT, got %s", doc.GetState())
	}

	// Publish to move to REVIEW
	err := doc.Publish()
	if err != nil {
		t.Errorf("Publish failed: %v", err)
	}

	if doc.GetState() != "REVIEW" {
		t.Errorf("Expected state REVIEW after publish, got %s", doc.GetState())
	}

	// Approve to move to PUBLISHED
	err = doc.Approve()
	if err != nil {
		t.Errorf("Approve failed: %v", err)
	}

	if doc.GetState() != "PUBLISHED" {
		t.Errorf("Expected state PUBLISHED after approve, got %s", doc.GetState())
	}
}

func TestDocument_Rejection(t *testing.T) {
	doc := NewDocument("test")

	// Submit for review
	doc.Publish()

	if doc.GetState() != "REVIEW" {
		t.Errorf("Expected state REVIEW, got %s", doc.GetState())
	}

	// Reject to go back to DRAFT
	err := doc.Reject()
	if err != nil {
		t.Errorf("Reject failed: %v", err)
	}

	if doc.GetState() != "DRAFT" {
		t.Errorf("Expected state DRAFT after reject, got %s", doc.GetState())
	}
}

func TestPlayer_Movement(t *testing.T) {
	player := NewPlayer()

	// Initial state should be IDLE
	if player.GetState() != "IDLE" {
		t.Errorf("Expected initial state IDLE, got %s", player.GetState())
	}

	// Move right transitions to WALKING
	err := player.MoveRight()
	if err != nil {
		t.Errorf("MoveRight failed: %v", err)
	}

	if player.GetState() != "WALKING" {
		t.Errorf("Expected state WALKING after move, got %s", player.GetState())
	}

	if player.x != 1 {
		t.Errorf("Expected x position 1, got %d", player.x)
	}
}

func TestPlayer_Jump(t *testing.T) {
	player := NewPlayer()

	// Jump from idle
	err := player.Jump()
	if err != nil {
		t.Errorf("Jump failed: %v", err)
	}

	if player.GetState() != "JUMPING" {
		t.Errorf("Expected state JUMPING, got %s", player.GetState())
	}

	if player.y != 10 {
		t.Errorf("Expected y position 10, got %d", player.y)
	}

	// Can't double jump
	err = player.Jump()
	if err == nil {
		t.Error("Expected error for double jump")
	}
}

func TestPlayer_Death(t *testing.T) {
	player := NewPlayer()
	initialHealth := player.health

	// Take fatal damage
	err := player.TakeDamage(initialHealth + 10)
	if err != nil {
		t.Errorf("TakeDamage failed: %v", err)
	}

	if player.GetState() != "DEAD" {
		t.Errorf("Expected state DEAD, got %s", player.GetState())
	}

	// Can't move when dead
	err = player.MoveRight()
	if err == nil {
		t.Error("Expected error when moving while dead")
	}

	// Can't jump when dead
	err = player.Jump()
	if err == nil {
		t.Error("Expected error when jumping while dead")
	}
}

func TestPlayer_CombatSequence(t *testing.T) {
	player := NewPlayer()

	// Walk
	player.MoveRight()
	if player.GetState() != "WALKING" {
		t.Errorf("Expected WALKING, got %s", player.GetState())
	}

	// Jump
	player.Jump()
	if player.GetState() != "JUMPING" {
		t.Errorf("Expected JUMPING, got %s", player.GetState())
	}

	// Take damage while jumping (non-fatal)
	player.TakeDamage(30)
	if player.GetState() != "IDLE" {
		t.Errorf("Expected IDLE after taking damage, got %s", player.GetState())
	}

	if player.health != 70 {
		t.Errorf("Expected health 70, got %d", player.health)
	}
}

// Benchmark tests
func BenchmarkStateTransition(b *testing.B) {
	conn := NewTCPConnection()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn.Open()
		conn.Close()
	}
}

func BenchmarkPlayerMovement(b *testing.B) {
	player := NewPlayer()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		player.MoveRight()
		player.Jump()
	}
}
