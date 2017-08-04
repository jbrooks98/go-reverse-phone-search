package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var scraperURL = "https://www.truepeoplesearch.com/results?phoneno="

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

type Address struct {
	Street1 string `json:"st1"`
	Street2 string `json:"st2"`
	Street3 string `json:"st3"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
}

type Person struct {
	FullName string `json:"fn"`
	Phone    PhoneNumber
	Address  Address
}

type PhoneNumber struct {
	Number string `json:"pn"`
}

func (a *App) Initialize(dbname string) {
	a.DB = NewSession(dbname)
	CreateDBTables(a.DB)
	a.DB.Close()
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

func (a *App) initializeRoutes() {
	// TODO update regex to accept a phone number
	a.Router.HandleFunc("/number/{phone:[0-9]+}", a.getPersonByNumber).Methods("GET")
}

func (a *App) getPersonByNumber(w http.ResponseWriter, r *http.Request) {
	phoneNumber := r.FormValue("pn")
	fullName := r.FormValue("fn")
	// TODO check that phone number is valid format
	//if err != nil {
	//	message := map[string]string{"error": "Invalid phone number"}
	//	createJSONResponse(w, http.StatusBadRequest, message)
	//}
	person, err := GetPersonByNumber(phoneNumber, a.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// TODO call scraper
			doc := urlDoc{scraperURL}.getDoc()
			person, err := GetPeopleFromScraper(doc, fullName)
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
		// TODO search by name
		log.Println(fullName)
	}
}

func createJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
