package main

import (
	"os"
)

func main() {
	a := App{}
	a.Initialize(os.Getenv("APP_DB_NAME"))
	a.Run("127.0.0.1:8000")
}
