// Package interceptingfilter demonstrates the Intercepting Filter pattern.
// It preprocesses and postprocesses requests through a chain of filters,
// commonly used for middleware in web applications.
package interceptingfilter

import (
	"fmt"
	"strings"
	"time"
)

// Request represents a request
type Request struct {
	Path      string
	Method    string
	Headers   map[string]string
	Body      string
	Timestamp time.Time
}

// Response represents a response
type Response struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// Filter interface for request filters
type Filter interface {
	Execute(req *Request, res *Response, chain *FilterChain)
}

// FilterChain manages the chain of filters
type FilterChain struct {
	filters []Filter
	index   int
	target  Target
}

// NewFilterChain creates a filter chain
func NewFilterChain(target Target) *FilterChain {
	return &FilterChain{
		filters: make([]Filter, 0),
		index:   0,
		target:  target,
	}
}

// AddFilter adds a filter to the chain
func (fc *FilterChain) AddFilter(filter Filter) {
	fc.filters = append(fc.filters, filter)
}

// Execute executes the filter chain
func (fc *FilterChain) Execute(req *Request, res *Response) {
	if fc.index < len(fc.filters) {
		filter := fc.filters[fc.index]
		fc.index++
		filter.Execute(req, res, fc)
	} else if fc.target != nil {
		fc.target.Execute(req, res)
	}
}

// Reset resets the chain for reuse
func (fc *FilterChain) Reset() {
	fc.index = 0
}

// Target is the final handler
type Target interface {
	Execute(req *Request, res *Response)
}

// Concrete Filters

// AuthenticationFilter checks authentication
type AuthenticationFilter struct{}

func (a *AuthenticationFilter) Execute(req *Request, res *Response, chain *FilterChain) {
	token := req.Headers["Authorization"]
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		res.StatusCode = 401
		res.Body = "Unauthorized: Missing or invalid token"
		return
	}
	fmt.Println("[AuthFilter] Authentication successful")
	chain.Execute(req, res)
}

// LoggingFilter logs requests
type LoggingFilter struct{}

func (l *LoggingFilter) Execute(req *Request, res *Response, chain *FilterChain) {
	start := time.Now()
	fmt.Printf("[LogFilter] %s %s started\n", req.Method, req.Path)

	chain.Execute(req, res)

	duration := time.Since(start)
	fmt.Printf("[LogFilter] %s %s completed in %v (status: %d)\n",
		req.Method, req.Path, duration, res.StatusCode)
}

// CompressionFilter compresses response
type CompressionFilter struct{}

func (c *CompressionFilter) Execute(req *Request, res *Response, chain *FilterChain) {
	chain.Execute(req, res)

	if res.StatusCode == 200 && len(res.Body) > 100 {
		res.Headers["Content-Encoding"] = "gzip"
		fmt.Println("[CompressionFilter] Response compressed")
	}
}

// RateLimitFilter limits request rate
type RateLimitFilter struct {
	requestCount int
	maxRequests  int
}

func NewRateLimitFilter(maxRequests int) *RateLimitFilter {
	return &RateLimitFilter{
		requestCount: 0,
		maxRequests:  maxRequests,
	}
}

func (r *RateLimitFilter) Execute(req *Request, res *Response, chain *FilterChain) {
	r.requestCount++
	if r.requestCount > r.maxRequests {
		res.StatusCode = 429
		res.Body = "Too Many Requests"
		fmt.Println("[RateLimitFilter] Rate limit exceeded")
		return
	}
	fmt.Printf("[RateLimitFilter] Request %d/%d\n", r.requestCount, r.maxRequests)
	chain.Execute(req, res)
}

// ValidationFilter validates request
type ValidationFilter struct{}

func (v *ValidationFilter) Execute(req *Request, res *Response, chain *FilterChain) {
	if req.Method == "" || req.Path == "" {
		res.StatusCode = 400
		res.Body = "Bad Request: Missing method or path"
		return
	}
	fmt.Println("[ValidationFilter] Request validated")
	chain.Execute(req, res)
}

// Example Target Handler

// APIHandler is the final target handler
type APIHandler struct{}

func (h *APIHandler) Execute(req *Request, res *Response) {
	res.StatusCode = 200
	res.Body = fmt.Sprintf("API Response for %s %s", req.Method, req.Path)
	res.Headers = map[string]string{
		"Content-Type": "application/json",
	}
	fmt.Println("[APIHandler] Request processed")
}

// FilterManager manages filters
type FilterManager struct {
	chain *FilterChain
}

// NewFilterManager creates a filter manager
func NewFilterManager(target Target) *FilterManager {
	return &FilterManager{
		chain: NewFilterChain(target),
	}
}

// AddFilter adds a filter
func (fm *FilterManager) AddFilter(filter Filter) {
	fm.chain.AddFilter(filter)
}

// Process processes a request through the filter chain
func (fm *FilterManager) Process(req *Request) *Response {
	res := &Response{
		StatusCode: 200,
		Headers:    make(map[string]string),
	}
	fm.chain.Reset()
	fm.chain.Execute(req, res)
	return res
}
