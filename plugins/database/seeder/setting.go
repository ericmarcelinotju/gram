package seeder

import (
	"github.com/ericmarcelinotju/gram/constant"
	"github.com/ericmarcelinotju/gram/model"
	"gorm.io/gorm"
)

type SettingSeederService struct {
	db *gorm.DB
}

func NewSettingSeederService(db *gorm.DB) *SettingSeederService {
	return &SettingSeederService{db: db}
}

func (s *SettingSeederService) Migrate() error {
	return s.db.AutoMigrate(&model.SettingEntity{})
}

func (s *SettingSeederService) Seed() error {
	seedDatas := []model.SettingEntity{
		{
			Name:  constant.SMTPHost,
			Value: "smtp.zoho.com",
		},
		{
			Name:  constant.SMTPPort,
			Value: "587",
		},
		{
			Name:  constant.SMTPEmail,
			Value: "noreply@mail.apilogik.com",
		},
		{
			Name:  constant.SMTPPassword,
			Value: "G3iSHExStmBd",
		},
		{
			Name:  constant.SFTPHost,
			Value: "10.39.105.80",
		},
		{
			Name:  constant.SFTPPort,
			Value: "22",
		},
		{
			Name:  constant.SFTPUsername,
			Value: "sftpuser",
		},
		{
			Name:  constant.SFTPPassword,
			Value: "P@ssw0rd1!",
		},
		{
			Name:  constant.SFTPStorageFolder,
			Value: "recording",
		},
	}

	for _, seedData := range seedDatas {
		if err := s.db.Create(&seedData).Error; err != nil {
			return err
		}
	}
	return nil
}
