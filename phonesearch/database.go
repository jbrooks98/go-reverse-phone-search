package phonesearch

type DatabaseSession struct {
	*mgo.Session
	databaseName string
}

/*
Connect to the local Db and set up the database.
*/
func NewSession(name string) *DatabaseSession {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	addIndexToSignatureEmails(session.DB(name))
	return &DatabaseSession{session, name}
}

func addIndexToPhoneNumbers(db *mgo.Database) {
	index := mgo.Index{
		Key:      []string{"email"},
		Unique:   true,
		DropDups: true,
	}
	indexErr := db.C("signatures").EnsureIndex(index)
	if indexErr != nil {
		panic(indexErr)
	}
}

func createPersonTable(db *Database) {
	// think about just storing the full name.  because it comes back from the scraper as just text we will have issues
	// parsing on suffix
	// TODO fields full name

}


func createAddressTable(db *Database) {
	// TODO street1-4, city, state, zip

}


func createPhoneNumberTable(db *Database) {
	// TODO 10 digit number

}


func createContactTable(db *Database) {
	// Contact Table (all are unique together)
	// - fk phone number
	// - fk address table
	// - fk person table
}

func addContactInfo(ci *ContactInfo) {
	// TODO add 
}