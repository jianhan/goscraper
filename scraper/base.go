package scraper

type Fetcher interface {
	fetchCategories(url string) ([]Category, error)
	fetchProducts(categories []Category) ([]Product, error)
}

type Base struct {
	name        string
	categoryURL string
	categories  []Category
	products    []Product
	currency    string
	homepageURL string
	fetcher     Fetcher
}

func (b *Base) Name() string {
	return b.name
}

func (b *Base) Categories() []Category {
	return b.categories
}

func (b *Base) Products() []Product {
	return b.products
}

func (b *Base) Scrape() (categories []Category, products []Product, err error) {
	// fetch categories
	if categories, err = b.fetcher.fetchCategories(b.categoryURL); err != nil {
		return
	}

	// fetch products
	if products, err = b.fetcher.fetchProducts(categories); err != nil {
		return
	}

	return
}
