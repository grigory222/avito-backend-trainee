package usecases

import (
	"github.com/grigory222/avito-backend-trainee/internal/models"
	"time"
)

type PVZService struct {
	pvzRepo  PVZRepository
	recRepo  ReceptionRepository
	prodRepo ProductRepository
}

func NewPVZService(pvzRepo PVZRepository, recRepo ReceptionRepository, prodRepo ProductRepository) *PVZService {
	return &PVZService{
		pvzRepo:  pvzRepo,
		recRepo:  recRepo,
		prodRepo: prodRepo,
	}
}

func (ps *PVZService) AddPVZ(city string) (models.PVZ, error) {
	pvz, err := ps.pvzRepo.AddPVZ(city)
	if err != nil {
		return models.PVZ{}, err
	}
	return pvz, nil
}

func (ps *PVZService) CloseLastReception(pvzId string) (models.Reception, error) {
	// проверить корректность id ПВЗ
	_, err := ps.pvzRepo.GetPVZById(pvzId)
	if err != nil {
		return models.Reception{}, err
	}

	// получить текущую(последнюю) приемку
	reception, err := ps.recRepo.GetLastReception(pvzId)
	if err != nil {
		return models.Reception{}, err
	}
	if reception.Status != models.INPROGRESS {
		return models.Reception{}, models.ErrNoActiveReception
	}

	// обновить статус
	newReception, err := ps.recRepo.UpdateReceptionStatus(reception.Id, models.CLOSE)
	if err != nil {
		return models.Reception{}, err
	}
	return newReception, nil
}

func (ps *PVZService) DeleteLastProduct(pvzId string) error {

	// проверить корректность id ПВЗ
	_, err := ps.pvzRepo.GetPVZById(pvzId)
	if err != nil {
		return err
	}

	// получить текущую(последнюю) приемку
	reception, err := ps.recRepo.GetLastReception(pvzId)
	if err != nil {
		return err
	}
	if reception.Status != models.INPROGRESS {
		return models.ErrNoActiveReception
	}

	// получить последний товар
	product, err := ps.prodRepo.GetLastProduct(reception.Id)
	if err != nil {
		return models.ErrNoProductsInReception
	}

	// удалить товар
	err = ps.prodRepo.DeleteProductById(product.Id)
	if err != nil {
		return err
	}
	return nil
}

// GetPVZWithPagination выводить только те ПВЗ и всю информацию по ним, которые в указанный диапазон времени проводили приёмы товаров
func (ps *PVZService) GetPVZWithPagination(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	offset := (page - 1) * limit
	rows, err := ps.pvzRepo.GetFlatPVZRows(*startDate, *endDate, offset, limit)
	if err != nil {
		return nil, err
	}
	pvzMap := make(map[string]*models.PVZWithReceptions)
	for _, row := range rows {
		// найти или создать ПВЗ
		pvz, ok := pvzMap[row.PVZId]
		if !ok {
			pvz = &models.PVZWithReceptions{
				PVZ: models.PVZ{
					Id:               row.PVZId,
					City:             row.City,
					RegistrationDate: row.RegistrationDate,
				},
			}
			pvzMap[row.PVZId] = pvz
		}

		// найти или создать приёмку
		var recPtr *models.ReceptionWithProducts
		for i := range pvz.Receptions {
			if pvz.Receptions[i].Reception.Id == row.ReceptionId {
				recPtr = &pvz.Receptions[i]
				break
			}
		}
		if recPtr == nil {
			newRec := models.ReceptionWithProducts{
				Reception: models.Reception{
					Id:       row.ReceptionId,
					DateTime: row.ReceptionDate,
					PVZId:    row.PVZId,
					Status:   row.Status,
				},
			}
			pvz.Receptions = append(pvz.Receptions, newRec)
			recPtr = &pvz.Receptions[len(pvz.Receptions)-1]
		}

		// добавить товар
		var rowProductId, rowProductDate, rowProductType string
		if row.ProductId.Valid && row.ProductDate.Valid && row.ProductType.Valid {
			rowProductId = row.ProductId.String
			rowProductDate = row.ProductDate.String
			rowProductType = row.ProductType.String
		}

		recPtr.Products = append(recPtr.Products, models.Product{
			Id:          rowProductId,
			DateTime:    rowProductDate,
			Type:        rowProductType,
			ReceptionId: row.ReceptionId,
		})
	}

	// собрать слайс
	var result []models.PVZWithReceptions
	for _, pvz := range pvzMap {
		result = append(result, *pvz)
	}

	return result, nil
}
