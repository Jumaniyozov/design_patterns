// Package memento demonstrates the Memento pattern.
// It captures and externalizes object state without violating encapsulation,
// enabling undo/redo and state restoration functionality.
package memento

import (
	"fmt"
	"time"
)

// EditorMemento stores editor state
type EditorMemento struct {
	content  string
	cursor   int
	saveTime time.Time
}

// TextEditor is the originator
type TextEditor struct {
	content string
	cursor  int
}

func NewTextEditor() *TextEditor {
	return &TextEditor{
		content: "",
		cursor:  0,
	}
}

// Write adds text at cursor position
func (e *TextEditor) Write(text string) {
	e.content = e.content[:e.cursor] + text + e.content[e.cursor:]
	e.cursor += len(text)
}

// SetContent replaces content
func (e *TextEditor) SetContent(content string) {
	e.content = content
	e.cursor = len(content)
}

// GetContent returns current content
func (e *TextEditor) GetContent() string {
	return e.content
}

// Save creates a memento
func (e *TextEditor) Save() *EditorMemento {
	return &EditorMemento{
		content:  e.content,
		cursor:   e.cursor,
		saveTime: time.Now(),
	}
}

// Restore restores from memento
func (e *TextEditor) Restore(m *EditorMemento) {
	e.content = m.content
	e.cursor = m.cursor
}

// History is the caretaker
type History struct {
	mementos []*EditorMemento
	current  int
}

func NewHistory() *History {
	return &History{
		mementos: make([]*EditorMemento, 0),
		current:  -1,
	}
}

// Save adds a memento to history
func (h *History) Save(m *EditorMemento) {
	// Remove any future states if we're not at the end
	h.mementos = h.mementos[:h.current+1]
	h.mementos = append(h.mementos, m)
	h.current++
}

// Undo returns previous state
func (h *History) Undo() *EditorMemento {
	if h.current > 0 {
		h.current--
		return h.mementos[h.current]
	}
	return nil
}

// Redo returns next state
func (h *History) Redo() *EditorMemento {
	if h.current < len(h.mementos)-1 {
		h.current++
		return h.mementos[h.current]
	}
	return nil
}

// CanUndo checks if undo is possible
func (h *History) CanUndo() bool {
	return h.current > 0
}

// CanRedo checks if redo is possible
func (h *History) CanRedo() bool {
	return h.current < len(h.mementos)-1
}

// Game state example

// GameMemento stores game state
type GameMemento struct {
	level      int
	health     int
	score      int
	checkpoint string
}

// GameState manages game state
type GameState struct {
	level      int
	health     int
	score      int
	checkpoint string
}

func NewGameState() *GameState {
	return &GameState{
		level:      1,
		health:     100,
		score:      0,
		checkpoint: "start",
	}
}

// Update modifies game state
func (g *GameState) Update(level, health, score int, checkpoint string) {
	g.level = level
	g.health = health
	g.score = score
	g.checkpoint = checkpoint
}

// SaveCheckpoint creates a memento
func (g *GameState) SaveCheckpoint() *GameMemento {
	return &GameMemento{
		level:      g.level,
		health:     g.health,
		score:      g.score,
		checkpoint: g.checkpoint,
	}
}

// LoadCheckpoint restores from memento
func (g *GameState) LoadCheckpoint(m *GameMemento) {
	g.level = m.level
	g.health = m.health
	g.score = m.score
	g.checkpoint = m.checkpoint
}

// GetStatus returns current status
func (g *GameState) GetStatus() string {
	return fmt.Sprintf("Level:%d HP:%d Score:%d Checkpoint:%s",
		g.level, g.health, g.score, g.checkpoint)
}

// CheckpointManager manages game checkpoints
type CheckpointManager struct {
	checkpoints map[string]*GameMemento
}

func NewCheckpointManager() *CheckpointManager {
	return &CheckpointManager{
		checkpoints: make(map[string]*GameMemento),
	}
}

// SaveCheckpoint saves a named checkpoint
func (cm *CheckpointManager) SaveCheckpoint(name string, memento *GameMemento) {
	cm.checkpoints[name] = memento
}

// LoadCheckpoint loads a named checkpoint
func (cm *CheckpointManager) LoadCheckpoint(name string) *GameMemento {
	return cm.checkpoints[name]
}

// ListCheckpoints returns all checkpoint names
func (cm *CheckpointManager) ListCheckpoints() []string {
	names := make([]string, 0, len(cm.checkpoints))
	for name := range cm.checkpoints {
		names = append(names, name)
	}
	return names
}
