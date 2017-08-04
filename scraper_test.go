package main

import (
	"testing"
	"runtime"
	"path"
)

func getAbsPath(p string) string {
	_, filename, _, _ := runtime.Caller(1)
	f := path.Join(path.Dir(filename), p)
	return f
}

func TestNumberScraping(t *testing.T) {
	filePath := getAbsPath("/testfiles/results.html")
	doc := FileDoc{filePath}.getDoc()
	results := scrapeNumber(doc)

	for _, r := range results {
		if r.DetailLink == "" || r.FullName == "" {
			t.Error("scraping phone number failed")
		}
	}
}

func TestAddressScraping(t *testing.T) {
	filePath := getAbsPath("/testfiles/results.html")
	doc := FileDoc{filePath}.getDoc()
	results := scrapeAddress(doc)
	if results.FullAddress == "" {
		t.Error("scraping full address failed")
	}
}

func TestAddressParser(t *testing.T) {
	expectedResult := &Address{}
	expectedResult.Street1 = "1 Blue Spruce Dr"
	expectedResult.City = "Bear"
	expectedResult.State = "DE"
	expectedResult.Zip = "19701-4125"

	result := &DetailResuts{}
	result.FullAddress = "1 Blue Spruce Dr<br/>Bear, DE 19701-4125"

	parsedAddress := result.parseFullAddress()

	if parsedAddress != expectedResult {
		t.Error("parsing full address failed")
	}
}
