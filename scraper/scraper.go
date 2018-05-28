package scraper

type Scraper interface {
	Scrape() error
	Categories() []Category
	Products() []Product
}

type Category struct {
	name string
	href string
}

type Product struct {
	name         string
	price        float64
	currency     string
	href         string
	categoryHref string
}
