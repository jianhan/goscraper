package scraper

import (
	"errors"
	"strings"
)

type Scraper interface {
	Name() string
	Scrape() error
	Categories() []Category
	Products() []Product
	FetchCategories() error
	FetchProducts() error
}

type Category struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Product struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Image       string  `json:"image"`
	CategoryURL string  `json:"category_url"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	Brand       string  `json:"brand"`
	URL         string  `json:"url"`
}

type ValidationTemplate struct {
	scraper Scraper
}

func (v *ValidationTemplate) Scrape() error {
	if err := v.validate(); err != nil {
		return err
	}

	if err := v.scraper.FetchCategories(); err != nil {
		return err
	}

	if err := v.scraper.FetchProducts(); err != nil {
		return err
	}

	return nil
}

func (v *ValidationTemplate) validate() error {
	if strings.Trim(v.scraper.Name(), " ") == "" {
		return errors.New("empty name")
	}

	if len(v.scraper.Categories()) == 0 {
		return errors.New("empty categories")
	}

	if len(v.scraper.Products()) == 0 {
		return errors.New("empty products")
	}

	return nil
}
