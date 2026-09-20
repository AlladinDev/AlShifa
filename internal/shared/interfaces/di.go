package interfaces

import "reflect"

type IDependencyInjection interface {
	AddService(key string, service any)
	GetService(key string) (any, reflect.Type)
	GetServiceType(key string) reflect.Type
	Instantiatemodules()
	ModuleExists(moduleName string) bool
	PrintservicesInDI()
	PrintModulesRegistered()
	RegisterModule(moduleName string, constructorFn func(diContainer IDependencyInjection), dependencies []string)
	ServiceExists(key string) bool
}
