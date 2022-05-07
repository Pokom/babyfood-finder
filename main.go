package main

import (
	"github.com/gocolly/colly"
	"log"
)

func main() {
	url := "https://www.costco.com/CatalogSearch?dept=All&keyword=enfamil"
	c := colly.NewCollector(
		colly.AllowedDomains("www.costco.com"),
	)

	c.OnHTML("div#tiles-body-attribute", func(e *colly.HTMLElement){
		log.Println(e)
	})
	c.OnHTML("div.product-tile-set", func(e *colly.HTMLElement) {
		log.Println(e)
	})
	c.OnHTML("body", func(e *colly.HTMLElement) {
		log.Println(e)
	})
	c.OnRequest(func(r *colly.Request){
		log.Printf("Visiting %s", r.URL)
	})

	c.OnResponse(func(r *colly.Response){
		log.Println(r)
	})

	c.Visit(url)
}