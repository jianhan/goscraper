package scraper

import (
	"fmt"
	"log"
	"net/http"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
)

type amazon struct {
	homepage   string
	name       string
	url        string
	categories []Category
	products   []Product
	currency   string
}

func NewAmazon() Scraper {
	return &amazon{
		homepage: "https://www.amazon.com",
		url:      "https://www.amazon.com/PC-Parts-Components/b/ref=nav_shopall_components?ie=UTF8&node=193870011",
		currency: "USD",
		name:     "Amazon",
	}
}

func (a *amazon) Name() string {
	return a.name
}

func (a *amazon) Categories() []Category {
	return a.categories
}

func (a *amazon) Products() []Product {
	return a.products
}

func (a *amazon) Scrape() error {
	// clear data first
	RemoveContents(a.name)
	// create dir for downloaded data
	if err := CreateDirIfNotExist(a.name); err != nil {
		return err
	}

	// start scraping
	a.fetchCategories(a.url)
	a.fetchProducts()

	return nil
}

func (a *amazon) fetchCategories(url string) error {
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
	doc.Find("h4").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "Computer Components" {
			//spew.Dump(s.Text())
			s.Parent().Parent().Next().Find("li").Each(func(ci int, cs *goquery.Selection) {
				cs.First().Find("a").Each(func(li int, ls *goquery.Selection) {
					href, ok := ls.Attr("href")
					if ok {
						a.categories = append(a.categories, Category{Name: ls.Text(), URL: a.homepage + href})
					}
				})
			})
		}
	})

	return nil
}

func (a *amazon) fetchProducts() error {
	for _, c := range a.categories {
		// Request the HTML page.
		spew.Dump(c.URL)
		res, err := http.Get(c.URL)
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
		doc.Find("div#mainResults").First().Each(func(i int, s *goquery.Selection) {
			//p := Product{CategoryURL: c.URL}
			s.Find("div.s-item-container").Each(func(divI int, divS *goquery.Selection) {

				// find text
				divS.First().Find("a.s-access-detail-page").Each(func(linkI int, linkS *goquery.Selection) {
					spew.Dump(linkS.Text())
				})

				// find image
				divS.First().Find("img.s-access-image").Each(func(imgI int, imgS *goquery.Selection) {
					imgSrc, ok := imgS.Attr("src")
					if ok {
						spew.Dump(imgSrc)
					}
				})

				// find brand
				divS.Find(".a-row .a-spacing-none").Each(func(brandI int, brandS *goquery.Selection) {
					// has two spans
					brandS.Find("span, sup").Each(func(byI int, byS *goquery.Selection) {
						if strings.ToLower(strings.Trim(byS.Text(), " ")) == "by" {
							spew.Dump(byS.Next().Text())
						}

						// find whole price at the same time
						if byS.HasClass("sx-price-whole") {
							spew.Dump(byS.Text())
						}

						// find fractional price
						if byS.HasClass("sx-price-fractional") {
							spew.Dump(byS.Text())
						}

						// find rating
						byS.Find("span.a-icon-alt").First().Each(func(ratingI int, ratingS *goquery.Selection) {
							ratingSlice := strings.Split(ratingS.Text(), " ")
							if len(ratingSlice) > 0 {
								spew.Dump(ratingSlice[0])
							}
						})
					})
				})
			})
		})
		break
	}

	return nil
}
