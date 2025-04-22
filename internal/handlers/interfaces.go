package handlers

import (
	"github.com/grigory222/avito-backend-trainee/internal/models"
	"time"
)

type ProductService interface {
	AddProduct(productType, pvzId string) (models.Product, error)
}

type PVZService interface {
	AddPVZ(city string) (models.PVZ, error)
	CloseLastReception(pvzId string) (models.Reception, error)
	DeleteLastProduct(pvzId string) error
	GetPVZWithPagination(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error)
}

type ReceptionService interface {
	AddReception(pvzId string) (models.Reception, error)
}
type UserService interface {
	UserRegister(email, password, role string) (models.User, error)
	UserLogin(email, password string) (models.User, error)
}
