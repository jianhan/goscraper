package main

import (
	"github.com/jianhan/goscraper/scraper"
	"github.com/sirupsen/logrus"
)

func main() {
	ncix := scraper.NewNCIXScrapper()
	if err := ncix.Scrape(); err != nil {
		logrus.Warn(err)
	}
	if err := scraper.OutputJSONData(ncix); err != nil {
		logrus.Warn(err)
	}
}
