// Package builder demonstrates the Builder pattern in Go.
// The Builder pattern separates the construction of a complex object from its
// representation, allowing step-by-step construction with optional parameters.
package builder

import (
	"fmt"
	"strings"
	"time"
)

// HTTPRequest represents a complex HTTP request object.
type HTTPRequest struct {
	url             string
	method          string
	headers         map[string]string
	body            string
	timeout         time.Duration
	retryCount      int
	followRedirects bool
}

// HTTPRequestBuilder builds HTTPRequest objects step by step.
type HTTPRequestBuilder struct {
	request HTTPRequest
}

// NewHTTPRequestBuilder creates a new builder with sensible defaults.
func NewHTTPRequestBuilder() *HTTPRequestBuilder {
	return &HTTPRequestBuilder{
		request: HTTPRequest{
			method:          "GET",
			headers:         make(map[string]string),
			timeout:         30 * time.Second,
			retryCount:      3,
			followRedirects: true,
		},
	}
}

// URL sets the request URL.
func (b *HTTPRequestBuilder) URL(url string) *HTTPRequestBuilder {
	b.request.url = url
	return b
}

// Method sets the HTTP method.
func (b *HTTPRequestBuilder) Method(method string) *HTTPRequestBuilder {
	b.request.method = strings.ToUpper(method)
	return b
}

// Header adds a header to the request.
func (b *HTTPRequestBuilder) Header(key, value string) *HTTPRequestBuilder {
	b.request.headers[key] = value
	return b
}

// Body sets the request body.
func (b *HTTPRequestBuilder) Body(body string) *HTTPRequestBuilder {
	b.request.body = body
	return b
}

// Timeout sets the request timeout.
func (b *HTTPRequestBuilder) Timeout(timeout time.Duration) *HTTPRequestBuilder {
	b.request.timeout = timeout
	return b
}

// RetryCount sets the number of retries.
func (b *HTTPRequestBuilder) RetryCount(count int) *HTTPRequestBuilder {
	b.request.retryCount = count
	return b
}

// FollowRedirects sets whether to follow redirects.
func (b *HTTPRequestBuilder) FollowRedirects(follow bool) *HTTPRequestBuilder {
	b.request.followRedirects = follow
	return b
}

// Build constructs and validates the final HTTPRequest.
func (b *HTTPRequestBuilder) Build() (*HTTPRequest, error) {
	// Validation
	if b.request.url == "" {
		return nil, fmt.Errorf("URL is required")
	}

	// Return a copy to ensure immutability
	headers := make(map[string]string)
	for k, v := range b.request.headers {
		headers[k] = v
	}

	return &HTTPRequest{
		url:             b.request.url,
		method:          b.request.method,
		headers:         headers,
		body:            b.request.body,
		timeout:         b.request.timeout,
		retryCount:      b.request.retryCount,
		followRedirects: b.request.followRedirects,
	}, nil
}

// String returns a string representation of the request.
func (r *HTTPRequest) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s\n", r.method, r.url))
	sb.WriteString(fmt.Sprintf("Timeout: %v\n", r.timeout))
	sb.WriteString(fmt.Sprintf("Retries: %d\n", r.retryCount))
	sb.WriteString(fmt.Sprintf("Follow Redirects: %v\n", r.followRedirects))
	if len(r.headers) > 0 {
		sb.WriteString("Headers:\n")
		for k, v := range r.headers {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}
	if r.body != "" {
		sb.WriteString(fmt.Sprintf("Body: %s\n", r.body))
	}
	return sb.String()
}

// SQLQuery represents a complex SQL query.
type SQLQuery struct {
	queryType string // SELECT, INSERT, UPDATE, DELETE
	table     string
	columns   []string
	where     []string
	orderBy   string
	limit     int
	offset    int
	joins     []string
}

// SQLQueryBuilder builds SQL queries step by step.
type SQLQueryBuilder struct {
	query SQLQuery
}

// NewQueryBuilder creates a new SQL query builder.
func NewQueryBuilder() *SQLQueryBuilder {
	return &SQLQueryBuilder{
		query: SQLQuery{
			columns: make([]string, 0),
			where:   make([]string, 0),
			joins:   make([]string, 0),
		},
	}
}

// Select specifies columns to select.
func (b *SQLQueryBuilder) Select(columns ...string) *SQLQueryBuilder {
	b.query.queryType = "SELECT"
	b.query.columns = append(b.query.columns, columns...)
	return b
}

// From specifies the table.
func (b *SQLQueryBuilder) From(table string) *SQLQueryBuilder {
	b.query.table = table
	return b
}

// Where adds a WHERE condition.
func (b *SQLQueryBuilder) Where(condition string) *SQLQueryBuilder {
	b.query.where = append(b.query.where, condition)
	return b
}

