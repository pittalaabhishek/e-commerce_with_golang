package models

import (
    "errors"
    "gorm.io/gorm"
)

var ErrInvalidRating = errors.New("rating must be between 1 and 5")

type Review struct {
    gorm.Model
    ProductID uint   `json:"product_id" gorm:"not null"`
    Name      string `json:"name" gorm:"default:'Anonymous'"`
    Review    string `json:"review"`
    Rating    int    `json:"rating" gorm:"check:rating >= 1 AND rating <= 5"`
}