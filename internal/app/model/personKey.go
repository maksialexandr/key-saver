package model

import (
	"fmt"
	"strings"
)

type PersonKey struct {
	KeyId     int
	DeviceId  int
	PanelCode int
	Value     string
	BuyerId   int
	Authentic bool
	CipherId  int
}

func (this *PersonKey) UnsetKeyId() {
	this.KeyId = 0
}

func (this *PersonKey) UnsetCipher() {
	this.CipherId = 0
}

func (this *PersonKey) PrepareValue() {
	this.Value = fmt.Sprintf("%014s", strings.TrimLeft(this.Value, "0"))
}
