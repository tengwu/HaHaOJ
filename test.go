package main

import (
	"fmt"
)

type A interface {
	fuck()
}

func gg(a A) {
	fmt.Println(a)
}

type B struct {
	x int
	y string
}

func (b *B) fuck() {

}

func main() {
	gg(&B{x: 10, y: "hello"})
}
