package scraper

type Scraper interface {
	Name() string
	Scrape() error
	Categories() []Category
	Products() []Product
}

type Category struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Product struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Image       string  `json:"image"`
	CategoryURL string  `json:"category_url"`
	Description string  `json:"description"`
}
