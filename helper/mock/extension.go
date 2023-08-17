package mock

import (
	"bytes"
	"errors"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
)

func (m *MockHTTPClient) ExpectHtmlPageReturn(url string, html string) {
	m.EXPECT().Do(gomock.Any()).DoAndReturn(func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == url {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(html)),
			}, nil
		}
		return nil, errors.New("unexpected url")
	})
}
