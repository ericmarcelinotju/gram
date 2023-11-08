package websocket

import (
	"errors"
	"log"

	"github.com/google/uuid"
)

// Message wraps the relevant information needed to broadcast a message
type Message struct {
	Channel   string // the channel name to broadcast on
	Operation string
	Data      interface{} // the data to broadcast
	Type      string
}

func (m Message) toResponseModel() (interface{}, error) {
	// if m.Type == "device" {
	// 	response, ok := m.Data.(*model.Client)
	// 	if !ok {
	// 		return nil, errors.New("invalid data")
	// 	}
	// 	return deviceResp.ToResponseModel(response), nil
	// } else if m.Type == "event" {
	// 	response, ok := m.Data.(*model.Event)
	// 	if !ok {
	// 		return nil, errors.New("invalid data")
	// 	}
	// 	return eventResp.ToResponseModel(response), nil
	// }
	return nil, errors.New("invalid data type")
}

// Dispatcher maintains the set of active clients and broadcasts messages to the clients.
type Dispatcher struct {
	// Broadcase messages to all client.
	broadcast chan *Message

	// Registered clients.
	Clients map[string]map[*Client]bool

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	OnHeartbeat func(string, bool)

	segment uuid.UUID
}

// NewDispatcher creates a new Dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[string]map[*Client]bool),
	}
}

// Broadcast returns the broadcast channel
func (d *Dispatcher) Broadcast() chan *Message {
	// NOTE: To scale this, a message queue must be used to "broadcast" messages
	// to all running webservers. This function would enqueue messages to the
	// the message queue system, while the dispatcher also reads messages from the queue
	// and sends them to the `broadcast` channel
	return d.broadcast
}

func (d *Dispatcher) sendMessage(channel, data string) {
	if clients, ok := d.Clients[channel]; ok {
		for client := range clients {
			select {
			case client.send <- data:
			default:
				close(client.send)
				delete(d.Clients[channel], client)
			}
		}
	}
}

// Run starts the dispatch loop
func (d *Dispatcher) Run() {
	for {
		select {
		case client := <-d.Register:
			channel := client.channel
			log.Printf("registered new client to '%s'", channel)
			if _, ok := d.Clients[channel]; !ok {
				d.Clients[channel] = make(map[*Client]bool)
			}
			d.Clients[channel][client] = true
			if channel == "heartbeat" && client.Key != nil {
				d.OnHeartbeat(*client.Key, true)
				d.sendMessage("gateway", *client.Key)
			}
		case client := <-d.unregister:
			channel := client.channel
			log.Printf("unregistered client from '%s'", channel)
			if _, ok := d.Clients[channel][client]; ok {
				delete(d.Clients[channel], client)
				close(client.send)
			}
			if channel == "heartbeat" && client.Key != nil {
				d.OnHeartbeat(*client.Key, false)
				d.sendMessage("gateway", *client.Key)
			}
		case message := <-d.broadcast:
			channel := message.Channel
			response, err := message.toResponseModel()
			if err != nil {
				log.Println("[Websocket] convert data error : " + err.Error())
				continue
			}
			if clients, ok := d.Clients[channel]; ok {
				for client := range clients {
					select {
					case client.send <- response:
					default:
						close(client.send)
						delete(d.Clients[channel], client)
					}
				}
			}
		}
	}
}
