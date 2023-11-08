package notifier

import "gitlab.com/firelogik/helios/domain/model"

// Notifier provides an abstraction on top of the notifier logic
type Notifier interface {
	Notify(title string, body interface{}, recipient *model.User) error
	Subscribe(string, string) error
	Unsubscribe(string, string) error
}
