package scraper

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/davecgh/go-spew/spew"
	"github.com/gosimple/slug"
	"github.com/joho/godotenv"
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

func OutputJSONData(scrapers ...Scraper) error {
	for _, scraper := range scrapers {
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

func UploadS3(scrapers ...Scraper) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s3Bucket := os.Getenv("S3_BUCKET")
	s3AccessKey := os.Getenv("S3_ACCESS_KEY")
	s3SecretKey := os.Getenv("S3_SECRET_KEY")
	s3Region := os.Getenv("S3_REGION")

	if s3Bucket == "" || s3AccessKey == "" || s3SecretKey == "" || s3Region == "" {
		return errors.New("invalid configs for s3")
	}

	creds := credentials.NewStaticCredentials(s3AccessKey, s3SecretKey, "")
	cfg := aws.NewConfig().WithRegion(s3Region).WithCredentials(creds)
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(cfg))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	filename := "umart/supplier.json"
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", filename, err)
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(filename),
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	spew.Dump(result)
	return nil
}
