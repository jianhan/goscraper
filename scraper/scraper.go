package scraper

type Scraper interface {
	Scrape() error
	Categories() []*Category
	Products() []*Product
}

type Category struct {
	name string `json:"name"`
	href string `json:"href"`
}

type Product struct {
	name         string  `json:"name"`
	price        float64 `json:"price"`
	currency     string  `json:"currency"`
	href         string  `json:"href"`
	categoryHref string  `json:"category_href"`
}
