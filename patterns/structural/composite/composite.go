package main

import "fmt"

// Когда использовать: Когда нужно работать с иерархической структурой объектов одинаково, независимо от того, является ли объект индивидуальным или составным.

// Компонент
type Component interface {
	search(string)
}

// Лист
type File struct {
	name string
}

func (f *File) search(keyword string) {
	fmt.Printf("Поиск ключевого слова %s в файле %s\n", keyword, f.name)
}

// Составной объект
type Folder struct {
	name       string
	components []Component
}

func (f *Folder) search(keyword string) {
	fmt.Printf("Рекурсивный поиск ключевого слова %s в папке %s\n", keyword, f.name)
	for _, composite := range f.components {
		composite.search(keyword)
	}
}

func (f *Folder) add(c Component) {
	f.components = append(f.components, c)
}

func main() {
	file1 := &File{name: "file1"}
	file2 := &File{name: "file2"}

	folder := &Folder{name: "folder1"}
	folder.add(file1)
	folder.add(file2)

	folder.search("test")
}
