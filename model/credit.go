package model

import (
	"awesomeProject/database"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

type CreditSchema struct {
	ID             uint32    `json:"id,omitempty"`
	Card           int       `json:"card"`
	TotalPrice     uint32    `json:"totalPrice,omitempty"`
	FeeAmount      uint32    `json:"feeAmount"`
	Fees           uint8     `json:"fees"`
	CurrentFee     uint8     `json:"currentFee"`
	CurrentFeePaid bool      `json:"currentFeePaid"`
	PurchaseDate   string    `json:"purchaseDate"`
	Completed      bool      `json:"completed,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

type CreditDue struct {
	Amount uint32 `json:"amount"`
}

func (c *CreditSchema) QGetCredit(db *sql.DB) ([]CreditSchema, error) {
	var rows *sql.Rows
	var err error
	fmt.Printf("ID: %d", c.ID)
	var query = fmt.Sprintf("SELECT * FROM credit WHERE id=%d", c.ID)
	if c.ID == 0 {
		query = `SELECT * FROM credit WHERE currentFeePaid=false
                       AND completed=false AND currentFee!=0`
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
		if err = rows.Scan(&c.ID, &c.Card, &c.TotalPrice,
			&c.FeeAmount, &c.Fees, &c.CurrentFee,
			&c.CurrentFeePaid, &c.PurchaseDate,
			&c.Completed, &c.CreatedAt); err != nil {
			return nil, err
		}
		credits = append(credits, *c)
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

func (c *CreditSchema) QInsertCredit(db *sql.DB) error {
	c.Completed = false
	c.CreatedAt = time.Now()
	var closeDay int
	rows, err := db.Query("SELECT closeDay FROM card WHERE id=$1", c.Card)
	if err != nil {
		return err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)
	for rows.Next() {
		err = rows.Scan(&closeDay)
		if err != nil {
			return err
		}
	}
	if c.CurrentFee == 1 && c.CurrentFeePaid == false &&
		len(c.PurchaseDate)-2 >= closeDay {
		c.CurrentFee = 0
	}
	err = db.QueryRow(database.CreditNewRecord,
		&c.Card,
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
		`SELECT * FROM credit LIMIT $1 OFFSET $2`,
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
		if err := rows.Scan(&c.ID, &c.Card, &c.TotalPrice,
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
	rows, err := db.Query("SELECT feeAmount FROM credit WHERE currentFeePaid!=true AND currentFee!=0")
	if err != nil {
		return debt, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	for rows.Next() {
		err = rows.Scan(&c.FeeAmount)
		if err != nil {
			return 0, err
		}
		debt += c.FeeAmount
	}
	return debt, err
}

func QDisplayDueCredit(db *sql.DB) ([]CreditDue, error) {
	var creditsDue []CreditDue
	rows, err := db.Query(`SELECT * FROM credit JOIN card ON credit.card = card.id 
         WHERE credit.completed=false`)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}(rows)

	var cardId int
	var cardType string
	var closeDay int

	for rows.Next() {
		var c CreditSchema
		err := rows.Scan(&c.ID, &c.Card, &c.TotalPrice,
			&c.FeeAmount, &c.Fees, &c.CurrentFee,
			&c.CurrentFeePaid, &c.PurchaseDate,
			&c.Completed, &c.CreatedAt, &cardId,
			&cardType, &closeDay)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		plus := 0
		if c.CurrentFeePaid == false && (c.CurrentFee != 0 ||
			closeDay > int(c.PurchaseDate[len(c.PurchaseDate)-2])) {
			plus = 1
		}
		for i := 0; i < int(c.Fees-c.CurrentFee)+plus; i++ {
			if i > len(creditsDue)-1 {
				cd := CreditDue{c.FeeAmount}
				creditsDue = append(creditsDue, cd)
			} else {
				creditsDue[i].Amount += c.FeeAmount
			}
		}
	}
	return creditsDue, nil
}

func QClearTableCredit(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE credit")
	return err
}
