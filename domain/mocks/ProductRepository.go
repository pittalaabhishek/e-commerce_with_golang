package mocks

import (
    "e-commerce_with_golang/models"
    "testing"
    "github.com/stretchr/testify/mock"
)

// ProductRepository is a mock implementation of domain.ProductRepository
type ProductRepository struct {
    mock.Mock
}

func (m *ProductRepository) GetAll() ([]models.Product, error) {
    args := m.Called()
    return args.Get(0).([]models.Product), args.Error(1)
}

func (m *ProductRepository) GetByID(id uint) (*models.Product, error) {
    args := m.Called(id)
    return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductRepository) Create(product *models.Product) error {
    args := m.Called(product)
    return args.Error(0)
}

// NewProductRepository creates a new instance of ProductRepository
func NewProductRepository(t *testing.T) *ProductRepository {
    mock := &ProductRepository{}
    mock.Mock.Test(t)

    t.Cleanup(func() { mock.AssertExpectations(t) })

    return mock
}
