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

func (u *umart) fetchProducts() error {
	i := 0
	for _, c := range u.categories {
		logrus.Infof("Start fetching %s : %s", c.Name, c.URL)
		if err := u.fetchProductsByURL(c.URL, c.URL); err != nil {
			return err
		}
		i++
		if i == 10 {
			break
		}
	}
	return nil
}

func (u *umart) fetchProductsByURL(url, categoryURL string) error {
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
	doc.Find("li.goods_info").Each(func(i int, s *goquery.Selection) {
		p := Product{CategoryURL: categoryURL}

		// find image
		s.First().Find("div.goods_img > a > img").Each(func(imgI int, imgS *goquery.Selection) {
			src, ok := imgS.Attr("src")
			if ok {
				p.Image = src
			}
		})

		// find product name
		s.First().Find("div.content_holder1 > div.goods_name > a").Each(func(nameI int, nameS *goquery.Selection) {

			// product url
			href, ok := nameS.Attr("href")
			if ok {
				p.URL = href
			}

			// product name
			p.Name = nameS.Text()

		})

		// find product price
		s.First().Find("span.goods_price").Each(func(priceI int, priceS *goquery.Selection) {
			priceRaw := strings.Replace(priceS.Text(), " ", "", -1)
			priceRaw = strings.Replace(priceRaw, ",", "", -1)
			priceRaw = strings.Replace(priceRaw, "$", "", -1)
			priceFloat, err := strconv.ParseFloat(priceRaw, 64)
			if err != nil {
				logrus.Warn(err)
			} else {
				p.Price = priceFloat
			}
		})
		p.Currency = u.currency
		u.products = append(u.products, p)
	})

	var nextPageURL string
	// find next page url
	doc.Find("ul.page li a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == ">" {
			href, ok := s.Attr("href")
			if ok {
				spew.Dump(href)
				nextPageURL = href
			}
		}
	})

	if nextPageURL != "" {
		u.fetchProductsByURL(u.homepageURL+"/"+nextPageURL, categoryURL)
	}

	return nil
}
