package main

import (
	"path"
	"runtime"
	"testing"
)

func getAbsPath(p string) string {
	_, filename, _, _ := runtime.Caller(1)
	f := path.Join(path.Dir(filename), p)
	return f
}

func TestNumberScraping(t *testing.T) {
	filePath := getAbsPath("/testfiles/results.templates")
	doc := FileDoc{filePath}.getDoc()
	results := scrapeNumber(doc)
	if results[0].DetailLink != "/results?phoneno=3027502606&rid=0x0" || results[0].FullName != "Donald James Willis Sr" {
		t.Error("scraping phone number failed")
	}
}

func TestAddressScraping(t *testing.T) {
	filePath := getAbsPath("/testfiles/details.templates")
	doc := FileDoc{filePath}.getDoc()
	result := ScrapeAddress(doc)
	if result.Street != "1 Blue Spruce Dr" {
		t.Error("scraping address street failed")
	}
	if result.State != "DE" {
		t.Error("scraping address state failed")
	}
	if result.City != "Bear" {
		t.Error("scraping address city failed")
	}
	if result.Zip != "19701-4125" {
		t.Error("scraping address zip failed")
	}

}
