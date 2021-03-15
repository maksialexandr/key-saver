package app

// тут хочу переделать, не использовать лучше
import (
	"github.com/sirupsen/logrus"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/form"
	"person-key-saver/internal/app/model"
	"person-key-saver/internal/app/repository"
)

type DataBuilder struct {
	container *container.Container
	payload   *form.Payload
	data      *model.DataExecute
}

func (this *DataBuilder) SetContainer(container *container.Container) {
	this.container = container
}

func (this *DataBuilder) setPayload(payload *form.Payload) {
	this.payload = payload
}

func (this *DataBuilder) setData(data *model.DataExecute) {
	this.data = data
}

func (this *DataBuilder) Build(payload *form.Payload) *model.DataExecute {
	this.setPayload(payload)
	this.setData(&model.DataExecute{})
	this.preBuild()
	this.data.CurrentKeys = this.prepareCurrent()
	return this.data
}

// подготовка ключей src и dst
func (this *DataBuilder) preBuild() {
	if !this.payload.IsEmptySourceMac() {
		this.data.SrcDevice, this.data.SrcKeys = this.prepareSource(this.payload.SrcMac)
	}
	if !this.payload.IsEmptyDestinationMac() {
		this.data.DstDevice, this.data.DstKeys = this.prepareDestination(this.payload.DstMac)
	}
}

func (this *DataBuilder) prepareSource(mac string) (*model.Device, map[string]model.PersonKey) {
	deviceRepository := this.container.Get(container.DEVICE_REPOSITORY).(*repository.DeviceRepository)
	keyRepository := this.container.Get(container.KEY_REPOSITORY).(*repository.PersonKeyRepository)
	logger := this.container.Get(container.LOGGER).(*logrus.Logger)
	keys := make(map[string]model.PersonKey)
	srcDevice, err := deviceRepository.GetItem(mac)

	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return srcDevice, keys
	}
	pKeys, err := keyRepository.GetItems(srcDevice.DeviceId)

	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return srcDevice, keys
	}

	for j := range pKeys {
		keys[pKeys[j].Value] = pKeys[j]
	}
	return srcDevice, keys
}

func (this *DataBuilder) prepareDestination(dst string) (*model.Device, map[string]model.PersonKey) {
	deviceRepository := this.container.Get(container.DEVICE_REPOSITORY).(*repository.DeviceRepository)
	keyRepository := this.container.Get(container.KEY_REPOSITORY).(*repository.PersonKeyRepository)
	logger := this.container.Get(container.LOGGER).(*logrus.Logger)
	keys := make(map[string]model.PersonKey)

	dstDevice, err := deviceRepository.GetItem(dst)
	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return dstDevice, keys
	}
	dstKeys, err := keyRepository.GetItems(dstDevice.DeviceId)

	if err != nil {
		logger.Error("Nil pointer or error " + err.Error())
		return dstDevice, keys
	}

	for j := range this.data.SrcKeys {
		key := this.data.SrcKeys[j]
		key.DeviceId = dstDevice.DeviceId
		key.UnsetKeyId()
		keys[j] = key
	}

	for j := range dstKeys {
		if dstKeys[j].Value == keys[dstKeys[j].Value].Value {
			key := keys[dstKeys[j].Value]
			key.KeyId = dstKeys[j].KeyId
			keys[dstKeys[j].Value] = key
		}
	}
	return dstDevice, keys
}

// подготовка ключей которыми будем оперировать
func (this *DataBuilder) prepareCurrent() map[string]model.PersonKey {
	result := make(map[string]model.PersonKey)
	var device *model.Device
	var keyMap map[string]model.PersonKey
	if this.data.DstDevice != nil {
		device = this.data.DstDevice
		keyMap = this.data.DstKeys
	} else {
		device = this.data.SrcDevice
		keyMap = this.data.SrcKeys
	}

	if !this.payload.IsEmptyKeys() {
		for i := range this.payload.Keys {
			keyId := 0
			if this.payload.Keys[i].Value == keyMap[this.payload.Keys[i].Value].Value {
				keyId = keyMap[this.payload.Keys[i].Value].KeyId
			}

			if device.SecurityMode != 0 {
				this.payload.Keys[i].CipherId = device.CipherId
			}
			this.payload.Keys[i].BuyerId = device.BuyerId
			this.payload.Keys[i].DeviceId = device.DeviceId
			this.payload.Keys[i].KeyId = keyId
			this.payload.Keys[i].Authentic = true
			result[this.payload.Keys[i].Value] = this.payload.Keys[i]
		}
	} else {
		for i := range keyMap {
			key := keyMap[i]
			key.DeviceId = device.DeviceId
			key.Authentic = true
			if device.SecurityMode != 0 {
				key.CipherId = device.CipherId
			}
			key.BuyerId = device.BuyerId
			result[i] = key
		}
	}
	return result
}
