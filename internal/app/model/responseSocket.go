package model

import "github.com/gorilla/websocket"

type ResponseSocket struct {
	MessageType int
	Conn        *websocket.Conn
}
