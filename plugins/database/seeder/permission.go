package seeder

import (
	"github.com/ericmarcelinotju/gram/model"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type PermissionSeederService struct {
	db *gorm.DB
}

func NewPermissionSeederService(db *gorm.DB) *PermissionSeederService {
	return &PermissionSeederService{db: db}
}

func (s *PermissionSeederService) Migrate() error {
	return s.db.AutoMigrate(&model.PermissionEntity{})
}

func (s *PermissionSeederService) Seed() error {
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
			entity := model.PermissionEntity{
				Model:       model.Model{Id: uuid.NewV4()},
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
