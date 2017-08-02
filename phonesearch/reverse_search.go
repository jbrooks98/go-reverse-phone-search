package phonesearch

type ContactInfo struct {
	FirstName string `json:"first_name"`
	LastName string  `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Number     int `json:"number"`
	Street1 string `json:"street1"`
	Street2 string `json:"street2"`
	Street3 string `json:"street3"`
	Street4 string `json:"street4"`
	City string `json:"city"`
	State string `json:"state"`
	Zip string `json:"zip"`
}


func fetchNumber(db *mgo.Database, number int) []ContactInfo {
	// TODO search the DB first if not found look up using the web scraper
	search := []PhoneSearch{}
	err := db.C("phonenumbers").Find(nil).All(&search)
	if err != nil {
		panic(err)
	}
	// if not found in DB get the results from the web scraper
    if not search {
		results := scraper(phoneNumber)

	}

	return search
}


func handleMultipleResults(results [][]ContactInfo, firstName, lastName string) ContactInfo {
	// TODO if search yields multiple people for a number find the specific person

	for item in range results {

	}
}


func storeResults(db *Database, ci *ContactInfo) {
	// TODO store results from scraper in the local db
}