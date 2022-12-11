package model

import (
	"awesomeProject/common"
	"database/sql"
	"log"
	"time"
)

type CreditSchema struct {
	ID             uint32    `json:"id"`
	TotalPrice     float64   `json:"totalPrice"`
	FeeAmount      uint32    `json:"feeAmount"`
	Fees           uint8     `json:"fees"`
	CurrentFee     uint8     `json:"currentFee"`
	CurrentFeePaid bool      `json:"currentFeePaid"`
	PurchaseDate   string    `json:"purchaseDate"`
	Completed      bool      `json:"completed"`
	CreatedAt      time.Time `json:"createdAt"`
}

func (c *CreditSchema) QGetCredit(db *sql.DB) error {
	return db.QueryRow("SELECT totalPrice FROM credit WHERE id=$1",
		&c.ID).Scan(&c.TotalPrice)
}

func (c *CreditSchema) QUpdateCredit(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE credit SET currentFeePaid=$1 WHERE id=$2",
			&c.CurrentFeePaid, &c.ID)
	return err
}

func (c *CreditSchema) QDeleteCredit(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM credit WHERE id=$1", &c.ID)
	return err
}

func (c *CreditSchema) QCreateCredit(db *sql.DB) error {
	err := db.QueryRow(common.CreditNewRecord,
		&c.TotalPrice,
		&c.FeeAmount,
		&c.Fees,
		&c.CurrentFee,
		&c.CurrentFeePaid,
		&c.PurchaseDate,
		&c.Completed,
		time.Now()).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func QGetCredits(db *sql.DB, start int, count int) ([]CreditSchema, error) {
	rows, err := db.Query(
		`SELECT id, totalPrice, feeAmount, fees, currentFee, 
       currentFeePaid, purchaseDate, completed, createdAt FROM credit LIMIT $1 OFFSET $2`,
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

	var credits []CreditSchema

	for rows.Next() {
		var c CreditSchema
		if err := rows.Scan(&c.ID, &c.TotalPrice,
			&c.FeeAmount, &c.Fees, &c.CurrentFee,
			&c.CurrentFeePaid, &c.PurchaseDate,
			&c.Completed, &c.CreatedAt); err != nil {
			return nil, err
		}
		credits = append(credits, c)
	}

	return credits, nil
}
