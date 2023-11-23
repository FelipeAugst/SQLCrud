package bank

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {

	stringConnection := "felipe:felipe/datacenter?charset=utf8&parseTime=True"
	db, err := sql.Open("mysql", stringConnection)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}
