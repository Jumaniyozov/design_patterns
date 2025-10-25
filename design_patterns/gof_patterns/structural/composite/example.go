package composite

import "fmt"

func Example1_FileSystem() {
	fmt.Println("\n=== Example 1: File System ===")

	root := NewComposite("root")
	home := NewComposite("home")
	user := NewComposite("user")

	file1 := NewLeaf("file1.txt")
	file2 := NewLeaf("file2.txt")
	file3 := NewLeaf("file3.txt")

	user.Add(file1)
	user.Add(file2)
	home.Add(user)
	home.Add(file3)
	root.Add(home)

	fmt.Println(root.Operation())
}

func Example2_OrganizationChart() {
	fmt.Println("\n=== Example 2: Organization Chart ===")

	ceo := NewComposite("CEO")
	cto := NewComposite("CTO")
	cfo := NewComposite("CFO")

	dev1 := NewLeaf("Developer 1")
	dev2 := NewLeaf("Developer 2")
	accountant := NewLeaf("Accountant")

	cto.Add(dev1)
	cto.Add(dev2)
	cfo.Add(accountant)
	ceo.Add(cto)
	ceo.Add(cfo)

	fmt.Println(ceo.Operation())
}
