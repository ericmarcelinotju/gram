package command

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ericmarcelinotju/gram/domain"
	permissionModule "github.com/ericmarcelinotju/gram/module/permission"
	roleModule "github.com/ericmarcelinotju/gram/module/role"
	userModule "github.com/ericmarcelinotju/gram/module/user"
	"github.com/ericmarcelinotju/gram/repository/database/seeder"
	"gorm.io/gorm"
)

func ProcessCommands(db *gorm.DB) {
	db = db.Session(&gorm.Session{SkipHooks: true})

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	cmdUser := flag.String("u", "", "Create super user")
	cmdMigrate := flag.Bool("m", false, "Migrate tables")
	cmdSeeding := flag.Bool("s", false, "Seeding Init Value")
	flag.Parse()

	if cmdUser != nil && len(*cmdUser) > 0 {
		userRepo := userModule.NewRepository(db, nil)
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
			[]domain.SeederService{
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
			[]domain.SeederService{
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
