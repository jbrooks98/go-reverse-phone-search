package main

import (
	"os"
	"testing"
)

func TestModels(t *testing.T) {
	app := &app{}
	tableNames := [4]string{"phone_number", "person", "Address"}
	dbName := "./unittest.db"

	app.initialize(dbName)
	defer os.Remove(dbName)
	defer app.DB.Close()

	// test tables created
	createDBTables(app.DB)
	for r := range tableNames {
		query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
		app.DB.Exec(query, r)
	}

	// test data
	number := "4112223333"
	name := "John Doe"

	person := &person{}
	person.FullName = name
	person.Address.State = "DE"
	person.Address.Zip = "19802"
	person.Address.Street = "123 Main St."
	person.Address.City = "Wilmington"
	person.phone.Number = number
	person.save(app.DB)

	pn := &phoneNumber{}
	pn.Number = number
	getPersonFromDb(pn, app.DB)

	// test that a person was added and we can retrieve them from the db
	if pn.Matches[0].FullName != person.FullName {
		t.Error("person name not found")
	}

	if pn.Matches[0].phone.Number != person.phone.Number {
		t.Error("person Number not found")
	}

	if pn.Matches[0].Address.Street != person.Address.Street {
		t.Error("person Address not found")
	}
}
