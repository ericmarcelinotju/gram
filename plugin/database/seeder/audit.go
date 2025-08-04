package seeder

import (
	"github.com/ericmarcelinotju/gram/model"
	"gorm.io/gorm"
)

type AuditSeederService struct {
	db *gorm.DB
}

func NewAuditSeederService(db *gorm.DB) *AuditSeederService {
	return &AuditSeederService{db: db}
}

func (s *AuditSeederService) Migrate() error {
	return s.db.AutoMigrate(&model.AuditEntity{})
}

func (s *AuditSeederService) Seed() error {
	return nil
}
