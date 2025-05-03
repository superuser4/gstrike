package relay

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var WSConn *websocket.Conn

func WSHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Websock upgrade error: %v\n", err)
		return
	}
	WSConn = conn
}
