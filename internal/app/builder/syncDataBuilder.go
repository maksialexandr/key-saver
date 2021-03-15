package builder

import (
	"github.com/sirupsen/logrus"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/model"
	"person-key-saver/internal/app/repository"
	"person-key-saver/internal/app/repository/api"
)

type SyncDataBuilder struct {
	BaseDate
	mac string
}

func NewSyncDataBuilder(con *container.Container, mac string) *SyncDataBuilder {
	return &SyncDataBuilder{
		BaseDate: BaseDate{
			container: con,
		},
		mac: mac,
	}
}

func (this *SyncDataBuilder) Build() *BaseDate {
	deviceRepository := this.container.Get(container.DEVICE_REPOSITORY).(*repository.DeviceRepository)
	deviceApiRepository := this.container.Get(container.DEVICE_API_REPOSITORY).(*api.DeviceRepository)
	keyRepository := this.container.Get(container.KEY_REPOSITORY).(*repository.PersonKeyRepository)
	logger := this.container.Get(container.LOGGER).(*logrus.Logger)

	device, err := deviceRepository.GetItem(this.mac)
	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return nil
	}

	pKeys, err := keyRepository.GetItems(device.DeviceId)
	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return nil
	}

	dKeys, err := deviceApiRepository.GetKeys(device)
	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return nil
	}

	mapKeys := map[string]model.PersonKey{}
	for i := range pKeys {
		mapKeys[pKeys[i].Value] = pKeys[i]
	}
	pKeys = nil

	for i := range dKeys {
		_, exist := mapKeys[dKeys[i].Value]
		if !exist {
			dKeys[i].DeviceId = device.DeviceId
			dKeys[i].BuyerId = device.BuyerId
			this.Keys = append(this.Keys, dKeys[i])
		}
	}
	mapKeys = nil
	dKeys = nil

	this.Device = device
	return &this.BaseDate
}
