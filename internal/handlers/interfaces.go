package handlers

import "github.com/grigory222/avito-backend-trainee/internal/models"

type ProductService interface {
	AddProduct(productType, pvzId string) (models.Product, error)
}
