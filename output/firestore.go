package output

import (
	"log"

	"github.com/jianhan/goscraper/scraper"

	"firebase.google.com/go"
	"github.com/davecgh/go-spew/spew"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type firestoreWriter struct {
	scrapers []scraper.Scraper
}

func NewFirestore(scrapers []scraper.Scraper) OutputWriter {
	return &firestoreWriter{
		scrapers: scrapers,
	}
}

func (f *firestoreWriter) Write() error {
	ctx := context.Background()
	app, err := firebase.NewApp(
		ctx,
		&firebase.Config{ProjectID: "reactfire-198405"},
		option.WithCredentialsFile("./output/serviceAccountKey.json"),
	)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// load all products
	productsMap := map[string]string{}
	products := client.Collection("products").Documents(ctx)
	for {
		doc, err := products.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		productsMap[doc.Data()["URL"].(string)] = doc.Ref.ID
	}

	// Get a new write batch.
	batch := client.Batch()

	for _, s := range f.scrapers {
		if err = s.Validate(); err != nil {
			return err
		}
		for _, p := range s.Products() {
			var id string
			if productsMap[p.URL] != "" {
				// exists
				id = productsMap[p.URL]
			} else {
				// not exists
				id = uuid.Must(uuid.NewV4()).String()

			}
			sfRef := client.Collection("products").Doc(id)
			batch.Set(sfRef, p)
		}
	}

	// Commit the batch.
	r, err := batch.Commit(ctx)
	if err != nil {
		return err
	}
	spew.Dump(r)
	return nil
}
