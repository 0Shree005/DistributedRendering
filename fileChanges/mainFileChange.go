package main

import (
	"fmt"
)

func main() {
	var stack Stack

	fmt.Println("Stack:", stack.Dump())

	stack.Push(1)
	stack.Push(6)
	stack.Push("one")
	stack.Push(10)
	stack.Push(10.01)
	stack.Push("two")

	fmt.Println("Stack:", stack.Dump())

	fmt.Println("Last item:", stack.Peek())

	fmt.Println("Stack is empty:", stack.IsEmpty())

	fmt.Println("Last removed item:", stack.Pop())

	fmt.Println("Stack:", stack.Dump())

	stack.Reset()

	fmt.Println("Stack is empty:", stack.IsEmpty())

	fmt.Println("Last item:", stack.Peek())
}
