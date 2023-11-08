package command

import (
	"context"
	"fmt"

	"gitlab.com/firelogik/helios/domain/model"
	"gitlab.com/firelogik/helios/domain/module/permission"
	"gitlab.com/firelogik/helios/domain/module/role"
	"gitlab.com/firelogik/helios/domain/module/user"
)

// UserCommandFactory create and returns a factory to create command line functions for user
func UserCommandFactory(permRepo permission.Repository, roleRepo role.Repository, userRepo user.Repository) func(context.Context, string) error {
	createSuperAdmin := func(ctx context.Context, username string) error {
		fmt.Printf("Creating user with username '%s'", username)

		var email string
		var password string

		fmt.Print("Email    :   ")
		fmt.Scanln(&email)
		fmt.Print("Password :   ")
		fmt.Scanln(&password)

		var role model.Role
		roles, _, err := roleRepo.SelectRole(ctx, &model.Role{
			Name: "superadmin",
		})
		if err != nil {
			return fmt.Errorf("error when reading roles %s", err)
		}

		permissions, _, err := permRepo.SelectPermission(ctx, nil)
		if err != nil {
			return fmt.Errorf("error when reading permissions %s", err)
		}

		if err != nil || len(roles) <= 0 {

			role = model.Role{
				Name:        "superadmin",
				Description: "Super Admin",
				Permissions: permissions,
			}
			err = roleRepo.InsertRole(ctx, &role)
			if err != nil {
				return fmt.Errorf("error when creating role %s", err)
			}
		} else {
			role.ID = roles[0].ID

			role = model.Role{
				ID:          roles[0].ID,
				Permissions: permissions,
			}

			err = roleRepo.UpdateRole(ctx, &role)
			if err != nil {
				return fmt.Errorf("error when updating role permissions %s", err)
			}
		}

		superAdminRoleID := role.ID

		err = userRepo.InsertUser(ctx, &model.User{
			Username: username,
			Email:    email,
			Password: password,
			RoleID:   superAdminRoleID,
		})
		if err != nil {
			return fmt.Errorf("error when creating user %s", err)
		}

		fmt.Printf("Created user with username '%s'\n", username)
		return nil
	}
	return createSuperAdmin
}
