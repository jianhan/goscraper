package scraper

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
)

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

	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	// get all links with class categoryLink
	doc.Find("a.categoryLink").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		href, ok := s.Attr("href")
		if ok && href != "" {
			m.categories = append(m.categories, Category{Name: s.Text(), Href: href})
		}
	})
	spew.Dump(m.categories)
	return nil
}

func (m *megabuyau) fetchProducts() error {
	return nil
}
