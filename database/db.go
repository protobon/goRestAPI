package database

import (
	"awesomeProject/common"
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

	CreateTableProduct(db)
	CreateTableCredit(db)
	return db
}

func CreateTableProduct(db *sql.DB) {
	if _, err := db.Exec(common.ProductTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func CreateTableCredit(db *sql.DB) {
	if _, err := db.Exec(common.CreditTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}
