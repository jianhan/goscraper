package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jianhan/goscraper/output"
	"github.com/jianhan/goscraper/scraper"
	"github.com/sirupsen/logrus"
)

func main() {
	scrapers := []scraper.Scraper{scraper.NewNCIX(false), scraper.NewMegabuyau(false), scraper.NewUmart(false)}
	if err := run(scrapers...); err != nil {
		logrus.Error(err)
	}

	if err := output.NewJSONWriter(scrapers).Write(); err != nil {
		spew.Dump(err)
		logrus.Error(err)
	}

	if err := output.NewS3Writer(scrapers).Write(); err != nil {
		spew.Dump(err)
		logrus.Error(err)
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
