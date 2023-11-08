package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
	domain "gitlab.com/firelogik/helios/domain/websocket"
	response "gitlab.com/firelogik/helios/utils/http"
)

// NewRoutesFactory create and returns a factory to create routes for the transactions
func NewRoutesFactory(group *gin.RouterGroup) func(svc domain.WebsocketService) {
	websocketRoutesFactory := func(svc domain.WebsocketService) {
		group.GET("/:channel", func(c *gin.Context) {
			channel := c.Param("channel")
			key, err := Bind(c)
			if err != nil {
				response.ResponseError(c, err, http.StatusInternalServerError)
				return
			}
			// upgrader upgrades the request to WS
			var upgrader = websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin:     func(r *http.Request) bool { return true },
			}
			// serveWs handles websocket requests from the peer.
			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}

			err = svc.Connect(conn, &domain.Channel{Channel: channel, Key: key})
			if err != nil {
				return
			}
		})
	}

	return websocketRoutesFactory
}
