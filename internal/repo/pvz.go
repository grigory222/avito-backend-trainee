package repo

import (
	"github.com/grigory222/avito-backend-trainee/internal/models"
	"github.com/jmoiron/sqlx"
)

type PVZRepository struct {
	db *sqlx.DB
}

func NewPVZRepository(db *sqlx.DB) *PVZRepository {
	return &PVZRepository{
		db: db,
	}
}

// AddPVZ создание ПВЗ
func (pvzRepo *PVZRepository) AddPVZ(city string) (models.PVZ, error) {
	sql := "INSERT INTO pvz (city) VALUES ($1) RETURNING id, registration_date, city"
	pvz := models.PVZ{}
	err := pvzRepo.db.QueryRow(sql, city).Scan(
		&pvz.Id,
		&pvz.RegistrationDate,
		&pvz.City)
	if err != nil {
		return pvz, models.ErrDBInsert
	}
	return pvz, nil
}

// GetPVZById получение ПВЗ по id
func (pvzRepo *PVZRepository) GetPVZById(pvzId string) (models.PVZ, error) {
	sql := "SELECT * FROM pvz WHERE id = $1"
	pvz := models.PVZ{}
	err := pvzRepo.db.QueryRow(sql, pvzId).Scan(
		&pvz.Id,
		&pvz.RegistrationDate,
		&pvz.City)
	if err != nil {
		return models.PVZ{}, models.ErrPVZNotFound
	}
	return pvz, nil
}

func (pvzRepo *PVZRepository) GetPVZsWithPagination(offset, limit int) ([]models.PVZ, error) {
	
}
