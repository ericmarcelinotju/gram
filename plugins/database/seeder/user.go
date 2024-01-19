package seeder

import (
	"github.com/ericmarcelinotju/gram/model"
	"github.com/ericmarcelinotju/gram/utils/crypt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type UserSeederService struct {
	db *gorm.DB
}

func NewUserSeederService(db *gorm.DB) *UserSeederService {
	return &UserSeederService{db: db}
}

func (s *UserSeederService) Migrate() error {
	return s.db.AutoMigrate(&model.UserEntity{})
}

func (s *UserSeederService) Seed() error {
	var roles []model.RoleEntity
	query := s.db.Model(&model.RoleEntity{}).Find(&roles)
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

	seedDatas := []model.UserEntity{
		{
			Model:    model.Model{Id: uuid.NewV4()},
			Name:     "super",
			Email:    "eric@datis.co.id",
			Password: superAdminPassword,
			RoleId:   roles[0].Id,
		},
		{
			Model:    model.Model{Id: uuid.NewV4()},
			Name:     "admin",
			Email:    "admin@admin.com",
			Password: adminPassword,
			RoleId:   roles[1].Id,
		},
	}
	for _, seedData := range seedDatas {
		if err := s.db.Create(&seedData).Error; err != nil {
			return err
		}
	}
	return nil
}
