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

	CreateTableDummy(db)
	CreateTableCreditCard(db)
	CreateTableCredit(db)
	CreateTableRent(db)
	CreateTableBill(db)
	CreateTableDebitCard(db)
	CreateTableDebit(db)
	return db
}

func CreateTableDummy(db *sql.DB) {
	if _, err := db.Exec(DummyTableCreate); err != nil {
		log.Fatal(err)
	}
}

func CreateTableCreditCard(db *sql.DB) {
	if _, err := db.Exec(CreditCardTableCreate); err != nil {
		log.Fatal(err)
	}
}

func CreateTableCredit(db *sql.DB) {
	if _, err := db.Exec(CreditTableCreate); err != nil {
		log.Fatal(err)
	}
}

func CreateTableRent(db *sql.DB) {
	if _, err := db.Exec(RentTableCreate); err != nil {
		log.Fatal(err)
	}
}

func CreateTableBill(db *sql.DB) {
	if _, err := db.Exec(BillTableCreate); err != nil {
		log.Fatal(err)
	}
}

func CreateTableDebitCard(db *sql.DB) {
	if _, err := db.Exec(DebitCardTableCreate); err != nil {
		log.Fatal(err)
	}
}

func CreateTableDebit(db *sql.DB) {
	if _, err := db.Exec(DebitTableCreate); err != nil {
		log.Fatal(err)
	}
}
