package scraper

type Scraper interface {
	Scrape() error
}
