package output

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/gosimple/slug"
	"github.com/jianhan/goscraper/scraper"
)

type jsonWriter struct {
	scrapers []scraper.Scraper
}

func NewJSONWriter(scrapers []scraper.Scraper) OutputWriter {
	return &jsonWriter{
		scrapers: scrapers,
	}
}

func (j *jsonWriter) Write() error {
	// simple validation
	if len(j.scrapers) == 0 {
		return errors.New("empty scrapers, can not write to json file")
	}
	for _, scraper := range j.scrapers {
		if err := scraper.Validate(); err != nil {
			return err
		}

		// create folder if not exists & clean folder
		folderName := slug.Make(scraper.Name())
		createDirIfNotExist(folderName)
		removeContents(folderName)

		products, categories := scraper.Products(), scraper.Categories()
		productsJSON, err := json.Marshal(products)
		if err != nil {
			return err
		}

		categoriesJSON, err := json.Marshal(categories)
		if err != nil {
			return err
		}

		// write products
		if err = ioutil.WriteFile(path.Join(folderName, "products.json"), productsJSON, 0644); err != nil {
			return err
		}

		// write categories
		if err = ioutil.WriteFile(path.Join(folderName, "categories.json"), categoriesJSON, 0644); err != nil {
			return err
		}

		// write supplier
		supplierJS, err := json.Marshal(struct {
			Name        string `json:"name"`
			HomepageURL string `json:"homepage_url"`
			Currency    string `json:"currency"`
		}{
			Name:        scraper.Name(),
			HomepageURL: scraper.HomepageURL(),
			Currency:    scraper.Currency(),
		})
		if err = ioutil.WriteFile(path.Join(folderName, "supplier.json"), supplierJS, 0644); err != nil {
			return err
		}
	}

	return nil
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	return nil
}
