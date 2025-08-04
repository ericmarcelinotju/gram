package net

import (
	"crypto/tls"
	"log"
	"math/rand"
	"time"

	"net/http"

	"github.com/gorilla/websocket"
)

func (h *NetClient) Connect(channel string) (c *websocket.Conn, err error) {

	log.Printf("[WebSocketClient] Connecting to %s", h.WsURL+channel)

	header := http.Header{}

	cookie := `Authorization="Bearer ` + *h.Token + `"`

	header.Set("Cookie", cookie)

	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	c, resp, err := dialer.Dial(h.WsURL+channel, header)
	if err != nil {
		if err == websocket.ErrBadHandshake {
			log.Printf("[WebSocketClient] Handshake failed with status %d", resp.StatusCode)
		}

		sleep := 5 * time.Second
		jitter := time.Duration(rand.Int63n(int64(sleep)))
		sleep = sleep + jitter/2

		time.Sleep(sleep)
		log.Printf("[WebSocketClient] Try Redial Websocket")

		return h.Connect(channel)
	}

	return c, nil
}
