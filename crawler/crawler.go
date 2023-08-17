package crawler

import (
	"errors"
	"fmt"
	"github.com/rayspock/go-web-crawler/helper"
	"io"
	"net/http"
	"sync"
)

type Crawler interface {
	Crawl(wr io.Writer, baseUrl string, depth int)
	Fetched() map[string]error
}

var errLoading = errors.New("url load in progress")

type crawler struct {
	fetched     map[string]error
	linkFetcher helper.LinkFetcher
	sync.Mutex
}

func New(linkFetcher helper.LinkFetcher) Crawler {
	return &crawler{
		fetched:     make(map[string]error),
		linkFetcher: linkFetcher,
	}
}

func (c *crawler) Fetched() map[string]error {
	return c.fetched
}

func (c *crawler) Crawl(wr io.Writer, baseUrl string, depth int) {
	if depth <= 0 {
		fmt.Printf("<- Done with %v with depth 0.\n", baseUrl)
		writeln(wr, baseUrl)
		return
	}

	c.Lock()
	if _, ok := c.fetched[baseUrl]; ok {
		c.Unlock()
		fmt.Printf("<- Done with %v, already fetched.\n", baseUrl)
		return
	}
	c.fetched[baseUrl] = errLoading
	c.Unlock()

	paths, err := c.linkFetcher(&http.Client{}, baseUrl)

	c.Lock()
	c.fetched[baseUrl] = err
	c.Unlock()

	if err != nil {
		fmt.Printf("<- Error on %v: %v\n", baseUrl, err)
		return
	}
	fmt.Printf("Found: %s\n", baseUrl)

	done := make(chan bool)
	for i, path := range paths {
		fmt.Printf("-> Crawling child %v/%v of %v : %v.\n", i, len(paths), baseUrl, path)
		go func(path string) {
			c.Crawl(wr, baseUrl+path, depth-1)
			done <- true
		}(path)
	}
	for i, path := range paths {
		fmt.Printf("<- [%v] %v/%v Waiting for child %v.\n", baseUrl, i, len(paths), path)
		<-done
	}
	fmt.Printf("<- Done with %v.\n", baseUrl)
	writeln(wr, baseUrl)
}

func writeln(wr io.Writer, str string) {
	line := fmt.Sprintf("%v\n", str)
	if _, err := wr.Write([]byte(line)); err != nil {
		fmt.Printf("Error writing to file: %v", err)
	}
}
