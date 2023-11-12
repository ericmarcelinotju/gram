package notifier

import "github.com/ericmarcelinotju/gram/dto"

// Notifier provides an abstraction on top of the notifier logic
type Notifier interface {
	Notify(title string, body interface{}, recipient *dto.UserDto) error
	Subscribe(string, string) error
	Unsubscribe(string, string) error
}
