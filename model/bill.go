package model

import (
	"database/sql"
	"log"
	"time"
)

type BillSchema struct {
	ID        int       `json:"id"`
	Amount    int       `json:"amount"`
	DueDate   time.Time `json:"due_date"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *BillSchema) QGetBill(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM bill WHERE id=$1",
		b.ID).Scan(&b.Amount, &b.DueDate, &b.Active, &b.CreatedAt, &b.UpdatedAt)
}

func (b *BillSchema) QUpdateBill(db *sql.DB) error {
	b.UpdatedAt = time.Now()
	_, err :=
		db.Exec(`UPDATE bill SET active=$1, updatedAt=$2 WHERE id=$3`,
			b.Active, b.UpdatedAt, b.ID)
	return err
}

func (b *BillSchema) QDeleteBill(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM bill WHERE id=$1", b.ID)
	return err
}

func (b *BillSchema) QCreateBill(db *sql.DB) error {
	b.CreatedAt = time.Now()
	err := db.QueryRow(
		`INSERT INTO bill(amount, dueDate, active, createdAt)
VALUES($1, $2, $3, $4) RETURNING id`,
		b.Amount, b.DueDate, b.Active, b.CreatedAt).Scan(&b.ID)

	if err != nil {
		return err
	}

	return nil
}

func (b *BillSchema) QGetBills(db *sql.DB, start int, count int) ([]BillSchema, error) {
	rows, err := db.Query(
		"SELECT * FROM bill LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var bills []BillSchema

	for rows.Next() {
		var bill BillSchema
		if err = rows.Scan(&bill.ID, &bill.Amount,
			&bill.DueDate, &bill.Active,
			&bill.CreatedAt, &bill.UpdatedAt); err != nil {
			return nil, err
		}
		bills = append(bills, bill)
	}

	return bills, nil
}
