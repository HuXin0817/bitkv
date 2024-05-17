package main

type ServiceContext struct {
	Config Config
}

func NewServiceContext(c Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
