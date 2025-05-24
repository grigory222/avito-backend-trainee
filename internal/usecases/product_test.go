//go:generate mockery --name=ProductRepository --output=mocks
//go:generate mockery --name=ReceptionRepository --output=mocks
//go:generate mockery --name=PVZRepository --output=mocks

package usecases_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grigory222/avito-backend-trainee/internal/models"
	"github.com/grigory222/avito-backend-trainee/internal/usecases"
	"github.com/grigory222/avito-backend-trainee/internal/usecases/mocks"
)

func TestProductService_AddProduct(t *testing.T) {
	tests := []struct {
		name            string
		setupMocks      func() (*mocks.ProductRepository, *mocks.ReceptionRepository, *mocks.PVZRepository)
		wantErr         bool
		wantErrIs       error
		wantProductType string
	}{
		{
			name: "success",
			setupMocks: func() (*mocks.ProductRepository, *mocks.ReceptionRepository, *mocks.PVZRepository) {
				mockProdRepo := new(mocks.ProductRepository)
				mockRecRepo := new(mocks.ReceptionRepository)
				mockPVZRepo := new(mocks.PVZRepository)

				mockPVZRepo.On("GetPVZById", "pvz1").Return(models.PVZ{Id: "pvz1"}, nil)
				mockRecRepo.On("GetLastReception", "pvz1").Return(models.Reception{Id: "rec1", Status: models.INPROGRESS}, nil)
				mockProdRepo.On("AddProduct", "Electronics", "rec1").Return(models.Product{Id: "prod1", Type: "Electronics", ReceptionId: "rec1"}, nil)

				return mockProdRepo, mockRecRepo, mockPVZRepo
			},
			wantErr:         false,
			wantProductType: "Electronics",
		},
		{
			name: "pvz not found",
			setupMocks: func() (*mocks.ProductRepository, *mocks.ReceptionRepository, *mocks.PVZRepository) {
				mockProdRepo := new(mocks.ProductRepository)
				mockRecRepo := new(mocks.ReceptionRepository)
				mockPVZRepo := new(mocks.PVZRepository)

				mockPVZRepo.On("GetPVZById", "pvz1").Return(models.PVZ{}, models.ErrDBInsert)

				return mockProdRepo, mockRecRepo, mockPVZRepo
			},
			wantErr: true,
		},
		{
			name: "no active reception",
			setupMocks: func() (*mocks.ProductRepository, *mocks.ReceptionRepository, *mocks.PVZRepository) {
				mockProdRepo := new(mocks.ProductRepository)
				mockRecRepo := new(mocks.ReceptionRepository)
				mockPVZRepo := new(mocks.PVZRepository)

				mockPVZRepo.On("GetPVZById", "pvz1").Return(models.PVZ{Id: "pvz1"}, nil)
				mockRecRepo.On("GetLastReception", "pvz1").Return(models.Reception{Id: "rec1", Status: models.CLOSE}, nil)

				return mockProdRepo, mockRecRepo, mockPVZRepo
			},
			wantErr:   true,
			wantErrIs: models.ErrNoActiveReception,
		},
		{
			name: "repo error on AddProduct",
			setupMocks: func() (*mocks.ProductRepository, *mocks.ReceptionRepository, *mocks.PVZRepository) {
				mockProdRepo := new(mocks.ProductRepository)
				mockRecRepo := new(mocks.ReceptionRepository)
				mockPVZRepo := new(mocks.PVZRepository)

				mockPVZRepo.On("GetPVZById", "pvz1").Return(models.PVZ{Id: "pvz1"}, nil)
				mockRecRepo.On("GetLastReception", "pvz1").Return(models.Reception{Id: "rec1", Status: models.INPROGRESS}, nil)
				mockProdRepo.On("AddProduct", "Electronics", "rec1").Return(models.Product{}, errors.New("insert error"))

				return mockProdRepo, mockRecRepo, mockPVZRepo
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProdRepo, mockRecRepo, mockPVZRepo := tt.setupMocks()
			service := usecases.NewProductService(mockProdRepo, mockRecRepo, mockPVZRepo)

			product, err := service.AddProduct("Electronics", "pvz1")

			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantProductType, product.Type)
			}

			mockProdRepo.AssertExpectations(t)
			mockRecRepo.AssertExpectations(t)
			mockPVZRepo.AssertExpectations(t)
		})
	}
}
