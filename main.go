package main

import "github.com/jianhan/goscraper/scraper"

func main() {
	ncix := scraper.NewNCIXScrapper()
	ncix.Scrape()
}
