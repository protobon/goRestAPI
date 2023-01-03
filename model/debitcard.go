package model

import (
	"database/sql"
	"log"
	"time"
)

type DebitCardSchema struct {
	ID        int       `json:"id"`
	Salary    int       `json:"salary"`
	Currency  string    `json:"currency"`
	Renew     bool      `json:"renew"`
	Money     int       `json:"money"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func (dc *DebitCardSchema) QGetCard(db *sql.DB) error {
	return db.QueryRow(`SELECT * FROM debit_card WHERE id=$1`,
		dc.ID).Scan(&dc.Salary, &dc.Currency,
		&dc.Renew, &dc.Money, &dc.CreatedAt, &dc.UpdatedAt)
}

func (dc *DebitCardSchema) QUpdateCard(db *sql.DB) error {
	dc.UpdatedAt = time.Now()
	_, err :=
		db.Exec(`UPDATE debit_card SET salary=$1, currency=$2,
                      renew=$3, money=$4, updatedAt=$5 WHERE id=$6`,
			dc.Salary, dc.Currency, dc.Renew, dc.Money, dc.UpdatedAt, dc.ID)
	return err
}

func (dc *DebitCardSchema) QDeleteCard(db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM debit_card WHERE id=$1`, dc.ID)
	return err
}

func (dc *DebitCardSchema) QCreateCard(db *sql.DB) error {
	dc.CreatedAt = time.Now()
	err := db.QueryRow(
		`INSERT INTO debit_card(salary, currency,
                       renew, money, createdAt) VALUES($1, $2, $3, $4, $5) RETURNING id`,
		dc.Salary, dc.Currency, dc.Renew, dc.Money, dc.CreatedAt).Scan(&dc.ID)

	if err != nil {
		return err
	}

	return nil
}

func (dc *DebitCardSchema) QGetCards(db *sql.DB, start int, count int) ([]DebitCardSchema, error) {
	rows, err := db.Query(
		"SELECT * FROM debit_card LIMIT $1 OFFSET $2",
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

	var cards []DebitCardSchema

	for rows.Next() {
		var card DebitCardSchema
		if err = rows.Scan(&card.ID, &card.Salary,
			&card.Currency, &card.Renew, &card.Money,
			&card.CreatedAt, &card.UpdatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}
