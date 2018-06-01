package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type base struct {
	name        string
	categoryURL string
	categories  []Category
	products    []Product
	currency    string
	homepageURL string
}

func (b *base) Name() string {
	return b.name
}

func (b *base) Categories() []Category {
	return b.categories
}

func (b *base) Products() []Product {
	return b.products
}

func (b *base) addCategory(c Category) {
	b.categories = append(b.categories, c)
}

func (b *base) addProduct(p Product) {
	b.products = append(b.products, p)
}

func (b *base) htmlDoc(url string) (*goquery.Document, func() error, error) {
	// get html page
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		return nil, nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, nil, err
	}

	return doc, res.Body.Close, nil
}

func (b *base) priceStrToFloat(priceStr string) (price float64, err error) {
	priceRaw := strings.Replace(priceStr, " ", "", -1)
	priceRaw = strings.Replace(priceRaw, ",", "", -1)
	priceRaw = strings.Replace(priceRaw, "$", "", -1)
	price, err = strconv.ParseFloat(priceRaw, 64)
	if err != nil {
		return
	}

	return
}

func (b *base) getLinkFullURL(url string) string {
	if strings.HasPrefix(url, b.homepageURL) {
		return url
	}
	url = strings.Replace(url, " ", "", -1)
	url = strings.Replace(url, "/", "", -1)

	return b.homepageURL + "/" + url
}

//func (b *base) OutputJSON() error {
//	if len(b.products) == 0 {
//		return errors.New("empty products")
//	}
//
//	if len(b.categories) == 0 {
//		return errors.New("empty categories")
//	}
//
//	productsJSON, err := json.Marshal(b.products)
//	if err != nil {
//		return err
//	}
//
//	categoriesJSON, err := json.Marshal(b.categories)
//	if err != nil {
//		return err
//	}
//
//	if err = ioutil.WriteFile(path.Join(b.Name(), "products.json"), productsJSON, 0644); err != nil {
//		return err
//	}
//
//	if err = ioutil.WriteFile(path.Join(b.Name(), "categories.json"), categoriesJSON, 0644); err != nil {
//		return err
//	}
//
//	return nil
//}
