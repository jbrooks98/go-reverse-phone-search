package phonesearch

import (
	// standard library packages
	"log"
    // 3rd party packages
	"github.com/PuerkitoBio/goquery"
)
type SearchResults struct {
	DetailLink string
	FullName string
}

type DetailResuts struct {
	FullAddress string

}

type Address struct {
	street1 string
	street2 string
	street3 string
	city    string
	state   string
	zip     string
}

func (d *DetailResuts) parseFullAddress() *Address {
	a := &Address{}
	return a
}
// var scraperURL = "https://www.truepeoplesearch.com/results?phoneno="

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
//
//}

func scrapeNumber(url string) []*SearchResults {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
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

func scrapeAddress(url string) *DetailResuts {
    d := &DetailResuts{}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".link-to-more").Each(func(i int, s *goquery.Selection) {
		d.FullAddress = s.Text()
	})
	return d
}
