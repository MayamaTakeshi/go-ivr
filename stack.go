package main

import (
	"github.com/beevik/etree"
)

type ElementList []*etree.Element

type Stack []ElementList

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(elemList ElementList) {
	*s = append(*s, elemList) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() (ElementList, bool) {
	if s.IsEmpty() {
		return nil, false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

func (s *Stack) Top() (ElementList) {
	if s.IsEmpty() {
		return nil
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		return element
	}
}
