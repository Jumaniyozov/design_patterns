package main

import (
	"fmt"
	"sync"
)

type singleton struct{}

func (s *singleton) Print() {
	fmt.Println("instance")
}

var instance *singleton
var once sync.Once

func GetInstance() *singleton {
	once.Do(func() {
		instance = &singleton{}
	})
	return instance
}

func main() {
	inst := GetInstance()
	inst.Print()
	inst.Print()
	inst.Print()
}
