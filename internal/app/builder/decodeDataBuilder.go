package builder

import (
	"github.com/sirupsen/logrus"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/form"
	"person-key-saver/internal/app/model"
	"person-key-saver/internal/app/repository"
)

type DecodeDataBuilder struct {
	BaseDate
}

func NewDecodeDataBuilder(con *container.Container, payload *form.Payload) *DecodeDataBuilder {
	return &DecodeDataBuilder{
		BaseDate: BaseDate{
			container: con,
			payload:   payload,
		},
	}
}

func (this *DecodeDataBuilder) Build() *BaseDate {
	deviceRepository := this.container.Get(container.DEVICE_REPOSITORY).(*repository.DeviceRepository)
	keyRepository := this.container.Get(container.KEY_REPOSITORY).(*repository.PersonKeyRepository)
	logger := this.container.Get(container.LOGGER).(*logrus.Logger)

	device, err := deviceRepository.GetItem(this.payload.SrcMac)
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
		if pKeys[j].CipherId != 0 {
			mapKeys[pKeys[j].Value] = pKeys[j]
		}
	}

	if !this.payload.IsEmptyKeys() {
		for i := range this.payload.Keys {
			if this.payload.Keys[i].Value == mapKeys[this.payload.Keys[i].Value].Value {
				this.payload.Keys[i].UnsetCipher()
				this.Keys = append(this.Keys, this.payload.Keys[i])
			}
		}
	} else {
		for _, value := range mapKeys {
			value.UnsetCipher()
			this.Keys = append(this.Keys, value)
		}
	}

	this.Device = device
	return &this.BaseDate
}
