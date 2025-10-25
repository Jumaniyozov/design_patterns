package mediator

import "fmt"

type Mediator interface {
	Notify(sender string, event string)
}

type ChatMediator struct {
	users map[string]*User
}

func NewChatMediator() *ChatMediator {
	return &ChatMediator{users: make(map[string]*User)}
}

func (m *ChatMediator) Register(user *User) {
	m.users[user.name] = user
}

func (m *ChatMediator) Notify(sender string, message string) {
	for name, user := range m.users {
		if name != sender {
			user.Receive(sender, message)
		}
	}
}

type User struct {
	name     string
	mediator *ChatMediator
}

func NewUser(name string, mediator *ChatMediator) *User {
	user := &User{name: name, mediator: mediator}
	mediator.Register(user)
	return user
}

func (u *User) Send(message string) {
	fmt.Printf("[%s] Sending: %s\n", u.name, message)
	u.mediator.Notify(u.name, message)
}

func (u *User) Receive(sender, message string) {
	fmt.Printf("[%s] Received from %s: %s\n", u.name, sender, message)
}