// OrderBy sets the ORDER BY clause.
func (b *SQLQueryBuilder) OrderBy(orderBy string) *SQLQueryBuilder {
	b.query.orderBy = orderBy
	return b
}

// Limit sets the LIMIT.
func (b *SQLQueryBuilder) Limit(limit int) *SQLQueryBuilder {
	b.query.limit = limit
	return b
}

// Offset sets the OFFSET.
func (b *SQLQueryBuilder) Offset(offset int) *SQLQueryBuilder {
	b.query.offset = offset
	return b
}

// Join adds a JOIN clause.
func (b *SQLQueryBuilder) Join(join string) *SQLQueryBuilder {
	b.query.joins = append(b.query.joins, join)
	return b
}

// Build constructs the final SQL query string.
func (b *SQLQueryBuilder) Build() (string, error) {
	if b.query.table == "" {
		return "", fmt.Errorf("table name is required")
	}

	var sb strings.Builder

	// SELECT clause
	sb.WriteString("SELECT ")
	if len(b.query.columns) == 0 {
		sb.WriteString("*")
	} else {
		sb.WriteString(strings.Join(b.query.columns, ", "))
	}

	// FROM clause
	sb.WriteString(fmt.Sprintf(" FROM %s", b.query.table))

	// JOIN clauses
	for _, join := range b.query.joins {
		sb.WriteString(fmt.Sprintf(" %s", join))
	}

	// WHERE clause
	if len(b.query.where) > 0 {
		sb.WriteString(fmt.Sprintf(" WHERE %s", strings.Join(b.query.where, " AND ")))
	}

	// ORDER BY clause
	if b.query.orderBy != "" {
		sb.WriteString(fmt.Sprintf(" ORDER BY %s", b.query.orderBy))
	}

	// LIMIT clause
	if b.query.limit > 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", b.query.limit))
	}

	// OFFSET clause
	if b.query.offset > 0 {
		sb.WriteString(fmt.Sprintf(" OFFSET %d", b.query.offset))
	}

	return sb.String(), nil
}

// EmailMessage represents a complex email message.
type EmailMessage struct {
	from        string
	to          []string
	cc          []string
	bcc         []string
	subject     string
	body        string
	htmlBody    string
	attachments []string
	priority    string
}

// EmailBuilder builds email messages.
type EmailBuilder struct {
	email EmailMessage
}

// NewEmailBuilder creates a new email builder.
func NewEmailBuilder() *EmailBuilder {
	return &EmailBuilder{
		email: EmailMessage{
			to:          make([]string, 0),
			cc:          make([]string, 0),
			bcc:         make([]string, 0),
			attachments: make([]string, 0),
			priority:    "normal",
		},
	}
}

// From sets the sender.
func (b *EmailBuilder) From(from string) *EmailBuilder {
	b.email.from = from
	return b
}

// To adds a recipient.
func (b *EmailBuilder) To(to ...string) *EmailBuilder {
	b.email.to = append(b.email.to, to...)
	return b
}

// CC adds a CC recipient.
func (b *EmailBuilder) CC(cc ...string) *EmailBuilder {
	b.email.cc = append(b.email.cc, cc...)
	return b
}

// BCC adds a BCC recipient.
func (b *EmailBuilder) BCC(bcc ...string) *EmailBuilder {
	b.email.bcc = append(b.email.bcc, bcc...)
	return b
}

// Subject sets the subject.
func (b *EmailBuilder) Subject(subject string) *EmailBuilder {
	b.email.subject = subject
	return b
}

// Body sets the plain text body.
func (b *EmailBuilder) Body(body string) *EmailBuilder {
	b.email.body = body
	return b
}

// HTMLBody sets the HTML body.
func (b *EmailBuilder) HTMLBody(html string) *EmailBuilder {
	b.email.htmlBody = html
	return b
}

// Attachment adds an attachment.
func (b *EmailBuilder) Attachment(path string) *EmailBuilder {
	b.email.attachments = append(b.email.attachments, path)
	return b
}

// Priority sets the priority.
func (b *EmailBuilder) Priority(priority string) *EmailBuilder {
	b.email.priority = priority
	return b
}

// Build constructs and validates the email.
func (b *EmailBuilder) Build() (*EmailMessage, error) {
	if b.email.from == "" {
		return nil, fmt.Errorf("sender (from) is required")
	}
	if len(b.email.to) == 0 {
		return nil, fmt.Errorf("at least one recipient (to) is required")
	}
	if b.email.subject == "" {
		return nil, fmt.Errorf("subject is required")
	}
	if b.email.body == "" && b.email.htmlBody == "" {
		return nil, fmt.Errorf("body or HTML body is required")
	}

	// Return copy
	return &EmailMessage{
		from:        b.email.from,
		to:          append([]string{}, b.email.to...),
		cc:          append([]string{}, b.email.cc...),
		bcc:         append([]string{}, b.email.bcc...),
		subject:     b.email.subject,
		body:        b.email.body,
		htmlBody:    b.email.htmlBody,
		attachments: append([]string{}, b.email.attachments...),
		priority:    b.email.priority,
	}, nil
}

