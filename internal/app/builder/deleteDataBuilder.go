package builder

import (
	"github.com/sirupsen/logrus"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/form"
	"person-key-saver/internal/app/model"
	"person-key-saver/internal/app/repository"
	"person-key-saver/internal/app/repository/api"
)

type DeleteDataBuilder struct {
	BaseDate
}

func NewDeleteDataBuilder(con *container.Container, payload *form.Payload) *DeleteDataBuilder {
	return &DeleteDataBuilder{
		BaseDate: BaseDate{
			container: con,
			payload:   payload,
		},
	}
}

func (this *DeleteDataBuilder) Build() *BaseDate {
	deviceRepository := this.container.Get(container.DEVICE_REPOSITORY).(*repository.DeviceRepository)
	deviceApiRepository := this.container.Get(container.DEVICE_API_REPOSITORY).(*api.DeviceRepository)
	keyRepository := this.container.Get(container.KEY_REPOSITORY).(*repository.PersonKeyRepository)
	logger := this.container.Get(container.LOGGER).(*logrus.Logger)

	device, err := deviceRepository.GetItem(this.payload.SrcMac)
	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return nil
	}

	this.PhysicalKeys, err = deviceApiRepository.GetKeys(device)
	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return nil
	}

	pKeys, err := keyRepository.GetItems(device.DeviceId)
	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return nil
	}

	mapKeys := make(map[string]model.PersonKey)
	for j := range pKeys {
		mapKeys[pKeys[j].Value] = pKeys[j]
	}

	if !this.payload.IsEmptyKeys() {
		for i := range this.payload.Keys {
			if this.payload.Keys[i].Value == mapKeys[this.payload.Keys[i].Value].Value {
				this.Keys = append(this.Keys, this.payload.Keys[i])
			}
		}
	} else {
		this.Keys = pKeys
	}

	this.Device = device
	return &this.BaseDate
}
