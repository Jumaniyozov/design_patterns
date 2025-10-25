package decorator

import "fmt"

// Example1_CoffeeDecorators demonstrates decorating coffee.
func Example1_CoffeeDecorators() {
	fmt.Println("\n=== Example 1: Coffee Decorators ===")

	coffee := &SimpleCoffee{}
	fmt.Printf("%s: $%.2f\n", coffee.Description(), coffee.Cost())

	coffeeWithMilk := NewMilkDecorator(coffee)
	fmt.Printf("%s: $%.2f\n", coffeeWithMilk.Description(), coffeeWithMilk.Cost())

	coffeeWithMilkAndSugar := NewSugarDecorator(coffeeWithMilk)
	fmt.Printf("%s: $%.2f\n", coffeeWithMilkAndSugar.Description(), coffeeWithMilkAndSugar.Cost())

	fancyCoffee := NewWhippedCreamDecorator(NewSugarDecorator(NewMilkDecorator(&SimpleCoffee{})))
	fmt.Printf("%s: $%.2f\n", fancyCoffee.Description(), fancyCoffee.Cost())
}

// Example2_DataSourceDecorators demonstrates file operations with decorators.
func Example2_DataSourceDecorators() {
	fmt.Println("\n=== Example 2: Data Source Decorators ===")

	file := NewFileDataSource("data.txt")
	file.WriteData("Hello, World!")
	file.ReadData()

	fmt.Println("\nWith encryption:")
	encrypted := NewEncryptionDecorator(NewFileDataSource("secure.txt"))
	encrypted.WriteData("Secret message")
	encrypted.ReadData()

	fmt.Println("\nWith encryption and compression:")
	encryptedCompressed := NewCompressionDecorator(NewEncryptionDecorator(NewFileDataSource("secure_compressed.txt")))
	encryptedCompressed.WriteData("Top secret data")
	encryptedCompressed.ReadData()
}

// Example3_DecoratorBenefits shows benefits.
func Example3_DecoratorBenefits() {
	fmt.Println("\n=== Example 3: Decorator Benefits ===")
	fmt.Println("✓ Add functionality dynamically")
	fmt.Println("✓ Combine multiple decorators")
	fmt.Println("✓ No class explosion")
	fmt.Println("✓ Single Responsibility Principle")
}
