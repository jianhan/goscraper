package scraper

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"
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

func (b *base) OutputJSON() error {
	if len(b.products) == 0 {
		return errors.New("empty products")
	}

	if len(b.categories) == 0 {
		return errors.New("empty categories")
	}

	productsJSON, err := json.Marshal(b.products)
	if err != nil {
		return err
	}

	categoriesJSON, err := json.Marshal(b.categories)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(path.Join(b.Name(), "products.json"), productsJSON, 0644); err != nil {
		return err
	}

	if err = ioutil.WriteFile(path.Join(b.Name(), "categories.json"), categoriesJSON, 0644); err != nil {
		return err
	}

	return nil
}
