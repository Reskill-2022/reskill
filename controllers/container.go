package controllers

type Container struct {
	UserController *UserController
}

func NewContainer() *Container {
	return &Container{
		UserController: NewUserController(),
	}
}
