// Code generated by mockery v2.53.4. DO NOT EDIT.

package mocks

import (
	models "github.com/grigory222/avito-backend-trainee/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// ProductRepository is an autogenerated mock type for the ProductRepository type
type ProductRepository struct {
	mock.Mock
}

// AddProduct provides a mock function with given fields: productType, receptionId
func (_m *ProductRepository) AddProduct(productType string, receptionId string) (models.Product, error) {
	ret := _m.Called(productType, receptionId)

	if len(ret) == 0 {
		panic("no return value specified for AddProduct")
	}

	var r0 models.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (models.Product, error)); ok {
		return rf(productType, receptionId)
	}
	if rf, ok := ret.Get(0).(func(string, string) models.Product); ok {
		r0 = rf(productType, receptionId)
	} else {
		r0 = ret.Get(0).(models.Product)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(productType, receptionId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteProductById provides a mock function with given fields: id
func (_m *ProductRepository) DeleteProductById(id string) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteProductById")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetLastProduct provides a mock function with given fields: recId
func (_m *ProductRepository) GetLastProduct(recId string) (models.Product, error) {
	ret := _m.Called(recId)

	if len(ret) == 0 {
		panic("no return value specified for GetLastProduct")
	}

	var r0 models.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (models.Product, error)); ok {
		return rf(recId)
	}
	if rf, ok := ret.Get(0).(func(string) models.Product); ok {
		r0 = rf(recId)
	} else {
		r0 = ret.Get(0).(models.Product)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(recId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewProductRepository creates a new instance of ProductRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProductRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProductRepository {
	mock := &ProductRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
