package proxy

import "fmt"

func Example1_LazyLoading() {
	fmt.Println("\n=== Example 1: Lazy Loading Proxy ===")
	image := NewProxyImage("photo.jpg")
	fmt.Println("Image created but not loaded yet")
	image.Display()
	fmt.Println("\nDisplaying again (already loaded):")
	image.Display()
}
