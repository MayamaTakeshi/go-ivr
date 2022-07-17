package main

import (
	"fmt"
	// "errors"
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

func Test_stack(t *testing.T) {
	assert := assert.New(t)

	doc := etree.NewDocument()
	s := `<IVR>
	<Play>test.wav</Play>
	<Speak voice="suzan" language="en-US">Hello</Speak>
</IVR>`

	err := doc.ReadFromString(s)
	assert.Nil(err)

	root := doc.Root()
	elemList := root.ChildElements()

	var stack Stack

	stack.Push(elemList)
	stack.Push(elemList)

	topBefore, ok := stack.Top()
	assert.True(ok)
	fmt.Println(topBefore)

	_, topBefore = topBefore[0], topBefore[1:]

	fmt.Println(topBefore)

	topAfter1stPop, ok := stack.Pop()
	assert.NotEqual(topBefore, topAfter1stPop)
	assert.True(ok)
	assert.False(stack.IsEmpty())

	fmt.Println(topAfter1stPop)

	topAfter2ndPop, ok := stack.Pop()
	assert.NotEqual(topBefore, topAfter2ndPop)
	assert.True(ok)
	assert.True(stack.IsEmpty())

	fmt.Println(topAfter2ndPop)
}
