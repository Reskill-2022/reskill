package repository

type Container struct {
	UserRepository *UserRepository
}

func NewContainer() *Container {
	return &Container{
		UserRepository: NewUserRepository(),
	}
}
