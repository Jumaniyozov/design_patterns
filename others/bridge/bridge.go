// Package bridge demonstrates the Bridge pattern.
// It decouples abstraction from implementation so both can vary independently,
// preventing class explosion when you have multiple dimensions of variation.
package bridge

import "fmt"

// Implementation Interface - Renderer

// Renderer defines the implementation interface for rendering operations
type Renderer interface {
	RenderCircle(radius float64) string
	RenderRectangle(width, height float64) string
}

// Concrete Implementations

// VectorRenderer renders shapes as vector graphics
type VectorRenderer struct {
	name string
}

// NewVectorRenderer creates a new vector renderer
func NewVectorRenderer() *VectorRenderer {
	return &VectorRenderer{name: "Vector"}
}

func (v *VectorRenderer) RenderCircle(radius float64) string {
	return fmt.Sprintf("[%s] Drawing circle with radius %.2f using SVG paths", v.name, radius)
}

func (v *VectorRenderer) RenderRectangle(width, height float64) string {
	return fmt.Sprintf("[%s] Drawing rectangle %.2fx%.2f using SVG paths", v.name, width, height)
}

// RasterRenderer renders shapes as raster/bitmap graphics
type RasterRenderer struct {
	name       string
	resolution string
}

// NewRasterRenderer creates a new raster renderer
func NewRasterRenderer(resolution string) *RasterRenderer {
	return &RasterRenderer{
		name:       "Raster",
		resolution: resolution,
	}
}

func (r *RasterRenderer) RenderCircle(radius float64) string {
	return fmt.Sprintf("[%s@%s] Drawing circle with radius %.2f as pixels",
		r.name, r.resolution, radius)
}

func (r *RasterRenderer) RenderRectangle(width, height float64) string {
	return fmt.Sprintf("[%s@%s] Drawing rectangle %.2fx%.2f as pixel grid",
		r.name, r.resolution, width, height)
}

// OpenGLRenderer renders shapes using OpenGL
type OpenGLRenderer struct {
	name string
	version string
}

// NewOpenGLRenderer creates a new OpenGL renderer
func NewOpenGLRenderer(version string) *OpenGLRenderer {
	return &OpenGLRenderer{
		name:    "OpenGL",
		version: version,
	}
}

func (o *OpenGLRenderer) RenderCircle(radius float64) string {
	return fmt.Sprintf("[%s %s] Rendering circle radius %.2f with GPU acceleration",
		o.name, o.version, radius)
}

func (o *OpenGLRenderer) RenderRectangle(width, height float64) string {
	return fmt.Sprintf("[%s %s] Rendering rectangle %.2fx%.2f with GPU acceleration",
		o.name, o.version, width, height)
}

// Abstraction - Shape

// Shape is the abstraction that uses a Renderer implementation
type Shape struct {
	renderer Renderer
}

// SetRenderer changes the rendering implementation at runtime
func (s *Shape) SetRenderer(renderer Renderer) {
	s.renderer = renderer
}

// Refined Abstractions

// Circle is a refined abstraction
type Circle struct {
	Shape
	radius float64
}

// NewCircle creates a circle with a specific renderer
func NewCircle(radius float64, renderer Renderer) *Circle {
	return &Circle{
		Shape:  Shape{renderer: renderer},
		radius: radius,
	}
}

// Draw renders the circle using the bridge to implementation
func (c *Circle) Draw() string {
	return c.renderer.RenderCircle(c.radius)
}

// Resize changes the circle's radius
func (c *Circle) Resize(radius float64) {
	c.radius = radius
}

// GetInfo returns circle information
func (c *Circle) GetInfo() string {
	return fmt.Sprintf("Circle (radius: %.2f)", c.radius)
}

// Rectangle is another refined abstraction
type Rectangle struct {
	Shape
	width  float64
	height float64
}

// NewRectangle creates a rectangle with a specific renderer
func NewRectangle(width, height float64, renderer Renderer) *Rectangle {
	return &Rectangle{
		Shape:  Shape{renderer: renderer},
		width:  width,
		height: height,
	}
}

// Draw renders the rectangle using the bridge to implementation
func (r *Rectangle) Draw() string {
	return r.renderer.RenderRectangle(r.width, r.height)
}

// Resize changes the rectangle's dimensions
func (r *Rectangle) Resize(width, height float64) {
	r.width = width
	r.height = height
}

// GetInfo returns rectangle information
func (r *Rectangle) GetInfo() string {
	return fmt.Sprintf("Rectangle (%.2fx%.2f)", r.width, r.height)
}

// Another example: Message abstraction with different senders

// MessageSender is the implementation interface for sending messages
type MessageSender interface {
	Send(recipient, message string) string
}

// EmailSender sends messages via email
type EmailSender struct {
	smtpServer string
}

// NewEmailSender creates an email sender
func NewEmailSender(smtpServer string) *EmailSender {
	return &EmailSender{smtpServer: smtpServer}
}

func (e *EmailSender) Send(recipient, message string) string {
	return fmt.Sprintf("[Email via %s] To: %s | Message: %s", e.smtpServer, recipient, message)
}

// SMSSender sends messages via SMS
type SMSSender struct {
	gateway string
}

// NewSMSSender creates an SMS sender
func NewSMSSender(gateway string) *SMSSender {
	return &SMSSender{gateway: gateway}
}

func (s *SMSSender) Send(recipient, message string) string {
	return fmt.Sprintf("[SMS via %s] To: %s | Message: %s", s.gateway, recipient, message)
}

// Message is the abstraction for different message types
type Message struct {
	sender MessageSender
}

// UrgentMessage is a refined abstraction for urgent messages
type UrgentMessage struct {
	Message
	priority string
}

// NewUrgentMessage creates an urgent message with a sender
func NewUrgentMessage(sender MessageSender) *UrgentMessage {
	return &UrgentMessage{
		Message:  Message{sender: sender},
		priority: "HIGH",
	}
}

// Send sends an urgent message with priority indicator
func (u *UrgentMessage) Send(recipient, message string) string {
	formattedMessage := fmt.Sprintf("[%s PRIORITY] %s", u.priority, message)
	return u.sender.Send(recipient, formattedMessage)
}

// RegularMessage is a refined abstraction for regular messages
type RegularMessage struct {
	Message
}

// NewRegularMessage creates a regular message with a sender
func NewRegularMessage(sender MessageSender) *RegularMessage {
	return &RegularMessage{
		Message: Message{sender: sender},
	}
}

// Send sends a regular message
func (r *RegularMessage) Send(recipient, message string) string {
	return r.sender.Send(recipient, message)
}
