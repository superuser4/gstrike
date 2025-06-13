package commgr

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)



var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r , nil)
	if err != nil {
		return
	}
	for {
		_, msg, err := conn.ReadMessage()
		
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		fmt.Printf("Received: %s", msg)
	}
}


