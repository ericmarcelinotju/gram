package websocket

import (
	"github.com/gorilla/websocket"
)

type WebsocketService interface {
	Connect(*websocket.Conn, *Channel) error
}

type Service struct {
	repo Repository
}

// NewService creates a new service struct
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (svc *Service) Connect(conn *websocket.Conn, channel *Channel) error {
	return svc.repo.Connect(conn, channel)
}
