package crawler_test

import (
	"github.com/rayspock/go-web-crawler/crawler"
	"github.com/rayspock/go-web-crawler/helper"
	"github.com/stretchr/testify/suite"
	"io"
	"sync"
	"testing"
)

func TestCrawlerSuite(t *testing.T) {
	suite.Run(t, new(CrawlerSuite))
}

type CrawlerSuite struct {
	suite.Suite

	client crawler.Crawler
	writer io.Writer
}

func (s *CrawlerSuite) SetupTest() {
	s.client = crawler.New(MockLinkFetcher([]string{"/", "/about", "/posts", "/posts/1"}, nil))
	s.writer = &MockWriter{}
}

func (s *CrawlerSuite) TestCrawl() {
	const baseUrl = "https://blog.rayspock.com"
	s.client.Crawl(s.writer, baseUrl, 2)
	s.Len(s.client.Fetched(), 5)
	s.Contains(s.client.Fetched(), baseUrl)
}

// MockWriter is a mock implementation of io.Writer
type MockWriter struct {
	mutex   sync.Mutex
	Written []byte
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Written = p
	return len(p), nil
}

// MockLinkFetcher is a mock implementation of helper.LinkFetcher
func MockLinkFetcher(urls []string, err error) helper.LinkFetcher {
	return func(client helper.HTTPClient, url string) ([]string, error) {
		return urls, err
	}
}
