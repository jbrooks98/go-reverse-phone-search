package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func NewSession(name string) *sql.DB {
	db, err := sql.Open("sqlite3", name)
	checkErr(err)

	return db
}

func CreateAddressTable(db *sql.DB) {
	sqlTable := `CREATE TABLE IF NOT EXISTS address(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	street VARCHAR(150) NOT NULL,
	city VARCHAR(150) NOT NULL,
	state VARCHAR(20) NOT NULL,
	zip VARCHAR(20) NOT NULL,
    CONSTRAINT unique_address UNIQUE (street, city, state, zip)
	)`
	_, err := db.Exec(sqlTable)
	checkErr(err)
}

func CreatePhoneNumberTable(db *sql.DB) {
	sqlTable := `CREATE TABLE IF NOT EXISTS phone_number(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		number VARCHAR(15) NOT NULL,
		CONSTRAINT unique_number UNIQUE (number)
	)`
	_, err := db.Exec(sqlTable)
	checkErr(err)
}

func CreatePersonTable(db *sql.DB) {
	sqlTable := `CREATE TABLE IF NOT EXISTS person(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fullname VARCHAR(200) NOT NULL,
		address_id INTEGER NOT NULL,
		phone_number_id INTEGER NOT NULL,
        FOREIGN KEY(address_id) REFERENCES address(id),
        FOREIGN KEY(phone_number_id) REFERENCES phone_number(id),
        CONSTRAINT unique_person UNIQUE (fullname, address_id, phone_number_id)
	)`
	_, err := db.Exec(sqlTable)
	checkErr(err)
}

func addPerson(name string, numberID, addressID int64, db *sql.DB) int64 {
	sqlAddPerson := `INSERT OR REPLACE INTO person(
		fullname,
		address_id,
		phone_number_id
	) values(?, ?, ?)`
	stmt, err1 := db.Prepare(sqlAddPerson)
	checkErr(err1)
	defer stmt.Close()

	r, err2 := stmt.Exec(name, addressID, numberID)
	checkErr(err2)

	id, err3 := r.LastInsertId()
	checkErr(err3)

	return id
}

func addNumber(number string, db *sql.DB) int64 {
	sqlAddNumber := `INSERT OR REPLACE INTO phone_number (number) values (?)`
	stmt, err1 := db.Prepare(sqlAddNumber)
	checkErr(err1)
	defer stmt.Close()

	r, err2 := stmt.Exec(number)
	checkErr(err2)

	id, err3 := r.LastInsertId()
	checkErr(err3)

	return id
}

func addAddress(a *Address, db *sql.DB) int64 {
	sqlAddNumber := `
	INSERT OR REPLACE INTO address(
		street,
		city,
		state,
		zip
	) values(?, ?, ?, ?)
	`
	stmt, err1 := db.Prepare(sqlAddNumber)
	checkErr(err1)
	defer stmt.Close()

	r, err2 := stmt.Exec(a.Street, a.City, a.State, a.Zip)
	checkErr(err2)

	id, err3 := r.LastInsertId()
	checkErr(err3)

	return id
}

func getPersonFromDb(pn *PhoneNumber, db *sql.DB) {
	query := `SELECT person.fullname, phone_number.number, address.street, address.city, address.state, address.zip
	FROM person
	JOIN phone_number ON person.phone_number_id = phone_number.id
	JOIN address ON person.address_id = address.id
	WHERE phone_number.number=?`

	row := &Person{}
	matches := []*Person{}
	err := db.QueryRow(query, pn.Number).Scan(
		&row.FullName, &row.Phone.Number, &row.Address.Street,
		&row.Address.City, &row.Address.State, &row.Address.Zip,
	)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No person with that number in DB.")
		pn.updateStatus()
	case err != nil:
		log.Panic(err)
	default:
		matches = append(matches, row)
		pn.Matches = matches
		pn.updateStatus()
	}
}

func CreateDBTables(db *sql.DB) {
	CreateAddressTable(db)
	CreatePhoneNumberTable(db)
	CreatePersonTable(db)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
