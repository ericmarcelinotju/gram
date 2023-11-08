package websocket

import (
	ws "gitlab.com/firelogik/helios/library/websocket"
)

func Init() (*ws.Dispatcher, error) {
	dispatcher := ws.NewDispatcher()
	go dispatcher.Run()

	return dispatcher, nil
}
