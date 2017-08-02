package phonesearch

import (
	"testing"
)

func TestNumberScraping(t *testing.T) {
	results := scrapeNumber("./tests/results.html")

	for _, r := range results {
		if r.DetailLink == "" || r.FullName == "" {
			t.Error("scraping phone number failed")
		}
	}
}

func TestAddressScraping(t *testing.T) {
	results := scrapeAddress("./tests/details.html")
	if results.FullAddress == "" {
		t.Error("scraping full address failed")
	}
}

func TestAddressParser(t *testing.T) {
	expectedResult := &Address{}
	expectedResult.street1 = "1 Blue Spruce Dr"
	expectedResult.city = "Bear"
	expectedResult.state = "DE"
	expectedResult.zip = "19701-4125"

	result := &DetailResuts{}
	result.FullAddress = "1 Blue Spruce Dr<br/>Bear, DE 19701-4125"

	parsedAddress := result.parseFullAddress()

	if parsedAddress != expectedResult {
		t.Error("parsing full address failed")
	}
}