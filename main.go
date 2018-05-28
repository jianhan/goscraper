package main

import (
	"fmt"

	"github.com/jianhan/goscraper/scraper"
)

func main() {
	ncix := scraper.NewNCIXScrapper()
	ncix.Scrape()
	fmt.Println(ncix.Categories())
}
