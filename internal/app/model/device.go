package model

type Device struct {
	DeviceId     int
	Mac          string
	Type         int
	Url          string
	Login        string
	Password     string
	SecurityMode int
	CipherId     int
	BuyerId      int
}
