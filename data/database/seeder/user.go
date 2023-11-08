package seeder

import (
	"github.com/google/uuid"
	"gitlab.com/firelogik/helios/data/schema"
	"gitlab.com/firelogik/helios/utils/crypt"
	"gorm.io/gorm"
)

type UserSeederService struct {
	db *gorm.DB
}

func NewUserSeederService(db *gorm.DB) *UserSeederService {
	return &UserSeederService{db: db}
}

func (s *UserSeederService) Migrate() error {
	return s.db.AutoMigrate(&schema.User{})
}

func (s *UserSeederService) Seeding() error {
	var roles []schema.Role
	query := s.db.Model(&schema.Role{}).Find(&roles)
	if err := query.Error; err != nil {
		return err
	}

	superAdminPassword, err := crypt.Hash("super")
	if err != nil {
		return err
	}
	adminPassword, err := crypt.Hash("admin")
	if err != nil {
		return err
	}

	seedDatas := []schema.User{
		{
			Model:    schema.Model{ID: uuid.New()},
			Username: "super",
			Email:    "eric@datis.co.id",
			Password: superAdminPassword,
			RoleID:   roles[0].ID,
		},
		{
			Model:    schema.Model{ID: uuid.New()},
			Username: "admin",
			Email:    "admin@admin.com",
			Password: adminPassword,
			RoleID:   roles[1].ID,
		},
	}
	for _, seedData := range seedDatas {
		if err := s.db.Create(&seedData).Error; err != nil {
			return err
		}
	}
	return nil
}
