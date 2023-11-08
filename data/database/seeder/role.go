package seeder

import (
	"github.com/google/uuid"
	"gitlab.com/firelogik/helios/data/schema"
	"gorm.io/gorm"
)

type RoleSeederService struct {
	db *gorm.DB
}

func NewRoleSeederService(db *gorm.DB) *RoleSeederService {
	return &RoleSeederService{db: db}
}

func (s *RoleSeederService) Migrate() error {
	if err := s.db.AutoMigrate(&schema.Role{}); err != nil {
		return err
	}
	return s.db.AutoMigrate(&schema.RolePermission{})
}

func (s *RoleSeederService) Seeding() error {
	var permissions []schema.Permission
	query := s.db.Model(&schema.Permission{}).Find(&permissions)

	if err := query.Error; err != nil {
		return err
	}

	var adminPermissionsMap map[string][]string = map[string][]string{
		"USER":    {"GET", "POST", "PUT", "DELETE"},
		"ROLE":    {"GET"},
		"SETTING": {"GET", "POST"},
	}

	var adminPermissions []schema.Permission

	for _, permission := range permissions {
		if isPermissionExistInMap(permission, adminPermissionsMap) {
			adminPermissions = append(adminPermissions, permission)
		}
	}
	seedDatas := []schema.Role{
		{
			Model:       schema.Model{ID: uuid.New()},
			Name:        "Super Admin",
			Description: "Super Administrator",
			Permissions: permissions,
		},
		{
			Model:       schema.Model{ID: uuid.New()},
			Name:        "Admin",
			Description: "Administrator",
			Permissions: adminPermissions,
		},
	}
	for _, seedData := range seedDatas {
		if err := s.db.Create(&seedData).Error; err != nil {
			return err
		}
	}
	return nil
}

func isPermissionExistInMap(permission schema.Permission, permissionsMap map[string][]string) bool {
	for module, methods := range permissionsMap {
		if permission.Module == module {
			for _, method := range methods {
				if permission.Method == method {
					return true
				}
			}
		}
	}
	return false
}
