// Package stranglerfig demonstrates the Strangler Fig pattern.
// It enables gradual migration from legacy systems by incrementally replacing
// functionality with new services while both systems run in parallel.
package stranglerfig

import (
	"fmt"
	"time"
)

// Request represents a service request
type Request struct {
	ID      string
	Feature string
	UserID  string
	Data    map[string]interface{}
}

// Response represents a service response
type Response struct {
	Data   interface{}
	Source string // "legacy" or "modern"
	Time   time.Duration
}

// Service interface for both legacy and modern implementations
type Service interface {
	Handle(req *Request) (*Response, error)
	Name() string
}

// LegacyService represents the old system
type LegacyService struct {
	name string
}

// NewLegacyService creates a legacy service
func NewLegacyService() *LegacyService {
	return &LegacyService{name: "LegacyService"}
}

func (l *LegacyService) Handle(req *Request) (*Response, error) {
	start := time.Now()

	// Simulate slow legacy processing
	time.Sleep(200 * time.Millisecond)

	result := fmt.Sprintf("Legacy result for %s", req.Feature)

	return &Response{
		Data:   result,
		Source: "legacy",
		Time:   time.Since(start),
	}, nil
}

func (l *LegacyService) Name() string {
	return l.name
}

// ModernService represents the new system
type ModernService struct {
	name string
}

// NewModernService creates a modern service
func NewModernService() *ModernService {
	return &ModernService{name: "ModernService"}
}

func (m *ModernService) Handle(req *Request) (*Response, error) {
	start := time.Now()

	// Simulate fast modern processing
	time.Sleep(50 * time.Millisecond)

	result := fmt.Sprintf("Modern result for %s (improved!)", req.Feature)

	return &Response{
		Data:   result,
		Source: "modern",
		Time:   time.Since(start),
	}, nil
}

func (m *ModernService) Name() string {
	return m.name
}

// RoutingStrategy determines which service to use
type RoutingStrategy interface {
	ShouldUseModern(req *Request) bool
}

// FeatureBasedRouter routes based on feature flags
type FeatureBasedRouter struct {
	migratedFeatures map[string]bool
}

// NewFeatureBasedRouter creates a feature-based router
func NewFeatureBasedRouter() *FeatureBasedRouter {
	return &FeatureBasedRouter{
		migratedFeatures: make(map[string]bool),
	}
}

// MigrateFeature marks a feature as migrated
func (r *FeatureBasedRouter) MigrateFeature(feature string) {
	r.migratedFeatures[feature] = true
}

func (r *FeatureBasedRouter) ShouldUseModern(req *Request) bool {
	return r.migratedFeatures[req.Feature]
}

// PercentageBasedRouter gradually increases traffic to modern service
type PercentageBasedRouter struct {
	percentage int // 0-100
	counter    int
}

// NewPercentageBasedRouter creates a percentage-based router
func NewPercentageBasedRouter(percentage int) *PercentageBasedRouter {
	return &PercentageBasedRouter{
		percentage: percentage,
		counter:    0,
	}
}

func (r *PercentageBasedRouter) ShouldUseModern(req *Request) bool {
	r.counter++
	return (r.counter % 100) < r.percentage
}

// SetPercentage updates the routing percentage
func (r *PercentageBasedRouter) SetPercentage(percentage int) {
	r.percentage = percentage
}

// UserBasedRouter routes based on user segments
type UserBasedRouter struct {
	betaUsers map[string]bool
}

// NewUserBasedRouter creates a user-based router
func NewUserBasedRouter() *UserBasedRouter {
	return &UserBasedRouter{
		betaUsers: make(map[string]bool),
	}
}

// AddBetaUser adds a user to beta program
func (r *UserBasedRouter) AddBetaUser(userID string) {
	r.betaUsers[userID] = true
}

func (r *UserBasedRouter) ShouldUseModern(req *Request) bool {
	return r.betaUsers[req.UserID]
}

// StranglerProxy routes requests between legacy and modern services
type StranglerProxy struct {
	legacy Service
	modern Service
	router RoutingStrategy
	stats  *MigrationStats
}

