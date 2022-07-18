package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_keyValueString2Map(t *testing.T) {
	assert := assert.New(t)

	m := make(map[string]string)

	err := keyValueString2Map(m, "xml_url=http://abc.com;domain=test.com", ";", "=")
	assert.Nil(err)

	assert.Equal(m["xml_url"], "http://abc.com")
	assert.Equal(m["domain"], "test.com")
}
