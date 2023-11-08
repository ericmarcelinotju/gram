package websocket

import (
	ws "github.com/ericmarcelinotju/gram/library/websocket"
)

func Init() (*ws.Dispatcher, error) {
	dispatcher := ws.NewDispatcher()
	go dispatcher.Run()

	return dispatcher, nil
}
