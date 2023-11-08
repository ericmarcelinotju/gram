package model

import "gorm.io/gorm"

type Pagination struct {
	Limit  int
	Offset int
}

func (p Pagination) Apply(db *gorm.DB) *gorm.DB {
	if p.Limit > 0 && p.Offset >= 0 {
		return db.
			Limit(p.Limit).
			Offset(p.Offset)
	}
	return db
}
