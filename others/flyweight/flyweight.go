// Package flyweight demonstrates the Flyweight pattern.
// It minimizes memory usage by sharing common state among many objects,
// separating intrinsic (shared) from extrinsic (unique) state.
package flyweight

import (
	"fmt"
	"sync"
)

// Example 1: Text Editor with Character Glyphs

// CharacterStyle represents intrinsic state (shared among characters)
type CharacterStyle struct {
	font   string
	size   int
	color  string
	bold   bool
	italic bool
}

// Render renders a character with this style at a specific position
func (cs *CharacterStyle) Render(char rune, x, y int) string {
	style := ""
	if cs.bold {
		style += "B"
	}
	if cs.italic {
		style += "I"
	}
	if style == "" {
		style = "Regular"
	}
	return fmt.Sprintf("'%c' at (%d,%d) [%s %dpt %s %s]",
		char, x, y, cs.font, cs.size, cs.color, style)
}

// GetInfo returns style information
func (cs *CharacterStyle) GetInfo() string {
	return fmt.Sprintf("%s %dpt %s", cs.font, cs.size, cs.color)
}

// StyleFactory manages shared CharacterStyle flyweights
type StyleFactory struct {
	styles map[string]*CharacterStyle
	mu     sync.RWMutex
}

// NewStyleFactory creates a style factory
func NewStyleFactory() *StyleFactory {
	return &StyleFactory{
		styles: make(map[string]*CharacterStyle),
	}
}

// GetStyle returns a shared style flyweight
func (sf *StyleFactory) GetStyle(font string, size int, color string, bold, italic bool) *CharacterStyle {
	key := fmt.Sprintf("%s-%d-%s-%v-%v", font, size, color, bold, italic)

	sf.mu.RLock()
	if style, exists := sf.styles[key]; exists {
		sf.mu.RUnlock()
		return style
	}
	sf.mu.RUnlock()

	sf.mu.Lock()
	defer sf.mu.Unlock()

	// Double-check after acquiring write lock
	if style, exists := sf.styles[key]; exists {
		return style
	}

	style := &CharacterStyle{
		font:   font,
		size:   size,
		color:  color,
		bold:   bold,
		italic: italic,
	}
	sf.styles[key] = style
	return style
}

// GetStyleCount returns number of unique styles created
func (sf *StyleFactory) GetStyleCount() int {
	sf.mu.RLock()
	defer sf.mu.RUnlock()
	return len(sf.styles)
}

// Character represents a character in the document with extrinsic state
type Character struct {
	char  rune
	x, y  int
	style *CharacterStyle // Reference to shared flyweight
}

// NewCharacter creates a character with a shared style
func NewCharacter(char rune, x, y int, style *CharacterStyle) *Character {
	return &Character{
		char:  char,
		x:     x,
		y:     y,
		style: style,
	}
}

// Render renders this character
func (c *Character) Render() string {
	return c.style.Render(c.char, c.x, c.y)
}

// Example 2: Game Forest with Shared Tree Types

// TreeType represents intrinsic state (shared among trees of same species)
type TreeType struct {
	name       string
	texture    string
	mesh       string
	leafColor  string
	barkColor  string
}

// Render renders a tree of this type at a specific location
func (tt *TreeType) Render(x, y float64, height float64) string {
	return fmt.Sprintf("Tree[%s] at (%.1f,%.1f) height:%.1fm | Texture:%s Leaves:%s Bark:%s",
		tt.name, x, y, height, tt.texture, tt.leafColor, tt.barkColor)
}

// GetInfo returns tree type information
func (tt *TreeType) GetInfo() string {
	return fmt.Sprintf("%s (Texture:%s)", tt.name, tt.texture)
}

// TreeFactory manages shared TreeType flyweights
type TreeFactory struct {
	treeTypes map[string]*TreeType
	mu        sync.RWMutex
}

// NewTreeFactory creates a tree factory
func NewTreeFactory() *TreeFactory {
	return &TreeFactory{
		treeTypes: make(map[string]*TreeType),
	}
}

// GetTreeType returns a shared tree type flyweight
func (tf *TreeFactory) GetTreeType(name, texture, mesh, leafColor, barkColor string) *TreeType {
	key := name

	tf.mu.RLock()
	if treeType, exists := tf.treeTypes[key]; exists {
		tf.mu.RUnlock()
		return treeType
	}
	tf.mu.RUnlock()

	tf.mu.Lock()
	defer tf.mu.Unlock()

	// Double-check
	if treeType, exists := tf.treeTypes[key]; exists {
		return treeType
	}

	treeType := &TreeType{
		name:      name,
		texture:   texture,
		mesh:      mesh,
		leafColor: leafColor,
		barkColor: barkColor,
	}
	tf.treeTypes[key] = treeType
	return treeType
}

// GetTreeTypeCount returns number of unique tree types
func (tf *TreeFactory) GetTreeTypeCount() int {
	tf.mu.RLock()
	defer tf.mu.RUnlock()
	return len(tf.treeTypes)
}

// Tree represents a tree instance with extrinsic state
type Tree struct {
	x, y     float64
	height   float64
	treeType *TreeType // Reference to shared flyweight
}

// NewTree creates a tree with a shared type
func NewTree(x, y, height float64, treeType *TreeType) *Tree {
	return &Tree{
		x:        x,
		y:        y,
		height:   height,
		treeType: treeType,
	}
}

// Render renders this tree
func (t *Tree) Render() string {
	return t.treeType.Render(t.x, t.y, t.height)
}

// Forest manages a collection of trees
type Forest struct {
	trees       []*Tree
	treeFactory *TreeFactory
}

// NewForest creates a forest with a tree factory
func NewForest() *Forest {
	return &Forest{
		trees:       make([]*Tree, 0),
		treeFactory: NewTreeFactory(),
	}
}

// PlantTree plants a tree in the forest
func (f *Forest) PlantTree(x, y, height float64, name, texture, mesh, leafColor, barkColor string) {
	treeType := f.treeFactory.GetTreeType(name, texture, mesh, leafColor, barkColor)
	tree := NewTree(x, y, height, treeType)
	f.trees = append(f.trees, tree)
}

// RenderAll renders all trees in the forest
func (f *Forest) RenderAll() []string {
	results := make([]string, len(f.trees))
	for i, tree := range f.trees {
		results[i] = tree.Render()
	}
	return results
}

// GetStats returns forest statistics
func (f *Forest) GetStats() string {
	return fmt.Sprintf("Forest: %d trees, %d unique species (%.1f%% memory saved)",
		len(f.trees),
		f.treeFactory.GetTreeTypeCount(),
		100.0*(1.0-float64(f.treeFactory.GetTreeTypeCount())/float64(len(f.trees))))
}
