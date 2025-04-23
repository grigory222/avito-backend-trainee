package repo

import (
	"github.com/grigory222/avito-backend-trainee/internal/models"
	"github.com/jmoiron/sqlx"
)

type ReceptionRepository struct {
	db *sqlx.DB
}

func NewReceptionRepository(db *sqlx.DB) *ReceptionRepository {
	return &ReceptionRepository{
		db: db,
	}
}

func (recRepo *ReceptionRepository) AddReception(pvzId string) (models.Reception, error) {
	sql := "INSERT INTO receptions (pvz_id) VALUES ($1) RETURNING id, date_time, pvz_id, status"
	reception := models.Reception{}
	err := recRepo.db.QueryRow(sql, pvzId).Scan(
		&reception.Id,
		&reception.DateTime,
		&reception.PVZId,
		&reception.Status,
	)
	if err != nil {
		return models.Reception{}, models.ErrDBInsert
	}
	return reception, nil

}

// получение последней приемки
func (recRepo *ReceptionRepository) GetLastReception(pvzId string) (models.Reception, error) {
	sql := "SELECT * FROM receptions WHERE pvz_id = $1 ORDER BY date_time DESC LIMIT 1"
	reception := models.Reception{}
	err := recRepo.db.QueryRow(sql, pvzId).Scan(
		&reception.Id,
		&reception.DateTime,
		&reception.PVZId,
		&reception.Status,
	)
	if err != nil {
		return models.Reception{}, models.ErrDBRead
	}
	return reception, nil
}

// обновление статуса
func (recRepo *ReceptionRepository) UpdateReceptionStatus(recId, status string) (models.Reception, error) {
	sql := "UPDATE receptions SET status = $1 WHERE id = $2 RETURNING id, date_time, pvz_id, status "
	reception := models.Reception{}
	err := recRepo.db.QueryRow(sql, status, recId).Scan(
		&reception.Id,
		&reception.DateTime,
		&reception.PVZId,
		&reception.Status,
	)
	if err != nil {
		return models.Reception{}, models.ErrDBUpdate
	}
	return reception, nil
}

// получение по PVZ ids
