package scraper

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveContents(dir string) error {
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

func OutputJSONData(scraper Scraper) error {
	products, categories := scraper.Products(), scraper.Categories()

	if len(products) == 0 {
		return errors.New("empty products")
	}

	if len(categories) == 0 {
		return errors.New("empty categories")
	}

	productsJSON, err := json.Marshal(products)
	if err != nil {
		return err
	}

	categoriesJSON, err := json.Marshal(categories)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(path.Join(scraper.Name(), "products.json"), productsJSON, 0644); err != nil {
		return err
	}

	if err = ioutil.WriteFile(path.Join(scraper.Name(), "categories.json"), categoriesJSON, 0644); err != nil {
		return err
	}

	return nil
}
