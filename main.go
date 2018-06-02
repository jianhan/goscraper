package main

import (
	"github.com/jianhan/goscraper/scraper"
	"github.com/sirupsen/logrus"
)

func main() {
	ncix := scraper.NewNCIX(true)
	if err := ncix.Scrape(); err != nil {
		panic(err)
	}
	if err := scraper.OutputJSONData(ncix); err != nil {
		panic(err)
	}

	megabuyAU := scraper.NewMegabuyau(true)
	if err := megabuyAU.Scrape(); err != nil {
		panic(err)
	}
	if err := scraper.OutputJSONData(megabuyAU); err != nil {
		logrus.Warn(err)
	}

	umart := scraper.NewUmart(true)
	if err := umart.Scrape(); err != nil {
		panic(err)
	}
	if err := scraper.OutputJSONData(umart); err != nil {
		logrus.Warn(err)
	}
}
