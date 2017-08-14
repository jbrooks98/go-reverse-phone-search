package main

func main() {
	a := app{}
	a.initialize("./phone-search.db")
	a.run("127.0.0.1:8000")
}
