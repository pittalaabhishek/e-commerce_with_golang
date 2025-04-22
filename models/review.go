package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	ProductID uint   `json:"product_id"`
	Name      string `json:"name"`
	Review    string `json:"review"`
	Rating    int    `json:"rating" gorm:"check:rating >= 1 AND rating <= 5"`
}