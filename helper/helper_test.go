package helper_test

import (
	"github.com/rayspock/go-web-crawler/helper"
	"github.com/rayspock/go-web-crawler/helper/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestHelperSuite(t *testing.T) {
	suite.Run(t, new(HelperSuite))
}

type HelperSuite struct {
	suite.Suite

	client *mock.MockHTTPClient
}

func (s *HelperSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.client = mock.NewMockHTTPClient(ctrl)
}

func (s *HelperSuite) TestFetch() {
	const (
		baseUrl  = "https://github.com"
		htmlBody = "<html><a href='https://github.com'></a><a href='/topics'></a><a href='https://rayspock.com'></a></html>"
	)
	s.client.ExpectHtmlPageReturn(baseUrl, htmlBody)
	urls, err := helper.Fetch(s.client, baseUrl)
	s.NoError(err)
	s.Len(urls, 1)
	s.Contains(urls, "/topics")
	s.NotContains(urls, "rayspock.com")
	s.NotContains(urls, "https")
	s.NotContains(urls, "github.com")
}
