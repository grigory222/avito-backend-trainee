package usecases_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/grigory222/avito-backend-trainee/internal/models"
	"github.com/grigory222/avito-backend-trainee/internal/usecases"
	"github.com/grigory222/avito-backend-trainee/internal/usecases/mocks"
)

func TestPVZService_AddPVZ(t *testing.T) {
	mockPVZRepo := new(mocks.PVZRepository)
	mockRecRepo := new(mocks.ReceptionRepository)
	mockProdRepo := new(mocks.ProductRepository)

	service := usecases.NewPVZService(mockPVZRepo, mockRecRepo, mockProdRepo)

	city := "Moscow"

	t.Run("success", func(t *testing.T) {
		mockPVZRepo.ExpectedCalls = nil

		mockPVZRepo.On("AddPVZ", city).Return(models.PVZ{Id: "pvz1", City: city}, nil)

		pvz, err := service.AddPVZ(city)
		assert.NoError(t, err)
		assert.Equal(t, city, pvz.City)

		mockPVZRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockPVZRepo.ExpectedCalls = nil

		mockPVZRepo.On("AddPVZ", city).Return(models.PVZ{}, errors.New("add error"))

		_, err := service.AddPVZ(city)
		assert.Error(t, err)
		assert.EqualError(t, err, "add error")

		mockPVZRepo.AssertExpectations(t)
	})
}

func TestPVZService_CloseLastReception(t *testing.T) {
	mockPVZRepo := new(mocks.PVZRepository)
	mockRecRepo := new(mocks.ReceptionRepository)
	mockProdRepo := new(mocks.ProductRepository)

	service := usecases.NewPVZService(mockPVZRepo, mockRecRepo, mockProdRepo)

	pvzId := "pvz1"
	recId := "rec1"

	tests := []struct {
		name           string
		setupMocks     func()
		expectedError  error
		expectedStatus string
	}{
		{
			name: "success",
			setupMocks: func() {
				mockPVZRepo.On("GetPVZById", pvzId).Return(models.PVZ{Id: pvzId}, nil)
				mockRecRepo.On("GetLastReception", pvzId).Return(models.Reception{Id: recId, Status: models.INPROGRESS}, nil)
				mockRecRepo.On("UpdateReceptionStatus", recId, models.CLOSE).Return(models.Reception{Id: recId, Status: models.CLOSE}, nil)
			},
			expectedError:  nil,
			expectedStatus: models.CLOSE,
		},
		{
			name: "pvz not found",
			setupMocks: func() {
				mockPVZRepo.On("GetPVZById", pvzId).Return(models.PVZ{}, errors.New("not found"))
			},
			expectedError: errors.New("not found"),
		},
		{
			name: "no active reception",
			setupMocks: func() {
				mockPVZRepo.On("GetPVZById", pvzId).Return(models.PVZ{Id: pvzId}, nil)
				mockRecRepo.On("GetLastReception", pvzId).Return(models.Reception{Id: recId, Status: models.CLOSE}, nil)
			},
			expectedError: models.ErrNoActiveReception,
		},
		{
			name: "error updating status",
			setupMocks: func() {
				mockPVZRepo.On("GetPVZById", pvzId).Return(models.PVZ{Id: pvzId}, nil)
				mockRecRepo.On("GetLastReception", pvzId).Return(models.Reception{Id: recId, Status: models.INPROGRESS}, nil)
				mockRecRepo.On("UpdateReceptionStatus", recId, models.CLOSE).Return(models.Reception{}, errors.New("update error"))
			},
			expectedError: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPVZRepo.ExpectedCalls = nil
			mockRecRepo.ExpectedCalls = nil
			mockProdRepo.ExpectedCalls = nil

			tt.setupMocks()

			received, err := service.CloseLastReception(pvzId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, received.Status)
			}

			mockPVZRepo.AssertExpectations(t)
			mockRecRepo.AssertExpectations(t)
			mockProdRepo.AssertExpectations(t)
		})
	}
}

