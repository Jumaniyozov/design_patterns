// Package mediator demonstrates the Mediator pattern.
// It encapsulates how objects interact, promoting loose coupling
// by preventing direct references between components.
package mediator

import "fmt"

// ChatMediator defines the mediator interface
type ChatMediator interface {
	SendMessage(message string, user User)
	AddUser(user User)
}

// User interface for chat participants
type User interface {
	Send(message string)
	Receive(message string)
	GetName() string
}

// ChatRoom is the concrete mediator
type ChatRoom struct {
	users []User
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{users: make([]User, 0)}
}

func (c *ChatRoom) AddUser(user User) {
	c.users = append(c.users, user)
}

func (c *ChatRoom) SendMessage(message string, sender User) {
	for _, user := range c.users {
		if user != sender {
			user.Receive(fmt.Sprintf("%s: %s", sender.GetName(), message))
		}
	}
}

// ChatUser is a concrete user
type ChatUser struct {
	name     string
	mediator ChatMediator
	messages []string
}

func NewChatUser(name string, mediator ChatMediator) *ChatUser {
	user := &ChatUser{
		name:     name,
		mediator: mediator,
		messages: make([]string, 0),
	}
	mediator.AddUser(user)
	return user
}

func (u *ChatUser) Send(message string) {
	u.mediator.SendMessage(message, u)
}

func (u *ChatUser) Receive(message string) {
	u.messages = append(u.messages, message)
}

func (u *ChatUser) GetName() string {
	return u.name
}

func (u *ChatUser) GetMessages() []string {
	return u.messages
}

// Air Traffic Control example

// ATCMediator coordinates aircraft
type ATCMediator interface {
	RequestLanding(aircraft Aircraft)
	RequestTakeoff(aircraft Aircraft)
	NotifyPositionChange(aircraft Aircraft)
}

// Aircraft interface
type Aircraft interface {
	GetID() string
	Land()
	Takeoff()
	GetAltitude() int
}

// AirTrafficControl is the concrete mediator
type AirTrafficControl struct {
	runway   bool // true if occupied
	aircraft []Aircraft
}

func NewAirTrafficControl() *AirTrafficControl {
	return &AirTrafficControl{
		runway:   false,
		aircraft: make([]Aircraft, 0),
	}
}

func (atc *AirTrafficControl) RegisterAircraft(aircraft Aircraft) {
	atc.aircraft = append(atc.aircraft, aircraft)
}

func (atc *AirTrafficControl) RequestLanding(aircraft Aircraft) string {
	if atc.runway {
		return fmt.Sprintf("ATC to %s: Runway occupied, hold position", aircraft.GetID())
	}
	atc.runway = true
	return fmt.Sprintf("ATC to %s: Cleared for landing", aircraft.GetID())
}

func (atc *AirTrafficControl) RequestTakeoff(aircraft Aircraft) string {
	if atc.runway {
		return fmt.Sprintf("ATC to %s: Runway occupied, wait", aircraft.GetID())
	}
	atc.runway = true
	return fmt.Sprintf("ATC to %s: Cleared for takeoff", aircraft.GetID())
}

func (atc *AirTrafficControl) NotifyPositionChange(aircraft Aircraft) {
	// Coordinate with other aircraft
}

func (atc *AirTrafficControl) RunwayCleared() {
	atc.runway = false
}

// Airplane is a concrete aircraft
type Airplane struct {
	id       string
	altitude int
	mediator *AirTrafficControl
}

func NewAirplane(id string, mediator *AirTrafficControl) *Airplane {
	plane := &Airplane{
		id:       id,
		altitude: 10000,
		mediator: mediator,
	}
	mediator.RegisterAircraft(plane)
	return plane
}

func (a *Airplane) GetID() string {
	return a.id
}

func (a *Airplane) GetAltitude() int {
	return a.altitude
}

func (a *Airplane) Land() {
	a.altitude = 0
	a.mediator.RunwayCleared()
}

func (a *Airplane) Takeoff() {
	a.altitude = 10000
	a.mediator.RunwayCleared()
}

func (a *Airplane) RequestLanding() string {
	return a.mediator.RequestLanding(a)
}

func (a *Airplane) RequestTakeoff() string {
	return a.mediator.RequestTakeoff(a)
}
