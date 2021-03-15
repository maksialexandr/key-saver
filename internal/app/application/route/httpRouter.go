package route

import (
	"github.com/gorilla/mux"
	"net/http"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/controllers"
)

type HttpRouter struct {
	Router     *mux.Router
	controller *controllers.HttpMainController
}

func NewHttpRouter() *HttpRouter {
	return &HttpRouter{}
}

func (this *HttpRouter) Configure(container *container.Container) {
	this.controller = controllers.NewHttpMainController(container)

	this.Router = mux.NewRouter()

	this.Router.HandleFunc("/key/save", func(w http.ResponseWriter, r *http.Request) {
		this.controller.SaveAction(w, r)
	}).Methods("POST")

	this.Router.HandleFunc("/key/encode", func(w http.ResponseWriter, r *http.Request) {
		this.controller.EncodeAction(w, r)
	}).Methods("PUT")

	this.Router.HandleFunc("/key/decode", func(w http.ResponseWriter, r *http.Request) {
		this.controller.DecodeAction(w, r)
	}).Methods("PUT")

	this.Router.HandleFunc("/key/delete", func(w http.ResponseWriter, r *http.Request) {
		this.controller.DeleteAction(w, r)
	}).Methods("DELETE")

	this.Router.HandleFunc("/key/copy", func(w http.ResponseWriter, r *http.Request) {
		this.controller.CopyAction(w, r)
	}).Methods("POST")
}
