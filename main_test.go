package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestKeyValueString2Map(t *testing.T) {
	assert := assert.New(t)

	m := KeyValueString2Map("xml_url=http://abc.com;domain=test.com", ";", "=")	

	assert.Equal(m["xml_url"], "http://abc.com")
	assert.Equal(m["domain"], "test.com")
}
