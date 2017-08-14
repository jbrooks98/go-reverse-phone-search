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
	doc := fileDoc{filePath}.getDoc()
	pn := &phoneNumber{}
	scrapeNumber(doc, pn)
	if pn.Matches[0].addressLink != "/results?phoneno=3027502606&rid=0x0" || pn.Matches[0].FullName != "Donald James Willis Sr" {
		t.Error("scraping phone Number failed")
	}
}

func TestAddressScraping(t *testing.T) {
	filePath := getAbsPath("/testfiles/details.html")
	doc := fileDoc{filePath}.getDoc()
	person := &person{}
	person = scrapeAddress(doc, person)
	if person.Address.Street != "1 Blue Spruce Dr" {
		t.Error("scraping Address Street failed")
	}
	if person.Address.State != "DE" {
		t.Error("scraping Address State failed")
	}
	if person.Address.City != "Bear" {
		t.Error("scraping Address City failed")
	}
	if person.Address.Zip != "19701-4125" {
		t.Error("scraping Address Zip failed")
	}

}

func TestCaptchaScraping(t *testing.T) {
	filePath := getAbsPath("/testfiles/captcha.html")
	doc := fileDoc{filePath}.getDoc()
	result := isCaptcha(doc)

	if result != true {
		t.Error("failed to scrape captcha")
	}
}
