# Web Crawler - Golang

Fetch URLs in parallel.

```bash
# Parameters
-depth int
    	Maximum of depth to crawl (default 1)
-website string
    Website URL (default "https://github.com")
```

## Quick Start

```bash
$ go run main.go -depth 2 -website https://github.com
```