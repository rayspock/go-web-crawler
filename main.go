package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"regexp"
	"sync"

	"golang.org/x/net/html"
)

var domainPtr *string
var depthPtr *int

func main() {

	domainPtr = flag.String("website", "https://github.com", "Website URL")
	depthPtr = flag.Int("depth", 1, "Maximum of depth to crawl")

	flag.Parse()

	crawl(*domainPtr, *depthPtr)
	log.Println("Fetching stats\n--------------")
	for url, err := range fetched.m {
		if err != nil {
			log.Printf("%v failed: %v\n", url, err)
		} else {
			log.Printf("%v was fetched\n", url)
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
		log.Printf("<- Done with %v, depth 0.\n", url)
		return
	}

	fetched.Lock()
	if _, ok := fetched.m[url]; ok {
		fetched.Unlock()
		log.Printf("<- Done with %v, already fetched.\n", url)
		return
	}
	fetched.m[url] = errLoading
	fetched.Unlock()

	urls, err := fetch(url)

	fetched.Lock()
	fetched.m[url] = err
	fetched.Unlock()

	if err != nil {
		log.Printf("<- Error on %v: %v\n", url, err)
		return
	}
	log.Printf("Found: %s\n", url)

	done := make(chan bool)
	for i, u := range urls {
		log.Printf("-> Crawling child %v/%v of %v : %v.\n", i, len(urls), url, u)
		go func(url string) {
			crawl(url, depth-1)
			done <- true
		}(u)
	}
	for i, u := range urls {
		log.Printf("<- [%v] %v/%v Waiting for child %v.\n", url, i, len(urls), u)
		<-done
	}
	log.Printf("<- Done with %v\n", url)
}

//fetch ... return a slice of URLs found on that page
func fetch(url string) (urls []string, err error) {
	// Make HTTP Get request
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	z := html.NewTokenizer(response.Body)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			if ok, href := getSameOriginHref(&t); ok {
				urls = append(urls, href)
			}
		}
	}
}

func getHref(t *html.Token) (ok bool, href string) {
	for _, a := range (*t).Attr {
		if a.Key == "href" {
			ok = true
			href = a.Val
		}
	}
	return
}

func getSameOriginHref(t *html.Token) (ok bool, href string) {
	const exp = `^\/`
	if _ok, _href := getHref(t); _ok {
		matched, err := regexp.Match(exp, []byte(_href))
		if err != nil {
			ok = false
			log.Println("Something went wrong: ", err)
			return
		}
		if matched {
			ok = true
			href = *domainPtr + _href
		}
	}
	return
}
