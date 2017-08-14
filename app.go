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

type App struct {
	DB *sql.DB
}

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
}

type Person struct {
	FullName    string `json:"fn"`
	AddressLink string
	MatchRank   int
	Phone       PhoneNumber
	Address     Address
}

func (p *Person) save(db *sql.DB) int64 {
	// TODO wrap in a transaction
	numberID := addNumber(p.Phone.Number, db)
	address := &Address{
		Street: p.Address.Street,
		City:   p.Address.City,
		State:  p.Address.State,
		Zip:    p.Address.Zip,
	}
	addressID := addAddress(address, db)

	id := addPerson(p.FullName, numberID, addressID, db)
	return id

}

type PhoneNumber struct {
	Number              string `json:"pn"`
	Name                string `json:"name"`
	Matches             []*Person
	matchLock           sync.Mutex
	multipleMatchesChan chan bool
	foundMatchChan      chan bool
	noMatchChan         chan bool
}

func (p *PhoneNumber) updateStatus() {
	p.matchLock.Lock()
	numOfMatches := len(p.Matches)
	p.matchLock.Unlock()

	go func() {
		if numOfMatches == 0 {
			log.Println("no match channel")
			p.noMatchChan <- true
		} else if numOfMatches == 1 {
			log.Println("1 Match")
			p.foundMatchChan <- true
		} else {
			log.Println("multiples")
			p.multipleMatchesChan <- true
		}
	}()
}

func (a *App) Initialize(dbName string) {
	a.DB = NewSession(dbName)
	CreateDBTables(a.DB)
	a.initializeRoutes()
}

func (a *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, nil))
}

func (a *App) initializeRoutes() {
	http.HandleFunc("/api/search/", a.getPersonByNumber)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
}

func (pn *PhoneNumber) handleMultipleResults() *PhoneNumber {
	matchedPerson := &Person{}
	matches := []*Person{}
	// set an initial rank
	matchedPerson.MatchRank = initalMatchRank
	for i := range pn.Matches {
		log.Println(pn.Name, pn.Matches[i].FullName)
		rank := fuzzy.RankMatch(pn.Name, pn.Matches[i].FullName)
		// no match
		if rank == -1 {
			continue
		}
		pn.Matches[i].MatchRank = rank
		if pn.Matches[i].MatchRank < matchedPerson.MatchRank {
			matchedPerson.FullName = pn.Matches[i].FullName
			matchedPerson.Phone.Number = pn.Number
			matchedPerson.MatchRank = pn.Matches[i].MatchRank
			matchedPerson.AddressLink = pn.Matches[i].AddressLink
		}
	}
	log.Println("matched person", matchedPerson)
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

func (a *App) getPersonByNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		createJSONErrorResponse(w, http.StatusOK, errors.New("Only POST method allowed").Error())
	}
	number := r.FormValue("pn")
	log.Println("num", number)
	if !isValidPhoneNumber(number) {
		msg := "Invalid phone number"
		log.Println(msg)
		createJSONErrorResponse(w, http.StatusOK, msg)
		return
	}
	pn := &PhoneNumber{}
	pn.Number = number
	pn.noMatchChan = make(chan bool)
	pn.foundMatchChan = make(chan bool)
	pn.multipleMatchesChan = make(chan bool)

	pn.Name = r.FormValue("fn")
	log.Println("fn and pn ", pn.Name, number)
	getPersonFromDb(pn, a.DB)

	for {
		select {
		case <-pn.foundMatchChan:
			createJSONResponse(w, http.StatusOK, pn)
			return
		case <-pn.noMatchChan:
			// get from scraper
			log.Println("recieved No match")
			url := scraperURL + pn.Number
			doc := urlDoc{url}.getDoc()
			if isCaptcha(doc) {
				pn.updateStatus()
				errTxt := "Please handle captcha " + url
				createJSONErrorResponse(w, http.StatusOK, errors.New(errTxt).Error())
				return
			}
			pn = scrapeNumber(doc, pn)
			if len(pn.Matches) < 1 {
				createJSONErrorResponse(w, http.StatusOK, "No results found")
				return
			}
			addressURL := baseURL + pn.Matches[0].AddressLink
			addressDoc := urlDoc{addressURL}.getDoc()
			person := scrapeAddress(addressDoc, pn.Matches[0])
			person.save(a.DB)
			pn.updateStatus()

		case <-pn.multipleMatchesChan:
			log.Println("multiple Matches")
			pn := pn.handleMultipleResults()
			if pn.Matches[0].MatchRank == initalMatchRank {
				createJSONErrorResponse(w, http.StatusOK, "No results found")
				return
			}
			addressURL := baseURL + pn.Matches[0].AddressLink
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
