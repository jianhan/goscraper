package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jianhan/goscraper/scraper"
)

func main() {
	//ncix := scraper.NewNCIX(true)
	//megabuyAU := scraper.NewMegabuyau(true)
	//umart := scraper.NewUmart(true)
	//if err := run(ncix, megabuyAU, umart); err != nil {
	//	logrus.Error(err)
	//}
	//if err := scraper.OutputJSONData(umart, megabuyAU, ncix); err != nil {
	//	logrus.Warn(err)
	//}

	err := scraper.UploadS3()
	spew.Dump(err)
}

func run(scrapers ...scraper.Scraper) error {
	for _, scraper := range scrapers {
		if err := scraper.Scrape(); err != nil {
			return err
		}
	}

	return nil
}
