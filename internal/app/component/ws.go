package component

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WS struct {
	test     int
	Upgrader websocket.Upgrader
}

func NewWs() *WS {
	return &WS{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (this *WS) Read(conn *websocket.Conn) (messageType int, p []byte, err error) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return 0, nil, err
		}

		return messageType, p, err
	}
}

func (this *WS) Write(conn *websocket.Conn, messageType int, msg *interface{}) {
	str, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	err = conn.WriteMessage(messageType, str)
	if err != nil {
		log.Println(err)
	}
}

func (this *WS) GetConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	this.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	return this.Upgrader.Upgrade(w, r, nil)
}
