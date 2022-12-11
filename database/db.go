package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func DBInit(user string, password string, dbname string) *sql.DB {
	connectionString :=
		fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
			user,
			password,
			"localhost",
			5432,
			dbname)

	var err error
	var db *sql.DB
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	CreateTableProducts(db)
	return db
}

func CreateTableProducts(db *sql.DB) {
	if _, err := db.Exec(productsTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

const productsTableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
   id SERIAL,
   name TEXT NOT NULL,
   price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
   CONSTRAINT products_pkey PRIMARY KEY (id)
)`
