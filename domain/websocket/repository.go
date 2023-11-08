package websocket

import "github.com/gorilla/websocket"

// Repository provides an abstraction on top of the driver data source
type Repository interface {
	Connect(*websocket.Conn, *Channel) error
}
