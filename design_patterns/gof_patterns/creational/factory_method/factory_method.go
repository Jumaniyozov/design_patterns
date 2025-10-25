// Package factory_method demonstrates the Factory Method pattern in Go.
// The Factory Method pattern defines an interface for creating objects,
// but lets factory functions decide which concrete type to instantiate.
package factory_method

import (
	"fmt"
	"strings"
	"time"
)

// Notifier defines the interface for all notification types.
type Notifier interface {
	Send(message string) error
	GetChannel() string
}

// EmailNotifier sends notifications via email.
type EmailNotifier struct {
	smtpServer string
	from       string
}

func (e *EmailNotifier) Send(message string) error {
	fmt.Printf("[EMAIL] Sending via %s from %s: %s\n", e.smtpServer, e.from, message)
	return nil
}

func (e *EmailNotifier) GetChannel() string {
	return "email"
}

// SMSNotifier sends notifications via SMS.
type SMSNotifier struct {
	gateway string
	from    string
}

func (s *SMSNotifier) Send(message string) error {
	fmt.Printf("[SMS] Sending via %s from %s: %s\n", s.gateway, s.from, message)
	return nil
}

func (s *SMSNotifier) GetChannel() string {
	return "sms"
}

// PushNotifier sends push notifications.
type PushNotifier struct {
	service string
	appID   string
}

func (p *PushNotifier) Send(message string) error {
	fmt.Printf("[PUSH] Sending via %s (App: %s): %s\n", p.service, p.appID, message)
	return nil
}

func (p *PushNotifier) GetChannel() string {
	return "push"
}

// SlackNotifier sends notifications to Slack.
type SlackNotifier struct {
	webhookURL string
	channel    string
}

func (s *SlackNotifier) Send(message string) error {
	fmt.Printf("[SLACK] Posting to #%s via webhook: %s\n", s.channel, message)
	return nil
}

func (s *SlackNotifier) GetChannel() string {
	return "slack"
}

// NotificationConfig holds configuration for creating notifiers.
type NotificationConfig struct {
	SMTPServer   string
	SMSGateway   string
	PushService  string
	SlackWebhook string
	From         string
	AppID        string
	SlackChannel string
}

// NewNotifier is the factory function that creates the appropriate notifier
// based on the channel type.
func NewNotifier(channel string, config NotificationConfig) (Notifier, error) {
	channel = strings.ToLower(channel)

	switch channel {
	case "email":
		return &EmailNotifier{
			smtpServer: config.SMTPServer,
			from:       config.From,
		}, nil
	case "sms":
		return &SMSNotifier{
			gateway: config.SMSGateway,
			from:    config.From,
		}, nil
	case "push":
		return &PushNotifier{
			service: config.PushService,
			appID:   config.AppID,
		}, nil
	case "slack":
		return &SlackNotifier{
			webhookURL: config.SlackWebhook,
			channel:    config.SlackChannel,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported notification channel: %s", channel)
	}
}

// Parser defines the interface for document parsers.
type Parser interface {
	Parse(data []byte) (map[string]interface{}, error)
	GetFormat() string
}

// JSONParser parses JSON documents.
type JSONParser struct{}

func (j *JSONParser) Parse(data []byte) (map[string]interface{}, error) {
	fmt.Printf("[JSON Parser] Parsing %d bytes of JSON data\n", len(data))
	// Simplified - real implementation would use encoding/json
	return map[string]interface{}{"format": "json", "size": len(data)}, nil
}

func (j *JSONParser) GetFormat() string {
	return "JSON"
}

// XMLParser parses XML documents.
type XMLParser struct{}

func (x *XMLParser) Parse(data []byte) (map[string]interface{}, error) {
	fmt.Printf("[XML Parser] Parsing %d bytes of XML data\n", len(data))
	return map[string]interface{}{"format": "xml", "size": len(data)}, nil
}

func (x *XMLParser) GetFormat() string {
	return "XML"
}

// YAMLParser parses YAML documents.
type YAMLParser struct{}

func (y *YAMLParser) Parse(data []byte) (map[string]interface{}, error) {
	fmt.Printf("[YAML Parser] Parsing %d bytes of YAML data\n", len(data))
	return map[string]interface{}{"format": "yaml", "size": len(data)}, nil
}

func (y *YAMLParser) GetFormat() string {
	return "YAML"
}

// NewParser creates a parser based on file extension.
func NewParser(filename string) (Parser, error) {
	// Extract extension
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("no file extension found in: %s", filename)
	}

	ext := strings.ToLower(parts[len(parts)-1])

	switch ext {
	case "json":
		return &JSONParser{}, nil
	case "xml":
		return &XMLParser{}, nil
	case "yaml", "yml":
		return &YAMLParser{}, nil
	default:
		return nil, fmt.Errorf("unsupported file format: .%s", ext)
	}
}

// HTTPClient defines the interface for HTTP clients.
type HTTPClient interface {
	Get(url string) (string, error)
	Post(url string, body string) (string, error)
	GetEnvironment() string
}

// ProductionClient is optimized for production use.
type ProductionClient struct {
	timeout time.Duration
}

func (p *ProductionClient) Get(url string) (string, error) {
	fmt.Printf("[PROD Client] GET %s (timeout: %v)\n", url, p.timeout)
	return "production response", nil
}

func (p *ProductionClient) Post(url string, body string) (string, error) {
	fmt.Printf("[PROD Client] POST %s (timeout: %v)\n", url, p.timeout)
	return "production response", nil
}

func (p *ProductionClient) GetEnvironment() string {
	return "production"
}

// DevelopmentClient has verbose logging and relaxed timeouts.
type DevelopmentClient struct {
	timeout time.Duration
	verbose bool
}

func (d *DevelopmentClient) Get(url string) (string, error) {
	if d.verbose {
		fmt.Printf("[DEV Client] GET %s (timeout: %v, verbose: true)\n", url, d.timeout)
	}
	return "development response", nil
}

func (d *DevelopmentClient) Post(url string, body string) (string, error) {
	if d.verbose {
		fmt.Printf("[DEV Client] POST %s with body: %s (timeout: %v)\n", url, body, d.timeout)
	}
	return "development response", nil
}

func (d *DevelopmentClient) GetEnvironment() string {
	return "development"
}

// MockClient is used for testing.
type MockClient struct {
	responses map[string]string
}

func (m *MockClient) Get(url string) (string, error) {
	fmt.Printf("[MOCK Client] GET %s\n", url)
	if response, ok := m.responses[url]; ok {
		return response, nil
	}
	return "mock response", nil
}

func (m *MockClient) Post(url string, body string) (string, error) {
	fmt.Printf("[MOCK Client] POST %s\n", url)
	return "mock response", nil
}

func (m *MockClient) GetEnvironment() string {
	return "test"
}

// NewHTTPClient creates an HTTP client based on environment.
func NewHTTPClient(environment string) HTTPClient {
	switch strings.ToLower(environment) {
	case "production", "prod":
		return &ProductionClient{
			timeout: 30 * time.Second,
		}
	case "development", "dev":
		return &DevelopmentClient{
			timeout: 5 * time.Minute,
			verbose: true,
		}
	case "test", "testing":
		return &MockClient{
			responses: make(map[string]string),
		}
	default:
		// Default to production settings
		return &ProductionClient{
			timeout: 30 * time.Second,
		}
	}
}
