package model

import (
	"database/sql"
	"log"
	"time"
)

type DebitSchema struct {
	ID        int       `json:"id"`
	Total     int       `json:"total"`
	Card      int       `json:"card"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func (d *DebitSchema) QGetDebit(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM debit WHERE id=$1",
		d.ID).Scan(&d.Total, &d.Card, &d.CreatedAt, &d.UpdatedAt)
}

func (d *DebitSchema) QDeleteDebit(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM debit WHERE id=$1", d.ID)
	return err
}

func (d *DebitSchema) QCreateDebit(db *sql.DB) error {
	d.CreatedAt = time.Now()
	err := db.QueryRow(
		`INSERT INTO debit(total, card, createdAt) VALUES($1, $2, $3) RETURNING id`,
		d.Total, d.Card, d.CreatedAt).Scan(&d.ID)

	if err != nil {
		return err
	}

	return nil
}

func (d *DebitSchema) QGetDebits(db *sql.DB, start int, count int) ([]DebitSchema, error) {
	rows, err := db.Query(
		"SELECT * FROM debit LIMIT $1 OFFSET $2",
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

	var debits []DebitSchema

	for rows.Next() {
		var debit DebitSchema
		if err = rows.Scan(&debit.ID, &debit.Total,
			&debit.Card, &debit.CreatedAt,
			&debit.UpdatedAt); err != nil {
			return nil, err
		}
		debits = append(debits, debit)
	}

	return debits, nil
}
