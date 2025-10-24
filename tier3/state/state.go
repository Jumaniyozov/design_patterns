// Package state implements the State pattern.
//
// The State pattern allows an object to alter its behavior when its internal
// state changes. The object will appear to change its class.
//
// Key components:
// - State interface: Defines state-specific behavior
// - Concrete states: Implement state-specific behavior
// - Context: Maintains current state and delegates to it
package state

import (
	"errors"
	"fmt"
)

// TCPState defines the interface for TCP connection states.
type TCPState interface {
	Open(conn *TCPConnection) error
	Close(conn *TCPConnection) error
	Send(conn *TCPConnection, data []byte) error
	Receive(conn *TCPConnection) ([]byte, error)
	GetStateName() string
}

// TCPConnection represents a TCP connection with state.
type TCPConnection struct {
	state TCPState
	data  []byte
}

// NewTCPConnection creates a new TCP connection.
func NewTCPConnection() *TCPConnection {
	return &TCPConnection{
		state: &ClosedState{},
	}
}

// SetState changes the connection state.
func (c *TCPConnection) SetState(state TCPState) {
	fmt.Printf("[Connection] State transition: %s -> %s\n",
		c.state.GetStateName(), state.GetStateName())
	c.state = state
}

// Open opens the connection.
func (c *TCPConnection) Open() error {
	return c.state.Open(c)
}

// Close closes the connection.
func (c *TCPConnection) Close() error {
	return c.state.Close(c)
}

// Send sends data.
func (c *TCPConnection) Send(data []byte) error {
	return c.state.Send(c, data)
}

// Receive receives data.
func (c *TCPConnection) Receive() ([]byte, error) {
	return c.state.Receive(c)
}

// GetState returns the current state name.
func (c *TCPConnection) GetState() string {
	return c.state.GetStateName()
}

// ClosedState represents a closed connection.
type ClosedState struct{}

func (s *ClosedState) GetStateName() string { return "CLOSED" }

func (s *ClosedState) Open(conn *TCPConnection) error {
	fmt.Println("[CLOSED] Opening connection...")
	conn.SetState(&EstablishedState{})
	return nil
}

func (s *ClosedState) Close(conn *TCPConnection) error {
	return errors.New("connection already closed")
}

func (s *ClosedState) Send(conn *TCPConnection, data []byte) error {
	return errors.New("cannot send: connection closed")
}

func (s *ClosedState) Receive(conn *TCPConnection) ([]byte, error) {
	return nil, errors.New("cannot receive: connection closed")
}

// EstablishedState represents an established connection.
type EstablishedState struct{}

func (s *EstablishedState) GetStateName() string { return "ESTABLISHED" }

func (s *EstablishedState) Open(conn *TCPConnection) error {
	return errors.New("connection already open")
}

func (s *EstablishedState) Close(conn *TCPConnection) error {
	fmt.Println("[ESTABLISHED] Closing connection...")
	conn.SetState(&ClosedState{})
	return nil
}

func (s *EstablishedState) Send(conn *TCPConnection, data []byte) error {
	fmt.Printf("[ESTABLISHED] Sending %d bytes\n", len(data))
	return nil
}

func (s *EstablishedState) Receive(conn *TCPConnection) ([]byte, error) {
	fmt.Println("[ESTABLISHED] Receiving data...")
	return conn.data, nil
}

// ListenState represents a listening connection.
type ListenState struct{}

func (s *ListenState) GetStateName() string { return "LISTEN" }

func (s *ListenState) Open(conn *TCPConnection) error {
	fmt.Println("[LISTEN] Accepting connection...")
	conn.SetState(&EstablishedState{})
	return nil
}

func (s *ListenState) Close(conn *TCPConnection) error {
	fmt.Println("[LISTEN] Closing listener...")
	conn.SetState(&ClosedState{})
	return nil
}

func (s *ListenState) Send(conn *TCPConnection, data []byte) error {
	return errors.New("cannot send: connection in listen state")
}

func (s *ListenState) Receive(conn *TCPConnection) ([]byte, error) {
	return nil, errors.New("cannot receive: connection in listen state")
}

// DocumentState defines the interface for document workflow states.
type DocumentState interface {
	Publish(doc *Document) error
	Approve(doc *Document) error
	Reject(doc *Document) error
	GetStateName() string
}

// Document represents a document with workflow state.
type Document struct {
	content string
	state   DocumentState
}

// NewDocument creates a new document in draft state.
func NewDocument(content string) *Document {
	return &Document{
		content: content,
		state:   &DraftState{},
	}
}

// SetState changes the document state.
func (d *Document) SetState(state DocumentState) {
	fmt.Printf("[Document] State: %s -> %s\n",
		d.state.GetStateName(), state.GetStateName())
	d.state = state
}

// Publish publishes the document.
func (d *Document) Publish() error {
	return d.state.Publish(d)
}

// Approve approves the document.
func (d *Document) Approve() error {
	return d.state.Approve(d)
}

// Reject rejects the document.
func (d *Document) Reject() error {
	return d.state.Reject(d)
}

// GetState returns the current state name.
func (d *Document) GetState() string {
	return d.state.GetStateName()
}

// DraftState represents a draft document.
type DraftState struct{}

func (s *DraftState) GetStateName() string { return "DRAFT" }