// NewStranglerProxy creates a strangler proxy
func NewStranglerProxy(legacy, modern Service, router RoutingStrategy) *StranglerProxy {
	return &StranglerProxy{
		legacy: legacy,
		modern: modern,
		router: router,
		stats:  NewMigrationStats(),
	}
}

// Handle routes the request to appropriate service
func (p *StranglerProxy) Handle(req *Request) (*Response, error) {
	if p.router.ShouldUseModern(req) {
		p.stats.IncrementModern()
		fmt.Printf("[Proxy] Routing %s to modern service\n", req.Feature)
		return p.modern.Handle(req)
	}

	p.stats.IncrementLegacy()
	fmt.Printf("[Proxy] Routing %s to legacy service\n", req.Feature)
	return p.legacy.Handle(req)
}

// GetStats returns migration statistics
func (p *StranglerProxy) GetStats() *MigrationStats {
	return p.stats
}

// MigrationStats tracks migration progress
type MigrationStats struct {
	legacyRequests int
	modernRequests int
}

// NewMigrationStats creates migration stats
func NewMigrationStats() *MigrationStats {
	return &MigrationStats{}
}

// IncrementLegacy increments legacy request count
func (m *MigrationStats) IncrementLegacy() {
	m.legacyRequests++
}

// IncrementModern increments modern request count
func (m *MigrationStats) IncrementModern() {
	m.modernRequests++
}

// GetMigrationPercentage returns percentage of requests on modern system
func (m *MigrationStats) GetMigrationPercentage() float64 {
	total := m.legacyRequests + m.modernRequests
	if total == 0 {
		return 0
	}
	return float64(m.modernRequests) / float64(total) * 100
}

// Summary returns a summary string
func (m *MigrationStats) Summary() string {
	return fmt.Sprintf("Legacy: %d, Modern: %d, Migration: %.1f%%",
		m.legacyRequests,
		m.modernRequests,
		m.GetMigrationPercentage())
}

// MigrationManager manages the migration process
type MigrationManager struct {
	proxy *StranglerProxy
	router *FeatureBasedRouter
}

// NewMigrationManager creates a migration manager
func NewMigrationManager(legacy, modern Service) *MigrationManager {
	router := NewFeatureBasedRouter()
	proxy := NewStranglerProxy(legacy, modern, router)

	return &MigrationManager{
		proxy:  proxy,
		router: router,
	}
}

// MigrateFeature migrates a feature to modern service
func (mm *MigrationManager) MigrateFeature(feature string) {
	fmt.Printf("[Migration] Migrating feature: %s\n", feature)
	mm.router.MigrateFeature(feature)
}

// HandleRequest handles a request through the proxy
func (mm *MigrationManager) HandleRequest(req *Request) (*Response, error) {
	return mm.proxy.Handle(req)
}

// GetProgress returns migration progress
func (mm *MigrationManager) GetProgress() string {
	return mm.proxy.GetStats().Summary()
}

// GradualMigration demonstrates gradual traffic shifting
type GradualMigration struct {
	proxy      *StranglerProxy
	router     *PercentageBasedRouter
	percentage int
}

// NewGradualMigration creates a gradual migration
func NewGradualMigration(legacy, modern Service) *GradualMigration {
	router := NewPercentageBasedRouter(0)
	proxy := NewStranglerProxy(legacy, modern, router)

	return &GradualMigration{
		proxy:      proxy,
		router:     router,
		percentage: 0,
	}
}

// IncreaseTraffic increases traffic to modern service
func (gm *GradualMigration) IncreaseTraffic(amount int) {
	gm.percentage += amount
	if gm.percentage > 100 {
		gm.percentage = 100
	}
	gm.router.SetPercentage(gm.percentage)
	fmt.Printf("[Migration] Traffic to modern service: %d%%\n", gm.percentage)
}

// HandleRequest handles a request
func (gm *GradualMigration) HandleRequest(req *Request) (*Response, error) {
	return gm.proxy.Handle(req)
}

// GetProgress returns current progress
func (gm *GradualMigration) GetProgress() string {
	return fmt.Sprintf("Target: %d%% | Actual: %s",
		gm.percentage,
		gm.proxy.GetStats().Summary())
}
