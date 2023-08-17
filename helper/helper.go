package helper

import (
	"fmt"
	"net/http"
	"regexp"

	"golang.org/x/net/html"
)

//go:generate mockgen -package=mock -destination=mock/helper.go -source=./helper.go

// HTTPClient ... interface for http.Client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type LinkFetcher func(client HTTPClient, url string) (urls []string, err error)

// Fetch ... fetches all the links from a given URL
func Fetch(client HTTPClient, url string) (urls []string, err error) {
	// Make HTTP Get request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	z := html.NewTokenizer(res.Body)
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
			fmt.Println("Something went wrong: ", err)
			return
		}
		if matched {
			ok = true
			href = _href
		}
	}
	return
}
