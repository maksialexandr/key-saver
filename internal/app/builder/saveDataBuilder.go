package builder

import (
	"github.com/sirupsen/logrus"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/form"
	"person-key-saver/internal/app/repository"
)

type SaveDataBuilder struct {
	BaseDate
}

func NewSaveDataBuilder(con *container.Container, payload *form.Payload) *SaveDataBuilder {
	return &SaveDataBuilder{
		BaseDate: BaseDate{
			container: con,
			payload:   payload,
		},
	}
}

func (this *SaveDataBuilder) Build() *BaseDate {
	deviceRepository := this.container.Get(container.DEVICE_REPOSITORY).(*repository.DeviceRepository)
	logger := this.container.Get(container.LOGGER).(*logrus.Logger)

	device, err := deviceRepository.GetItem(this.payload.SrcMac)
	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return nil
	}

	for i := range this.payload.Keys {
		if device.SecurityMode != 0 {
			this.payload.Keys[i].CipherId = device.CipherId
		}
		this.payload.Keys[i].BuyerId = device.BuyerId
		this.payload.Keys[i].DeviceId = device.DeviceId
		this.payload.Keys[i].Authentic = true
	}
	this.Keys = this.payload.Keys
	this.Device = device
	return &this.BaseDate
}
