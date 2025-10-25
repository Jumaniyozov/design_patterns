package flyweight

import "fmt"

type TreeType struct {
	name    string
	color   string
	texture string
}

func (t *TreeType) Draw(x, y int) {
	fmt.Printf("Drawing %s tree at (%d,%d) with color %s\n", t.name, x, y, t.color)
}

type TreeFactory struct {
	treeTypes map[string]*TreeType
}

func NewTreeFactory() *TreeFactory {
	return &TreeFactory{treeTypes: make(map[string]*TreeType)}
}

func (f *TreeFactory) GetTreeType(name, color, texture string) *TreeType {
	key := name + color + texture
	if treeType, exists := f.treeTypes[key]; exists {
		return treeType
	}
	treeType := &TreeType{name: name, color: color, texture: texture}
	f.treeTypes[key] = treeType
	return treeType
}

type Tree struct {
	x        int
	y        int
	treeType *TreeType
}

func (t *Tree) Draw() {
	t.treeType.Draw(t.x, t.y)
}
