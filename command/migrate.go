package command

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/ericmarcelinotju/gram/domain"
)

// MigrateCommandFactory create and returns a factory to create command line functions for migration
func MigrateCommandFactory(services []domain.SeederService) func(context.Context) error {
	migrate := func(ctx context.Context) error {
		var migrateErrors map[string]error = make(map[string]error)
		var err error

		for _, service := range services {
			var serviceName string
			if t := reflect.TypeOf(service); t.Kind() == reflect.Ptr {
				serviceName = strings.Split(t.Elem().String(), ".")[1]
			} else {
				serviceName = t.Name()
			}

			fmt.Printf("Migrating %s\n", serviceName)
			err = service.Migrate()
			migrateErrors[serviceName] = nil
			if err != nil {
				migrateErrors[serviceName] = err
				fmt.Printf("Error when migrating %ss %s\n", serviceName, err)
			}
		}

		var successModules []string
		var errorModules []string

		for module, err := range migrateErrors {
			if err == nil {
				successModules = append(successModules, module)
			} else {
				errorModules = append(errorModules, module)
			}
		}

		if len(successModules) <= 0 {
			fmt.Println("All migrations failed!")
		} else {
			fmt.Printf("Successful migration: '%s'\n", strings.Join(successModules, ", "))
		}

		if len(errorModules) <= 0 {
			fmt.Println("all migrations successful!")
			return nil
		} else {
			return fmt.Errorf("failed migration: '%s'\n", strings.Join(errorModules, ", "))
		}
	}
	return migrate
}
