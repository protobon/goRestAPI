package main

import (
	"goRestAPI/api"
	"os"
)

func main() {
	app := api.App{}

	app.Initialize(
		os.Getenv("TEST_POSTGRES_USER"),
		os.Getenv("TEST_POSTGRES_PASSWORD"),
		os.Getenv("TEST_POSTGRES_DB"),
	)

	app.Run(":8080")
}
