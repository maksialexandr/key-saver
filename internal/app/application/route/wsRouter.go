package route

import (
	"github.com/gorilla/mux"
	"net/http"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/controllers"
)

type WSRouter struct {
	Router     *mux.Router
	controller *controllers.WsMainController
}

func NewWSRouter() *WSRouter {
	return &WSRouter{}
}

func (this *WSRouter) Configure(container *container.Container) {
	this.controller = controllers.NewWsMainController(container)

	this.Router = mux.NewRouter()

	this.Router.HandleFunc("/keys/save", func(w http.ResponseWriter, r *http.Request) {
		this.controller.SaveAction(w, r)
	})

	this.Router.HandleFunc("/keys/encode", func(w http.ResponseWriter, r *http.Request) {
		this.controller.EncodeAction(w, r)
	})

	this.Router.HandleFunc("/keys/decode", func(w http.ResponseWriter, r *http.Request) {
		this.controller.DecodeAction(w, r)
	})

	this.Router.HandleFunc("/keys/delete", func(w http.ResponseWriter, r *http.Request) {
		this.controller.DeleteAction(w, r)
	})

	this.Router.HandleFunc("/keys/copy", func(w http.ResponseWriter, r *http.Request) {
		this.controller.CopyAction(w, r)
	})
}
