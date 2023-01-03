package model

import (
	"database/sql"
	"errors"
	"fmt"
	"goRestAPI/database"
	"log"
	"strconv"
	"time"
)

type CreditSchema struct {
	ID             uint32    `json:"id"`
	Card           int       `json:"card"`
	TotalPrice     uint32    `json:"totalPrice,omitempty"`
	FeeAmount      uint32    `json:"feeAmount"`
	ExpiredAmount  uint32    `json:"expiredAmount,omitempty"`
	Fees           uint8     `json:"fees"`
	CurrentFee     uint8     `json:"currentFee"`
	CurrentFeePaid bool      `json:"currentFeePaid"`
	PurchaseDate   string    `json:"purchaseDate"`
	Completed      bool      `json:"completed"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
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
		query = database.DueCreditThisMonth
	}

	rows, err = db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var credits []CreditSchema

	for rows.Next() {
		if err = rows.Scan(&c.ID, &c.Card, &c.TotalPrice,
			&c.FeeAmount, &c.ExpiredAmount, &c.Fees,
			&c.CurrentFee, &c.CurrentFeePaid, &c.PurchaseDate,
			&c.Completed, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		credits = append(credits, *c)
	}

	return credits, nil
}

func (c *CreditSchema) QInsertCredit(db *sql.DB) error {
	if c.Fees <= 0 || c.CurrentFee <= 0 {
		return nil
	}
	c.ExpiredAmount = 0
	c.Completed = false
	c.CreatedAt = time.Now()
	var closeDay int
	rows, err := db.Query("SELECT closeDay FROM card WHERE id=$1", c.Card)
	if err != nil {
		return err
	}

	var purchaseDay int
	purchaseDay, err = strconv.Atoi(c.PurchaseDate[len(c.PurchaseDate)-2:])
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
		purchaseDay >= closeDay {
		c.CurrentFee = 0
	}
	err = db.QueryRow(database.CreditNewRecord,
		&c.Card,
		&c.TotalPrice,
		&c.FeeAmount,
		&c.ExpiredAmount,
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

func (c *CreditSchema) QDeleteCredit(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM credit WHERE id=$1", &c.ID)
	return err
}

func (c *CreditSchema) QPayCredit(db *sql.DB) error {
	var err error
	err = db.QueryRow(database.GetCreditToPay,
		&c.ID).Scan(&c.Fees, &c.CurrentFee, &c.CurrentFeePaid, &c.Completed)
	if err != nil {
		return err
	}
	if c.Completed == true {
		return errors.New("payment_already_completed")
	}
	if c.CurrentFeePaid == true {
		return errors.New("quota_already_fulfilled")
	}

	c.CurrentFeePaid = true
	if c.CurrentFee == c.Fees {
		_, err = db.Exec(database.CreditCompletePayment,
			&c.ID)
		if err != nil {
			return err
		}
	}
	_, err = db.Exec(database.PayCredit,
		&c.ID)
	if err != nil {
		return err
	}
	return nil
}

func QCreditNextQuotaAll(db *sql.DB) error {
	credits, err := QGetAllCredits(db)
	for i := 0; i < len(credits); i++ {
		if credits[i].CurrentFeePaid == false {
			credits[i].ExpiredAmount += credits[i].FeeAmount
			_, err = db.Exec(database.AddToExpired, credits[i].ExpiredAmount, credits[i].ID)
		}
	}
	_, err = db.Exec(database.NextQuota)
	if err != nil {
		return err
	}
	return nil
}

func QGetAllCredits(db *sql.DB) ([]CreditSchema, error) {
	rows, err := db.Query(database.AllCredits)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var credits []CreditSchema

	for rows.Next() {
		var c CreditSchema
		if err = rows.Scan(&c.ID, &c.Card, &c.TotalPrice,
			&c.FeeAmount, &c.ExpiredAmount, &c.Fees,
			&c.CurrentFee, &c.CurrentFeePaid, &c.PurchaseDate,
			&c.Completed, &c.CreatedAt); err != nil {
			return nil, err
		}
		credits = append(credits, c)
	}

	return credits, nil
}

func QThisMonthDebtCredit(db *sql.DB) (uint32, error) {
	var debt uint32 = 0
	var c CreditSchema
	rows, err := db.Query(database.CreditDebtThisMonth)
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
		err = rows.Scan(&c.FeeAmount, &c.ExpiredAmount)
		if err != nil {
			return 0, err
		}
		debt += c.FeeAmount + c.ExpiredAmount
	}
	return debt, err
}

func QDisplayDueCredit(db *sql.DB) ([]CreditDue, error) {
	var creditsDue = []CreditDue{{0}}
	rows, err := db.Query(database.DueCreditAllTime)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}(rows)

	var c CreditSchema
	var cardId int
	var cardType string
	var closeDay int

	for rows.Next() {
		err = rows.Scan(&c.ID, &c.Card, &c.TotalPrice,
			&c.FeeAmount, &c.ExpiredAmount, &c.Fees,
			&c.CurrentFee, &c.CurrentFeePaid, &c.PurchaseDate,
			&c.Completed, &c.CreatedAt, &cardId,
			&cardType, &closeDay)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		if c.ExpiredAmount > 0 {
			creditsDue[0].Amount += c.ExpiredAmount
		}

		plus := 0
		if c.CurrentFee == 0 {
			plus = 1
		} else {
			if c.CurrentFeePaid == false {
				creditsDue[0].Amount += c.FeeAmount
			}
		}
		for i := 1; i <= int(c.Fees-c.CurrentFee)+plus; i++ {
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
	if err != nil {
		return err
	}
	_, err = db.Exec("ALTER SEQUENCE credit_id_seq RESTART WITH 1")
	if err != nil {
		return err
	}
	return nil
}
