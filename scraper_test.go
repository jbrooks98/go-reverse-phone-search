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
	filePath := getAbsPath("/testfiles/results.html")
	doc := FileDoc{filePath}.getDoc()
	pn := &PhoneNumber{}
	scrapeNumber(doc, pn)
	if pn.Matches[0].AddressLink != "/results?phoneno=3027502606&rid=0x0" || pn.Matches[0].FullName != "Donald James Willis Sr" {
		t.Error("scraping phone number failed")
	}
}

func TestAddressScraping(t *testing.T) {
	filePath := getAbsPath("/testfiles/details.html")
	doc := FileDoc{filePath}.getDoc()
	person := &Person{}
	person = scrapeAddress(doc, person)
	if person.Address.Street != "1 Blue Spruce Dr" {
		t.Error("scraping address street failed")
	}
	if person.Address.State != "DE" {
		t.Error("scraping address state failed")
	}
	if person.Address.City != "Bear" {
		t.Error("scraping address city failed")
	}
	if person.Address.Zip != "19701-4125" {
		t.Error("scraping address zip failed")
	}

}

func TestCaptchaScraping(t *testing.T) {
	filePath := getAbsPath("/testfiles/captcha.html")
	doc := FileDoc{filePath}.getDoc()
	result := isCaptcha(doc)

	if result != true {
		t.Error("failed to scrape captcha")
	}
}
