package main

import (
	"database/sql"
	"encoding/json"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"log"
	"net/http"
	"strconv"
	"sync"
	"os"
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
	log.Println(numOfMatches)

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

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (a *App) Initialize(dbName string) {
	dbName = "./test.db"
	a.DB = NewSession(dbName)
	CreateDBTables(a.DB)
	a.initializeRoutes()
}

func (a *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, nil))
}

func (a *App) initializeRoutes() {
	// TODO update regex to accept a phone number
	http.HandleFunc("/search/", a.getPersonByNumber)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
}

func (pn *PhoneNumber) handleMultipleResults() *Person {
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
			matchedPerson = pn.Matches[i]
		}
	}
	matches = append(matches, matchedPerson)
	pn.Matches = matches

	return matchedPerson
}

func isValidPhoneNumber(number string) bool {
	if len(number) > 10 {
		return false
	}
	_, err := strconv.Atoi(number)
	if err != nil {
		return false
	}
	return true
}

func (a *App) getPersonByNumber(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		msg := "Request method not allowed"
		log.Println(msg)
		createJSONErrorResponse(w, http.StatusOK, msg)
		return
	}

	number := r.FormValue("pn")
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
			doc := urlDoc{scraperURL + pn.Number}.getDoc()
			_, err := scrapeNumber(doc, pn, a.DB)
			if err != nil {
				createJSONErrorResponse(w, http.StatusOK, err.Error())
				return
			}
		case <-pn.multipleMatchesChan:
			log.Println("multiple Matches")
			personMatch := pn.handleMultipleResults()
			if personMatch.MatchRank == initalMatchRank {
				createJSONErrorResponse(w, http.StatusOK, "No match found")
				return
			}
			scrapeAddress(personMatch, pn, a.DB)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
