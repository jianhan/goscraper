package scraper

import (
	"fmt"
	"log"
	"net/http"

	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
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
			m.categories = append(m.categories, Category{Name: s.Text(), Href: href})
		}
	})

	return nil
}

func (m *megabuyau) fetchProducts() error {
	// TODO: uncomment this code when deploy to production
	m.categories = append(m.categories[:0], m.categories[0+1:]...)
	for _, c := range m.categories {
		// Request the HTML page.
		res, err := http.Get(c.Href)
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
		doc.Find("div.productListing div.productListingRow").Each(func(i int, s *goquery.Selection) {
			p := Product{CategoryHref: c.Href}

			// find image
			s.Find("div.image > a > img").Each(func(ii int, is *goquery.Selection) {
				src, ok := is.Attr("src")
				if ok {
					p.Image = src
				}
			})

			// find name
			s.Find("div.nameDescription > a").Each(func(ni int, ns *goquery.Selection) {
				p.Name = ns.Text()
			})

			// find price
			s.Find("div.price > span").Each(func(ni int, ns *goquery.Selection) {
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
		spew.Dump(m.products)
		break
	}
	return nil
}
