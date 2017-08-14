package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var initalMatchRank = 99999
var baseURL = "https://www.truepeoplesearch.com"
var scraperURL = baseURL + "/results?phoneno="

type app struct {
	DB *sql.DB
}

type address struct {
	Street string `json:"Street"`
	City   string `json:"City"`
	State  string `json:"State"`
	Zip    string `json:"Zip"`
}

type person struct {
	FullName    string `json:"fn"`
	addressLink string
	matchRank   int
	phone       phoneNumber
	Address     address
}

func (p *person) save(db *sql.DB) int64 {
	numberID := addNumber(p.phone.Number, db)
	address := &address{
		Street: p.Address.Street,
		City:   p.Address.City,
		State:  p.Address.State,
		Zip:    p.Address.Zip,
	}
	addressID := addAddress(address, db)

	id := addPerson(p.FullName, numberID, addressID, db)
	return id

}

type phoneNumber struct {
	Number              string `json:"pn"`
	name                string
	Matches             []*person
	matchLock           sync.Mutex
	multipleMatchesChan chan bool
	foundMatchChan      chan bool
	noMatchChan         chan bool
	getAddressChan      chan bool
}

func (p *phoneNumber) updateStatus() {
	p.matchLock.Lock()
	numOfMatches := len(p.Matches)
	p.matchLock.Unlock()

	go func() {
		if numOfMatches == 1 && p.Matches[0].Address.Street == "" {
			p.getAddressChan <- true
		} else if numOfMatches == 0 {
			p.noMatchChan <- true
		} else if numOfMatches == 1 {
			p.foundMatchChan <- true
		} else {
			p.multipleMatchesChan <- true
		}
	}()
}

func (a *app) initialize(dbName string) {
	a.DB = newSession(dbName)
	createDBTables(a.DB)
	a.initializeRoutes()
}

func (a *app) run(address string) {
	log.Fatal(http.ListenAndServe(address, nil))
}

func (a *app) initializeRoutes() {
	http.HandleFunc("/api/search/", a.getPersonByNumber)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
}

func (pn *phoneNumber) handleMultipleResults() *phoneNumber {
	matchedPerson := &person{}
	matches := []*person{}
	// set an initial rank
	matchedPerson.matchRank = initalMatchRank
	for i := range pn.Matches {
		rank := fuzzy.RankMatch(pn.name, pn.Matches[i].FullName)
		// no match
		if rank == -1 {
			continue
		}
		pn.Matches[i].matchRank = rank
		if pn.Matches[i].matchRank < matchedPerson.matchRank {
			matchedPerson.FullName = pn.Matches[i].FullName
			matchedPerson.phone.Number = pn.Number
			matchedPerson.matchRank = pn.Matches[i].matchRank
			matchedPerson.addressLink = pn.Matches[i].addressLink
		}
	}
	matches = append(matches, matchedPerson)
	pn.Matches = matches

	return pn
}

func isValidPhoneNumber(number string) bool {
	num, err := strconv.Atoi(number)
	if err != nil {
		return false
	}
	if len(strconv.Itoa(num)) != 10 {
		return false
	}
	return true
}

func (a *app) getPersonByNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		createJSONErrorResponse(w, http.StatusOK, errors.New("Only POST method allowed").Error())
	}
	number := r.FormValue("pn")
	if !isValidPhoneNumber(number) {
		msg := "Invalid phone Number"
		createJSONErrorResponse(w, http.StatusOK, msg)
		return
	}
	pn := &phoneNumber{}
	pn.Number = number
	pn.noMatchChan = make(chan bool)
	pn.foundMatchChan = make(chan bool)
	pn.multipleMatchesChan = make(chan bool)
	pn.getAddressChan = make(chan bool)

	pn.name = r.FormValue("fn")
	getPersonFromDb(pn, a.DB)

	for {
		select {
		case <-pn.foundMatchChan:
			createJSONResponse(w, http.StatusOK, pn)
			return
		case <-pn.noMatchChan:
			// get from scraper
			url := scraperURL + pn.Number
			doc := urlDoc{url}.getDoc()
			if isCaptcha(doc) {
				pn.updateStatus()
				errTxt := "Please handle captcha @" + url
				createJSONErrorResponse(w, http.StatusOK, errors.New(errTxt).Error())
				return
			}
			pn = scrapeNumber(doc, pn)
			if len(pn.Matches) < 1 {
				createJSONErrorResponse(w, http.StatusOK, "No results found")
				return
			}
			pn.updateStatus()

		case <-pn.multipleMatchesChan:
			pn := pn.handleMultipleResults()
			if pn.Matches[0].matchRank == initalMatchRank {
				createJSONErrorResponse(w, http.StatusOK, "No results found")
				return
			}
			pn.updateStatus()

		case <-pn.getAddressChan:
			addressURL := baseURL + pn.Matches[0].addressLink
			doc := urlDoc{addressURL}.getDoc()
			person := scrapeAddress(doc, pn.Matches[0])

			person.save(a.DB)
			pn.updateStatus()
		}
	}
	createJSONResponse(w, http.StatusOK, pn)
}

func createJSONErrorResponse(w http.ResponseWriter, code int, errMsg string) {
	msg := map[string]string{"error": errMsg}
	createJSONResponse(w, code, msg)
}

func createJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.WriteHeader(code)
	w.Write(response)
}
