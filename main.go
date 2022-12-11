package main

import (
	"awesomeProject/app"
	"fmt"
	"os"
)

func main() {
	a := app.App{}

	fmt.Println("user startup: ", os.Getenv("APP_DB_USERNAME"))
	fmt.Println("password startup: ", os.Getenv("APP_DB_PASSWORD"))
	fmt.Println("dbname startup: ", os.Getenv("APP_DB_NAME"))

	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	a.Run(":8080")
}
