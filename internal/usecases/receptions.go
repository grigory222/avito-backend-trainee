package usecases

import "github.com/grigory222/avito-backend-trainee/internal/models"

type ReceptionService struct {
	pvzRepo PVZRepository
	recRepo ReceptionRepository
}

func NewReceptionService(pvzRepo PVZRepository, recRepo ReceptionRepository) *ReceptionService {
	return &ReceptionService{
		pvzRepo: pvzRepo,
		recRepo: recRepo,
	}
}

func (rs *ReceptionService) AddReception(pvzId string) (models.Reception, error) {
	// проверить корректность id ПВЗ
	_, err := rs.pvzRepo.GetPVZById(pvzId)
	if err != nil {
		return models.Reception{}, err
	}

	// проверить статус последней приёмки
	reception, err := rs.recRepo.GetLastReception(pvzId)
	if reception.Id != "" && reception.Status == models.INPROGRESS {
		return models.Reception{}, models.ErrPVZAlreadyHasReception
	}

	// добавить приемку
	newReception, err := rs.recRepo.AddReception(pvzId)
	if err != nil {
		return models.Reception{}, err
	}
	return newReception, nil
}
