package helper

import (
	"bytes"
	"io/ioutil"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ClientMock struct {
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	var body = []byte("<html><a href='https://github.com'></a><a href='/topics'></a><a href='https://rayspock.com'></a></html>")
	return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader(body))}, nil
}

func TestFetch(t *testing.T) {
	const domain = "https://github.com"
	urls, err := Fetch(&ClientMock{}, domain)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		panic(err)
	}
	for _, u := range urls {
		fmt.Printf("Found: %s\n", u)
		assert.NotContains(t, u, "https")
		assert.NotContains(t, u, "rayspock.com")
	}
}
