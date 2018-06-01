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

type umart struct {
	name        string
	url         string
	categories  []Category
	products    []Product
	currency    string
	homepageURL string
}

func NewUmart() Scraper {
	return &scraper{
		name:              "Umart",
		categoryURL:       "https://www.umart.com.au/all-categories.html",
		homepageURL:       "https://www.umart.com.au",
		currency:          "AUD",
		categoriesFetcher: uMartCategoriesFetcher,
	}
}

func uMartCategoriesFetcher(url, homepageURL string) ([]Category, error) {
	categories := []Category{}

	// fetch category html page
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	// get all links with class categoryLink
	doc.Find("div.ovhide.productsIn.productText > a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok && href != "" {
			categories = append(categories, Category{Name: s.Text(), URL: homepageURL + "/" + href})
		}
	})

	return categories, nil
}

func uMartProductsFetcher(categories []Category, currency string) (products []Product, err error) {
	for _, c := range categories {
		if err := fetchUMartProductsByURL(c.URL, c.URL, currency); err != nil {
			return nil, err
		}
	}
	return nil
}

func fetchUMartProductsByURL(url, categoryURL, currency, homepageURL string) ([]Product, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	products := []Product{}

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
		p.Currency = currency
		products = append(products, p)
	})

	var nextPageURL string
	// find next page url
	doc.Find("ul.page li a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == ">" {
			href, ok := s.Attr("href")
			if ok {
				nextPageURL = href
			}
		}
	})

	if nextPageURL != "" {
		fetchUMartProductsByURL(u.homepageURL+"/"+nextPageURL, categoryURL)
	}

	return nil
}
