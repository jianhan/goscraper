package main

import (
	"github.com/jianhan/goscraper/scraper"
	"github.com/sirupsen/logrus"
)

func main() {
	ncix := scraper.NewNCIX(true)
	megabuyAU := scraper.NewMegabuyau(true)
	umart := scraper.NewUmart(true)
	if err := run(ncix, megabuyAU, umart); err != nil {
		logrus.Error(err)
	}
	if err := scraper.OutputJSONData(umart, megabuyAU, ncix); err != nil {
		logrus.Warn(err)
	}
}

func run(scrapers ...scraper.Scraper) error {
	for _, scraper := range scrapers {
		if err := scraper.Scrape(); err != nil {
			return err
		}
	}

	return nil
}
