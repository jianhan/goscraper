package main

import (
	"github.com/jianhan/goscraper/scraper"
	"github.com/sirupsen/logrus"
)

func main() {
	//ncix := scraper.NewNCIXScrapper()
	//if err := ncix.Scrape(); err != nil {
	//	panic(err)
	//}
	//if err := scraper.OutputJSONData(ncix); err != nil {
	//	panic(err)
	//}

	megabuyAU := scraper.NewMegabuyau()
	if err := megabuyAU.Scrape(); err != nil {
		panic(err)
	}
	if err := scraper.OutputJSONData(megabuyAU); err != nil {
		logrus.Warn(err)
	}
}
