package usecases

import "github.com/grigory222/avito-backend-trainee/internal/models"

type ProductService struct {
	prodRepo ProductRepository
	recRepo  ReceptionRepository
	pvzRepo  PVZRepository
}

func NewProductService(prodRepo ProductRepository, recRepo ReceptionRepository, pvzRepo PVZRepository) *ProductService {
	return &ProductService{
		prodRepo: prodRepo,
		recRepo:  recRepo,
		pvzRepo:  pvzRepo,
	}
}

func (ps *ProductService) AddProduct(productType, pvzId string) (models.Product, error) {
	// проверить наличие ПВЗ с таким id
	_, err := ps.pvzRepo.GetPVZById(pvzId)
	if err != nil {
		return models.Product{}, err
	}

	// проверить статус последней приёмки
	reception, err := ps.recRepo.GetLastReception(pvzId)
	if err != nil {
		return models.Product{}, err
	}
	if reception.Status != models.INPROGRESS {
		return models.Product{}, models.ErrNoActiveReception
	}

	// добавить продукт
	product, err := ps.prodRepo.AddProduct(productType, reception.Id)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}
