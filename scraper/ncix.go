package scraper

import (
	"log"
	"net/http"

	"fmt"

	"strings"

	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

type ncix struct {
	url        string
	categories []*Category
	products   []*Product
	currency   string
}

func NewNCIXScrapper() Scraper {
	return &ncix{
		// every scrapper follow different algorithm, therefore do not needed to pass as param
		url:      "https://www.ncix.com/categories/",
		currency: "CAD",
	}
}

func (n *ncix) Categories() []*Category {
	return n.categories
}

func (n *ncix) Products() []*Product {
	return n.products
}

func (n *ncix) Scrape() error {
	n.fetchCategories(n.url)
	n.fetchProducts()
	return nil
}

func (n *ncix) fetchCategories(url string) error {
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

	// find categories
	doc.Find("div#sublinks a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		href, ok := s.Attr("href")
		if ok {
			n.categories = append(n.categories, &Category{name: s.Text(), href: href})
		}
	})

	return nil
}

func (n *ncix) fetchProducts() error {
	for _, c := range n.categories {
		// Request the HTML page.
		res, err := http.Get(c.href)
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

		// find products
		doc.Find("span.listing a").Each(func(i int, s *goquery.Selection) {
			p := &Product{}
			// For each item found, get the band and title
			href, ok := s.Attr("href")

			// find image
			s.Parent().Parent().Prev().Find("img").Each(func(j int, js *goquery.Selection) {
				imageSrc, ok := js.Attr("src")
				if ok {
					p.href = imageSrc
				}
			})

			// find price
			s.Parent().Parent().Next().Next().Find("strong").Each(func(j int, js *goquery.Selection) {
				priceRaw := strings.Replace(strings.TrimLeft(js.Text(), "$"), ",", "", -1)
				priceFloat, err := strconv.ParseFloat(priceRaw, 64)
				if err != nil {
					logrus.Warn(err)
				} else {
					p.price = priceFloat
				}
			})
			if ok {
				p.currency, p.name, p.href = s.Text(), n.currency, href
				n.products = append(n.products, p)
			}
		})
		break
	}

	return nil
}
