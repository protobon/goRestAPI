package model

import (
	"database/sql"
	"time"
)

type RentSchema struct {
	ID        int       `json:"id,omitempty"`
	Amount    int       `json:"amount"`
	DueDate   time.Time `json:"dueDate"`
	Active    bool      `json:"active,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

func (r *RentSchema) QGetRent(db *sql.DB) error {
	return db.QueryRow(`SELECT * FROM rent WHERE id=$1`,
		r.ID).Scan(&r.Amount, &r.DueDate, &r.Active, &r.CreatedAt, &r.UpdatedAt)
}

func (r *RentSchema) QUpdateRent(db *sql.DB) error {
	r.UpdatedAt = time.Now()
	_, err :=
		db.Exec(`UPDATE rent SET amount=$1, dueDate=$2, updatedAt=$3 WHERE id=$4`,
			r.Amount, r.DueDate, r.UpdatedAt, r.ID)
	return err
}

func (r *RentSchema) QDeleteRent(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM rent WHERE id=$1", r.ID)
	return err
}

func (r *RentSchema) QCreateRent(db *sql.DB) error {
	r.CreatedAt = time.Now()
	r.Active = true
	err := db.QueryRow(
		`INSERT INTO rent(amount, dueDate,
                 active, createdAt) VALUES($1, $2, $3, $4) RETURNING id`,
		r.Amount, r.DueDate, r.Active, r.CreatedAt).Scan(&r.ID)

	if err != nil {
		return err
	}

	return nil
}
