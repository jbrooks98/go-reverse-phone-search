package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func newSession(name string) *sql.DB {
	db, err := sql.Open("sqlite3", name)
	checkErr(err)

	return db
}

func createAddressTable(db *sql.DB) {
	sqlTable := `CREATE TABLE IF NOT EXISTS Address(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	Street VARCHAR(150) NOT NULL,
	City VARCHAR(150) NOT NULL,
	State VARCHAR(20) NOT NULL,
	Zip VARCHAR(20) NOT NULL,
	CONSTRAINT unique_address UNIQUE (Street, City, State, Zip)
	)`
	_, err := db.Exec(sqlTable)
	checkErr(err)
}

func createPhoneNumberTable(db *sql.DB) {
	sqlTable := `CREATE TABLE IF NOT EXISTS phone_number(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	Number VARCHAR(15) NOT NULL,
	CONSTRAINT unique_number UNIQUE (Number)
	)`
	_, err := db.Exec(sqlTable)
	checkErr(err)
}

func createPersonTable(db *sql.DB) {
	sqlTable := `CREATE TABLE IF NOT EXISTS person(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	fullname VARCHAR(200) NOT NULL,
	address_id INTEGER NOT NULL,
	phone_number_id INTEGER NOT NULL,
	FOREIGN KEY(address_id) REFERENCES Address(id),
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
	sqlAddNumber := `INSERT OR REPLACE INTO phone_number (Number) values (?)`
	stmt, err1 := db.Prepare(sqlAddNumber)
	checkErr(err1)
	defer stmt.Close()

	r, err2 := stmt.Exec(number)
	checkErr(err2)

	id, err3 := r.LastInsertId()
	checkErr(err3)

	return id
}

func addAddress(a *address, db *sql.DB) int64 {
	sqlAddNumber := `
	INSERT OR REPLACE INTO Address(
	Street,
	City,
	State,
	Zip
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

func getPersonFromDb(pn *phoneNumber, db *sql.DB) {
	query := `SELECT person.fullname, phone_number.Number, Address.Street, Address.City, Address.State, Address.Zip
	FROM person
	JOIN phone_number ON person.phone_number_id = phone_number.id
	JOIN Address ON person.address_id = Address.id
	WHERE phone_number.Number=?`

	row := &person{}
	matches := []*person{}
	err := db.QueryRow(query, pn.Number).Scan(
		&row.FullName, &row.phone.Number, &row.Address.Street,
		&row.Address.City, &row.Address.State, &row.Address.Zip,
	)
	switch {
	case err == sql.ErrNoRows:
		pn.updateStatus()
	case err != nil:
		log.Panic(err)
	default:
		matches = append(matches, row)
		pn.Matches = matches
		pn.updateStatus()
	}
}

func createDBTables(db *sql.DB) {
	createAddressTable(db)
	createPhoneNumberTable(db)
	createPersonTable(db)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
