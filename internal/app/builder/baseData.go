package builder

import (
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/form"
	"person-key-saver/internal/app/model"
)

type BaseDate struct {
	Keys         []model.PersonKey
	PhysicalKeys []model.PersonKey
	Device       *model.Device
	container    *container.Container
	payload      *form.Payload
}
