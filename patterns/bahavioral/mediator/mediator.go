package main

import "fmt"

// Когда использовать: Когда нужно уменьшить количество связей между объектами.

// Посредник
type ChatRoom struct{}

func (c *ChatRoom) showMessage(user *User, message string) {
	fmt.Printf("[%s]: %s\n", user.name, message)
}

// Коллега
type User struct {
	name string
	chat *ChatRoom
}

func (u *User) sendMessage(message string) {
	u.chat.showMessage(u, message)
}

func main() {
	chat := &ChatRoom{}
	user1 := &User{name: "Алиса", chat: chat}
	user2 := &User{name: "Боб", chat: chat}

	user1.sendMessage("Привет, Боб!")
	user2.sendMessage("Привет, Алиса!")
}
