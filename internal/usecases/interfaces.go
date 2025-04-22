package usecases

import (
	"github.com/grigory222/avito-backend-trainee/internal/models"
	"time"
)

type ProductRepository interface {
	AddProduct(productType, receptionId string) (models.Product, error)
	GetLastProduct(recId string) (models.Product, error)
	DeleteProductById(id string) error
	//GetProductsByReceptionIds(recIds []string, startDate, endDate *time.Time) ([]models.Product, error)
}

type PVZRepository interface {
	AddPVZ(city string) (models.PVZ, error)
	GetPVZById(pvzId string) (models.PVZ, error)
	GetFlatPVZRows(startDate, endDate time.Time, offset, limit int) ([]models.FlatRow, error)
	//GetAllPVZs(ctx context.Context) ([]models.PVZ, error)
}

type ReceptionRepository interface {
	GetLastReception(pvzId string) (models.Reception, error)
	AddReception(pvzId string) (models.Reception, error)
	UpdateReceptionStatus(recId, status string) (models.Reception, error)
	//GetReceptionsByPVZIds(pvzIds []string, startDate, endDate *time.Time) ([]models.Reception, error)
}

type UserRepository interface {
	AddNewUser(email, passwordHash, role string) (models.User, error)
	GetUserByEmailAndPassword(email, passwordHash string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
}
