package internals

import (
	"fmt"
	"log"
	"reflect"
)

type ServiceEntry struct {
	Value any
	Type  reflect.Type
}

type DI struct {
	Services map[string]*ServiceEntry
}

//helper functions ------

func (d *DI) GetServiceType(key string) reflect.Type {
	return d.Services[key].Type
}

func (d *DI) PrintServicesInDI() {
	fmt.Println("Printing services in di container total : ", len(d.Services))
	for key := range d.Services {
		fmt.Println(key)
	}
}
func NewDI() *DI {
	return &DI{Services: make(map[string]*ServiceEntry)}
}

func (d *DI) AddService(key string, service any) {
	if reflect.TypeOf(service).Kind() != reflect.Ptr {
		log.Fatal("Pass Service instance as pointer not raw ")
	}
	d.Services[key] = &ServiceEntry{
		Value: service,
		Type:  reflect.TypeOf(service),
	}
}

func (d *DI) GetService(key string) (any, reflect.Type) {
	entry, ok := d.Services[key]
	if !ok {
		return nil, nil
	}
	return entry.Value, entry.Type
}
