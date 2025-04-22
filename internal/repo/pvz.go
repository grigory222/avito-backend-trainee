package repo

import (
	"github.com/grigory222/avito-backend-trainee/internal/models"
	"github.com/jmoiron/sqlx"
	"time"
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

func (pvzRepo *PVZRepository) GetFlatPVZRows(startDate, endDate time.Time, offset, limit int) ([]models.FlatRow, error) {
	sql := `
		WITH filtered_pvzs AS (
			SELECT DISTINCT pvzs.id
			FROM pvzs
			JOIN receptions ON pvzs.id = receptions.pvz_id
			JOIN products ON receptions.id = products.reception_id
			WHERE products.date_time BETWEEN $1 AND $2
			ORDER BY pvzs.id
			OFFSET $3
			LIMIT $4
		)
		SELECT
			p.id AS pvz_id,
			p.city,
			p.registration_date,
			r.id AS reception_id,
			r.date_time AS reception_date,
			r.status,
			pr.id AS product_id,
			pr.date_time AS product_date,
			pr.product_type
		FROM filtered_pvzs f
		JOIN pvzs p ON p.id = f.id
		JOIN receptions r ON p.id = r.pvz_id
		JOIN products pr ON r.id = pr.reception_id
		WHERE pr.date_time BETWEEN $1 AND $2
		ORDER BY p.id, r.id, pr.id;
	`

	rows, err := pvzRepo.db.Query(sql, startDate, endDate, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.FlatRow
	for rows.Next() {
		var row models.FlatRow
		err := rows.Scan(
			&row.PVZId, &row.City, &row.RegistrationDate,
			&row.ReceptionId, &row.ReceptionDate, &row.Status,
			&row.ProductId, &row.ProductDate, &row.ProductType,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, rows.Err()
}
