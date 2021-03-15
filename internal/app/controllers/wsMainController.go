package controllers

import (
	"encoding/json"
	"net/http"
	"person-key-saver/internal/app/component"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/form"
	"person-key-saver/internal/app/model"
	"person-key-saver/internal/app/service"
)

type WsMainController struct {
	container *container.Container
}

func NewWsMainController(container *container.Container) *WsMainController {
	return &WsMainController{
		container: container,
	}
}

func (this *WsMainController) SaveAction(w http.ResponseWriter, r *http.Request) {
	var wsComponent = this.container.Get(container.WS_COMPONENT).(*component.WS)

	conn, err := wsComponent.GetConnection(w, r)
	if err != nil {
		return
	}
	messageType, p, err := wsComponent.Read(conn)
	if err != nil {
		return
	}

	var payload form.Payload
	if err := json.Unmarshal(p, &payload); err != nil {
		return
	}

	executor := service.Executor{
		Container: this.container,
		Payload:   &payload,
		ResponseSocket: &model.ResponseSocket{
			MessageType: messageType,
			Conn:        conn,
		},
	}
	executor.Save()
}

func (this *WsMainController) EncodeAction(w http.ResponseWriter, r *http.Request) {
	var wsComponent = this.container.Get(container.WS_COMPONENT).(*component.WS)
	conn, err := wsComponent.GetConnection(w, r)
	if err != nil {
		return
	}
	messageType, p, err := wsComponent.Read(conn)
	if err != nil {
		return
	}

	var payload form.Payload
	if err := json.Unmarshal(p, &payload); err != nil {
		return
	}

	executor := service.Executor{
		Container: this.container,
		Payload:   &payload,
		ResponseSocket: &model.ResponseSocket{
			MessageType: messageType,
			Conn:        conn,
		},
	}
	executor.Encode()
}

func (this *WsMainController) DecodeAction(w http.ResponseWriter, r *http.Request) {
	var wsComponent = this.container.Get(container.WS_COMPONENT).(*component.WS)
	conn, err := wsComponent.GetConnection(w, r)
	if err != nil {
		return
	}
	messageType, p, err := wsComponent.Read(conn)
	if err != nil {
		return
	}

	var payload form.Payload
	if err := json.Unmarshal(p, &payload); err != nil {
		return
	}

	executor := service.Executor{
		Container: this.container,
		Payload:   &payload,
		ResponseSocket: &model.ResponseSocket{
			MessageType: messageType,
			Conn:        conn,
		},
	}
	executor.Decode()
}

func (this *WsMainController) DeleteAction(w http.ResponseWriter, r *http.Request) {
	var wsComponent = this.container.Get(container.WS_COMPONENT).(*component.WS)
	conn, err := wsComponent.GetConnection(w, r)
	if err != nil {
		return
	}
	messageType, p, err := wsComponent.Read(conn)
	if err != nil {
		return
	}

	var payload form.Payload
	if err := json.Unmarshal(p, &payload); err != nil {
		return
	}

	executor := service.Executor{
		Container: this.container,
		Payload:   &payload,
		ResponseSocket: &model.ResponseSocket{
			MessageType: messageType,
			Conn:        conn,
		},
	}
	executor.Delete()
}

func (this *WsMainController) CopyAction(w http.ResponseWriter, r *http.Request) {
	var wsComponent = this.container.Get(container.WS_COMPONENT).(*component.WS)
	conn, err := wsComponent.GetConnection(w, r)
	if err != nil {
		return
	}
	messageType, p, err := wsComponent.Read(conn)
	if err != nil {
		return
	}

	var payload form.Payload
	if err := json.Unmarshal(p, &payload); err != nil {
		return
	}

	executor := service.Executor{
		Container: this.container,
		Payload:   &payload,
		ResponseSocket: &model.ResponseSocket{
			MessageType: messageType,
			Conn:        conn,
		},
	}
	executor.Copy()
}