func (s *DraftState) Publish(doc *Document) error {
	fmt.Println("[DRAFT] Submitting for review...")
	doc.SetState(&ReviewState{})
	return nil
}

func (s *DraftState) Approve(doc *Document) error {
	return errors.New("cannot approve: document is in draft")
}

func (s *DraftState) Reject(doc *Document) error {
	return errors.New("cannot reject: document is in draft")
}

// ReviewState represents a document under review.
type ReviewState struct{}

func (s *ReviewState) GetStateName() string { return "REVIEW" }

func (s *ReviewState) Publish(doc *Document) error {
	return errors.New("document already in review")
}

func (s *ReviewState) Approve(doc *Document) error {
	fmt.Println("[REVIEW] Approving document...")
	doc.SetState(&PublishedState{})
	return nil
}

func (s *ReviewState) Reject(doc *Document) error {
	fmt.Println("[REVIEW] Rejecting document...")
	doc.SetState(&DraftState{})
	return nil
}

// PublishedState represents a published document.
type PublishedState struct{}

func (s *PublishedState) GetStateName() string { return "PUBLISHED" }

func (s *PublishedState) Publish(doc *Document) error {
	return errors.New("document already published")
}

func (s *PublishedState) Approve(doc *Document) error {
	return errors.New("document already published")
}

func (s *PublishedState) Reject(doc *Document) error {
	return errors.New("cannot reject: document is published")
}

// PlayerState defines the interface for player character states.
type PlayerState interface {
	MoveLeft(player *Player) error
	MoveRight(player *Player) error
	Jump(player *Player) error
	TakeDamage(player *Player, damage int) error
	GetStateName() string
}

// Player represents a game character with state.
type Player struct {
	state  PlayerState
	health int
	x      int
	y      int
}

// NewPlayer creates a new player.
func NewPlayer() *Player {
	return &Player{
		state:  &IdleState{},
		health: 100,
		x:      0,
		y:      0,
	}
}

// SetState changes the player state.
func (p *Player) SetState(state PlayerState) {
	fmt.Printf("[Player] State: %s -> %s\n",
		p.state.GetStateName(), state.GetStateName())
	p.state = state
}

// MoveLeft moves the player left.
func (p *Player) MoveLeft() error {
	return p.state.MoveLeft(p)
}

// MoveRight moves the player right.
func (p *Player) MoveRight() error {
	return p.state.MoveRight(p)
}

// Jump makes the player jump.
func (p *Player) Jump() error {
	return p.state.Jump(p)
}

// TakeDamage deals damage to the player.
func (p *Player) TakeDamage(damage int) error {
	return p.state.TakeDamage(p, damage)
}

// GetState returns the current state name.
func (p *Player) GetState() string {
	return p.state.GetStateName()
}

// IdleState represents an idle player.
type IdleState struct{}

func (s *IdleState) GetStateName() string { return "IDLE" }

func (s *IdleState) MoveLeft(player *Player) error {
	player.x--
	player.SetState(&WalkingState{})
	return nil
}

func (s *IdleState) MoveRight(player *Player) error {
	player.x++
	player.SetState(&WalkingState{})
	return nil
}

func (s *IdleState) Jump(player *Player) error {
	player.y += 10
	player.SetState(&JumpingState{})
	return nil
}

func (s *IdleState) TakeDamage(player *Player, damage int) error {
	player.health -= damage
	if player.health <= 0 {
		player.SetState(&DeadState{})
	}
	return nil
}

// WalkingState represents a walking player.
type WalkingState struct{}

func (s *WalkingState) GetStateName() string { return "WALKING" }

func (s *WalkingState) MoveLeft(player *Player) error {
	player.x--
	return nil
}

func (s *WalkingState) MoveRight(player *Player) error {
	player.x++
	return nil
}

func (s *WalkingState) Jump(player *Player) error {
	player.y += 10
	player.SetState(&JumpingState{})
	return nil
}

func (s *WalkingState) TakeDamage(player *Player, damage int) error {
	player.health -= damage
	if player.health <= 0 {
		player.SetState(&DeadState{})
	}
	return nil
}

// JumpingState represents a jumping player.
type JumpingState struct{}

func (s *JumpingState) GetStateName() string { return "JUMPING" }

func (s *JumpingState) MoveLeft(player *Player) error {
	player.x--
	return nil
}

func (s *JumpingState) MoveRight(player *Player) error {
	player.x++
	return nil
}

func (s *JumpingState) Jump(player *Player) error {
	return errors.New("already jumping")
}

func (s *JumpingState) TakeDamage(player *Player, damage int) error {
	player.health -= damage
	if player.health <= 0 {
		player.SetState(&DeadState{})
	} else {
		// Land after taking damage
		player.y = 0
		player.SetState(&IdleState{})
	}
	return nil
}

// DeadState represents a dead player.
type DeadState struct{}

func (s *DeadState) GetStateName() string { return "DEAD" }

func (s *DeadState) MoveLeft(player *Player) error {
	return errors.New("cannot move: player is dead")
}

func (s *DeadState) MoveRight(player *Player) error {
	return errors.New("cannot move: player is dead")
}

func (s *DeadState) Jump(player *Player) error {
	return errors.New("cannot jump: player is dead")
}

func (s *DeadState) TakeDamage(player *Player, damage int) error {
	return nil // Already dead
}
