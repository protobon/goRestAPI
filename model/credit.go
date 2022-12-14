package model

import (
	"awesomeProject/common"
	"database/sql"
	"errors"
	"fmt"
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
	Completed      bool      `json:"completed,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

func (c *CreditSchema) QGetCredit(db *sql.DB) ([]CreditSchema, error) {
	var rows *sql.Rows
	var err error
	fmt.Printf("ID: %d", c.ID)
	var query = fmt.Sprintf("SELECT * FROM credit WHERE id=%d", c.ID)
	if c.ID == 0 {
		query = "SELECT * FROM credit WHERE currentFeePaid=false AND completed=false"
	}

	rows, err = db.Query(query)
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
	c.Completed = false
	c.CreatedAt = time.Now()
	err := db.QueryRow(common.CreditNewRecord,
		&c.TotalPrice,
		&c.FeeAmount,
		&c.Fees,
		&c.CurrentFee,
		&c.CurrentFeePaid,
		&c.PurchaseDate,
		false,
		&c.CreatedAt).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func (c *CreditSchema) QPayCredit(db *sql.DB) error {
	var err error
	err = db.QueryRow(`SELECT * FROM credit WHERE id=$1`,
		&c.ID).Scan(&c.ID, &c.TotalPrice, &c.FeeAmount, &c.Fees,
		&c.CurrentFee, &c.CurrentFeePaid, &c.PurchaseDate,
		&c.Completed, &c.CreatedAt)
	if err != nil {
		return err
	}
	if c.Completed == true {
		return errors.New("full_payment_already_fulfilled")
	}
	if c.CurrentFeePaid == true {
		return errors.New("quota_already_fulfilled")
	}

	c.CurrentFeePaid = true
	if c.CurrentFee == c.Fees {
		_, err =
			db.Exec("UPDATE credit SET currentFeePaid=true, completed=true WHERE id=$1",
				&c.ID)
		return err
	}
	_, err =
		db.Exec("UPDATE credit SET currentFeePaid=true WHERE id=$1",
			&c.ID)
	return err
}

func (c *CreditSchema) QNextQuota(db *sql.DB) error {
	_, err :=
		db.Exec(`UPDATE credit SET currentFee=currentFee+1, currentFeePaid=false
              WHERE completed!=true AND currentFee<fees`)
	return err
}

func QGetAllCredits(db *sql.DB, start int, count int) ([]CreditSchema, error) {
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

func QCalcDebtCredit(db *sql.DB) (uint32, error) {
	var debt uint32 = 0
	var c CreditSchema
	rows, err := db.Query("SELECT * FROM credit WHERE currentFeePaid!=true")
	if err != nil {
		return debt, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	for rows.Next() {
		err := rows.Scan(&c.ID, &c.TotalPrice,
			&c.FeeAmount, &c.Fees, &c.CurrentFee,
			&c.CurrentFeePaid, &c.PurchaseDate,
			&c.Completed, &c.CreatedAt)
		if err != nil {
			return 0, err
		}
		debt += c.FeeAmount
	}
	return debt, err
}

func QClearTableCredit(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE credit")
	return err
}
