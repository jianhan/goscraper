package scraper

type Scraper interface {
	Name() string
	Scrape() (categories []Category, products []Product, err error)
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
	Rating      float64 `json:"rating"`
	Brand       string  `json:"brand"`
	URL         string  `json:"url"`
}

type Fetcher interface {
	fetchCategories(url string) ([]Category, error)
	fetchProducts(categories []Category) ([]Product, error)
}

type scraper struct {
	name              string
	categoryURL       string
	categories        []Category
	products          []Product
	currency          string
	homepageURL       string
	categoriesFetcher func(url, homepageURL string) ([]Category, error)
	productsFetcher   func(categories []Category) ([]Product, error)
}

func (b *scraper) Name() string {
	return b.name
}

func (b *scraper) Categories() []Category {
	return b.categories
}

func (b *scraper) Products() []Product {
	return b.products
}

func (b *scraper) Scrape() (categories []Category, products []Product, err error) {
	// fetch categories
	if categories, err = b.categoriesFetcher(b.categoryURL, b.homepageURL); err != nil {
		return
	}

	// fetch products
	if products, err = b.productsFetcher(categories); err != nil {
		return
	}

	return
}
