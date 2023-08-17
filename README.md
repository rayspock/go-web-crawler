# Web Crawler

A web crawler to fetch all the links from a given website via go routines. 

## Getting Started

To get started with the repository, simply clone or download the code examples and open them in your favorite text
editor or IDE.

Before running the examples, make sure you have Go installed on your machine. You can download and install the latest
version of Go from the official website at https://golang.org.

To run the code, simply navigate to the directory containing the code example and run the following command:
```bash
# Run the code
$ go run cmd/main.go -u https://github.com -d 2 -o output.txt
```
```bash
# Parameters
  -d, --depth int       Maximum of depth to crawl (default 1)
  -o, --output string   Output file name
  -u, --url string      Website URL (default "https://github.com")
```


## Development 

```bash
# Run the tests
$ make test

# Generate mock
$ make generate
```
