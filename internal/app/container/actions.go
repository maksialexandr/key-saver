package container

type Action = string

const (
	ACTION_LOAD   = Action("save")
	ACTION_DELETE = Action("delete")
	ACTION_ENCODE = Action("encode")
	ACTION_DECODE = Action("decode")
	ACTION_COPY   = Action("copy")
	ACTION_MOVE   = Action("move")
)
