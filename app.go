package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/renstrom/fuzzysearch/fuzzy"
)

var baseURL = "https://www.truepeoplesearch.com"
var scraperURL = baseURL + "/results?phoneno="

type App struct {
	DB     *sql.DB
}

type Address struct {
	Street string `json:"st"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
}

type Person struct {
	FullName    string `json:"fn"`
	AddressLink string `json:"address_link"`
	MatchRank   int    `json:"rank"`
	Phone       PhoneNumber
	Address     Address
}

type PhoneNumber struct {
	Number string `json:"pn"`
}

func (a *App) Initialize(dbName string) {
	a.DB = NewSession(dbName)
	CreateDBTables(a.DB)
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func (a *App) initializeRoutes() {
	// TODO update regex to accept a phone number
	http.HandleFunc("/search/", a.getPersonByNumber)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
}

func rankResults(userInput string, value string) int {
	return fuzzy.RankMatch(userInput, value)
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
	log.Println("request heard")
	if r.Method != "GET" {
		createJSONResponse(w, http.StatusMethodNotAllowed, "Request method not accepted")
	}
	phoneNumber := "3027502606" // r.FormValue("pn")
	log.Println(phoneNumber)
	if !isValidPhoneNumber(phoneNumber) {
		message := map[string]string{"error": "Invalid phone number"}
		createJSONResponse(w, http.StatusBadRequest, message)
	}
	fullName := "Jason E Brooks"// r.FormValue("fn")
	log.Println(fullName)
	people, err := GetPeopleByNumber(phoneNumber, a.DB)
	if err != nil {
		createJSONResponse(w, http.StatusInternalServerError, err)
	}
	if len(people) < 1 {
		// number not found in the DB so get from the scraper
		doc := urlDoc{scraperURL + phoneNumber}.getDoc()
		person, err := GetPersonFromScraper(doc, fullName)
		if err != nil {
			message := map[string]string{"error": err.Error()}
			createJSONResponse(w, http.StatusNotFound, message)
		}
        log.Println("came back from scraper")
		//TODO If link is not populated due to not finding the matching name will throw http panic - indoex out of range
		addressURL := baseURL + person.AddressLink
		log.Println(addressURL)
		addressDoc := urlDoc{addressURL}.getDoc()
		address := ScrapeAddress(addressDoc)
		log.Println("scraped the address")
		person.Address.State = address.State
		person.Address.City = address.City
		person.Address.Street = address.Street
		person.Address.Zip = address.Zip
		person.Phone.Number = phoneNumber

		person.Save(a.DB)
		log.Println("saved to db")
		createJSONResponse(w, http.StatusOK, person)
	}
	if len(people) > 1 {
		winner := &Person{}
		winner.MatchRank = 99999
		for i := range people {
			people[i].MatchRank = rankResults(fullName, people[i].FullName)
			if people[i].MatchRank < winner.MatchRank {
				winner = people[i]
			}
		}
		if winner.FullName != "" {
			createJSONResponse(w, http.StatusOK, winner)
		} else {
			createJSONResponse(w, http.StatusNotFound, "No results found")
		}
	} else {
		log.Println("found in db")
		createJSONResponse(w, http.StatusOK, people)
	}

}

func createJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
