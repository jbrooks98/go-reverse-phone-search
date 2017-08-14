package main

import (
	"os"
	"testing"
)

func TestModels(t *testing.T) {
	app := &App{}
	tableNames := [4]string{"phone_number", "person", "address"}
	dbName := "./unittest.db"

	app.Initialize(dbName)
	defer os.Remove(dbName)
	defer app.DB.Close()

	// test tables created
	CreateDBTables(app.DB)
	for r := range tableNames {
		query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
		app.DB.Exec(query, r)
	}

	// test data
	number := "4112223333"
	name := "John Doe"

	person := &Person{}
	person.FullName = name
	person.Address.State = "DE"
	person.Address.Zip = "19802"
	person.Address.Street = "123 Main St."
	person.Address.City = "Wilmington"
	person.Phone.Number = number
	person.save(app.DB)

	pn := &PhoneNumber{}
	pn.Number = number
	getPersonFromDb(pn, app.DB)

	// test that a person was added and we can retrieve them from the db
	if pn.Matches[0].FullName != person.FullName {
		t.Error("Person name not found")
	}

	if pn.Matches[0].Phone.Number != person.Phone.Number {
		t.Error("Person number not found")
	}

	if pn.Matches[0].Address.Street != person.Address.Street {
		t.Error("Person address not found")
	}
}
