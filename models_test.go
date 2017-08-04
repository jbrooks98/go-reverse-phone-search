package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestModels(t *testing.T) {
	app := &App{}
	tableNames := [4]string{"phone_number", "person", "address"}
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dbPath := dir + "testdb.db"
	app.Initialize(dbPath)
	defer os.Remove(dbPath)
	defer app.DB.Close()

	// test tables created
	CreateDBTables(app.DB)
	for r := range tableNames {
		query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
		app.DB.Exec(query, r)
	}

	// add a person
	person := &Person{}
	person.FullName = "John Doe"
	person.Address.State = "DE"
	person.Address.Zip = "19802"
	person.Address.Street = "123 Main St."
	person.Address.City = "Wilmington"
	person.Phone.Number = "1112223333"
	person.Save(app.DB)

	createdPerson, err := GetPeopleByNumber(person.Phone.Number, app.DB)
	if err != nil {
		t.Error(err)
	}
	// test that a person was added and we can retrieve them from the db
	if createdPerson[0].FullName != person.FullName {
		t.Error("Person name not found")
	}

	if createdPerson[0].Phone.Number != person.Phone.Number {
		t.Error("Person number not found")
	}

	if createdPerson[0].Address.Street != person.Address.Street {
		t.Error("Person address not found")
	}
}
