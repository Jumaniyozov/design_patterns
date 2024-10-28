package main

import "fmt"

//Когда использовать: Когда нужно контролировать доступ к объекту, добавляя дополнительные действия при доступе.

// Интерфейс объекта
type Image interface {
	display()
}

// Реальный объект
type RealImage struct {
	fileName string
}

func (r *RealImage) display() {
	fmt.Println("Отображение", r.fileName)
}

func (r *RealImage) loadFromDisk() {
	fmt.Println("Загрузка", r.fileName)
}

func NewRealImage(fileName string) *RealImage {
	image := &RealImage{fileName: fileName}
	image.loadFromDisk()
	return image
}

// Прокси
type ProxyImage struct {
	realImage *RealImage
	fileName  string
}

func (p *ProxyImage) display() {
	if p.realImage == nil {
		p.realImage = NewRealImage(p.fileName)
	}
	p.realImage.display()
}

func NewProxyImage(fileName string) *ProxyImage {
	return &ProxyImage{fileName: fileName}
}

func main() {
	image := NewProxyImage("test.jpg")
	// Загрузка произойдет здесь
	image.display()
	// Повторная загрузка не нужна
	image.display()
}
