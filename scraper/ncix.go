package scraper

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type ncix struct {
	url        string
	categories []Category
}

func NewNCIXScrapper() Scraper {
	return &ncix{
		// every scrapper follow different algorithm, therefore do not needed to pass as param
		url: "https://www.ncix.com/categories/",
	}
}

func (n *ncix) Categories() []Category {
	return n.categories
}

func (n *ncix) Scrape() error {
	// Request the HTML page.
	res, err := http.Get(n.url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	n.fetchCategories(doc)
	return nil
}

func (n *ncix) fetchCategories(doc *goquery.Document) {
	// Find the review items
	doc.Find("div#sublinks a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		href, ok := s.Attr("href")
		if ok {
			n.categories = append(n.categories, Category{text: s.Text(), href: href})
		}
	})
}
