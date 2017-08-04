package main

import (
	// standard library packages
	"log"
	// 3rd party packages
	"errors"
	"github.com/PuerkitoBio/goquery"
	"os"
)

type SearchResults struct {
	DetailLink string
	FullName   string
}

type DetailResuts struct {
	FullAddress string
}

func (d *DetailResuts) parseFullAddress() *Address {
	a := &Address{}
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

func (u urlDoc ) getDoc() *goquery.Document {
	doc := &goquery.Document{}
	doc, err := goquery.NewDocument(u.urlPath)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

//func parseFullName(fn []string) SearchResults {
//	firstName := fullName[0]
//	if len(fn) == 2 {
//		lastName := fullName[1]
//	}
//	// TODO this scenario could be lastname and suffix
//	else if len(fn) == 3 {
//		middleName := fullName[1]
//		lastName := fullName[2] if
//	}
//	else if len(fn) == 4 {
//		middleName := fullName[1]
//		lastName := fullName[2]
//		suffix := fullName[3]
//	}
////}

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
		detailSelector := s.Find("data-detail-link")
		sr.DetailLink = detailSelector.Text()

		s.Find(".h4").Each(func(j int, h *goquery.Selection) {
			sr.FullName = h.Text()
		})
		r = append(r, sr)
	})
	return r
}

func scrapeAddress(doc *goquery.Document) *DetailResuts {
	d := &DetailResuts{}
	doc.Find(".link-to-more").Each(func(i int, s *goquery.Selection) {
		d.FullAddress = s.Text()
	})
	return d
}
