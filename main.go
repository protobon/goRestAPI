package main

import (
	"goRestAPI/api"
	"os"
)

func main() {
	app := api.App{}

	app.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	app.Run(":8080")
}
