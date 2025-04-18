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
	err = ps.prodRepo.DeleteProduct(product.Id)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PVZService) GetPVZWithPagination(startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	// выводить только те ПВЗ и всю информацию по ним, которые в указанный диапазон времени проводили приёмы товаров

	// сделать тройной джойн
	// сделать map[pvz]reception

	// либо 2.
	// Получить ПВЗ
	// Для каждого ПВЗ получить приёмки
	// Для каждой приёмки получить товары

	// 1. получить товары в указанном диапазоне и ПВЗ id
	// 2. сделать список уникальных ресепшн ид
	// 3. join pvz + receptions С ПАГИНАЦЕЙ
	// 4. для каждого ПВЗ вставить reception
	// 5. пройтись по всем товарам и вставить в соотв. пвз и ресепшн

	// 1. products := GetProductsWithTime(startDate, endDate) ([]models.ProductWithPVZId, error)
	// ProductWithPVZId {
	// 	Id          string
	//	DateTime    string
	//	Type        string
	//	ReceptionId string
	//  PVZId		string
	//}

	// 2.
	// receptionsMap := make(map[reception_id] int)
	// for product := range products {
	//	receptionsMap[product.ReceptionId] = product.PVZId
	//}
	// uniqueReceptions := make([]int, 0, len(receptionsMap))
	// for key, _ := range receptionsMap {
	//	uniqueReceptions = append(uniqueReceptions, key)
	//}

	// 3.
	// GetPVZsByReceptionIds(receptionIds []int) ([]models.PVZAndReception, error)
	// pvzsAndReceptions := GetPVZsByReceptions(uniqueReceptions)
	// type PVZAndReception struct {
	//	PVZ        PVZ
	//	Reception  Reception
	//}

	// 4.
	//type PVZWithReceptions struct {
	//	PVZ        PVZ
	//	Receptions []ReceptionWithProducts
	//}
	// pvzWithReception := make([]PVZWithReceptions, 0)
	// pvzMap := make(map[pvzId]PVZWithReceptions, 0)
	//for p := range pvzsAndReceptions {
	//  if _, ok := pvzMap[p.PVZ.Id]; !ok{
	// pvzMap[p.PVZ.Id] := PVZWithReceptions{PVZ{...}, []ReceptionWithProducts{{...}}}
	//}  else {
	//	pvzMap[p.PVZ.Id].Receptions = append(pvzMap[p.PVZ.Id].Receptions, ReceptionWithProducts{...})
	//}
	//
	//}

	// 5.
	// for p := range products {
	//	pvzMap[p.PVZId].ReceptionWithProducts.Products = append(pvzMap[p.PVZId].ReceptionWithProducts.Products, Product{p...., p....})
	//}

}
