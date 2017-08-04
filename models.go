package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func NewSession(name string) *sql.DB {
	db, err := sql.Open("sqlite3", name)
	checkErr(err)

	return db
}

func CreateAddressTable(db *sql.DB) {
	sqlTable := `CREATE TABLE IF NOT EXISTS address(
	id INTEGER PRIMARY KEY,
	street VARCHAR(150),
	city VARCHAR(150),
	state VARCHAR(20),
	zip VARCHAR(20)
	)`
	_, err := db.Exec(sqlTable)
	checkErr(err)
}

func CreatePhoneNumberTable(db *sql.DB) {
	sqlTable := `CREATE TABLE IF NOT EXISTS phone_number(
		id INTEGER PRIMARY KEY,
		number VARCHAR(15) NOT NULL
	)`
	_, err := db.Exec(sqlTable)
	checkErr(err)
}

func CreatePersonTable(db *sql.DB) {
	sqlTable := `CREATE TABLE IF NOT EXISTS person(
		id INTEGER PRIMARY KEY,
		fullname VARCHAR(200),
		address_id INTEGER,
		phone_number_id INTEGER,
        FOREIGN KEY(address_id) REFERENCES address(id),
        FOREIGN KEY(phone_number_id) REFERENCES phone_number(id)
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
	sqlAddNumber := `INSERT OR REPLACE INTO phone_number(
		number
	) values(?)`
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

func GetPeopleByNumber(phnum string, db *sql.DB) ([]*Person, error) {
	query := `SELECT person.fullname, phone_number.number, address.street, address.city, address.state, address.zip
	FROM person
	JOIN phone_number ON person.phone_number_id = phone_number.id
	JOIN address ON person.address_id = address.id
	WHERE phone_number.number = ?`
	rows, err := db.Query(query, phnum)
	checkErr(err)
	defer rows.Close()

	person := []*Person{}
	for rows.Next() {
		p := &Person{}
		err = rows.Scan(&p.FullName	, &p.Phone.Number, &p.Address.Street, &p.Address.City, &p.Address.State, &p.Address.Zip)
		checkErr(err)

		person = append(person, p)
	}
	return person, nil

}

func (p *Person) Save(db *sql.DB) int64 {
	numberID := addNumber(p.Phone.Number, db)
	address := &Address{
		Street: p.Address.Street,
		City:   p.Address.City,
		State:  p.Address.State,
		Zip:    p.Address.Zip,
	}
	addressID := addAddress(address, db)
	return addPerson(p.FullName, numberID, addressID, db)

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
