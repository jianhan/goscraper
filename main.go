package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jianhan/goscraper/scraper"
	"github.com/sirupsen/logrus"
)

func main() {
	ncix := scraper.NewNCIXScrapper()
	if err := ncix.Scrape(); err != nil {
		logrus.Warn(err)
	}

	spew.Dump(ncix.Products())
}
