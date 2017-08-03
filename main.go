package go_reverse_phone_search


import (
	"os"
)

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_NAME"))
	a.Run(":8080")
}

