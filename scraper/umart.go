package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

type umart struct {
	base
}

func NewUmart() Scraper {
	b := base{
		name:        "Umart",
		categoryURL: "https://www.umart.com.au/all-categories.html",
		homepageURL: "https://www.umart.com.au",
		currency:    "AUD",
	}

	return &umart{b}
}

func (u *umart) Scrape() error {
	if err := u.FetchCategories(); err != nil {
		return err
	}
	if err := u.FetchProducts(); err != nil {
		return err
	}

	return nil
}

func (u *umart) FetchCategories() error {
	doc, fn, err := u.htmlDoc(u.categoryURL)
	if err != nil {
		return err
	}
	defer fn()

	// get all links with class categoryLink
	doc.Find("div.ovhide.productsIn.productText > a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok && href != "" {
			u.categories = append(u.categories, Category{Name: s.Text(), URL: u.getLinkFullURL(href)})
		}
	})

	return nil
}

func (u *umart) FetchProducts() error {
	for _, c := range u.categories {
		if err := u.fetchProductsByURL(c.URL, u.categoryURL); err != nil {
			return err
		}
		break
	}

	return nil
}

func (u *umart) fetchProductsByURL(url, categoryURL string) error {
	doc, fn, err := u.htmlDoc(url)
	if err != nil {
		return err
	}
	defer fn()

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
			if p.Price, err = u.priceStrToFloat(priceS.Text()); err != nil {
				logrus.Warn(err)
			}
		})
		p.Currency = u.currency
		u.products = append(u.products, p)
	})

	// find next page url
	var nextPageURL string
	doc.Find("ul.page li a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == ">" {
			href, ok := s.Attr("href")
			if ok {
				nextPageURL = u.getLinkFullURL(href)
			}
		}
	})
	if nextPageURL != "" {
		u.fetchProductsByURL(nextPageURL, categoryURL)
	}

	return nil
}