func TestPVZService_DeleteLastProduct(t *testing.T) {
	mockPVZRepo := new(mocks.PVZRepository)
	mockRecRepo := new(mocks.ReceptionRepository)
	mockProdRepo := new(mocks.ProductRepository)

	service := usecases.NewPVZService(mockPVZRepo, mockRecRepo, mockProdRepo)

	pvzId := "pvz1"
	recId := "rec1"
	prodId := "prod1"

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError error
	}{
		{
			name: "success",
			setupMocks: func() {
				mockPVZRepo.On("GetPVZById", pvzId).Return(models.PVZ{Id: pvzId}, nil)
				mockRecRepo.On("GetLastReception", pvzId).Return(models.Reception{Id: recId, Status: models.INPROGRESS}, nil)
				mockProdRepo.On("GetLastProduct", recId).Return(models.Product{Id: prodId}, nil)
				mockProdRepo.On("DeleteProductById", prodId).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "pvz not found",
			setupMocks: func() {
				mockPVZRepo.On("GetPVZById", pvzId).Return(models.PVZ{}, errors.New("not found"))
			},
			expectedError: errors.New("not found"),
		},
		{
			name: "no active reception",
			setupMocks: func() {
				mockPVZRepo.On("GetPVZById", pvzId).Return(models.PVZ{Id: pvzId}, nil)
				mockRecRepo.On("GetLastReception", pvzId).Return(models.Reception{Id: recId, Status: models.CLOSE}, nil)
			},
			expectedError: models.ErrNoActiveReception,
		},
		{
			name: "no products in reception",
			setupMocks: func() {
				mockPVZRepo.On("GetPVZById", pvzId).Return(models.PVZ{Id: pvzId}, nil)
				mockRecRepo.On("GetLastReception", pvzId).Return(models.Reception{Id: recId, Status: models.INPROGRESS}, nil)
				mockProdRepo.On("GetLastProduct", recId).Return(models.Product{}, errors.New("no products"))
			},
			expectedError: models.ErrNoProductsInReception,
		},
		{
			name: "error deleting product",
			setupMocks: func() {
				mockPVZRepo.On("GetPVZById", pvzId).Return(models.PVZ{Id: pvzId}, nil)
				mockRecRepo.On("GetLastReception", pvzId).Return(models.Reception{Id: recId, Status: models.INPROGRESS}, nil)
				mockProdRepo.On("GetLastProduct", recId).Return(models.Product{Id: prodId}, nil)
				mockProdRepo.On("DeleteProductById", prodId).Return(errors.New("delete error"))
			},
			expectedError: errors.New("delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPVZRepo.ExpectedCalls = nil
			mockRecRepo.ExpectedCalls = nil
			mockProdRepo.ExpectedCalls = nil

			tt.setupMocks()

			err := service.DeleteLastProduct(pvzId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockPVZRepo.AssertExpectations(t)
			mockRecRepo.AssertExpectations(t)
			mockProdRepo.AssertExpectations(t)
		})
	}
}

func TestPVZService_GetPVZWithPagination(t *testing.T) {
	mockPVZRepo := new(mocks.PVZRepository)
	mockRecRepo := new(mocks.ReceptionRepository)
	mockProdRepo := new(mocks.ProductRepository)

	service := usecases.NewPVZService(mockPVZRepo, mockRecRepo, mockProdRepo)

	startDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	page := 1
	limit := 10
	offset := (page - 1) * limit

	mockRows := []models.FlatRow{
		{
			PVZId:            "pvz1",
			City:             "Moscow",
			RegistrationDate: "2023-01-01T10:00:00Z",
			ReceptionId:      "rec1",
			ReceptionDate:    "2023-05-01T12:00:00Z",
			Status:           "INPROGRESS",
			ProductId:        sql.NullString{String: "prod1", Valid: true},
			ProductDate:      sql.NullString{String: "2023-05-01T12:10:00Z", Valid: true},
			ProductType:      sql.NullString{String: "Electronics", Valid: true},
		},
		{
			PVZId:            "pvz1",
			City:             "Moscow",
			RegistrationDate: "2023-01-01T10:00:00Z",
			ReceptionId:      "rec1",
			ReceptionDate:    "2023-05-01T12:00:00Z",
			Status:           "INPROGRESS",
			ProductId:        sql.NullString{String: "prod2", Valid: true},
			ProductDate:      sql.NullString{String: "2023-05-01T12:15:00Z", Valid: true},
			ProductType:      sql.NullString{String: "Books", Valid: true},
		},
		{
			PVZId:            "pvz2",
			City:             "Saint Petersburg",
			RegistrationDate: "2023-02-01T11:00:00Z",
			ReceptionId:      "rec2",
			ReceptionDate:    "2023-06-01T13:00:00Z",
			Status:           "CLOSE",
			ProductId:        sql.NullString{Valid: false},
			ProductDate:      sql.NullString{Valid: false},
			ProductType:      sql.NullString{Valid: false},
		},
	}

	mockPVZRepo.
		On("GetFlatPVZRows", startDate, endDate, offset, limit).
		Return(mockRows, nil)

	result, err := service.GetPVZWithPagination(&startDate, &endDate, page, limit)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	expected := map[string]struct {
		City             string
		RegistrationDate string
		Receptions       map[string]struct {
			DateTime string
			Status   string
			Products []struct {
				Id          string
				DateTime    string
				Type        string
				ReceptionId string
			}
		}
	}{
		"pvz1": {
			City:             "Moscow",
			RegistrationDate: "2023-01-01T10:00:00Z",
			Receptions: map[string]struct {
				DateTime string
				Status   string
				Products []struct {
					Id          string
					DateTime    string
					Type        string
					ReceptionId string
				}
			}{
				"rec1": {
					DateTime: "2023-05-01T12:00:00Z",
					Status:   "INPROGRESS",
					Products: []struct {
						Id          string
						DateTime    string
						Type        string
						ReceptionId string
					}{
						{"prod1", "2023-05-01T12:10:00Z", "Electronics", "rec1"},
						{"prod2", "2023-05-01T12:15:00Z", "Books", "rec1"},
					},
				},
			},
		},
		"pvz2": {
			City:             "Saint Petersburg",
			RegistrationDate: "2023-02-01T11:00:00Z",
			Receptions: map[string]struct {
				DateTime string
				Status   string
				Products []struct {
					Id          string
					DateTime    string
					Type        string
					ReceptionId string
				}
			}{
				"rec2": {
					DateTime: "2023-06-01T13:00:00Z",
					Status:   "CLOSE",
					Products: []struct {
						Id          string
						DateTime    string
						Type        string
						ReceptionId string
					}{
						{"", "", "", "rec2"},
					},
				},
			},
		},
	}

	for _, pvz := range result {
		expPVZ, ok := expected[pvz.PVZ.Id]
		assert.True(t, ok, "unexpected PVZ id: %s", pvz.PVZ.Id)
		assert.Equal(t, expPVZ.City, pvz.PVZ.City)
		assert.Equal(t, expPVZ.RegistrationDate, pvz.PVZ.RegistrationDate)

		assert.Len(t, pvz.Receptions, len(expPVZ.Receptions))
		for _, rec := range pvz.Receptions {
			expRec, ok := expPVZ.Receptions[rec.Reception.Id]
			assert.True(t, ok, "unexpected Reception id: %s", rec.Reception.Id)
			assert.Equal(t, expRec.DateTime, rec.Reception.DateTime)
			assert.Equal(t, expRec.Status, rec.Reception.Status)
			assert.Equal(t, pvz.PVZ.Id, rec.Reception.PVZId)

			assert.Len(t, rec.Products, len(expRec.Products))
			for i, prod := range rec.Products {
				expProd := expRec.Products[i]
				assert.Equal(t, expProd.Id, prod.Id)
				assert.Equal(t, expProd.DateTime, prod.DateTime)
				assert.Equal(t, expProd.Type, prod.Type)
				assert.Equal(t, expProd.ReceptionId, prod.ReceptionId)
			}
		}
	}

	mockPVZRepo.AssertExpectations(t)
}
