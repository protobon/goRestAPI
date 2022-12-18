package model

import (
	"database/sql"
	"log"
)

type ProductSchema struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (p *ProductSchema) QGetProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM product WHERE id=$1",
		p.ID).Scan(&p.Name, &p.Price)
}

func (p *ProductSchema) QUpdateProduct(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE product SET name=$1, price=$2 WHERE id=$3",
			p.Name, p.Price, p.ID)
	return err
}

func (p *ProductSchema) QDeleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM product WHERE id=$1", p.ID)
	return err
}

func (p *ProductSchema) QCreateProduct(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO product(name, price) VALUES($1, $2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func QGetProducts(db *sql.DB, start int, count int) ([]ProductSchema, error) {
	rows, err := db.Query(
		"SELECT id, name,  price FROM product LIMIT $1 OFFSET $2",
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

	var products []ProductSchema

	for rows.Next() {
		var p ProductSchema
		if err = rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
