package demo

type ComponentY struct{}

type ComponentZ struct{}

type ServiceA struct {
	y *ComponentY
	z *ComponentZ
}

func NewComponentY() *ComponentY {
	return &ComponentY{}
}

func NewComponentZ() *ComponentZ {
	return &ComponentZ{}
}

func NewServiceA(y *ComponentY, z *ComponentZ) *ServiceA {
	return &ServiceA{y, z}
}
