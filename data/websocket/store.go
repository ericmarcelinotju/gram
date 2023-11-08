package websocket

import (
	"github.com/gorilla/websocket"
	domain "gitlab.com/firelogik/helios/domain/websocket"
	ws "gitlab.com/firelogik/helios/library/websocket"
)

type Store struct {
	dispatcher *ws.Dispatcher
}

// New creates a new Store struct
func New(dispatcher *ws.Dispatcher) (*Store, error) {
	return &Store{dispatcher: dispatcher}, nil
}

func (s *Store) Connect(conn *websocket.Conn, channel *domain.Channel) error {
	client := ws.NewClient(s.dispatcher, conn, channel.Channel, channel.Key)
	s.dispatcher.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WriteDispatch()
	go client.ReadDispatch()

	return nil
}
