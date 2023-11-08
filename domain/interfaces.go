package domain

// Service provides an abstraction on top of the service data source
type SeederService interface {
	Migrate() error
	Seeding() error
}
