package go_reverse_phone_search

import (
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
	"encoding/json"
	"log"
	"strings"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

type Address struct {
	Street1   string  `json:"st1"`
	Street2   string  `json:"st2"`
	Street3   string  `json:"st3"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Zip       string  `json:"zip"`
}

type Person struct {
	FullName  string       `json:"fn"`
	Number    PhoneNumber
	Address   Address
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
	// TODO check that phone is valid format
	//if err != nil {
	//	message := map[string]string{"error": "Invalid phone number"}
	//	createJSONResponse(w, http.StatusBadRequest, message)
	//}
	people, err :=  GetPeopleByNumber(phoneNumber, a.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// TODO call scraper
			if err := ph.getPeopleFromScraper(); err != nil {
				message := map[string]string{"error": "Person not found"}
				createJSONResponse(w, http.StatusNotFound, message)
			} else {
				// TODO save info in DB
				personInfo.save()

			}
		default:
			createJSONResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
}

func createJSONResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
