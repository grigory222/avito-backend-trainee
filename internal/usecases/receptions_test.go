package usecases_test

import (
	"errors"
	"testing"

	"github.com/grigory222/avito-backend-trainee/internal/models"
	"github.com/grigory222/avito-backend-trainee/internal/usecases"
	"github.com/grigory222/avito-backend-trainee/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
)

func TestReceptionService_AddReception(t *testing.T) {
	const pvzId = "pvz123"

	tests := []struct {
		name          string
		setupMocks    func(pvzRepo *mocks.PVZRepository, recRepo *mocks.ReceptionRepository)
		wantReception models.Reception
		wantErr       error
	}{
		{
			name: "success_add_reception",
			setupMocks: func(pvzRepo *mocks.PVZRepository, recRepo *mocks.ReceptionRepository) {
				pvzRepo.On("GetPVZById", pvzId).
					Return(models.PVZ{Id: pvzId}, nil).Once()

				recRepo.On("GetLastReception", pvzId).
					Return(models.Reception{}, nil).Once() // нет активной приемки

				recRepo.On("AddReception", pvzId).
					Return(models.Reception{Id: "rec123", PVZId: pvzId, Status: models.INPROGRESS}, nil).Once()
			},
			wantReception: models.Reception{Id: "rec123", PVZId: pvzId, Status: models.INPROGRESS},
			wantErr:       nil,
		},
		{
			name: "error_pvz_not_found",
			setupMocks: func(pvzRepo *mocks.PVZRepository, recRepo *mocks.ReceptionRepository) {
				pvzRepo.On("GetPVZById", pvzId).
					Return(models.PVZ{}, errors.New("not found")).Once()
			},
			wantReception: models.Reception{},
			wantErr:       errors.New("not found"),
		},
		{
			name: "error_already_has_inprogress_reception",
			setupMocks: func(pvzRepo *mocks.PVZRepository, recRepo *mocks.ReceptionRepository) {
				pvzRepo.On("GetPVZById", pvzId).
					Return(models.PVZ{Id: pvzId}, nil).Once()

				recRepo.On("GetLastReception", pvzId).
					Return(models.Reception{Id: "rec999", PVZId: pvzId, Status: models.INPROGRESS}, nil).Once()
			},
			wantReception: models.Reception{},
			wantErr:       models.ErrPVZAlreadyHasReception,
		},
		{
			name: "error_adding_reception_failed",
			setupMocks: func(pvzRepo *mocks.PVZRepository, recRepo *mocks.ReceptionRepository) {
				pvzRepo.On("GetPVZById", pvzId).
					Return(models.PVZ{Id: pvzId}, nil).Once()

				recRepo.On("GetLastReception", pvzId).
					Return(models.Reception{}, nil).Once()

				recRepo.On("AddReception", pvzId).
					Return(models.Reception{}, errors.New("db error")).Once()
			},
			wantReception: models.Reception{},
			wantErr:       errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPVZRepo := new(mocks.PVZRepository)
			mockRecRepo := new(mocks.ReceptionRepository)

			tt.setupMocks(mockPVZRepo, mockRecRepo)

			service := usecases.NewReceptionService(mockPVZRepo, mockRecRepo)
			gotRec, err := service.AddReception(pvzId)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr.Error())
				assert.Equal(t, models.Reception{}, gotRec)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantReception, gotRec)
			}

			mockPVZRepo.AssertExpectations(t)
			mockRecRepo.AssertExpectations(t)
		})
	}
}
