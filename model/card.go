package model

import (
	"database/sql"
	"log"
)

type CardSchema struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	CloseDay int    `json:"closeDay"`
}

func (c *CardSchema) GetCard(db *sql.DB) error {
	return db.QueryRow("SELECT type, closeDay FROM card WHERE id=$1",
		c.ID).Scan(&c.Type, &c.CloseDay)
}

func (c *CardSchema) UpdateCard(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE card SET type=$1, closeDay=$2 WHERE id=$3",
			c.Type, c.CloseDay, c.ID)
	return err
}

func (c *CardSchema) DeleteCard(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM card WHERE id=$1", c.ID)
	return err
}

func (c *CardSchema) CreateCard(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO card(type, closeDay) VALUES($1, $2) RETURNING id",
		c.Type, c.CloseDay).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetCards(db *sql.DB, start int, count int) ([]CardSchema, error) {
	rows, err := db.Query(
		"SELECT id, type,  closeDay FROM card LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var cards []CardSchema

	for rows.Next() {
		var c CardSchema
		if err := rows.Scan(&c.ID, &c.Type, &c.CloseDay); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}

	return cards, nil
}
