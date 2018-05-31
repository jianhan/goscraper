package scraper

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type umart struct {
	name        string
	url         string
	categories  []Category
	products    []Product
	currency    string
	homepageURL string
}

func NewUmart() Scraper {
	return &umart{
		url:         "https://www.umart.com.au/all-categories.html",
		currency:    "AUD",
		name:        "umart",
		homepageURL: "https://www.umart.com.au",
	}
}

func (u *umart) Name() string {
	return u.name
}

func (u *umart) Categories() []Category {
	return u.categories
}

func (u *umart) Products() []Product {
	return u.products
}

func (u *umart) Scrape() error {
	// clear data first
	RemoveContents(u.name)
	// create dir for downloaded data
	if err := CreateDirIfNotExist(u.name); err != nil {
		return err
	}

	// start scraping
	u.fetchCategories(u.url)
	u.fetchProducts()
	return nil
}

func (u *umart) fetchCategories(url string) error {
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
	doc.Find("div.ovhide.productsIn.productText > a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok && href != "" {
			u.categories = append(u.categories, Category{Name: s.Text(), URL: u.homepageURL + "/" + href})
		}
	})

	return nil
}

func (u *umart) fetchProducts() {

}
