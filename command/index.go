package command

import (
	"context"
	"flag"
	"fmt"
	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/plugins/database"
	"os"
	"time"

	permissionModule "github.com/ericmarcelinotju/gram/module/permission"
	roleModule "github.com/ericmarcelinotju/gram/module/role"
	userModule "github.com/ericmarcelinotju/gram/module/user"
	"github.com/ericmarcelinotju/gram/plugins/database/seeder"
	"gorm.io/gorm"
)

type SeederService interface {
	Seed() error
	Migrate() error
}

func ProcessCommands(db *gorm.DB) {
	db = db.Session(&gorm.Session{SkipHooks: true})

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	cmdUser := flag.String("u", "", "Create super user")
	cmdMigrate := flag.Bool("m", false, "Migrate tables")
	cmdSeeding := flag.Bool("s", false, "Seeding Init Value")
	cmdIsTesting := flag.Bool("t", false, "Is Testing Environtment")
	flag.Parse()
	if *cmdIsTesting {
		config.ChangeEnv(".env.test")
		configuration := config.Reload()
		newDB, err := database.Connect(configuration.Database)
		if err == nil {
			db = newDB
		}
	}

	if cmdUser != nil && len(*cmdUser) > 0 {
		userRepo := userModule.NewRepository(db, nil, nil)
		roleRepo := roleModule.NewRepository(db)
		permRepo := permissionModule.NewRepository(db)

		createSuperAdmin := UserCommandFactory(permRepo, roleRepo, userRepo)
		err := createSuperAdmin(ctx, *cmdUser)
		if err != nil {
			cancel()
			fmt.Println(err)
			os.Exit(1)
			return
		}
		cancel()
		os.Exit(0)
	} else if cmdMigrate != nil && *cmdMigrate {
		migrate := MigrateCommandFactory(
			[]SeederService{
				seeder.NewAuditSeederService(db),
				seeder.NewSettingSeederService(db),
				seeder.NewPermissionSeederService(db),
				seeder.NewRoleSeederService(db),
				seeder.NewUserSeederService(db),
			},
		)
		err := migrate(ctx)
		if err != nil {
			cancel()
			fmt.Println(err)
			os.Exit(1)
			return
		}
		cancel()
		os.Exit(0)
	} else if cmdSeeding != nil && *cmdSeeding {
		seeding := SeedingCommandFactory(
			[]SeederService{
				seeder.NewSettingSeederService(db),
				seeder.NewPermissionSeederService(db),
				seeder.NewRoleSeederService(db),
				seeder.NewUserSeederService(db),
			},
		)
		err := seeding(ctx)
		if err != nil {
			cancel()
			fmt.Println(err)
			os.Exit(1)
			return
		}
		cancel()
		os.Exit(0)
	}
	cancel()
}
