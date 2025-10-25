// Package prototype demonstrates the Prototype pattern.
// It creates new objects by cloning existing prototypes, which is efficient
// when object creation is expensive or complex.
package prototype

import (
	"fmt"
	"time"
)

// Prototype interface defines the Clone method
type Prototype interface {
	Clone() Prototype
	GetInfo() string
}

// Stats represents character statistics
type Stats struct {
	Health    int
	Attack    int
	Defense   int
	Speed     int
	MagicPower int
}

// Equipment represents character equipment
type Equipment struct {
	Weapon string
	Armor  string
	Items  []string
}

// DeepCopy creates a deep copy of equipment
func (e *Equipment) DeepCopy() *Equipment {
	items := make([]string, len(e.Items))
	copy(items, e.Items)
	return &Equipment{
		Weapon: e.Weapon,
		Armor:  e.Armor,
		Items:  items,
	}
}

// Character represents a game character (prototype)
type Character struct {
	Name      string
	Class     string
	Stats     Stats
	Equipment *Equipment
	CreatedAt time.Time
}

// Clone creates a deep copy of the character
func (c *Character) Clone() Prototype {
	return &Character{
		Name:  c.Name,
		Class: c.Class,
		Stats: Stats{ // Structs are copied by value
			Health:     c.Stats.Health,
			Attack:     c.Stats.Attack,
			Defense:    c.Stats.Defense,
			Speed:      c.Stats.Speed,
			MagicPower: c.Stats.MagicPower,
		},
		Equipment: c.Equipment.DeepCopy(), // Deep copy for pointer
		CreatedAt: time.Now(),             // New timestamp for clone
	}
}

// GetInfo returns character information
func (c *Character) GetInfo() string {
	return fmt.Sprintf("%s (%s) - HP:%d ATK:%d DEF:%d SPD:%d MAG:%d | Weapon:%s",
		c.Name, c.Class,
		c.Stats.Health, c.Stats.Attack, c.Stats.Defense, c.Stats.Speed, c.Stats.MagicPower,
		c.Equipment.Weapon)
}

// ConfigurationProfile demonstrates cloning complex configurations
type ConfigurationProfile struct {
	Name     string
	Settings map[string]interface{}
	Metadata map[string]string
}

// Clone creates a deep copy of the configuration
func (cp *ConfigurationProfile) Clone() Prototype {
	// Deep copy maps
	settings := make(map[string]interface{})
	for k, v := range cp.Settings {
		settings[k] = v
	}

	metadata := make(map[string]string)
	for k, v := range cp.Metadata {
		metadata[k] = v
	}

	return &ConfigurationProfile{
		Name:     cp.Name,
		Settings: settings,
		Metadata: metadata,
	}
}

// GetInfo returns configuration information
func (cp *ConfigurationProfile) GetInfo() string {
	return fmt.Sprintf("Config: %s (Settings: %d, Metadata: %d)",
		cp.Name, len(cp.Settings), len(cp.Metadata))
}

// PrototypeRegistry manages a collection of prototypes
type PrototypeRegistry struct {
	prototypes map[string]Prototype
}

// NewPrototypeRegistry creates a new prototype registry
func NewPrototypeRegistry() *PrototypeRegistry {
	return &PrototypeRegistry{
		prototypes: make(map[string]Prototype),
	}
}

// Register adds a prototype to the registry
func (pr *PrototypeRegistry) Register(key string, prototype Prototype) {
	pr.prototypes[key] = prototype
}

// Get retrieves and clones a prototype
func (pr *PrototypeRegistry) Get(key string) Prototype {
	if prototype, exists := pr.prototypes[key]; exists {
		return prototype.Clone()
	}
	return nil
}

// List returns all registered prototype keys
func (pr *PrototypeRegistry) List() []string {
	keys := make([]string, 0, len(pr.prototypes))
	for k := range pr.prototypes {
		keys = append(keys, k)
	}
	return keys
}

// Pre-configured prototypes

// NewWarriorPrototype creates a warrior character prototype
func NewWarriorPrototype() *Character {
	return &Character{
		Name:  "Warrior",
		Class: "Warrior",
		Stats: Stats{
			Health:     100,
			Attack:     20,
			Defense:    15,
			Speed:      10,
			MagicPower: 5,
		},
		Equipment: &Equipment{
			Weapon: "Sword",
			Armor:  "Heavy Armor",
			Items:  []string{"Health Potion", "Shield"},
		},
		CreatedAt: time.Now(),
	}
}

// NewMagePrototype creates a mage character prototype
func NewMagePrototype() *Character {
	return &Character{
		Name:  "Mage",
		Class: "Mage",
		Stats: Stats{
			Health:     60,
			Attack:     10,
			Defense:    5,
			Speed:      12,
			MagicPower: 30,
		},
		Equipment: &Equipment{
			Weapon: "Staff",
			Armor:  "Robes",
			Items:  []string{"Mana Potion", "Spell Book"},
		},
		CreatedAt: time.Now(),
	}
}

// NewRoguePrototype creates a rogue character prototype
func NewRoguePrototype() *Character {
	return &Character{
		Name:  "Rogue",
		Class: "Rogue",
		Stats: Stats{
			Health:     80,
			Attack:     18,
			Defense:    8,
			Speed:      20,
			MagicPower: 8,
		},
		Equipment: &Equipment{
			Weapon: "Daggers",
			Armor:  "Light Armor",
			Items:  []string{"Poison Vial", "Lockpick"},
		},
		CreatedAt: time.Now(),
	}
}

// NewDevelopmentConfig creates a development configuration prototype
func NewDevelopmentConfig() *ConfigurationProfile {
	return &ConfigurationProfile{
		Name: "Development",
		Settings: map[string]interface{}{
			"debug":    true,
			"loglevel": "debug",
			"port":     8080,
			"cache":    false,
		},
		Metadata: map[string]string{
			"environment": "dev",
			"region":      "local",
		},
	}
}

// NewProductionConfig creates a production configuration prototype
func NewProductionConfig() *ConfigurationProfile {
	return &ConfigurationProfile{
		Name: "Production",
		Settings: map[string]interface{}{
			"debug":    false,
			"loglevel": "error",
			"port":     443,
			"cache":    true,
		},
		Metadata: map[string]string{
			"environment": "prod",
			"region":      "us-east-1",
		},
	}
}

// InitializeRegistry creates and populates a prototype registry
func InitializeRegistry() *PrototypeRegistry {
	registry := NewPrototypeRegistry()

	// Register character prototypes
	registry.Register("warrior", NewWarriorPrototype())
	registry.Register("mage", NewMagePrototype())
	registry.Register("rogue", NewRoguePrototype())

	// Register configuration prototypes
	registry.Register("dev-config", NewDevelopmentConfig())
	registry.Register("prod-config", NewProductionConfig())

	return registry
}
