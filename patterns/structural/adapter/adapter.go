package main

import "fmt"

// Когда использовать: Когда нужно использовать существующий класс, но его интерфейс не соответствует потребностям.

// Целевой интерфейс
type Computer interface {
	InsertIntoLightningPort()
}

// Адаптируемый класс
type Android struct{}

func (a *Android) InsertIntoMicroUSBPort() {
	fmt.Println("Кабель MicroUSB подключен к Android устройству.")
}

// Адаптер
type MicroUSBToLightningAdapter struct {
	androidPhone *Android
}

func (m *MicroUSBToLightningAdapter) InsertIntoLightningPort() {
	fmt.Println("Адаптер преобразует сигнал Lightning в MicroUSB.")
	m.androidPhone.InsertIntoMicroUSBPort()
}

// Клиентский код
type Client struct{}

func (c *Client) InsertLightningConnectorIntoComputer(com Computer) {
	fmt.Println("Клиент вставляет коннектор Lightning в компьютер.")
	com.InsertIntoLightningPort()
}

func main() {
	client := &Client{}
	androidPhone := &Android{}
	adapter := &MicroUSBToLightningAdapter{androidPhone}
	client.InsertLightningConnectorIntoComputer(adapter)
}
