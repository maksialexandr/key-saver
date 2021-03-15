package api

import (
	"*****/intercom/intercom-go-lib"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"key-saver/internal/app/model"
	"strconv"
	"strings"
	"time"
)

const (
	TypeIMD1B = iota
	TypeSputnik
	TypeBeward
	TypeBewardSd06m
)

// Сервис для работы с устройствами
type DeviceRepository struct {
	httpClient intercom.ClientInterface
	mqttClient mqtt.Client
}

func NewDeviceRepository(httpClient intercom.ClientInterface, mqttClient mqtt.Client) *DeviceRepository {
	return &DeviceRepository{
		httpClient: httpClient,
		mqttClient: mqttClient,
	}
}

func (s *DeviceRepository) AddKey(device *model.Device, settings map[string]string) (bool, error) {
	apiDevice := intercom.CreateDevice(
		intercom.DeviceType(device.Type), device.Mac, device.Url, device.Login, device.Password, s.httpClient, s.mqttClient,
	)

	if device.Type == TypeIMD1B {
		s.sleep(1000)
	}

	if apiDevice != nil {
		if settings["CipherIndex"] != "" {
			return apiDevice.AddMifareKey(settings)
		} else {
			return apiDevice.AddRfidKey(settings)
		}
	}
	return true, nil
}

func (s *DeviceRepository) DeleteKey(device *model.Device, settings map[string]string) (bool, error) {
	apiDevice := intercom.CreateDevice(
		intercom.DeviceType(device.Type), device.Mac, device.Url, device.Login, device.Password, s.httpClient, s.mqttClient,
	)

	if device.Type == TypeIMD1B {
		s.sleep(1000)
	}

	if apiDevice != nil {
		if settings["CipherIndex"] != "" {
			return apiDevice.DeleteMifareKey(settings)
		} else {
			return apiDevice.DeleteRfidKey(settings)
		}
	}
	return true, nil
}

func (s *DeviceRepository) EncodeKey(device model.Device, settings map[string]string) (bool, error) {
	apiDevice := intercom.CreateDevice(
		intercom.DeviceType(device.Type), device.Mac, device.Url, device.Login, device.Password, s.httpClient, s.mqttClient,
	)

	settings["Destination"] = "0"
	if apiDevice != nil {
		return apiDevice.EncodeKey(settings)
	}
	return true, nil
}

func (s *DeviceRepository) DecodeKey(device model.Device, settings map[string]string) (bool, error) {
	apiDevice := intercom.CreateDevice(
		intercom.DeviceType(device.Type), device.Mac, device.Url, device.Login, device.Password, s.httpClient, s.mqttClient,
	)

	if apiDevice != nil {
		return apiDevice.DecodeKey(settings)
	}
	return true, nil
}

func (s *DeviceRepository) GetKeys(device *model.Device) ([]model.PersonKey, error) {
	var pKeys []model.PersonKey

	apiDevice := intercom.CreateDevice(
		intercom.DeviceType(device.Type), device.Mac, device.Url, device.Login, device.Password, s.httpClient, s.mqttClient,
	)
	data, _ := apiDevice.GetRfidKeys()
	params := map[string]string{}
	splitFunc := func(r rune) bool {
		return strings.ContainsRune("=", r)
	}
	words := strings.Fields(string(data))
	for _, word := range words {
		words := strings.FieldsFunc(word, splitFunc)
		if len(words) == 2 {
			params[words[0]] = words[1]
		}
	}

	i := 1
	for {
		if value, ok := params["KeyValue"+strconv.Itoa(i)]; ok {
			panelCode, _ := strconv.Atoi(params["KeyApartment"+strconv.Itoa(i)])
			pKey := model.PersonKey{
				Value:     value,
				PanelCode: panelCode,
			}
			pKeys = append(pKeys, pKey)
			i++
		} else {
			break
		}
	}

	data, _ = apiDevice.GetMifareKeys()
	params = map[string]string{}
	words = strings.Fields(string(data))
	for _, word := range words {
		words := strings.FieldsFunc(word, splitFunc)
		if len(words) == 2 {
			params[words[0]] = words[1]
		}
	}

	i = 1
	for {
		if value, ok := params["Key"+strconv.Itoa(i)]; ok {
			panelCode, _ := strconv.Atoi(params["Apartment"+strconv.Itoa(i)])
			cipherId, _ := strconv.Atoi(params["CipherIndex"+strconv.Itoa(i)])
			pKey := model.PersonKey{
				Value:     value,
				PanelCode: panelCode,
				CipherId:  cipherId,
			}
			pKeys = append(pKeys, pKey)
			i++
		} else {
			break
		}
	}
	return pKeys, nil
}

// Рандомный слип в интервале
// Для того чтобы плата Imd1b нормально работала, нужно передышки между запросами
func (d *DeviceRepository) sleep(millisecond int) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(millisecond)
	time.Sleep(time.Duration(n) * time.Millisecond)
}
