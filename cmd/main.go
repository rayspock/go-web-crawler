package main

import (
	"fmt"
	"github.com/rayspock/go-web-crawler/crawler"
	"github.com/rayspock/go-web-crawler/helper"
	flag "github.com/spf13/pflag"
	"os"
)

func main() {
	baseUrl := flag.StringP("url", "u", "", "Website URL")
	maximumOfDepth := flag.IntP("depth", "d", 1, "Maximum of depth to crawl")
	outputFileName := flag.StringP("output", "o", "", "Output file name")

	flag.Parse()

	if *baseUrl == "" {
		flag.Usage()
		fmt.Println("Please provide a website URL")
		os.Exit(1)
	}
	if *outputFileName == "" {
		flag.Usage()
		fmt.Println("Please provide an output file name")
		os.Exit(1)
	}

	// Open output file
	f, err := os.Create(*outputFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	c := crawler.New(helper.Fetch)
	c.Crawl(f, *baseUrl, *maximumOfDepth)

	fmt.Println("Fetching stats\n--------------")
	out := os.Stdout
	for url, err := range c.Fetched() {
		if err != nil {
			fmt.Fprintf(out, "%v failed: %v\n", url, err)
		} else {
			fmt.Fprintf(out, "%v was fetched\n", url)
		}
	}
}
