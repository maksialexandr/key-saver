package container

type Token = string

type Container struct {
	items map[Token]interface{}
}

func New() *Container {
	return &Container{
		items: make(map[Token]interface{}),
	}
}

func (this *Container) Get(token Token) interface{} {
	if val, ok := this.items[token]; ok {
		return val
	}
	return nil
}

func (this *Container) Set(token Token, item interface{}) {
	this.items[token] = item
}
