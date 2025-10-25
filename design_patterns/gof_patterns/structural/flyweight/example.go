package flyweight

import "fmt"

func Example1_TreeForest() {
	fmt.Println("\n=== Example 1: Tree Forest ===")
	factory := NewTreeFactory()

	oak := factory.GetTreeType("Oak", "Green", "Rough")
	pine := factory.GetTreeType("Pine", "Dark Green", "Smooth")

	trees := []*Tree{
		{10, 20, oak},
		{30, 40, oak},
		{50, 60, pine},
		{70, 80, oak},
	}

	for _, tree := range trees {
		tree.Draw()
	}

	fmt.Printf("\nShared %d tree types instead of %d objects\n", len(factory.treeTypes), len(trees))
}
