package main

import (
	"errors"
	"testing"

	"github.com/beevik/etree"
	"github.com/stretchr/testify/assert"
)

func Test_validate_ok(t *testing.T) {
	assert := assert.New(t)

	doc := etree.NewDocument()
	s := `<IVR>
	<Play>test.wav</Play>
	<Speak voice="suzan" language="en-US">Hello</Speak>
</IVR>`

	err := doc.ReadFromString(s)
	assert.Nil(err)

	root := doc.Root()
	err = validate(root, 0)
	assert.Nil(err)
}

func Test_validate_invalid_child(t *testing.T) {
	assert := assert.New(t)

	doc := etree.NewDocument()
	s := `<IVR>
	<Play>test.wav</Play>
	<Speak voice="suzan" language="en-US"><Play/></Speak>
</IVR>`

	err := doc.ReadFromString(s)
	assert.Nil(err)

	root := doc.Root()
	err = validate(root, 0)
	assert.Equal(err, errors.New("Element Speak cannot contain Element Play"))
}

func Test_validate_unknown_element(t *testing.T) {
	assert := assert.New(t)

	doc := etree.NewDocument()
	s := `<IVR>
	<Play>test.wav</Play>
	<Speak voice="suzan" language="en-US">Hello</Speak>
	<BlaBla/>
</IVR>
	`
	err := doc.ReadFromString(s)
	assert.Nil(err)

	root := doc.Root()
	err = validate(root, 0)
	assert.Equal(err, errors.New("Unknown Element BlaBla"))
}

func Test_validate_unknown_attribute(t *testing.T) {
	assert := assert.New(t)

	doc := etree.NewDocument()
	s := `<IVR>
	<Play>test.wav</Play>
	<Speak voice="suzan" language="en-US">Hello</Speak>
	<GetDigits maxTries="5" dummyAttribute="10">
		<Play>menu.wav</Play>
	</GetDigits>
</IVR>
	`
	err := doc.ReadFromString(s)
	assert.Nil(err)

	root := doc.Root()
	err = validate(root, 0)
	assert.Equal(err, errors.New("Element GetDigits unknown Attribute dummyAttribute"))
}
