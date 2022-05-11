package main

import (
	"fmt"
	"github.com/twilio/twilio-go"
	"log"
	"os"

	"github.com/playwright-community/playwright-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var authToken string

func sendSMS(msg string) error {
	accountSid := "ACe59edba87a888fbfbf2ce38ba33d03eb"

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo("+17324067063")
	params.SetFrom("+19803755389")
	params.SetBody(msg)

	resp, err := client.ApiV2010.CreateMessage(params)
	if err != nil {
		return err
	} else {
		fmt.Println("Message Sid: " + *resp.Sid)
	}
	return nil
}
func init() {
	err := playwright.Install()
	if err != nil {
		log.Fatal(err)
	}
	authToken = os.Getenv("TWILIO_API_TOKEN")
	fmt.Println(authToken)
}

type SearchResult struct {
	message string
	results []*Result
}

func (sr *SearchResult) String() string {
	msg := ""
	for _, result := range sr.results {
		msg += fmt.Sprintf("Product: %surl:%s\n",result.name, result.url)
	}
	return msg
}

func searchCostco(page playwright.Page) (*SearchResult, error) {
	_, err := page.Goto("https://www.costco.com")
	if err != nil {
		return nil, err
	}
	search, err := page.WaitForSelector("[placeholder=\"Search\"]")
	if err != nil {
		return nil, fmt.Errorf("could not find search bar: %v", err.Error())
	}
	if err := search.Fill("enfamil"); err != nil {
		return nil, fmt.Errorf("Could not fill placeholder for search: %v", err.Error())
	}

	if err := page.Press("[placeholder=\"Search\"]", "Enter"); err != nil {
		return nil, fmt.Errorf("could not click search button")
	}

	_, err = page.WaitForNavigation()
	if err != nil {
		return nil, fmt.Errorf("could not load page")
	}

	searchResults, err := page.WaitForSelector("div#search-results")
	if err != nil {
		return nil, fmt.Errorf("Could not load searchResults: %s", err.Error())
	}

	entries, err := searchResults.QuerySelectorAll("div.product")
	if err != nil {
		return nil, fmt.Errorf("could not get entries: %v", err)
	}
	costcoResults := NewSearchResults()
	for _, entry := range entries {
		productTile, err := entry.QuerySelector("span.description > a")
		if err != nil {
			log.Printf("Could not load product name: %v", err.Error())
			continue
		}

		productName, err := productTile.TextContent()
		if err != nil {
			log.Printf("Could not get product text: %v", err.Error())
			continue
		}
		link, err := productTile.GetAttribute("href")
		if err != nil {
			log.Printf("Could not get href: %v", err.Error())
			continue
		}
		costcoResults.AddResult(NewResult(productName, link))
	}
	return costcoResults, nil
}

func (sr *SearchResult) AddResult(result *Result) {
	sr.results = append(sr.results, result)
}

func NewResult(name string, link string) *Result {
	return &Result{
		name,
		link,
	}
}

type Result struct {
	name string
	url string
}

func NewSearchResults() *SearchResult {
	return &SearchResult{}
}

func main() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: playwright.String("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36"),
	})
	if err != nil {
		log.Fatalf("Could not setup new context: %v", err)
	}

	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	res, err := searchCostco(page)
	err = sendSMS(res.String())
	if err != nil {
		log.Println(err)
	}


	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}