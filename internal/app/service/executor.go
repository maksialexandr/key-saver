package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"person-key-saver/internal/app/builder"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/form"
	"person-key-saver/internal/app/model"
)

type Executor struct {
	Container      *container.Container
	Payload        *form.Payload
	ResponseSocket *model.ResponseSocket
}

func (this *Executor) Delete() {
	kService := this.Container.Get(container.KEY_SERVICE).(*KeyService)
	dataBuilder := builder.NewDeleteDataBuilder(this.Container, this.Payload)
	data := dataBuilder.Build()

	keyMap := map[string]model.PersonKey{}
	for i := range data.PhysicalKeys {
		keyMap[data.PhysicalKeys[i].Value] = data.PhysicalKeys[i]
	}

	for i := range data.Keys {
		_, exist := keyMap[data.Keys[i].Value]
		result, err := kService.Delete(data.Device, data.Keys[i], !exist)
		if this.ResponseSocket != nil {
			this.sendSocket(this.ResponseSocket.Conn, this.ResponseSocket.MessageType, &model.PersonKeyResponse{
				Value:  data.Keys[i].Value,
				Result: result && err == nil,
			})
		}
	}
	kService.DeleteAllKeys(data.Device)
	if this.ResponseSocket != nil {
		this.sendSocket(this.ResponseSocket.Conn, this.ResponseSocket.MessageType, &model.PersonKeyResponse{Complete: true})
	}
}

func (this *Executor) Copy() {
	kService := this.Container.Get(container.KEY_SERVICE).(*KeyService)
	dataBuilder := builder.NewCopyDataBuilder(this.Container, this.Payload)
	data := dataBuilder.Build()

	for i := range data.Keys {
		result, err := kService.Load(data.Device, data.Keys[i])
		if this.ResponseSocket != nil {
			this.sendSocket(this.ResponseSocket.Conn, this.ResponseSocket.MessageType, &model.PersonKeyResponse{
				Value:  data.Keys[i].Value,
				Result: result && err == nil,
			})
		}
	}
	if this.ResponseSocket != nil {
		this.sendSocket(this.ResponseSocket.Conn, this.ResponseSocket.MessageType, &model.PersonKeyResponse{Complete: true})
	}
}

func (this *Executor) Decode() {
	kService := this.Container.Get(container.KEY_SERVICE).(*KeyService)
	dataBuilder := builder.NewDecodeDataBuilder(this.Container, this.Payload)
	data := dataBuilder.Build()

	for i := range data.Keys {
		result, err := kService.Decode(*data.Device, data.Keys[i])
		if this.ResponseSocket != nil {
			this.sendSocket(this.ResponseSocket.Conn, this.ResponseSocket.MessageType, &model.PersonKeyResponse{
				Value:  data.Keys[i].Value,
				Result: result && err == nil,
			})
		}
	}
	if this.ResponseSocket != nil {
		this.sendSocket(this.ResponseSocket.Conn, this.ResponseSocket.MessageType, &model.PersonKeyResponse{Complete: true})
	}
}

func (this *Executor) Encode() {
	kService := this.Container.Get(container.KEY_SERVICE).(*KeyService)
	dataBuilder := builder.NewEncodeDataBuilder(this.Container, this.Payload)
	data := dataBuilder.Build()
	// временный костыль, изза старых прошивок
	for i := range data.Keys {
		cipherId := data.Keys[i].CipherId
		data.Keys[i].UnsetCipher()
		result, err := kService.Delete(data.Device, data.Keys[i], false)
		if result {
			data.Keys[i].CipherId = cipherId
			result, err = kService.Load(data.Device, data.Keys[i])
		}

		//result, err := kService.Encode(*data.Device, data.Keys[i])
		if this.ResponseSocket != nil {
			this.sendSocket(this.ResponseSocket.Conn, this.ResponseSocket.MessageType, &model.PersonKeyResponse{
				Value:  data.Keys[i].Value,
				Result: result && err == nil,
			})
		}
	}
	if this.ResponseSocket != nil {
		this.sendSocket(this.ResponseSocket.Conn, this.ResponseSocket.MessageType, &model.PersonKeyResponse{Complete: true})
	}
}

func (this *Executor) Save() {
	dataBuilder := builder.NewSaveDataBuilder(this.Container, this.Payload)
	kService := this.Container.Get(container.KEY_SERVICE).(*KeyService)
	data := dataBuilder.Build()

	if this.Payload.Delete {
		kService.DeleteAllKeys(data.Device)
	}

	for i := range data.Keys {
		result, err := kService.Load(data.Device, data.Keys[i])
		if this.ResponseSocket != nil {
			this.sendSocket(this.ResponseSocket.Conn, this.ResponseSocket.MessageType, &model.PersonKeyResponse{
				Value:  data.Keys[i].Value,
				Result: result && err == nil,
			})
		}
	}
	if this.ResponseSocket != nil {
		this.sendSocket(this.ResponseSocket.Conn, this.ResponseSocket.MessageType, &model.PersonKeyResponse{Complete: true})
	}
}

func (this *Executor) sendSocket(conn *websocket.Conn, messageType int, msg *model.PersonKeyResponse) {
	str, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	err = conn.WriteMessage(messageType, str)
	if err != nil {
		log.Println(err)
	}
}
