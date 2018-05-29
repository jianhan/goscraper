package scraper

type megabuyau struct {
	name       string
	url        string
	categories []Category
	products   []Product
	currency   string
}

func NewMegabuyau() Scraper {
	return &megabuyau{
		url:      "https://www.megabuy.com.au/computer-components-c1160.html",
		currency: "AUD",
		name:     "Mega Buy Australia",
	}
}

func (m *megabuyau) Name() string {
	return m.name
}

func (m *megabuyau) Categories() []Category {
	return m.categories
}

func (m *megabuyau) Products() []Product {
	return m.products
}

func (m *megabuyau) Scrape() error {
	// clear data first
	RemoveContents(m.name)
	// create dir for downloaded data
	if err := CreateDirIfNotExist(m.name); err != nil {
		return err
	}

	// start scraping
	m.fetchCategories(m.url)
	m.fetchProducts()
	return nil
}

func (m *megabuyau) fetchCategories(url string) error {
	return nil
}

func (m *megabuyau) fetchProducts() error {
	return nil
}
