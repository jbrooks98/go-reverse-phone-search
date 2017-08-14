package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"regexp"
	"strings"
)

type scraperDoc interface {
	getDoc() *goquery.Document
}

type fileDoc struct {
	FilePath string
}

type urlDoc struct {
	urlPath string
}

// used for testing
func (f fileDoc) getDoc() *goquery.Document {
	doc := &goquery.Document{}
	d, e := os.Open(f.FilePath)
	if e != nil {
		log.Println(e)
	}
	defer d.Close()
	doc, err := goquery.NewDocumentFromReader(d)
	if err != nil {
		log.Println(err)
	}
	return doc
}

func (u urlDoc) getDoc() *goquery.Document {
	doc := &goquery.Document{}
	doc, err := goquery.NewDocument(u.urlPath)
	if err != nil {
		log.Println(err)
	}
	return doc
}

func parseFullAddress(fullAdress string) *address {
	// TODO come up with better parser
	a := &address{}
	result := strings.Split(fullAdress, "\n")
	a.Street = cleanAddressField(result[1])

	result = strings.Split(result[2], ",")
	a.City = cleanAddressField(result[0])

	result = strings.Split(cleanAddressField(result[1]), " ")
	a.State = cleanAddressField(result[0])
	a.Zip = cleanAddressField(result[1])

	return a
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

func scrapeNumber(doc *goquery.Document, pn *phoneNumber) *phoneNumber {
	doc.Find(".card.card-block.shadow-form.card-summary").Each(func(i int, s *goquery.Selection) {
		person := &person{}
		person.phone.Number = pn.Number
		person.addressLink = s.AttrOr("data-detail-link", "")

		s.Find(".h4").Each(func(j int, h *goquery.Selection) {
			person.FullName = strings.TrimSpace(h.Text())
			pn.Matches = append(pn.Matches, person)
		})
	})
	return pn
}

func cleanAddressField(s string) string {
	addressStr := strings.TrimSpace(s)
	// removes multiple whitespaces inside the string
	re_inside_whtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

	return re_inside_whtsp.ReplaceAllString(addressStr, " ")
}

func scrapeAddress(doc *goquery.Document, person *person) *person {
	fullAddress := doc.Find(".link-to-more").First().Text()
	address := parseFullAddress(fullAddress)

	person.Address.State = address.State
	person.Address.City = address.City
	person.Address.Street = address.Street
	person.Address.Zip = address.Zip

	return person
}
