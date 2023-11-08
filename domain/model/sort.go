package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Sort struct {
	Column string
	Order  string
}

func (s Sort) Apply(db *gorm.DB) *gorm.DB {
	if len(s.Column) > 0 && len(s.Order) > 0 {
		order := fmt.Sprintf("%s %s", s.Column, s.Order)
		return db.Order(order)
	}
	return db.Order("created_at DESC")
}
