package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"sort"
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
	http.HandleFunc("/reverse/[a-z]{10}/", a.getPersonByNumber)
	// http.Handle("/", http.FileServer(http.Dir("./static")))
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
	if r.Method != "GET" {
		createJSONResponse(w, http.StatusMethodNotAllowed, "Request method not accepted")
	}
	phoneNumber := r.FormValue("pn")
	if !isValidPhoneNumber(phoneNumber) {
		message := map[string]string{"error": "Invalid phone number"}
		createJSONResponse(w, http.StatusBadRequest, message)
	}
	fullName := r.FormValue("fn")
	person, err := GetPeopleByNumber(phoneNumber, a.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// number not found in the DB so get from the scraper
			doc := urlDoc{scraperURL}.getDoc()
			person, err := GetPeopleFromScraper(doc, fullName)

			addressURL := baseURL + person.AddressLink
			addressDoc := urlDoc{addressURL}.getDoc()
			address := ScrapeAddress(addressDoc)
			person.Address.State = address.State
			person.Address.City = address.City
			person.Address.Street = address.Street
			person.Address.Zip = address.Zip
			person.Phone.Number = phoneNumber
			if err != nil {
				message := map[string]string{"error": err.Error()}
				createJSONResponse(w, http.StatusNotFound, message)
			} else {
				// TODO save info in DB
				person.Save(a.DB)
			}
		default:
			createJSONResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if len(person) > 1 {
		winner := &Person{}
		winner.MatchRank = 99999
		for i := range person {
			person[i].MatchRank = rankResults(fullName, person[i].FullName)
			if person[i].MatchRank < winner.MatchRank {
				winner = person[i]
			}
		}
		log.Println(winner.FullName)
	}
}

func createJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
