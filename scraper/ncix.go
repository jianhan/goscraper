package scraper

type ncix struct {
	url string
}

func NewNCIXScrapper() Scraper{
	return &ncix{
		// every scrapper follow different algorithm, therefore do not needed to pass as param
		url: "https://www.ncix.com/categories/",
	}
}

func (n *ncix) Scrape() error {

	return nil
}
