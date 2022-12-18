package model

import (
	"database/sql"
	"log"
)

type VariableServiceSchema struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (s *VariableServiceSchema) GetService(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM service_variable WHERE id=$1",
		s.ID).Scan(&s.Name, &s.Price)
}

func (s *VariableServiceSchema) QUpdateService(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE service_variable SET name=$1, price=$2 WHERE id=$3",
			s.Name, s.Price, s.ID)
	return err
}

func (s *VariableServiceSchema) QDeleteService(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM service_variable WHERE id=$1", s.ID)
	return err
}

func (s *VariableServiceSchema) QCreateService(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO service_variable(name, price) VALUES($1, $2) RETURNING id",
		s.Name, s.Price).Scan(&s.ID)

	if err != nil {
		return err
	}

	return nil
}

func QGetVariableServices(db *sql.DB, start int, count int) ([]VariableServiceSchema, error) {
	rows, err := db.Query(
		"SELECT id, name,  price FROM service_variable LIMIT $1 OFFSET $2",
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

	var services []VariableServiceSchema

	for rows.Next() {
		var s VariableServiceSchema
		if err = rows.Scan(&s.ID, &s.Name, &s.Price); err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	return services, nil
}