// String returns a string representation of the email.
func (e *EmailMessage) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("From: %s\n", e.from))
	sb.WriteString(fmt.Sprintf("To: %s\n", strings.Join(e.to, ", ")))
	if len(e.cc) > 0 {
		sb.WriteString(fmt.Sprintf("CC: %s\n", strings.Join(e.cc, ", ")))
	}
	if len(e.bcc) > 0 {
		sb.WriteString(fmt.Sprintf("BCC: %s\n", strings.Join(e.bcc, ", ")))
	}
	sb.WriteString(fmt.Sprintf("Subject: %s\n", e.subject))
	sb.WriteString(fmt.Sprintf("Priority: %s\n", e.priority))
	if e.body != "" {
		sb.WriteString(fmt.Sprintf("Body: %s\n", e.body))
	}
	if e.htmlBody != "" {
		sb.WriteString(fmt.Sprintf("HTML Body: %s\n", e.htmlBody))
	}
	if len(e.attachments) > 0 {
		sb.WriteString(fmt.Sprintf("Attachments: %s\n", strings.Join(e.attachments, ", ")))
	}
	return sb.String()
}

// Computer represents a complex product built by builder.
type Computer struct {
	cpu      string
	ram      int // GB
	storage  int // GB
	gpu      string
	os       string
	monitor  string
	keyboard string
	mouse    string
}

// ComputerBuilder builds Computer objects.
type ComputerBuilder struct {
	computer Computer
}

// NewComputerBuilder creates a new computer builder.
func NewComputerBuilder() *ComputerBuilder {
	return &ComputerBuilder{
		computer: Computer{
			os: "Linux", // default
		},
	}
}

// CPU sets the CPU.
func (b *ComputerBuilder) CPU(cpu string) *ComputerBuilder {
	b.computer.cpu = cpu
	return b
}

// RAM sets the RAM in GB.
func (b *ComputerBuilder) RAM(gb int) *ComputerBuilder {
	b.computer.ram = gb
	return b
}

// Storage sets the storage in GB.
func (b *ComputerBuilder) Storage(gb int) *ComputerBuilder {
	b.computer.storage = gb
	return b
}

// GPU sets the GPU.
func (b *ComputerBuilder) GPU(gpu string) *ComputerBuilder {
	b.computer.gpu = gpu
	return b
}

// OS sets the operating system.
func (b *ComputerBuilder) OS(os string) *ComputerBuilder {
	b.computer.os = os
	return b
}

// Monitor sets the monitor.
func (b *ComputerBuilder) Monitor(monitor string) *ComputerBuilder {
	b.computer.monitor = monitor
	return b
}

// Keyboard sets the keyboard.
func (b *ComputerBuilder) Keyboard(keyboard string) *ComputerBuilder {
	b.computer.keyboard = keyboard
	return b
}

// Mouse sets the mouse.
func (b *ComputerBuilder) Mouse(mouse string) *ComputerBuilder {
	b.computer.mouse = mouse
	return b
}

// Build constructs the computer.
func (b *ComputerBuilder) Build() (*Computer, error) {
	if b.computer.cpu == "" {
		return nil, fmt.Errorf("CPU is required")
	}
	if b.computer.ram == 0 {
		return nil, fmt.Errorf("RAM is required")
	}
	if b.computer.storage == 0 {
		return nil, fmt.Errorf("storage is required")
	}

	return &Computer{
		cpu:      b.computer.cpu,
		ram:      b.computer.ram,
		storage:  b.computer.storage,
		gpu:      b.computer.gpu,
		os:       b.computer.os,
		monitor:  b.computer.monitor,
		keyboard: b.computer.keyboard,
		mouse:    b.computer.mouse,
	}, nil
}

// String returns a string representation of the computer.
func (c *Computer) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CPU: %s\n", c.cpu))
	sb.WriteString(fmt.Sprintf("RAM: %dGB\n", c.ram))
	sb.WriteString(fmt.Sprintf("Storage: %dGB\n", c.storage))
	if c.gpu != "" {
		sb.WriteString(fmt.Sprintf("GPU: %s\n", c.gpu))
	}
	sb.WriteString(fmt.Sprintf("OS: %s\n", c.os))
	if c.monitor != "" {
		sb.WriteString(fmt.Sprintf("Monitor: %s\n", c.monitor))
	}
	if c.keyboard != "" {
		sb.WriteString(fmt.Sprintf("Keyboard: %s\n", c.keyboard))
	}
	if c.mouse != "" {
		sb.WriteString(fmt.Sprintf("Mouse: %s\n", c.mouse))
	}
	return sb.String()
}
