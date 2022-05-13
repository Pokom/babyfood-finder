package main

import (
	"flag"
	"fmt"
	"github.com/twilio/twilio-go"
	"log"
	"os"
	"strings"

	"github.com/playwright-community/playwright-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type SMS struct {
	TwilioAccountSid string
	TwilioAuthToken  string
	TwilioFromNumber string
}

type SMSConfig struct {
	TwilioAuthToken  string
	TwilioAccountSid string
	TwilioFromNumber string
}

func NewSMS(config *SMSConfig) *SMS {
	return &SMS{
		TwilioAuthToken:  config.TwilioAuthToken,
		TwilioAccountSid: config.TwilioAccountSid,
		TwilioFromNumber: config.TwilioFromNumber,
	}
}

func (s *SMS) Send(msg string, to string) error {
	log.Printf("Config sid=(%s)\nSending message(%s) to %s ", s.TwilioAccountSid, msg, to)
	accountSid := s.TwilioAccountSid

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: s.TwilioAuthToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(s.TwilioFromNumber)
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
}

type SearchResult struct {
	message string
	results []*Result
}

func (sr *SearchResult) String() string {
	msg := ""
	for _, result := range sr.results {
		msg += fmt.Sprintf("Product: %surl:%s\n", result.name, result.url)
	}
	return msg
}

func (r *Result) String() string {
	return fmt.Sprintf("Name: %surl:%s\n", strings.Trim(r.name, "\n"), r.url)
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
	url  string
}

func (r Result) Contains(search string) bool {
	log.Printf("Does %s contain: %s", r.name, search)
	return strings.Contains(strings.ToLower(r.name), strings.ToLower(search))
}

func NewSearchResults() *SearchResult {
	return &SearchResult{}
}

func main() {
	var search, to string
	flag.StringVar(&search, "search", "enfamil", "Search term to filter results by.")
	flag.StringVar(&to, "to", "", "Phone number to send sms message to. *Needs to be validated first*")
	flag.Parse()
	if len(to) == 0 {
		log.Fatalf("to needs to be set.")
	}
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

	searchResult, err := searchCostco(page)
	var filtered []*Result
	for _, result := range searchResult.results {
		if result.Contains(search) {
			filtered = append(filtered, result)
		}
	}
	twilioSsd := os.Getenv("TWILIO_API_SSD")
	twilioApiToken := os.Getenv("TWILIO_API_TOKEN")
	twilioFromNumber := os.Getenv("TWILIO_FROM_NUMBER")
	sms := NewSMS(&SMSConfig{
		TwilioAccountSid: twilioSsd,
		TwilioAuthToken: twilioApiToken,
		TwilioFromNumber: twilioFromNumber,
	})

	if len(filtered) > 0 {
		msg := fmt.Sprintf("Count %d %s", len(filtered), filtered)
		if err := sms.Send(msg, to); err != nil {
			log.Fatalf("Could not send sms: %s", err.Error())
		}
	} else {
		log.Printf("Searching for %s returned zero results.\n", search)
	}

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}
