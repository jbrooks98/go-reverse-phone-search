package main

import (
	"database/sql"
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

func parseFullAddress(fullAdress string) *Address {
	// TODO come up with better parser
	a := &Address{}
	result := strings.Split(fullAdress, "\n")
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

func isCaptcha(doc *goquery.Document) bool {
	foundCaptcha := false
	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		pageTitle := s.Find("title").Text()
		if strings.ToLower(pageTitle) == "captcha" {
			foundCaptcha = true
		}
	})
	return foundCaptcha
}

func scrapeNumber(doc *goquery.Document, pn *PhoneNumber, db *sql.DB) (*PhoneNumber, error) {
	if isCaptcha(doc) {
		pn.updateStatus()
		return pn, errors.New("Please handle captcha")
	}
	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		pageTitle := s.Find("title").Text()
		if strings.ToLower(pageTitle) == "captcha" {

		}
	})
	doc.Find(".card.card-block.shadow-form.card-summary").Each(func(i int, s *goquery.Selection) {
		person := &Person{}
		person.Phone.Number = pn.Number
		person.AddressLink = s.AttrOr("data-detail-link", "")

		s.Find(".h4").Each(func(j int, h *goquery.Selection) {
			person.FullName = strings.TrimSpace(h.Text())
			pn.Matches = append(pn.Matches, person)
			log.Println("got fullName")
		})
	})
	pn.updateStatus()
	return pn, nil
}

// TODO - possibly put into a Go routine and listen on the chaneel from the scrape Number to do its work
func cleanAddressField(s string) string {
	addressStr := strings.TrimSpace(s)
	// removes multiple whitespaces inside the string
	re_inside_whtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

	return re_inside_whtsp.ReplaceAllString(addressStr, " ")
}

func scrapeAddress(p *Person, pn *PhoneNumber, db *sql.DB) {
	addressURL := baseURL + p.AddressLink
	doc := urlDoc{addressURL}.getDoc()

	fullAddress := doc.Find(".link-to-more").First().Text()
	address := parseFullAddress(fullAddress)

	p.Address.State = address.State
	p.Address.City = address.City
	p.Address.Street = address.Street
	p.Address.Zip = address.Zip
	p.Save(db)
	pn.updateStatus()
}
