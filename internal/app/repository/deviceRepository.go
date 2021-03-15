package repository

import (
	"database/sql"
	"person-key-saver/internal/app/model"
)

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{
		db: db,
	}
}

func (this *DeviceRepository) GetItem(mac string) (*model.Device, error) {
	device := new(model.Device)
	if err := this.db.QueryRow("SELECT DISTINCT d.DEVICE_ID AS deviceId, d.MAC, d.TYPE, "+
		"COALESCE(url.VALUE_CHAR, '') AS url, COALESCE(login.VALUE_CHAR, '') AS login, "+
		"COALESCE(password.VALUE_CHAR, '') AS password, "+
		"COALESCE(securityMode.VALUE_NUMBER, 0) AS securityMode, "+
		"COALESCE(dtc.CIPHER_ID, 0) AS cipherId, d.BUYER_ID as buyerId FROM td.device d "+
		"LEFT JOIN td.deviceOption url ON url.DEVICE_ID = d.DEVICE_ID AND url.OPTION_NAME = 'url' "+
		"LEFT JOIN td.deviceOption login ON login.DEVICE_ID = d.DEVICE_ID AND login.OPTION_NAME = 'login' "+
		"LEFT JOIN td.deviceOption password ON password.DEVICE_ID = d.DEVICE_ID AND password.OPTION_NAME = 'password' "+
		"LEFT JOIN td.deviceOption securityMode ON securityMode.DEVICE_ID = d.DEVICE_ID AND securityMode.OPTION_NAME = 'isMifare' "+
		"LEFT JOIN td.deviceToCipher dtc ON dtc.DEVICE_ID = d.DEVICE_ID "+
		"WHERE  d.MAC = ?",
		mac).Scan(&device.DeviceId, &device.Mac, &device.Type, &device.Url,
		&device.Login, &device.Password, &device.SecurityMode, &device.CipherId, &device.BuyerId); err != nil {
		return nil, err
	}
	return device, nil
}

func (this *DeviceRepository) GetItems() ([]*model.Device, error) {
	var pKeys []*model.Device
	rows, err := this.db.Query("SELECT DISTINCT d.DEVICE_ID AS deviceId, d.MAC, d.TYPE, " +
		"COALESCE(url.VALUE_CHAR, '') AS url, COALESCE(login.VALUE_CHAR, '') AS login, " +
		"COALESCE(password.VALUE_CHAR, '') AS password, " +
		"COALESCE(securityMode.VALUE_NUMBER, 0) AS securityMode, " +
		"COALESCE(dtc.CIPHER_ID, 0) AS cipherId, d.BUYER_ID as buyerId FROM td.device d " +
		"LEFT JOIN td.deviceOption url ON url.DEVICE_ID = d.DEVICE_ID AND url.OPTION_NAME = 'url' " +
		"LEFT JOIN td.deviceOption login ON login.DEVICE_ID = d.DEVICE_ID AND login.OPTION_NAME = 'login' " +
		"LEFT JOIN td.deviceOption password ON password.DEVICE_ID = d.DEVICE_ID AND password.OPTION_NAME = 'password' " +
		"LEFT JOIN td.deviceOption securityMode ON securityMode.DEVICE_ID = d.DEVICE_ID AND securityMode.OPTION_NAME = 'isMifare' " +
		"LEFT JOIN td.deviceToCipher dtc ON dtc.DEVICE_ID = d.DEVICE_ID WHERE d.TYPE = 2 and d.STATUS_CODE = 2")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		pKey := new(model.Device)
		err = rows.Scan(&pKey.DeviceId,
			&pKey.Mac,
			&pKey.Type,
			&pKey.Url,
			&pKey.Login,
			&pKey.Password,
			&pKey.SecurityMode,
			&pKey.CipherId,
			&pKey.BuyerId)
		if err != nil {
			return nil, err
		}
		pKeys = append(pKeys, pKey)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return pKeys, nil
}
