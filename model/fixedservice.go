package model

import (
	"database/sql"
	"log"
)

type FixedServiceSchema struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (s *FixedServiceSchema) GetService(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM service_fixed WHERE id=$1",
		s.ID).Scan(&s.Name, &s.Price)
}

func (s *FixedServiceSchema) QUpdateService(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE service_fixed SET name=$1, price=$2 WHERE id=$3",
			s.Name, s.Price, s.ID)
	return err
}

func (s *FixedServiceSchema) QDeleteService(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM service_fixed WHERE id=$1", s.ID)
	return err
}

func (s *FixedServiceSchema) QCreateService(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO service_fixed(name, price) VALUES($1, $2) RETURNING id",
		s.Name, s.Price).Scan(&s.ID)

	if err != nil {
		return err
	}

	return nil
}

func QGetFixedServices(db *sql.DB, start int, count int) ([]FixedServiceSchema, error) {
	rows, err := db.Query(
		"SELECT id, name,  price FROM service_fixed LIMIT $1 OFFSET $2",
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

	var services []FixedServiceSchema

	for rows.Next() {
		var s FixedServiceSchema
		if err = rows.Scan(&s.ID, &s.Name, &s.Price); err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	return services, nil
}
