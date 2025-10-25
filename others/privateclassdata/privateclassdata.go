// Package privateclassdata demonstrates the Private Class Data pattern.
// It protects object state by making fields private and providing
// read-only access, ensuring immutability after construction.
package privateclassdata

import "fmt"

// Person demonstrates private class data with immutability
type Person struct {
	name string
	age  int
	ssn  string
}

// NewPerson creates an immutable person
func NewPerson(name string, age int, ssn string) *Person {
	return &Person{
		name: name,
		age:  age,
		ssn:  ssn,
	}
}

// Name returns the name (read-only)
func (p *Person) Name() string {
	return p.name
}

// Age returns the age (read-only)
func (p *Person) Age() int {
	return p.age
}

// String returns string representation (SSN hidden)
func (p *Person) String() string {
	return fmt.Sprintf("%s (age %d)", p.name, p.age)
}

// Configuration with private data
type Configuration struct {
	data configData
}

type configData struct {
	host     string
	port     int
	username string
	password string
}

// NewConfiguration creates an immutable configuration
func NewConfiguration(host string, port int, username, password string) *Configuration {
	return &Configuration{
		data: configData{
			host:     host,
			port:     port,
			username: username,
			password: password,
		},
	}
}

// Host returns host (read-only)
func (c *Configuration) Host() string {
	return c.data.host
}

// Port returns port (read-only)
func (c *Configuration) Port() int {
	return c.data.port
}

// ConnectionString returns connection string (password hidden)
func (c *Configuration) ConnectionString() string {
	return fmt.Sprintf("%s@%s:%d", c.data.username, c.data.host, c.data.port)
}

// Immutable point
type Point struct {
	x, y float64
}

// NewPoint creates an immutable point
func NewPoint(x, y float64) *Point {
	return &Point{x: x, y: y}
}

// X returns x coordinate
func (p *Point) X() float64 {
	return p.x
}

// Y returns y coordinate
func (p *Point) Y() float64 {
	return p.y
}

// Translate returns a new point (original unchanged)
func (p *Point) Translate(dx, dy float64) *Point {
	return NewPoint(p.x+dx, p.y+dy)
}
