package repository

import (
	"database/sql"
	"person-key-saver/internal/app/model"
)

type PersonKeyRepository struct {
	db *sql.DB
}

func NewPersonKeyRepository(db *sql.DB) *PersonKeyRepository {
	return &PersonKeyRepository{
		db: db,
	}
}

func (this *PersonKeyRepository) GetItems(deviceId int) ([]model.PersonKey, error) {
	var pKeys []model.PersonKey
	rows, err := this.db.Query("SELECT pk.KEY_ID as keyId, pk.DEVICE_ID as deviceId, ifnull(pk.FLAT_NUM, 0) as flatNum, pk.VALUE as value, "+
		"pk.BUYER_ID as buyerId, ifnull(pk.AUTHENTIC, 0) as authentic, ifnull(pk.CIPHER_ID, 0) as cipherId "+
		"FROM td.personKey pk WHERE  pk.DEVICE_ID = ?", deviceId)

	if err != nil {
		return nil, err
	}

	pKey := new(model.PersonKey)
	for rows.Next() {
		err = rows.Scan(&pKey.KeyId, &pKey.DeviceId, &pKey.PanelCode, &pKey.Value, &pKey.BuyerId, &pKey.Authentic, &pKey.CipherId)
		if err != nil {
			return nil, err
		}
		pKeys = append(pKeys, *pKey)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return pKeys, nil
}

func (this *PersonKeyRepository) Save(pk *model.PersonKey) error {
	stmt, err := this.db.Prepare("INSERT INTO td.personKey (KEY_ID, BUYER_ID, VALUE, DEVICE_ID, FLAT_NUM, AUTHENTIC, CIPHER_ID) " +
		"VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE FLAT_NUM=?, AUTHENTIC=?, CIPHER_ID=?")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(pk.KeyId, pk.BuyerId, pk.Value, pk.DeviceId, pk.PanelCode, pk.Authentic, pk.CipherId, pk.PanelCode, pk.Authentic, pk.CipherId)

	if err != nil {
		return err
	}

	return nil
}

func (this *PersonKeyRepository) Delete(pk *model.PersonKey) error {
	stmt, err := this.db.Prepare("DELETE FROM td.personKey WHERE KEY_ID=?")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(pk.KeyId)

	if err != nil {
		return err
	}

	return nil
}
