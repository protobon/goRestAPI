package model

import (
	"database/sql"
	"log"
	"time"
)

type DummySchema struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *DummySchema) QGetDummy(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM dummy WHERE id=$1",
		d.ID).Scan(&d.ID, &d.Name, &d.Price, &d.CreatedAt, &d.UpdatedAt)
}

func (d *DummySchema) QUpdateDummy(db *sql.DB) error {
	d.UpdatedAt = time.Now()
	_, err :=
		db.Exec("UPDATE dummy SET name=$1, price=$2, updated_at=$3 WHERE id=$4",
			d.Name, d.Price, d.UpdatedAt, d.ID)
	return err
}

func (d *DummySchema) QDeleteDummy(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM dummy WHERE id=$1", d.ID)
	return err
}

func (d *DummySchema) QCreateDummy(db *sql.DB) error {
	d.CreatedAt = time.Now()
	d.UpdatedAt = d.CreatedAt
	err := db.QueryRow(
		"INSERT INTO dummy(name, price, created_at, updated_at) VALUES($1, $2, $3, $4) RETURNING id",
		d.Name, d.Price, d.CreatedAt, d.UpdatedAt).Scan(&d.ID)

	if err != nil {
		return err
	}

	return nil
}

func (d *DummySchema) QGetDummies(db *sql.DB, start int, count int) ([]DummySchema, error) {
	rows, err := db.Query(
		"SELECT * FROM dummy LIMIT $1 OFFSET $2",
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

	var dummies []DummySchema

	for rows.Next() {
		var dummy DummySchema
		if err = rows.Scan(&d.ID, &d.Name, &d.Price,
			&d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		dummies = append(dummies, dummy)
	}

	return dummies, nil
}
