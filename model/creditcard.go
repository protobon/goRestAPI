package model

import (
	"database/sql"
	"log"
	"time"
)

type CreditCardSchema struct {
	ID          int       `json:"id"`
	Brand       string    `json:"brand"`
	CloseDay    int       `json:"closeDay"`
	AmountLimit int       `json:"amountLimit"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

func (cc *CreditCardSchema) QGetCard(db *sql.DB) error {
	return db.QueryRow(`SELECT * FROM credit_card WHERE id=$1`,
		cc.ID).Scan(&cc.Brand, &cc.CloseDay, &cc.AmountLimit,
		&cc.CreatedAt, &cc.UpdatedAt)
}

func (cc *CreditCardSchema) QUpdateLimit(db *sql.DB) error {
	cc.UpdatedAt = time.Now()
	_, err :=
		db.Exec("UPDATE credit_card SET amountLimit=$1, updatedAt=$2 WHERE id=$3",
			cc.AmountLimit, cc.UpdatedAt, cc.ID)
	return err
}

func (cc *CreditCardSchema) DeleteCard(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM credit_card WHERE id=$1", cc.ID)
	return err
}

func (cc *CreditCardSchema) CreateCard(db *sql.DB) error {
	cc.CreatedAt = time.Now()
	err := db.QueryRow(
		`INSERT INTO credit_card(brand,
                        closeDay, amountLimit,
                        createdAt) VALUES($1, $2, $3, $4) RETURNING id`,
		cc.Brand, cc.CloseDay, cc.AmountLimit, cc.CreatedAt).Scan(&cc.ID)

	if err != nil {
		return err
	}

	return nil
}

func (cc *CreditCardSchema) GetCards(db *sql.DB, start int, count int) ([]CreditCardSchema, error) {
	rows, err := db.Query(
		`SELECT * FROM credit_card LIMIT $1 OFFSET $2`,
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

	var cards []CreditCardSchema

	for rows.Next() {
		var card CreditCardSchema
		if err = rows.Scan(&card.ID,
			&card.Brand, &card.CloseDay,
			&card.AmountLimit, &card.CreatedAt,
			&card.UpdatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}
