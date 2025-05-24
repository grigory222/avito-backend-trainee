package models

import "database/sql"

type PVZ struct {
	Id               string `db:"id"`
	RegistrationDate string `db:"registration_date"`
	City             string `db:"city"`
}

type PVZWithReceptions struct {
	PVZ        PVZ
	Receptions []ReceptionWithProducts
}

type FlatRow struct {
	PVZId            string
	City             string
	RegistrationDate string
	ReceptionId      string
	ReceptionDate    string
	Status           string
	ProductId        sql.NullString
	ProductDate      sql.NullString
	ProductType      sql.NullString
}
