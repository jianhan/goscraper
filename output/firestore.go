package output

import (
	"log"

	"github.com/jianhan/goscraper/scraper"

	"firebase.google.com/go"
	"github.com/satori/go.uuid"
	"golang.org/x/net/context"
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

	// Get a new write batch.
	batch := client.Batch()

	for _, s := range f.scrapers {
		if err = s.Validate(); err != nil {
			return err
		}
		for _, p := range s.Products() {
			id := uuid.Must(uuid.NewV4())
			sfRef := client.Collection("products").Doc(id.String())
			batch.Set(sfRef, p)
		}
	}

	// Commit the batch.
	_, err = batch.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
