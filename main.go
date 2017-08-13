package main

func main() {
	a := App{}
	dbName := "./test.db"
	a.Initialize(dbName)
	a.Run("127.0.0.1:8000")
}
