package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

// struct to store scraped data
type Product struct {
	Title  string
	Price  string
	Review string
}

func main() {
	// a new Colly collector
	c := colly.NewCollector(
		// Allow redirected URLs
		colly.AllowURLRevisit(),
	)

	// a CSV file to store the scraped data
	file, err := os.Create("scraped_data.csv")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// header to CSV file
	writer.Write([]string{"Title", "Price", "Review"})

	// the URLs to scrape
	urls := []string{
		"https://www.fnp.com/gift/black-forest-cake-half-kg?widgettype=TrendingNow&widgetdisplayname=Gifts%20In%20Trend&pos=2&page=HP",
	}

	// the scraping logic
	c.OnHTML("div.product", func(e *colly.HTMLElement) {
		// Extract product details from HTML structure
		title := strings.TrimSpace(e.ChildText("h2"))
		price := strings.TrimSpace(e.ChildText("span.price"))
		review := strings.TrimSpace(e.ChildText("div.review"))

		// Create a new product object
		product := Product{
			Title:  title,
			Price:  price,
			Review: review,
		}

		// Write product details to CSV file
		err := writer.Write([]string{product.Title, product.Price, product.Review})
		if err != nil {
			log.Printf("Failed to write record to CSV: %v", err)
		}
	})

	// Error handling
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s\nFailed with response: %v\nError: %v", r.Request.URL, r, err)
	})

	// Start scraping each URL
	for _, url := range urls {
		err := c.Visit(url)
		if err != nil {
			log.Printf("Failed to scrape URL %s: %v", url, err)
		}
	}

	// Notify when scraping is done
	log.Println("Scraping finished. Check scraped_data.csv for results.")
}
