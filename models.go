package go_reverse_phone_search

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// Connect to the local Db and set up the database.

func NewSession(name string) *sql.DB {
	db, err := sql.Open("sqlite3", name)
	checkErr(err)

	return db
}

func CreatePersonTable(db *sql.DB) {
	// think about just storing the full name.  because it comes back from the scraper as just text we will have issues
	// parsing on suffix
	sqlTable := `
	CREATE TABLE IF NOT EXISTS person(
		id INTEGER PRIMARY KEY,
		fullname varchar(255) NOT NULL,
	);
	`
	_, err := db.Exec(sqlTable)
	checkErr(err)
}

func CreateAddressTable(db *sql.DB) {
	sqlTable := `
	CREATE TABLE IF NOT EXISTS address(
	id INTEGER PRIMARY KEY,
	street1 varchar(150),
	street2 varchar(150),
	street3 varchar(150),
	city varchar(150),
	state varchar(20),
	zip varchar(20),
	);
	`
	_, err := db.Exec(sqlTable)
	checkErr(err)
}


func CreatePhoneNumberTable(db *sql.DB) {
	sqlTable := `
	CREATE TABLE IF NOT EXISTS phone_number(
		id INTEGER PRIMARY KEY,
		number varchar(15) NOT NULL,
	);
	`
	_, err := db.Exec(sqlTable)
	checkErr(err)
}

func CreatePersonAddressTable(db *sql.DB) {
	sqlTable := `
	CREATE TABLE IF NOT EXISTS person_address(
		id INTEGER PRIMARY KEY,
		person_id INTEGER PRIMARY KEY,
		address_id INTEGER PRIMARY KEY
	);
	`
	_, err := db.Exec(sqlTable)
	checkErr(err)
}
func CreatePersonPhoneNumberTable(db *sql.DB) {
	sqlTable := `
	CREATE TABLE IF NOT EXISTS person_phone_number(
		id INTEGER PRIMARY KEY,
		person_id INTEGER PRIMARY KEY,
		phone_number_id INTEGER PRIMARY KEY
	);
	`
	_, err := db.Exec(sqlTable)
	checkErr(err)
}

func AddPerson(name string, db *sql.DB) {
	sqlAddPerson := `
	INSERT OR REPLACE INTO person(
		id,
		fullname
	) values(?)
	`
	stmt, err1 := db.Prepare(sqlAddPerson)
	checkErr(err1)
	defer stmt.Close()

	_, err2 := stmt.Exec(name)
	checkErr(err2)
}

func AddNumber(number string, db *sql.DB) {
	sqlAddNumber := `
	INSERT OR REPLACE INTO phone_number(
		id,
		number
	) values(?)
	`
	stmt, err1 := db.Prepare(sqlAddNumber)
	checkErr(err1)
	defer stmt.Close()

	_, err2 := stmt.Exec(number)
	checkErr(err2)
}

func AddAddress(a *Address, db *sql.DB) {
	sqlAddNumber := `
	INSERT OR REPLACE INTO phone_number(
		iscurrent,
		street1,
		street2,
		street3,
		city,
		state,
		zip
	) values(?, ?, ?, ?, ?, ?, ?)
	`
	stmt, err1 := db.Prepare(sqlAddNumber)
	checkErr(err1)
	defer stmt.Close()

	_, err2 := stmt.Exec(a.street1, a.street2, a.street3, a.city, a.state, a.zip)
	checkErr(err2)
}

func AddPersonAddress(db *sql.DB, person_id, address_id int) {
	sqlAddPersonAddress := `
	INSERT OR REPLACE INTO person_address(
		person_id,
		address_id
	) values(?, ?)
	`
	stmt, err1 := db.Prepare(sqlAddPersonAddress)
	checkErr(err1)
	defer stmt.Close()

	_, err2 := stmt.Exec(person_id, address_id)
	checkErr(err2)
}

type people struct {
	Number    string  `json:"pn"`
	FullName  string  `json:"fn"`
	Street1   string  `json:"st1"`
	Street2   string  `json:"st2"`
	Street3   string  `json:"st3"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Zip       string  `json:"zip"`
}

func GetPeopleByNumber(number string, db *sql.DB)  ([]people, error) {
	rows, err := db.Query("SELECT phone_number.phonenumber, fullname FROM person_phone_number WHERE phonenumber=?", number)
	checkErr(err)
	defer rows.Close()

	people := []people{}
	for rows.Next() {
		var p people
		err = rows.Scan(&p.Number, &p.FullName, &p.Street1, &p.Street2, &p.Street3, &p.City, &p.State, &p.Zip)
		checkErr(err)

		people = append(people, p)
	}
	return people, nil

}

func CreateDBTables(db *sql.DB) {
	CreatePersonTable(db)
	CreateAddressTable(db)
	CreatePhoneNumberTable(db)
	CreatePersonPhoneNumberTable(db)
	CreatePersonAddressTable(db)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}