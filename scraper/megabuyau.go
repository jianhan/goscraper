package scraper

import (
	"fmt"
	"log"
	"net/http"

	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
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
			m.categories = append(m.categories, Category{Name: s.Text(), URL: href})
		}
	})

	return nil
}

func (m *megabuyau) fetchProducts() error {
	for _, c := range m.categories {
		if err := m.fetchProductsByURL(c.URL, c.URL); err != nil {
			return err
		}
	}
	return nil
}

func (m *megabuyau) fetchProductsByURL(url, categoryURL string) error {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		return err
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

	// find products
	doc.Find("div.productListing div.productListingRow, div.productListing div.productListingRowAlt").Each(func(i int, s *goquery.Selection) {
		p := Product{CategoryURL: categoryURL}

		// find image
		s.Find("div.image > a > img").First().Each(func(ii int, is *goquery.Selection) {
			src, ok := is.Attr("src")
			if ok {
				p.Image = src
			}
		})

		// find name
		s.Find("div.nameDescription > a").First().Each(func(ni int, ns *goquery.Selection) {
			p.Name = ns.Text()
		})

		// find price
		s.Find("div.price > span").First().Each(func(ni int, ns *goquery.Selection) {
			priceRaw := strings.Replace(strings.TrimLeft(ns.Text(), "$"), ",", "", -1)
			priceFloat, err := strconv.ParseFloat(priceRaw, 64)
			if err != nil {
				logrus.Warn(err)
			} else {
				p.Price = priceFloat
			}
		})
		if p.Price > 0 {
			// append product into products
			m.products = append(m.products, p)
		}
	})

	// find next page url
	doc.Find("div.pagination").First().Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(ai int, as *goquery.Selection) {
			title, ok := as.Attr("title")
			if ok {
				if strings.ToLower(strings.Trim(title, " ")) == "next page" {
					nextPageHref, ok := as.Attr("href")
					if ok {
						m.fetchProductsByURL(nextPageHref, categoryURL)
					}
				}
			}
		})
	})

	return nil
}
