package form

import (
	"person-key-saver/internal/app/model"
)

type HttpPayload struct {
	SrcMac string
	DstMac string
	Key    model.PersonKey

	Action   string
	BuyerId  int
	CipherId int
}

func (this *HttpPayload) GetPayload() *Payload {
	var payload Payload
	payload.DstMac = this.DstMac
	payload.SrcMac = this.SrcMac
	payload.BuyerId = this.BuyerId
	payload.CipherId = this.CipherId
	payload.Keys = append(payload.Keys, this.Key)
	return &payload
}
