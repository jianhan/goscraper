package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

type ncix struct {
	base
}

func NewNCIX(testMode bool) Scraper {
	b := base{
		homepageURL: "https://www.ncix.com",
		name:        "NCIX",
		categoryURL: "https://www.ncix.com/categories/",
		currency:    "CAD",
		testMode:    testMode,
	}

	return &ncix{b}
}

func (n *ncix) Scrape() error {
	// start scraping
	if err := n.fetchCategories(); err != nil {
		return err
	}
	if err := n.fetchProducts(); err != nil {
		return err
	}

	return nil
}

func (n *ncix) fetchCategories() error {
	doc, fn, err := n.htmlDoc(n.categoryURL)
	if err != nil {
		return err
	}
	defer fn()

	// find categories
	doc.Find("div#sublinks a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		href, ok := s.Attr("href")
		if ok && href != "" {
			n.addCategory(Category{Name: s.Text(), URL: href})
		}
	})

	return nil
}

func (n *ncix) fetchProducts() error {
	for _, c := range n.categories {
		doc, fn, err := n.htmlDoc(c.URL)
		if err != nil {
			return err
		}
		defer fn()

		// find products
		doc.Find("span.listing a").Each(func(i int, s *goquery.Selection) {
			p := Product{CategoryURL: c.URL, Currency: n.currency}
			// find Href and Name
			href, ok := s.Attr("href")
			if ok {
				p.Currency, p.Name, p.URL = n.currency, s.Text(), href
			}

			// find image
			s.Parent().Parent().Prev().Find("img").First().Each(func(j int, js *goquery.Selection) {
				imageSrc, ok := js.Attr("src")
				if ok {
					p.Image = imageSrc
				}
			})

			// find Price
			s.Parent().Parent().Next().Next().Find("strong").First().Each(func(j int, js *goquery.Selection) {
				if p.Price, err = n.priceStrToFloat(js.Text()); err != nil {
					logrus.Warn(err)
				}
			})
			if p.Price > 0 {
				// append product into products
				n.addProduct(p)
			}
		})
		// test mode checking
		if n.testMode && len(n.products) > 0 {
			break
		}
	}

	return nil
}
