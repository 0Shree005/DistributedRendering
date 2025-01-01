// package filechangedetection

package main

import (
	"sync"
)

type Item interface{}

type Stack struct {
	items []Item
	mutex sync.Mutex
}

func (stack *Stack) Push(item Item) {
	stack.mutex.Lock()
	stack.items = append(stack.items, item)
	stack.mutex.Unlock()
}

func (stack *Stack) Pop() Item {
	stack.mutex.Lock()
	if len(stack.items) == 0 {
		return nil
	}
	lastItem := stack.items[len(stack.items)-1]
	stack.items = stack.items[:len(stack.items)-1]
	stack.mutex.Unlock()
	return lastItem
}

func (stack *Stack) Peek() Item {
	stack.mutex.Lock()
	if len(stack.items) == 0 {
		return nil
	}
	stack.mutex.Unlock()
	return stack.items[len(stack.items)-1]
}

func (stack *Stack) IsEmpty() bool {
	stack.mutex.Lock()
	stack.mutex.Unlock()
	return len(stack.items) == 0
}

func (stack *Stack) Reset() {
	stack.mutex.Lock()
	stack.items = nil
	stack.mutex.Unlock()
}

func (stack *Stack) Dump() []Item {
	stack.mutex.Lock()
	var copiedStack = make([]Item, len(stack.items))
	copy(copiedStack, stack.items)
	stack.mutex.Unlock()
	return copiedStack
}
