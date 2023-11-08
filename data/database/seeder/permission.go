package seeder

import (
	"github.com/google/uuid"
	"gitlab.com/firelogik/helios/data/schema"
	"gorm.io/gorm"
)

type PermissionSeederService struct {
	db *gorm.DB
}

func NewPermissionSeederService(db *gorm.DB) *PermissionSeederService {
	return &PermissionSeederService{db: db}
}

func (s *PermissionSeederService) Migrate() error {
	return s.db.AutoMigrate(&schema.Permission{})
}

func (s *PermissionSeederService) Seeding() error {
	permissionsMap := map[string][]string{
		"STATISTIC":  {"GET"},
		"LOG":        {"GET", "POST", "DELETE"},
		"PERMISSION": {"GET", "POST", "PUT", "DELETE"},
		"ROLE":       {"GET", "POST", "PUT", "DELETE"},
		"SETTING":    {"GET", "POST"},
		"USER":       {"GET", "POST", "PUT", "DELETE"},
	}

	for module, methods := range permissionsMap {
		for _, method := range methods {
			entity := schema.Permission{
				Model:       schema.Model{ID: uuid.New()},
				Module:      module,
				Method:      method,
				Description: "Seeded permissions",
			}
			if err := s.db.Create(&entity).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
