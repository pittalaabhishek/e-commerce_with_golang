package mocks

import (
    "e-commerce_with_golang/models"
    "testing"
    "github.com/stretchr/testify/mock"
)

// ReviewRepository is a mock implementation of domain.ReviewRepository
type ReviewRepository struct {
    mock.Mock
}

func (m *ReviewRepository) GetByProductID(productID uint) ([]models.Review, error) {
    args := m.Called(productID)
    return args.Get(0).([]models.Review), args.Error(1)
}

func (m *ReviewRepository) Create(review *models.Review) error {
    args := m.Called(review)
    return args.Error(0)
}

// NewReviewRepository creates a new instance of ReviewRepository
func NewReviewRepository(t *testing.T) *ReviewRepository {
    mock := &ReviewRepository{}
    mock.Mock.Test(t)

    t.Cleanup(func() { mock.AssertExpectations(t) })

    return mock
}
