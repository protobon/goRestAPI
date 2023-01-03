package model

import (
	"database/sql"
	"log"
)

type DummySchema struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (d *DummySchema) QGetDummy(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM dummy WHERE id=$1",
		d.ID).Scan(&d.Name, &d.Price)
}

func (d *DummySchema) QUpdateDummy(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE dummy SET name=$1, price=$2 WHERE id=$3",
			d.Name, d.Price, d.ID)
	return err
}

func (d *DummySchema) QDeleteDummy(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM dummy WHERE id=$1", d.ID)
	return err
}

func (d *DummySchema) QCreateDummy(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO dummy(name, price) VALUES($1, $2) RETURNING id",
		d.Name, d.Price).Scan(&d.ID)

	if err != nil {
		return err
	}

	return nil
}

func (d *DummySchema) QGetDummies(db *sql.DB, start int, count int) ([]DummySchema, error) {
	rows, err := db.Query(
		"SELECT id, name, price FROM dummy LIMIT $1 OFFSET $2",
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
		if err = rows.Scan(&dummy.ID, &dummy.Name, &dummy.Price); err != nil {
			return nil, err
		}
		dummies = append(dummies, dummy)
	}

	return dummies, nil
}
