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

	CreateTableProduct(db)
	CreateTableCard(db)
	CreateTableCredit(db)
	CreateTableServiceFixed(db)
	CreateTableServiceVariable(db)
	return db
}

func CreateTableProduct(db *sql.DB) {
	if _, err := db.Exec(ProductTableCreate); err != nil {
		log.Fatal(err)
	}
}

func CreateTableCard(db *sql.DB) {
	if _, err := db.Exec(CardTableCreate); err != nil {
		log.Fatal(err)
	}
}

func CreateTableCredit(db *sql.DB) {
	if _, err := db.Exec(CreditTableCreate); err != nil {
		log.Fatal(err)
	}
}

func CreateTableServiceFixed(db *sql.DB) {
	if _, err := db.Exec(FixedServiceTableCreate); err != nil {
		log.Fatal(err)
	}
}

func CreateTableServiceVariable(db *sql.DB) {
	if _, err := db.Exec(VariableServiceTableCreate); err != nil {
		log.Fatal(err)
	}
}
