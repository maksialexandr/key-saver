package service

import (
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/model"
	"person-key-saver/internal/app/repository"
	"person-key-saver/internal/app/repository/api"
	"strconv"
)

type KeyService struct {
	container *container.Container
}

func NewKeyService(container *container.Container) *KeyService {
	return &KeyService{
		container: container,
	}
}

func (this *KeyService) Load(device *model.Device, key model.PersonKey) (bool, error) {
	key.PrepareValue()
	deviceApiRepository := this.container.Get(container.DEVICE_API_REPOSITORY).(*api.DeviceRepository)
	kRep := this.container.Get(container.KEY_REPOSITORY).(*repository.PersonKeyRepository)

	prepareKey := this.prepareHardParams(&model.IntercomSettings{
		Key:         key.Value,
		Apartment:   key.PanelCode,
		CipherIndex: key.CipherId,
	})

	if key.KeyId != 0 {
		_, err := deviceApiRepository.DeleteKey(device, prepareKey)
		if err = kRep.Delete(&key); err != nil {
			return false, err
		}
	}

	//logger := this.container.Get(container.LOGGER).(*logrus.Logger)

	res, err := deviceApiRepository.AddKey(device, prepareKey)
	if res {
		if kRep.Save(&key) != nil {
			return false, err
		}
		return true, nil
	}
	return false, err
}

func (this *KeyService) Encode(device model.Device, key model.PersonKey) (bool, error) {
	key.PrepareValue()
	deviceApiRepository := this.container.Get(container.DEVICE_API_REPOSITORY).(*api.DeviceRepository)
	keyRepository := this.container.Get(container.KEY_REPOSITORY).(*repository.PersonKeyRepository)
	//logger := this.container.Get(container.LOGGER).(*logrus.Logger)

	res, err := deviceApiRepository.EncodeKey(device, this.prepareHardParams(&model.IntercomSettings{
		Key:         key.Value,
		Apartment:   key.PanelCode,
		CipherIndex: key.CipherId,
	}))

	if res {
		if keyRepository.Save(&key) != nil {
			return false, err
		}
		return true, nil
	}
	return false, err
}

func (this *KeyService) Decode(device model.Device, key model.PersonKey) (bool, error) {
	key.PrepareValue()
	deviceApiRepository := this.container.Get(container.DEVICE_API_REPOSITORY).(*api.DeviceRepository)
	kRep := this.container.Get(container.KEY_REPOSITORY).(*repository.PersonKeyRepository)

	res, err := deviceApiRepository.DecodeKey(device, this.prepareHardParams(&model.IntercomSettings{
		Key:         key.Value,
		Apartment:   key.PanelCode,
		CipherIndex: key.CipherId,
	}))

	if res {
		key.UnsetCipher()
		if kRep.Save(&key) != nil {
			return false, err
		}
		return true, nil
	}
	return false, err
}

func (this *KeyService) Delete(device *model.Device, key model.PersonKey, onlyDb bool) (bool, error) {
	key.PrepareValue()
	deviceApiRepository := this.container.Get(container.DEVICE_API_REPOSITORY).(*api.DeviceRepository)
	kRep := this.container.Get(container.KEY_REPOSITORY).(*repository.PersonKeyRepository)

	if !onlyDb {
		res, err := deviceApiRepository.DeleteKey(device, this.prepareHardParams(&model.IntercomSettings{
			Key:         key.Value,
			CipherIndex: key.CipherId,
		}))
		if !res || err != nil {
			return false, err
		}
	}

	if err := kRep.Delete(&key); err != nil {
		return false, err
	}
	return true, nil
}

func (this *KeyService) DeleteAllKeys(device *model.Device) {
	deviceApiRepository := this.container.Get(container.DEVICE_API_REPOSITORY).(*api.DeviceRepository)
	keyRepository := this.container.Get(container.KEY_REPOSITORY).(*repository.PersonKeyRepository)

	pKeys, err := keyRepository.GetItems(device.DeviceId)
	if err != nil {
		return
	}
	for i := range pKeys {
		_, _ = this.Delete(device, pKeys[i], false)
	}

	keys, err := deviceApiRepository.GetKeys(device)
	if err != nil {
		return
	}
	for i := range keys {
		_, _ = this.Delete(device, keys[i], false)
	}
}

func (this *KeyService) Copy(dstDevice *model.Device, key model.PersonKey) (bool, error) {
	key.PrepareValue()
	deviceApiRepository := this.container.Get(container.DEVICE_API_REPOSITORY).(*api.DeviceRepository)
	kRep := this.container.Get(container.KEY_REPOSITORY).(*repository.PersonKeyRepository)

	res, err := deviceApiRepository.AddKey(dstDevice, this.prepareHardParams(&model.IntercomSettings{
		Key:         key.Value,
		Apartment:   key.PanelCode,
		CipherIndex: key.CipherId,
	}))

	if res {
		if kRep.Save(&key) != nil {
			return false, err
		}
		return true, nil
	}
	return false, err
}

func (w *KeyService) prepareHardParams(setting *model.IntercomSettings) map[string]string {
	params := map[string]string{}
	params["Key"] = setting.Key

	if setting.Apartment != 0 {
		params["Apartment"] = strconv.Itoa(setting.Apartment)
	}

	if setting.CipherIndex != 0 {
		params["CipherIndex"] = strconv.Itoa(setting.CipherIndex)
	}
	return params
}
