package proxy

import "fmt"

type Image interface {
	Display()
}

type RealImage struct {
	filename string
}

func (r *RealImage) loadFromDisk() {
	fmt.Printf("Loading image: %s\n", r.filename)
}

func (r *RealImage) Display() {
	fmt.Printf("Displaying image: %s\n", r.filename)
}

type ProxyImage struct {
	filename  string
	realImage *RealImage
}

func NewProxyImage(filename string) *ProxyImage {
	return &ProxyImage{filename: filename}
}

func (p *ProxyImage) Display() {
	if p.realImage == nil {
		p.realImage = &RealImage{filename: p.filename}
		p.realImage.loadFromDisk()
	}
	p.realImage.Display()
}
