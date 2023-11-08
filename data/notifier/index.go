package notifier

import "github.com/ericmarcelinotju/gram/domain/model"

// Notifier provides an abstraction on top of the notifier logic
type Notifier interface {
	Notify(title string, body interface{}, recipient *model.User) error
	Subscribe(string, string) error
	Unsubscribe(string, string) error
}
