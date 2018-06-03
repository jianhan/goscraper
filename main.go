package main

import (
	"github.com/jianhan/goscraper/output"
	"github.com/jianhan/goscraper/scraper"
	"github.com/sirupsen/logrus"
)

func main() {
	//scrapers := []scraper.Scraper{scraper.NewNCIX(true), scraper.NewMegabuyau(true), scraper.NewUmart(true)}
	scrapers := []scraper.Scraper{scraper.NewNCIX(true)}
	if err := run(scrapers...); err != nil {
		logrus.Error(err)
	}

}

func run(scrapers ...scraper.Scraper) error {
	for _, scraper := range scrapers {
		if err := scraper.Scrape(); err != nil {
			return err
		}
	}

	if err := output.NewFirestore(scrapers).Write(); err != nil {
		return err
	}

	//if err := output.NewJSONWriter(scrapers).Write(); err != nil {
	//	return err
	//}

	//if err := output.NewS3Writer(scrapers).Write(); err != nil {
	//	return err
	//}

	return nil
}
