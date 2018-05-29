package scraper

type Scraper interface {
	Name() string
	Scrape() error
	Categories() []Category
	Products() []Product
}

type Category struct {
	Name string `json:"name"`
	Href string `json:"href"`
}

type Product struct {
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	Image        string  `json:"image"`
	CategoryHref string  `json:"category_href"`
	Description  string  `json:"description"`
}
