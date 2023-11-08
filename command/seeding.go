package command

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/ericmarcelinotju/gram/domain"
)

// SeedingCommandFactory create and returns a factory to create command line functions for seeding
func SeedingCommandFactory(services []domain.SeederService) func(context.Context) error {
	seed := func(ctx context.Context) error {
		var seedErrors map[string]error = make(map[string]error)
		var err error

		for _, service := range services {
			var serviceName string
			if t := reflect.TypeOf(service); t.Kind() == reflect.Ptr {
				serviceName = strings.Split(t.Elem().String(), ".")[1]
			} else {
				serviceName = t.Name()
			}

			fmt.Printf("Seeding %s\n", serviceName)
			err = service.Seeding()
			seedErrors[serviceName] = nil
			if err != nil {
				seedErrors[serviceName] = err
				fmt.Printf("Error when seeding %ss %s\n", serviceName, err)
			}
		}

		var successModules []string
		var errorModules []string

		for module, err := range seedErrors {
			if err == nil {
				successModules = append(successModules, module)
			} else {
				errorModules = append(errorModules, module)
			}
		}

		if len(successModules) <= 0 {
			fmt.Println("All seeding failed!")
		} else {
			fmt.Printf("Successful seeding: '%s'\n", strings.Join(successModules, ", "))
		}

		if len(errorModules) <= 0 {
			fmt.Println("all seeding successful!")
			return nil
		} else {
			return fmt.Errorf("failed seeding: '%s'\n", strings.Join(errorModules, ", "))
		}
	}
	return seed
}
