package output

import (
	"errors"
	"fmt"
	"log"
	"os"

	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gosimple/slug"
	"github.com/jianhan/goscraper/scraper"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type s3Writer struct {
	scrapers []scraper.Scraper
}

func NewS3Writer(scrapers []scraper.Scraper) OutputWriter {
	return &s3Writer{
		scrapers: scrapers,
	}
}
func (s *s3Writer) Write() error {
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

	for _, s := range s.scrapers {
		folderName := slug.Make(s.Name())

		productsFile, productsFilePath, err := getFileAndPath(folderName, "products")
		if err != nil {
			return err
		}
		r, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(productsFilePath),
			Body:   productsFile,
		})
		if err != nil {
			return fmt.Errorf("failed to upload file, %v", err)
		}
		logrus.WithField("result", r).Info("finished uploaded products")

		categoriesFile, categoriesFilePath, err := getFileAndPath(folderName, "categories")
		if err != nil {
			return err
		}
		r, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(categoriesFilePath),
			Body:   categoriesFile,
		})
		if err != nil {
			return fmt.Errorf("failed to upload file, %v", err)
		}
		logrus.WithField("result", r).Info("finished uploaded categories")

		supplierFile, supplierFilePath, err := getFileAndPath(folderName, "supplier")
		if err != nil {
			return err
		}
		r, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(supplierFilePath),
			Body:   supplierFile,
		})
		if err != nil {
			return fmt.Errorf("failed to upload file, %v", err)
		}
		logrus.WithField("result", r).Info("finished uploaded supplier")
	}

	return nil
}

func getFileAndPath(folderName, jsonFileName string) (*os.File, string, error) {
	path := path.Join(folderName, jsonFileName+".json")
	f, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open file %q, %v", path, err)
	}

	return f, path, nil
}
