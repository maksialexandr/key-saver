package model

type DataExecute struct {
	Action      string
	SrcDevice   *Device
	DstDevice   *Device
	SrcKeys     map[string]PersonKey // ключи источника
	DstKeys     map[string]PersonKey // ключи назначения
	CurrentKeys map[string]PersonKey // ключи которыми будем оперировать
}
