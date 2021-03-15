package builder

import (
	"github.com/sirupsen/logrus"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/form"
	"person-key-saver/internal/app/repository"
)

type CopyDataBuilder struct {
	BaseDate
}

func NewCopyDataBuilder(con *container.Container, payload *form.Payload) *CopyDataBuilder {
	return &CopyDataBuilder{
		BaseDate: BaseDate{
			container: con,
			payload:   payload,
		},
	}
}

func (this *CopyDataBuilder) Build() *BaseDate {
	deviceRepository := this.container.Get(container.DEVICE_REPOSITORY).(*repository.DeviceRepository)
	keyRepository := this.container.Get(container.KEY_REPOSITORY).(*repository.PersonKeyRepository)
	logger := this.container.Get(container.LOGGER).(*logrus.Logger)

	deviceSrc, err := deviceRepository.GetItem(this.payload.SrcMac)
	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return nil
	}

	deviceDst, err := deviceRepository.GetItem(this.payload.DstMac)
	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return nil
	}

	pKeys, err := keyRepository.GetItems(deviceSrc.DeviceId)
	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return nil
	}

	for i := range pKeys {
		key := pKeys[i]
		key.UnsetKeyId()
		key.DeviceId = deviceDst.DeviceId
		this.Keys = append(this.Keys, key)
	}
	this.Device = deviceDst
	return &this.BaseDate
}
