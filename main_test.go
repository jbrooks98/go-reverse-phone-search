package main

//
//import (
//	"os"
//	"testing"
//
//	"log"
//)
//
//var a App
//
//func TestMain(m *testing.M) {
//	a = App{}
//	a.Initialize(os.Getenv("APP_DB_NAME"))
//
//	ensureTableExists()
//
//	code := m.Run()
//
//	clearTable()
//
//	os.Exit(code)
//}
//
//func ensureTableExists() {
//	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
//		log.Fatal(err)
//	}
//}
//
//func clearTable() {
//	a.DB.Exec("DELETE FROM phone_number")
//	a.DB.Exec("ALTER SEQUENCE number_id_seq RESTART WITH 1")
//}
//
//const tableCreationQuery = `
//	CREATE TABLE IF NOT EXISTS phone_number(
//		id INTEGER PRIMARY KEY,
//		number varchar(15) NOT NULL,
//	);
//	`
