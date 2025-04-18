package db

import (
	"fmt"
	"github.com/grigory222/avito-backend-trainee/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func CreateConnection(config config.DB) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Host, config.User, config.Password, config.Name, config.Port)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
