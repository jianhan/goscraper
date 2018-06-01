package scraper

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gosimple/slug"
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
	if err := scraper.Validate(); err != nil {
		return err
	}

	// create folder if not exists & clean folder
	folderName := slug.Make(scraper.Name())
	CreateDirIfNotExist(folderName)
	RemoveContents(folderName)

	products, categories := scraper.Products(), scraper.Categories()
	productsJSON, err := json.Marshal(products)
	if err != nil {
		return err
	}

	categoriesJSON, err := json.Marshal(categories)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(path.Join(folderName, "products.json"), productsJSON, 0644); err != nil {
		return err
	}

	if err = ioutil.WriteFile(path.Join(folderName, "categories.json"), categoriesJSON, 0644); err != nil {
		return err
	}

	return nil
}
