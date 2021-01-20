package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/rayspock/go-web-crawler/helper"
)

var domainPtr *string
var depthPtr *int
var out io.Writer = os.Stdout

func main() {

	domainPtr = flag.String("website", "https://github.com", "Website URL")
	depthPtr = flag.Int("depth", 1, "Maximum of depth to crawl")

	flag.Parse()

	crawl(*domainPtr, *depthPtr)
	fmt.Println("Fetching stats\n--------------")
	for url, err := range fetched.m {
		if err != nil {
			fmt.Fprintf(out, "%v failed: %v\n", url, err)
		} else {
			fmt.Fprintf(out, "%v was fetched\n", url)
		}
	}
}

var fetched = struct {
	m map[string]error
	sync.Mutex
}{m: make(map[string]error)}

var errLoading = errors.New("url load in progress")

func crawl(url string, depth int) {
	if depth <= 0 {
		fmt.Printf("<- Done with %v, depth 0.\n", url)
		return
	}

	fetched.Lock()
	if _, ok := fetched.m[url]; ok {
		fetched.Unlock()
		fmt.Printf("<- Done with %v, already fetched.\n", url)
		return
	}
	fetched.m[url] = errLoading
	fetched.Unlock()

	urls, err := helper.Fetch(&http.Client{}, url)

	fetched.Lock()
	fetched.m[url] = err
	fetched.Unlock()

	if err != nil {
		fmt.Printf("<- Error on %v: %v\n", url, err)
		return
	}
	fmt.Printf("Found: %s\n", url)

	done := make(chan bool)
	for i, u := range urls {
		fmt.Printf("-> Crawling child %v/%v of %v : %v.\n", i, len(urls), url, u)
		go func(url string) {
			crawl(*domainPtr + url, depth-1)
			done <- true
		}(u)
	}
	for i, u := range urls {
		fmt.Printf("<- [%v] %v/%v Waiting for child %v.\n", url, i, len(urls), u)
		<-done
	}
	fmt.Printf("<- Done with %v\n", url)
}
