package main

import "fmt"

// Когда использовать: Когда нужно создавать новые объекты путем копирования существующих.

type Inode interface {
	print(string)
	clone() Inode
}

type File struct {
	name string
}

func (f *File) print(indentation string) {
	fmt.Println(indentation + f.name)
}

func (f *File) clone() Inode {
	return &File{name: f.name + "_clone"}
}

type Folder struct {
	children []Inode
	name     string
}

func (f *Folder) print(indentation string) {
	fmt.Println(indentation + f.name)
	for _, i := range f.children {
		i.print(indentation + indentation)
	}
}

func (f *Folder) clone() Inode {
	cloneFolder := &Folder{name: f.name + "_clone"}
	var tempChildren []Inode
	for _, i := range f.children {
		tempChildren = append(tempChildren, i.clone())
	}
	cloneFolder.children = tempChildren
	return cloneFolder
}

func main() {
	file1 := &File{name: "File1"}
	folder1 := &Folder{
		children: []Inode{file1},
		name:     "Folder1",
	}

	cloneFolder := folder1.clone()
	cloneFolder.print("  ")
}
