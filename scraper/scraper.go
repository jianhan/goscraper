package scraper

type Scraper interface {
	Scrape() error
	Categories() []Category
}

type Category struct {
	text string
	href string
}
