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

func CreatePersonTable(db *sql.DB) {
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
	street varchar(150),
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

func addPerson(name string, db *sql.DB) int64 {
	sqlAddPerson := `
	INSERT OR REPLACE INTO person(
		id,
		fullname
	) values(?)
	`
	stmt, err1 := db.Prepare(sqlAddPerson)
	checkErr(err1)
	defer stmt.Close()

	r, err2 := stmt.Exec(name)
	checkErr(err2)

	id, err3 := r.LastInsertId()
	checkErr(err3)

	return id
}

func addNumber(number string, db *sql.DB) int64 {
	sqlAddNumber := `
	INSERT OR REPLACE INTO phone_number(
		id,
		number
	) values(?)
	`
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
	INSERT OR REPLACE INTO phone_number(
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

func addPersonAddress(person_id, address_id int64, db *sql.DB) int64 {
	sqlAddPersonAddress := `
	INSERT OR REPLACE INTO person_address(
		person_id,
		address_id
	) values(?, ?)
	`
	stmt, err1 := db.Prepare(sqlAddPersonAddress)
	checkErr(err1)
	defer stmt.Close()

	r, err2 := stmt.Exec(person_id, address_id)
	checkErr(err2)

	id, err3 := r.LastInsertId()
	checkErr(err3)

	return id
}

func addPersonNumber(person_id, number_id int64, db *sql.DB) int64 {
	sqlAddPersonNumber := `
	INSERT OR REPLACE INTO person_number(
		person_id,
		number_id
	) values(?, ?)
	`
	stmt, err1 := db.Prepare(sqlAddPersonNumber)
	checkErr(err1)
	defer stmt.Close()

	r, err2 := stmt.Exec(person_id, number_id)
	checkErr(err2)

	id, err3 := r.LastInsertId()
	checkErr(err3)

	return id
}

func GetPersonByNumber(number string, db *sql.DB) ([]Person, error) {
	selectQuery := `
	SELECT n.number from phone_number n
	left Join person_number p on p.id = n.person_id.
	left join person_address a on a.person_id = p.id
	WHERE phone_number.number = ?
	`
	rows, err := db.Query(selectQuery, number)
	checkErr(err)
	defer rows.Close()

	person := []Person{}
	for rows.Next() {
		var p Person
		err = rows.Scan(&p.Phone.Number, &p.FullName, &p.Address.Street, &p.Address.City, &p.Address.State, &p.Address.Zip)
		checkErr(err)

		person = append(person, p)
	}
	return person, nil

}

func (p *Person) Save(db *sql.DB) {
	personID := addPerson(p.FullName, db)
	numberID := addNumber(p.Phone.Number, db)
	address := &Address{
		Street: p.Address.Street,
		City:   p.Address.City,
		State:  p.Address.State,
		Zip:    p.Address.Zip,
	}
	addressID := addAddress(address, db)
	_ = addPersonAddress(personID, addressID, db)
	_ = addPersonNumber(personID, numberID, db)
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
