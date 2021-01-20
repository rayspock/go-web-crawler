package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebCrawler(t *testing.T) {
	os.Args = []string{"cmd","-depth", "2", "-website", "https://github.com"}
	out = bytes.NewBuffer(nil)
	main()

	actual := out.(*bytes.Buffer).String()
	assert.Contains(t, actual, "https://github.com/topics")
}