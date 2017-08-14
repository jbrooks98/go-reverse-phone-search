package main

func main() {
	a := App{}
	dbName := "./test.db"
	a.initialize(dbName)
	a.run("127.0.0.1:8000")
}
