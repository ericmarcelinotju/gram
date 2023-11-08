package seeder

import (
	"github.com/ericmarcelinotju/gram/data/schema"
	"gorm.io/gorm"
)

type AuditSeederService struct {
	db *gorm.DB
}

func NewAuditSeederService(db *gorm.DB) *AuditSeederService {
	return &AuditSeederService{db: db}
}

func (s *AuditSeederService) Migrate() error {
	return s.db.AutoMigrate(&schema.Audit{})
}

func (s *AuditSeederService) Seeding() error {
	return nil
}
