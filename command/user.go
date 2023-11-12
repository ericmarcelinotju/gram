package command

import (
	"context"
	"fmt"

	"github.com/ericmarcelinotju/gram/dto"
	permissionModule "github.com/ericmarcelinotju/gram/module/permission"
	roleModule "github.com/ericmarcelinotju/gram/module/role"
	userModule "github.com/ericmarcelinotju/gram/module/user"
)

// UserCommandFactory create and returns a factory to create command line functions for user
func UserCommandFactory(permRepo permissionModule.Repository, roleRepo roleModule.Repository, userRepo userModule.Repository) func(context.Context, string) error {
	createSuperAdmin := func(ctx context.Context, username string) error {
		fmt.Printf("Creating user with username '%s'", username)

		var email string
		var password string

		fmt.Print("Email    :   ")
		fmt.Scanln(&email)
		fmt.Print("Password :   ")
		fmt.Scanln(&password)

		var role dto.RoleDto
		roles, _, err := roleRepo.Select(ctx, &dto.RoleDto{Name: "superadmin"}, nil, nil)
		if err != nil {
			return fmt.Errorf("error when reading roles %s", err)
		}

		permissions, _, err := permRepo.Select(ctx, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("error when reading permissions %s", err)
		}

		if err != nil || len(roles) <= 0 {

			role = dto.RoleDto{
				Name:        "superadmin",
				Description: "Super Admin",
				Permissions: permissions,
			}
			err = roleRepo.Insert(ctx, &role)
			if err != nil {
				return fmt.Errorf("error when creating role %s", err)
			}
		} else {
			role.Id = roles[0].Id

			role = dto.RoleDto{
				Id:          role.Id,
				Permissions: permissions,
			}

			err = roleRepo.Update(ctx, &role)
			if err != nil {
				return fmt.Errorf("error when updating role permissions %s", err)
			}
		}

		superAdminRoleID := role.Id

		err = userRepo.Insert(ctx, &dto.UserDto{
			Username: username,
			Email:    email,
			Password: password,
			RoleId:   superAdminRoleID,
		})
		if err != nil {
			return fmt.Errorf("error when creating user %s", err)
		}

		fmt.Printf("Created user with username '%s'\n", username)
		return nil
	}
	return createSuperAdmin
}
