package form

import (
	"person-key-saver/internal/app/model"
)

type Payload struct {
	SrcMac string
	DstMac string
	Keys   []model.PersonKey

	Action   string
	CipherId int
	BuyerId  int
	Delete   bool
	OnlyDb   bool
}

func (this *Payload) IsEmptySourceMac() bool {
	return this.SrcMac == ""
}

func (this *Payload) IsEmptyDestinationMac() bool {
	return this.DstMac == ""
}

func (this *Payload) IsEmptyCipher() bool {
	return this.CipherId == 0
}

func (this *Payload) IsEmptyBuyer() bool {
	return this.BuyerId == 0
}

func (this *Payload) IsEmptyKeys() bool {
	return len(this.Keys) == 0
}
