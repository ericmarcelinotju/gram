package seeder

import (
	"github.com/ericmarcelinotju/gram/model"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type RoleSeederService struct {
	db *gorm.DB
}

func NewRoleSeederService(db *gorm.DB) *RoleSeederService {
	return &RoleSeederService{db: db}
}

func (s *RoleSeederService) Migrate() error {
	if err := s.db.AutoMigrate(&model.RoleEntity{}); err != nil {
		return err
	}
	return s.db.AutoMigrate(&model.RolePermissionEntity{})
}

func (s *RoleSeederService) Seed() error {
	var permissions []model.PermissionEntity
	query := s.db.Model(&model.PermissionEntity{}).Find(&permissions)

	if err := query.Error; err != nil {
		return err
	}

	var adminPermissionsMap map[string][]string = map[string][]string{
		"USER":    {"GET", "POST", "PUT", "DELETE"},
		"ROLE":    {"GET"},
		"SETTING": {"GET", "POST"},
	}

	var adminPermissions []model.PermissionEntity

	for _, permission := range permissions {
		if isPermissionExistInMap(permission, adminPermissionsMap) {
			adminPermissions = append(adminPermissions, permission)
		}
	}
	seedDatas := []model.RoleEntity{
		{
			Model:       model.Model{Id: uuid.NewV4()},
			Name:        "Super Admin",
			Description: "Super Administrator",
			Permissions: permissions,
		},
		{
			Model:       model.Model{Id: uuid.NewV4()},
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

func isPermissionExistInMap(permission model.PermissionEntity, permissionsMap map[string][]string) bool {
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
