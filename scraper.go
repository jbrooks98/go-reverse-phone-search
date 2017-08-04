package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"regexp"
	"strings"
)

type SearchResults struct {
	DetailLink string
	FullName   string
}

type DetailResults struct {
	FullAddress string
}

func (d *DetailResults) parseFullAddress() *Address {
	// TODO come up with better parser
	a := &Address{}
	result := strings.Split(d.FullAddress, "\n")
	a.Street = cleanAddressField(result[1])

	result = strings.Split(result[2], ",")
	a.City = cleanAddressField(result[0])

	result = strings.Split(cleanAddressField(result[1]), " ")
	a.State = cleanAddressField(result[0])
	a.Zip = cleanAddressField(result[1])

	return a
}

type scraperDoc interface {
	getDoc() *goquery.Document
}

type FileDoc struct {
	FilePath string
}

type urlDoc struct {
	urlPath string
}

// used for testing
func (f FileDoc) getDoc() *goquery.Document {
	doc := &goquery.Document{}
	d, e := os.Open(f.FilePath)
	if e != nil {
		log.Fatal(e)
	}
	defer d.Close()
	doc, err := goquery.NewDocumentFromReader(d)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func (u urlDoc) getDoc() *goquery.Document {
	doc := &goquery.Document{}
	doc, err := goquery.NewDocument(u.urlPath)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func findNameMatch(inputName, correctName string) bool {
	return inputName == correctName
}

func GetPeopleFromScraper(doc *goquery.Document, fullName string) (*Person, error) {
	person := &Person{}
	peopleResults := scrapeNumber(doc)

	for i := 0; i <= len(peopleResults); i++ {
		if findNameMatch(peopleResults[i].FullName, fullName) {
			return person, nil
		}
	}
	return person, errors.New("Cannot find a name match")
}

func scrapeNumber(doc *goquery.Document) []*SearchResults {
	r := []*SearchResults{}
	doc.Find(".card.card-block.shadow-form.card-summary").Each(func(i int, s *goquery.Selection) {
		sr := &SearchResults{}
		sr.DetailLink = s.AttrOr("data-detail-link", "")
		s.Find(".h4").Each(func(j int, h *goquery.Selection) {
			sr.FullName = strings.TrimSpace(h.Text())
		})
		r = append(r, sr)
	})
	return r
}

func cleanAddressField(s string) string {
	addressStr := strings.TrimSpace(s)
	// remove multiple whitespaces inside the string
	re_inside_whtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

	return re_inside_whtsp.ReplaceAllString(addressStr, " ")
}

func scrapeAddress(doc *goquery.Document) *Address {
	d := &DetailResults{}
	d.FullAddress = doc.Find(".link-to-more").First().Text()
	return d.parseFullAddress()
}
