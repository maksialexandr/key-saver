package controllers

import (
	"encoding/json"
	"net/http"
	"person-key-saver/internal/app"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/form"
	"person-key-saver/internal/app/service"
)

type HttpMainController struct {
	container *container.Container
}

func NewHttpMainController(container *container.Container) *HttpMainController {
	return &HttpMainController{
		container: container,
	}
}

func (this *HttpMainController) SaveAction(w http.ResponseWriter, r *http.Request) {
	kService := this.container.Get(container.KEY_SERVICE).(*service.KeyService)
	var httpPayload form.HttpPayload
	err := json.NewDecoder(r.Body).Decode(&httpPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var dataBuilder app.DataBuilder
	dataBuilder.SetContainer(this.container)
	dataExecute := dataBuilder.Build(httpPayload.GetPayload())

	for i := range dataExecute.CurrentKeys {
		result, err := kService.Load(dataExecute.SrcDevice, dataExecute.CurrentKeys[i])
		if !result || err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func (this *HttpMainController) EncodeAction(w http.ResponseWriter, r *http.Request) {
	kService := this.container.Get(container.KEY_SERVICE).(*service.KeyService)
	var httpPayload form.HttpPayload
	err := json.NewDecoder(r.Body).Decode(&httpPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var dataBuilder app.DataBuilder
	dataBuilder.SetContainer(this.container)
	dataExecute := dataBuilder.Build(httpPayload.GetPayload())

	for i := range dataExecute.CurrentKeys {
		result, err := kService.Encode(*dataExecute.SrcDevice, dataExecute.CurrentKeys[i])
		if !result || err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusAccepted)
	return
}

func (this *HttpMainController) DecodeAction(w http.ResponseWriter, r *http.Request) {
	kService := this.container.Get(container.KEY_SERVICE).(*service.KeyService)
	var httpPayload form.HttpPayload
	err := json.NewDecoder(r.Body).Decode(&httpPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var dataBuilder app.DataBuilder
	dataBuilder.SetContainer(this.container)
	dataExecute := dataBuilder.Build(httpPayload.GetPayload())

	for i := range dataExecute.CurrentKeys {
		result, err := kService.Decode(*dataExecute.SrcDevice, dataExecute.CurrentKeys[i])
		if !result || err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (this *HttpMainController) DeleteAction(w http.ResponseWriter, r *http.Request) {
	kService := this.container.Get(container.KEY_SERVICE).(*service.KeyService)
	var httpPayload form.HttpPayload
	err := json.NewDecoder(r.Body).Decode(&httpPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var dataBuilder app.DataBuilder
	dataBuilder.SetContainer(this.container)
	dataExecute := dataBuilder.Build(httpPayload.GetPayload())

	for i := range dataExecute.CurrentKeys {
		result, err := kService.Delete(dataExecute.SrcDevice, dataExecute.CurrentKeys[i], false)
		if !result || err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	return
}

func (this *HttpMainController) CopyAction(w http.ResponseWriter, r *http.Request) {
	kService := this.container.Get(container.KEY_SERVICE).(*service.KeyService)
	var httpPayload form.HttpPayload
	err := json.NewDecoder(r.Body).Decode(&httpPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var dataBuilder app.DataBuilder
	dataBuilder.SetContainer(this.container)
	dataExecute := dataBuilder.Build(httpPayload.GetPayload())

	for i := range dataExecute.CurrentKeys {
		result, err := kService.Load(dataExecute.DstDevice, dataExecute.CurrentKeys[i])
		if !result || err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	return
}
