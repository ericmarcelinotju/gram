package seeder

import (
	"github.com/ericmarcelinotju/gram/data/schema"
	"gorm.io/gorm"
)

type LogSeederService struct {
	db *gorm.DB
}

func NewLogSeederService(db *gorm.DB) *LogSeederService {
	return &LogSeederService{db: db}
}

func (s *LogSeederService) Migrate() error {
	return s.db.AutoMigrate(&schema.Log{})
}

func (s *LogSeederService) Seeding() error {
	return nil
}
