package controllers

type Container struct{}

func NewContainer() *Container {
	return &Container{}
}
